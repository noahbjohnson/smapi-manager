package backend

import (
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

var AppFs = afero.NewOsFs()

var sugar *zap.SugaredLogger

const (
	smapiSubfolder      = "smapi-internal"
	configDirFolder     = "smapi_manager"
	configFileName      = "config"
	configFileExtension = "json"

	gameDirKey  = "GameDir"
	firstRunKey = "FirstRun"
	hasSmapiKey = "HasSmapi"
)

func getConfigDir() string {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(userConfigDir, configDirFolder)
}

func getOrCreateConfigDir(fs afero.Fs) (directory afero.File, err error) {
	dirName := getConfigDir()
	hasConfigDir, err := afero.DirExists(fs, dirName)
	if err != nil {
		return
	}
	if !hasConfigDir {
		err = fs.Mkdir(dirName, 0777)
		if err != nil {
			return
		}
	}
	return fs.Open(dirName)
}

func getGameDirectory() (path string) {
	userConfigDir, _ := os.UserConfigDir()
	userDir, _ := os.UserHomeDir()
	switch runtime.GOOS {
	case "darwin":
		if exists, _ := afero.Exists(AppFs, "/Applications/Stardew Valley.app"); exists {
			path = "/Applications/Stardew Valley.app/Contents/MacOS"
		} else {
			path = filepath.Join(userConfigDir, "Steam/SteamApps/common/Stardew Valley/Contents/MacOS")
		}
	case "android":
		path = "storage/emulated/0/StardewValley/"
	case "linux":
		if exists, _ := afero.Exists(AppFs, "GOGGames/StardewValley"); exists {
			path = filepath.Join(userDir, "GOGGames/StardewValley/game")
		} else {
			path = filepath.Join(userDir, ".local/share/Steam/steamapps/common/Stardew Valley")
		}
	case "windows":
		if exists, _ := afero.Exists(AppFs, "/GOG Games"); exists {
			path = "/GOG Games/Stardew Valley"
		} else if exists, _ := afero.Exists(AppFs, "/Program Files (x86)/GOG Galaxy/Games/Stardew Valley"); exists {
			path = "/Program Files (x86)/GOG Galaxy/Games/Stardew Valley"
		} else {
			path = "/Program Files (x86)/Steam/steamapps/common/Stardew Valley"
		}
	}
	return path
}

func hasSMAPI() (bool, error) {
	smapiFolder := filepath.Join(getGameDirectory(), smapiSubfolder)
	sugar.Debug("checking for smapi at: ", smapiFolder)
	return afero.Exists(AppFs, smapiFolder)
}

func initializeViper() (err error) {
	configDir := getConfigDir()
	configFilePath := filepath.Join(configDir, configFileName) + "." + configFileExtension
	viper.SetConfigName(configFileName)
	viper.SetConfigType(configFileExtension)
	viper.AddConfigPath(configDir)

	defaultGameDir := getGameDirectory()

	smapi, err := hasSMAPI()
	if err != nil {
		return err
	}

	viper.Set(hasSmapiKey, smapi)
	viper.SetDefault(gameDirKey, defaultGameDir)

	if configExists, _ := afero.Exists(AppFs, configFilePath); configExists {
		sugar.Debug("Config exists, reading from file")
		viper.SetDefault(firstRunKey, false)
		err = viper.ReadInConfig()
		if err != nil {
			return err
		}
		sugar.Debug("values loaded: ", viper.AllSettings())
		err = viper.WriteConfig()
		if err != nil {
			return err
		}
	} else {
		sugar.Debug("Config file not found, creating file")
		viper.SetDefault(firstRunKey, true)
		err = viper.SafeWriteConfig()
		if err != nil {
			return err
		}
	}
	return
}

// setup the config and such
func Initialize() string {
	log.Println("Initializing the logger")
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalln("Could not start logger", err)
	}
	defer logger.Sync()
	sugar = logger.Sugar()

	sugar.Info("Loading config directory")
	configDir, err := getOrCreateConfigDir(AppFs)
	if err != nil {
		sugar.Fatal("Could not open config directory", err)
	}

	sugar.Info("Initializing config")
	err = initializeViper()
	if err != nil {
		sugar.Fatal("Fatal error reading config file", err)
	}
	return configDir.Name()
}
package backend

import (
	"fmt"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

var AppFs = afero.NewOsFs()

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
	Sugar.Debug("checking for smapi at: ", smapiFolder)
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
		Sugar.Debug("Config exists, reading from file")
		viper.SetDefault(firstRunKey, false)
		err = viper.ReadInConfig()
		if err != nil {
			return err
		}
		Sugar.Debug("values loaded: ", viper.AllSettings())
		err = viper.WriteConfig()
		if err != nil {
			return err
		}
	} else {
		Sugar.Debug("Config file not found, creating file")
		viper.SetDefault(firstRunKey, true)
		err = viper.SafeWriteConfig()
		if err != nil {
			return err
		}
	}
	return
}

// Initialize sets up the config and such
func Initialize() string {
	log.Println("Initializing the logger")
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalln("Could not start logger", err)
	}
	defer logger.Sync()
	Sugar = logger.Sugar()

	Sugar.Info("Loading config directory")
	configDir, err := getOrCreateConfigDir(AppFs)
	if err != nil {
		Sugar.Fatal("Could not open config directory", err)
	}

	Sugar.Info("Initializing config")
	err = initializeViper()
	if err != nil {
		Sugar.Fatal("Fatal error reading config file", err)
	}
	return configDir.Name()
}

// HasSMAPI returns if SMAPI is installed
func HasSMAPI() bool {
	return viper.GetBool(hasSmapiKey)
}

// GameDir returns the game dir from the config
func GameDir() string {
	return viper.GetString(gameDirKey)
}

func GetSMAPI(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	var responseCode uint8 = 0
	if HasSMAPI() {
		responseCode = 1
	}
	Sugar.Infof("responding to smapi call with status %d", responseCode)
	_, _ = fmt.Fprintf(w, "%d", responseCode)
}

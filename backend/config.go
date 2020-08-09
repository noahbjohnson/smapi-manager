package backend

import (
	"encoding/json"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"net/http"
	"path/filepath"
)

func initializeViper() (err error) {
	configDir := getConfigPathString()
	configFilePath := filepath.Join(configDir, configFileName) + "." + configFileExtension
	viper.SetConfigName(configFileName)
	viper.SetConfigType(configFileExtension)
	viper.AddConfigPath(configDir)

	defaultGameDir := getGameDirectory()

	smapi, err := HasSMAPI()
	if err != nil {
		return err
	}

	viper.Set(hasSmapiKey, smapi)
	viper.SetDefault(gameDirKey, defaultGameDir)
	viper.SetDefault(modsKey, []Mod{})

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

// GameDir returns the game dir from the config
func GameDir() string {
	return viper.GetString(gameDirKey)
}

func EnumerateMods(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")

	err := json.NewEncoder(w).Encode(LoadMods(GameDir()))
	if err != nil {
		panic(err)
	}
}

// InitializeConfig sets up the config and such
func InitializeConfig() string {
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

	//Sugar.Info("Loading Mods")
	//mods := LoadMods(GameDir())

	return configDir.Name()
}

package backend

import (
	"github.com/spf13/afero"
	"go.uber.org/zap"
	"net/http"
	"path/filepath"
)

const (
	smapiSubfolder          = "smapi-internal"
	configDirFolder         = "smapi_manager"
	configFileName          = "config"
	configFileExtension     = "json"
	gameDirKey              = "GameDir"
	modsKey                 = "Mods"
	firstRunKey             = "FirstRun"
	hasSmapiKey             = "HasSmapi"
	ModsSubPath             = "Mods"
	modSubDir               = "smapi_manager"
	manifestFileName        = "manifest.json"
	mapTrailingCommaRegex   = ",(\\s*})"
	arrayTrailingCommaRegex = ",(\\s*])"
	leadingSpaceRegex       = "\\A(\\s*)"
	dosNewLineRegex         = "\r\n"
	endsWithDigitRegex      = "^.*\\d$"
	disablePrefix           = "."
	disableSuffix           = "."
	nbsp                    = "\ufeff"
	unixNewLine             = "\n"
	emptyString             = ``
	firstCaptureGroup       = `$1`
)

var (
	Sugar *zap.SugaredLogger
	AppFs = afero.NewOsFs()
)

func OpenSmapiInstall() error {
	return UrlOpener("https://stardewvalleywiki.com/Modding:Player_Guide/Getting_Started#Install_SMAPI")()
}

func HasSMAPI() (bool, error) {
	smapiFolder := filepath.Join(getGameDirectory(), smapiSubfolder)
	Sugar.Debug("checking for smapi at: ", smapiFolder)
	return afero.Exists(AppFs, smapiFolder)
}

func setupRoutes() error {
	http.HandleFunc("/upload", UploadFile)
	http.HandleFunc("/mods", EnumerateMods)
	//http.Handle("/", http.FileServer(frontendFS))
	return nil
}

// StartAPI sets up the api routes and starts API server in background
func StartAPI(addr string) error {
	err := setupRoutes()
	if err != nil {
		return err
	}

	go func() {
		err := http.ListenAndServe(addr, nil)
		if err != nil {
			panic(err)
		}
	}()

	return nil
}

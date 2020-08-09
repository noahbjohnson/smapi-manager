package backend

import (
	"github.com/rakyll/statik/fs"
	"github.com/spf13/afero"
	"github.com/webview/webview"
	"go.uber.org/zap"
	"log"
	"net/http"
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

func BindFunctions(wv webview.WebView) (err error) {
	err = wv.Bind("openSmapiInstall", UrlOpener("https://stardewvalleywiki.com/Modding:Player_Guide/Getting_Started#Install_SMAPI"))
	if err != nil {
		return err
	}
	err = wv.Bind("hasSmapi", HasSMAPI)
	return err
}

func setupRoutes() {
	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/upload", UploadFile)
	http.HandleFunc("/mods", EnumerateMods)
	http.Handle("/", http.FileServer(statikFS))
}

func StartAPI(addr string) {

	Initialize()
	setupRoutes()

	// Start API server in background
	go func() {
		err := http.ListenAndServe(addr, nil)
		if err != nil {
			panic(err)
		}
	}()
}

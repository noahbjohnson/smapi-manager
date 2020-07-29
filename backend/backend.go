package backend

import (
	"github.com/spf13/afero"
	"go.uber.org/zap"
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

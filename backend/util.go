package backend

import (
	"github.com/spf13/afero"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

// fixJSON fixes trailing commas and encoding issues in JSON
func fixJSON(inJson []byte) []byte {
	var mapTrailer = regexp.MustCompile(mapTrailingCommaRegex)
	var arrayTrailer = regexp.MustCompile(arrayTrailingCommaRegex)
	var leadingSpace = regexp.MustCompile(leadingSpaceRegex)
	var newline = regexp.MustCompile(dosNewLineRegex)

	jsonString := string(inJson[:])
	jsonString = mapTrailer.ReplaceAllString(jsonString, firstCaptureGroup)
	jsonString = arrayTrailer.ReplaceAllString(jsonString, firstCaptureGroup)
	jsonString = leadingSpace.ReplaceAllString(jsonString, emptyString)
	// DOS newlines
	jsonString = newline.ReplaceAllString(jsonString, unixNewLine)
	// byte order mark
	jsonString = strings.TrimLeft(jsonString, nbsp)
	return []byte(jsonString)
}

func getConfigPathString() string {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(userConfigDir, configDirFolder)
}

func getOrCreateConfigDir(fs afero.Fs) (directory afero.File, err error) {
	dirName := getConfigPathString()
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
	zipDir := filepath.Join(dirName, "zips/")
	hasZipDir, err := afero.DirExists(fs, zipDir)
	if err != nil {
		return
	}
	if !hasZipDir {
		err = fs.Mkdir(zipDir, 0777)
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

func appendManifestFilePath(dir string) string {
	return filepath.Join(dir, manifestFileName)
}

func isModDir(dirName string) bool {
	if !endsWithNumber(dirName) {
		dir, _ := afero.ReadDir(AppFs, dirName)
		for i := range dir {
			if dir[i].Name() == manifestFileName {
				return true
			}
		}
	}
	return false
}

func endsWithNumber(dirName string) bool {
	digitTest, err := regexp.Compile(endsWithDigitRegex)
	if err != nil {
		log.Fatal(err)
	}
	return digitTest.MatchString(dirName)
}

func isEnabled(dirName string) bool {
	dir := filepath.Dir(dirName)
	if strings.HasPrefix(dir, disablePrefix) {
		return false
	}
	return true
}

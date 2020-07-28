package backend

import (
	"encoding/json"
	"github.com/karrick/godirwalk"
	"github.com/spf13/afero"
	"log"
	"path/filepath"
	"regexp"
	"strings"
)

const manifestFileName = "manifest.json"

// For a content pack, ContentPackFor specifies which mod can read it.
type ContactPackRef struct {
	UniqueID       string `json:"UniqueID"`
	MinimumVersion string `json:"MinimumVersion"`
}

// The Dependencies field specifies other mods required to use this mod.
type DependencyRef struct {
	UniqueID string `json:"UniqueID"`
	// optional. If specified, older versions won't meet the requirement.
	MinimumVersion string `json:"MinimumVersion"`
	// You can mark a dependency as optional. It will be loaded first if it's installed, otherwise it'll be ignored.
	IsRequired bool `json:"IsRequired"`
}

// Every SMAPI mod or content pack must have a manifest.json file in its folder.
// SMAPI uses this to identify and load the mod, perform update checks, etc.
type ModMetadata struct {
	// The mod name. SMAPI uses this in player messages, logs, and errors.
	Name string `json:"Name"`
	// The name of the person who created the mod. Ideally this should include the username used to publish mods.
	Author string `json:"Author"`
	// The mod's semantic version.
	Version string `json:"Version"`
	// A short explanation of what your mod does (one or two sentences), shown in the SMAPI log
	Description string `json:"Description"`
	// A unique identifier for your mod. The recommended format is <your name>.<mod name>, with no spaces or
	//  special characters. SMAPI uses this for update checks, mod dependencies, and compatibility blacklists
	//  (if the mod breaks in a future version of the game). When another mod needs to reference this mod,
	//  it uses the unique ID.
	UniqueID string `json:"UniqueID"`
	// All mods must specify either EntryDll (for a SMAPI mod) or ContentPackFor (for a content pack).
	//  These are mutually exclusive — you can't specify both.
	//  For a SMAPI mod, EntryDll is the mod's compiled DLL filename in its mod folder.
	EntryDll string `json:"EntryDll"` // not present in content packs
	// For a content pack, ContentPackFor specifies which mod can read it. The MinimumVersion is optional.
	ContentPackFor ContactPackRef `json:"ContentPackFor"` // not present in mods
	// The MinimumApiVersion fields sets the minimum SMAPI version needed to use this mod.
	//  If a player tries to use the mod with an older SMAPI version, they'll see a friendly message
	//  saying they need to update SMAPI. This also serves as a proxy for the minimum game version,
	//  since SMAPI itself enforces a minimum game version.
	MinimumApiVersion string `json:"MinimumApiVersion"`
	// The Dependencies field specifies other mods required to use this mod.
	//  If a player tries to use the mod and the dependencies aren't installed,
	//  the mod won't be loaded and they'll see a friendly message saying they need to
	//  install those
	Dependencies []DependencyRef `json:"Dependencies"`
	UpdateKeys   []string        `json:"UpdateKeys"`
}

type Mod struct {
	Directory string      `json:"directory"`
	Enabled   bool        `json:"enabled"`
	Metadata  ModMetadata `json:"metadata"`
}

func (m *Mod) enable() error {
	currentDir := filepath.Dir(m.Directory)
	if strings.HasPrefix(currentDir, ".") {
		newDir := strings.TrimLeft(currentDir, ".")
		if strings.HasSuffix(currentDir, ".") {
			newDir = strings.TrimRight(newDir, ".")
		}
		parent := strings.TrimRight(m.Directory, currentDir)
		newPath := filepath.Join(parent, newDir)
		err := AppFs.Rename(m.Directory, newPath)
		if err != nil {
			return err
		}
		m.Enabled = true
	}
	return nil
}

func (m *Mod) disable() error {
	currentDir := filepath.Dir(m.Directory)
	if !strings.HasPrefix(currentDir, ".") {
		newPath := filepath.Join(strings.TrimRight(m.Directory, currentDir), "."+currentDir)
		err := AppFs.Rename(m.Directory, newPath)
		if err != nil {
			return err
		}
		m.Enabled = false
	}
	return nil
}

func appendManifestFilePath(dir string) string {
	return filepath.Join(dir, manifestFileName)
}

func isModDir(dirName string) bool {
	dir, _ := afero.ReadDir(AppFs, dirName)
	if endsWithNumber(dirName) {
		return false
	}
	for i := range dir {
		if dir[i].Name() == manifestFileName {
			return true
		}
	}
	return false
}

func endsWithNumber(dirName string) bool {
	digitTest, err := regexp.Compile("^.*\\d$")
	if err != nil {
		log.Fatal(err)
	}
	return digitTest.MatchString(dirName)
}

func isEnabled(dirName string) bool {
	dir := filepath.Dir(dirName)
	if strings.HasPrefix(dir, ".") {
		return false
	}
	return true
}

func getModDirs(modsDir string) []string {
	var modDirs []string
	err := godirwalk.Walk(modsDir, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			if de.IsDir() && isModDir(osPathname) {
				modDirs = append(modDirs, osPathname)
				return filepath.SkipDir
			}
			return nil
		},
		Unsorted: true,
	})
	if err != nil {
		log.Fatal(err)
	}
	return modDirs
}

func fixJSON(inJson []byte) []byte {
	var mapTrailer = regexp.MustCompile(",(\\s*})")
	var arrayTrailer = regexp.MustCompile(",(\\s*])")
	var leadingSpace = regexp.MustCompile("\\A(\\s*)")
	var newline = regexp.MustCompile("\r\n")

	jsonString := string(inJson[:])
	jsonString = mapTrailer.ReplaceAllString(jsonString, `$1`)
	jsonString = arrayTrailer.ReplaceAllString(jsonString, `$1`)
	jsonString = leadingSpace.ReplaceAllString(jsonString, ``)
	// DOS newlines
	jsonString = newline.ReplaceAllString(jsonString, "\n")
	// byte order mark
	jsonString = strings.TrimLeft(jsonString, "\ufeff")
	return []byte(jsonString)
}

func loadMods(modDirs []string) (mods []Mod) {
	for _, dir := range modDirs {
		jsonFilePath := appendManifestFilePath(dir)
		jsonFile, err := afero.ReadFile(AppFs, jsonFilePath)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("parsing ", jsonFilePath, " with length ", len(jsonFile))
		cleanJSON := fixJSON(jsonFile)
		metadata := ModMetadata{}
		err = json.Unmarshal(cleanJSON, &metadata)
		if err != nil {
			log.Printf("error decoding json manifest: %v", err)
			if e, ok := err.(*json.SyntaxError); ok {
				log.Printf("syntax error at byte offset %d", e.Offset)
			}
			panic(err)
		}
		mods = append(mods, Mod{
			Directory: dir,
			Enabled:   isEnabled(dir),
			Metadata:  metadata,
		})
	}
	return mods
}

func LoadMods(gameDir string) []Mod {
	modDir := filepath.Join(gameDir, "Mods/")
	modDirs := getModDirs(modDir)
	mods := loadMods(modDirs)
	return mods
}
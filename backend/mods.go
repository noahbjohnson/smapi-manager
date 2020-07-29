package backend

import (
	"encoding/json"
	"github.com/karrick/godirwalk"
	"github.com/spf13/afero"
	"log"
	"path/filepath"
	"strings"
)

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

// Mod is a a SMAPI mod
type Mod struct {
	Directory string      `json:"directory"`
	Enabled   bool        `json:"enabled"`
	Managed   bool        `json:"managed"`
	Metadata  ModMetadata `json:"metadata"`
}

func (m *Mod) enable() error {
	currentDir := filepath.Dir(m.Directory)
	if strings.HasPrefix(currentDir, disablePrefix) && m.Managed {
		newDir := strings.TrimLeft(currentDir, disablePrefix)
		if strings.HasSuffix(currentDir, disableSuffix) {
			newDir = strings.TrimRight(newDir, disableSuffix)
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
	if !strings.HasPrefix(currentDir, disablePrefix) && m.Managed {
		newPath := filepath.Join(strings.TrimRight(m.Directory, currentDir), disablePrefix+currentDir)
		err := AppFs.Rename(m.Directory, newPath)
		if err != nil {
			return err
		}
		m.Enabled = false
	}
	return nil
}

func (m *Mod) RefreshEnabled() {
	m.Enabled = isEnabled(m.Directory)
}

type ModManager struct {
	Mods    []Mod  `json:"mods"`
	GameDir string `json:"game_dir"`
}

func (m *ModManager) LoadModsFromFS() {
	m.Mods = LoadMods(m.GameDir)
}

func (m *ModManager) RefreshEnabled() {
	for _, mod := range m.Mods {
		mod.RefreshEnabled()
	}
}

func (m *ModManager) modsDir() {
	for _, mod := range m.Mods {
		mod.RefreshEnabled()
	}
}

func loadMods(modsDir string) (mods []Mod) {
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

// LoadMods searches the provided game directory for mods and parses them into Mod structures
func LoadMods(gameDir string) []Mod {
	modDir := filepath.Join(gameDir, ModsSubPath)
	mods := loadMods(modDir)
	return mods
}

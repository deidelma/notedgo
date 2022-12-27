package noted

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

const defaultConfigurationFile = ".config/noted/noted.json"

// configuration state
var Configuration struct {
	Path string
	Data *ConfigData
}

// initialize the configuration data and path to their defaults
func init() {
	home, _ := os.UserHomeDir()
	Configuration.Path = filepath.Join(home, defaultConfigurationFile)
	Configuration.Data = MakeDefaultConfiguration()
}

type ConfigData struct {
	NotesPath    string `json:"notes_path"`
	DatabasePath string `json:"database_path"`
	Initialized  bool   `json:"initialized"`
	Version      string `json:"version"`
	Autosave     bool   `json:"autosave"`
	UseGui       bool   `json:"use_gui"`
}

func MakeDefaultConfiguration() *ConfigData {
	cwd, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("unable to determine working directory:%v", err))
	}
	dbpath := cwd + string(os.PathSeparator) + "noted.sqlite3"
	return &ConfigData{NotesPath: cwd, DatabasePath: dbpath, Initialized: true, Version: "0.5.0", Autosave: true, UseGui: false}
}

func ConfigurationFileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		panic("unexpected error in test of configuration")
	}
	return true
}

func WriteConfigurationFile(path string, config *ConfigData) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	// does the directory exist, if not create it
	dir := filepath.Dir(path)
	_, err = os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(dir, os.ModeDir)
			if err != nil {
				panic("Unexpected error while trying to create directories: " + dir)
			}
		}
	}
	// write the file
	err = os.WriteFile(path, data, 0644)
	return err
}

func ReadConfigurationFile(path string) (*ConfigData, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	c := &ConfigData{}
	err = json.Unmarshal(data, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// LoadConfiguration loads the data from the current configuration path
func LoadConfiguration() {
	LoadConfigurationFromPath(Configuration.Path)
}

// LoadConfigurationFromPath loads the configuration data from the provided path
func LoadConfigurationFromPath(path string) {
	c, err := ReadConfigurationFile(path)
	if err != nil {
		log.Err(err).Msg("unable to read configuration file")
		os.Exit(2)
	}
	Configuration.Data = c
}

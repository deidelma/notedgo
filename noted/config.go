package noted

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
)

const defaultConfigurationFile = ".config/noted/noted.json"

// State holds the configuration information for the project
var State struct {
	Path  string
	Data  *ConfigData
	Debug bool
}

// initialize the configuration data and path to their defaults
func init() {
	home, _ := os.UserHomeDir()
	State.Path = filepath.Join(home, defaultConfigurationFile)
	State.Data = MakeDefaultConfiguration()
	State.Debug = false
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
	return &ConfigData{
		NotesPath:    cwd,
		DatabasePath: dbpath,
		Initialized:  true,
		Version:      "0.5.0",
		Autosave:     true,
		UseGui:       false,
	}
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
// or uses custom path based on command line and environment variables.
//
// Returns err if invalid data passed via the command line or environment variables.
func LoadConfiguration() error {
	debug := false
	configPath := ""
	ProcessEnvironment(&debug, &configPath)
	ParseCommandLine(&debug, &configPath)
	if debug {
		handleDebug()
	}
	if configPath != "" {
		log.Info().Msg("using custom configuration file")
		LoadConfigurationFromPath(configPath)
	} else {
		LoadConfigurationFromPath(State.Path)
	}
	return nil
}

// ParseCommandLine checks command line parameters for overrides of environment variables and defaults.
func ParseCommandLine(debug *bool, path *string) {
	flag.BoolVar(debug, "debug", false, "set debug mode")
	flag.StringVar(path, "config", "", "custom path for config file")
	flag.Parse()
}

// ProcessEnvironment checks for program specific environment variables.  If found, the provided
// pointers will be set to the appropriate values.
//
//	debug -- switch to debug mode for logging and other features.
//	path -- set a custom path for the configuration file.
func ProcessEnvironment(debug *bool, path *string) {
	for _, variables := range os.Environ() {
		pair := strings.SplitN(variables, "=", 2)
		key := strings.ToLower(pair[0])
		log.Info().Msg(fmt.Sprintf("%s : %s", pair[0], pair[1]))
		if key == "NOTED_CONFIG" {
			*path = pair[1]
		} else if key == "NOTED_DEBUG" || key == "DEBUG" {
			*debug = true
		}
	}
}

func handleDebug() {
	log.Info().Msg("handling debug")
}

// LoadConfigurationFromPath loads the configuration data from the provided path
func LoadConfigurationFromPath(path string) {
	c, err := ReadConfigurationFile(path)
	if err != nil {
		log.Err(err).Msg("unable to read configuration file")
		os.Exit(2)
	}
	State.Data = c
}

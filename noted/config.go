package noted

import (
	"bytes"
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

type CmdLineData struct {
	Debug      bool
	ConfigPath string
	HelpMsg    string
	Args       []string
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
			err := os.MkdirAll(dir, 0777)
			if err != nil {
				panic("Unexpected error while trying to create directories: " + dir)
			}
		}
	}
	// write the file
	err = os.WriteFile(path, data, 0777)
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
	cmdData := CmdLineData{}
	ProcessEnvironment(&cmdData.Debug, &cmdData.ConfigPath)
	//err := ParseCmdLine(os.Args[0], os.Args[1:], &debug, &configPath, &output)
	err := ParseCmdLine(os.Args[0], os.Args[1:], &cmdData)
	if err != nil {
		log.Err(err).Msg("unable to parse command line")
		return err
	}
	if cmdData.Debug {
		handleDebug()
	}
	if cmdData.ConfigPath != "" {
		log.Info().Msg("using custom configuration file")
		LoadConfigurationFromPath(cmdData.ConfigPath)
	} else {
		LoadConfigurationFromPath(State.Path)
	}
	return nil
}

// ParseCmdLine checks command line parameters for overrides of environment variables and defaults.
//
// Handles parsing using the approach proposed by
// eliben (see https://github.com/eliben/code-for-blog/blob/master/2020/go-testing-flags/main.go)
// This layer over flag helps make testing feasible
func ParseCmdLine(progname string, args []string, cmdData *CmdLineData) (err error) {
	flags := flag.NewFlagSet(progname, flag.ContinueOnError)
	var buf bytes.Buffer
	flags.SetOutput(&buf)

	flags.BoolVar(&cmdData.Debug, "debug", false, "set debug mode")
	flags.StringVar(&cmdData.ConfigPath, "config", "", "custom path for config file")
	err = flags.Parse(args)
	cmdData.HelpMsg = buf.String()
	cmdData.Args = flags.Args()
	if err != nil {
		return err
	}
	return nil
}

// ProcessEnvironment checks for program specific environment variables.  If found, the provided
// pointers will be set to the appropriate values.
//
//	debug -- switch to debug mode for logging and other features.
//	path -- set a custom path for the configuration file.
func ProcessEnvironment(debug *bool, path *string) {
	*path = os.Getenv("NOTED_CONFIG")
	*debug = strings.ToLower(os.Getenv("NOTED_DEBUG")) == "true"
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

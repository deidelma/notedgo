package noted

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const defaultConfigurationFile = ".config/noted/noted.json"

// ConfigState holds the configuration information for the project
var ConfigState struct {
	Path  string
	Data  *ConfigData
	Debug bool
}

// initialize the configuration data and path to their defaults
func init() {
	home, _ := os.UserHomeDir()
	ConfigState.Path = filepath.Join(home, defaultConfigurationFile)
	ConfigState.Data = MakeDefaultConfiguration()
	ConfigState.Debug = false
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

// ConfigData holds configuration settings loaded from the configuration file.
type ConfigData struct {
	NotesPath    string `json:"notes_path"`
	DatabasePath string `json:"database_path"`
	Initialized  bool   `json:"initialized"`
	Version      string `json:"version"`
	Autosave     bool   `json:"autosave"`
	UseGui       bool   `json:"use_gui"`
}

// CmdLineData holds the data returned from flag.Process()
type CmdLineData struct {
	Debug      bool
	ConfigPath string
	HelpMsg    string
	Args       []string
}

// MakeDefaultConfiguration creates a pointer to a new instance of ConfigData
// initialized to the hardwired system defaults.
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

// ConfigurationFileExists returns true if there is a
// configuration file at the provided location.  Assumes that
// the path is fully qualified and includes the name of the
// file itself.
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

// WriteConfigurationFile attempts to write the configuration file
// at the provided path using the data provide.  Currently,
// writes a json file.
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

// ReadConfigurationFile attempts to load the configuration file
// from the provided path.
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
// programName is usually os.Args[0]
//
// args is usually os.Args[1:]
//
// Returns err if invalid data passed via the command line or environment variables.
func LoadConfiguration(programName string, args []string) error {
	cmdData := CmdLineData{}
	// first load environment variables
	ProcessEnvironment(&cmdData.Debug, &cmdData.ConfigPath)
	// then check the command line for overrides
	err := ParseCmdLine(programName, args, &cmdData)
	if err != nil {
		log.Err(err).Msg("unable to parse command line")
		return err
	}
	// set up debug mode if necessary
	handleDebug(cmdData.Debug)

	// check for custom config path
	if cmdData.ConfigPath != "" {
		// use custom configuration path if provided
		log.Debug().Msg(fmt.Sprintf("using custom configuration file: %s", cmdData.ConfigPath))
		LoadConfigurationFromPath(cmdData.ConfigPath)
	} else {
		LoadConfigurationFromPath(ConfigState.Path)
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

	flags.BoolVar(&cmdData.Debug, "debug", cmdData.Debug, "set debug mode")
	flags.StringVar(&cmdData.ConfigPath, "config", cmdData.ConfigPath, "custom path for config file")
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

func handleDebug(debugOn bool) {
	if debugOn {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}

// LoadConfigurationFromPath loads the configuration data from the provided path
func LoadConfigurationFromPath(path string) {
	c, err := ReadConfigurationFile(path)
	if err != nil {
		log.Err(err).Msg("unable to read configuration file")
		os.Exit(2)
	}
	ConfigState.Data = c
}

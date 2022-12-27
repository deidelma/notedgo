package noted

import (
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var mockDirectory string

//goland:noinspection ALL
func setupMockDirectory(t *testing.T) func(t *testing.T) {
	mock, _ := os.Getwd()

	mockDirectory = filepath.Join(mock, "mock") //mock + string(os.PathSeparator) + "mock"
	_, err := os.Stat(mockDirectory)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(mockDirectory, 0755)
			if err != nil {
				panic("unable to create mock directory")
			}
		} else {
			panic("unexpected error for mockDirectory")
		}
	}
	return func(t *testing.T) {
		err := os.RemoveAll(mockDirectory)
		if err != nil {
			log.Err(err)
			t.Fail()
		}
	}
}

func TestMakeDefault(t *testing.T) {

	cwd, _ := os.Getwd()
	c := MakeDefaultConfiguration()
	assert.Equal(t, cwd, c.NotesPath)
	assert.True(t, strings.HasSuffix(c.DatabasePath, "noted.sqlite3"))
	assert.True(t, c.Initialized)
	assert.False(t, c.UseGui)
	assert.True(t, c.Autosave)
}

func TestConfigurationFileExists(t *testing.T) {
	assert.False(t, ConfigurationFileExists("./notaconfig.json"))
}

func TestWriteConfigurationFile(t *testing.T) {
	takeDownMockDirectory := setupMockDirectory(t)
	defer takeDownMockDirectory(t)

	path := filepath.Join(mockDirectory, "noted.json")
	conf := MakeDefaultConfiguration()
	err := WriteConfigurationFile(path, conf)
	if err != nil {
		log.Err(err)
		t.Fail()
	}
	assert.True(t, ConfigurationFileExists(path))
}
func TestWriteConfigurationPath(t *testing.T) {
	takeDownMockDirectory := setupMockDirectory(t)
	defer takeDownMockDirectory(t)

	path := filepath.Join(mockDirectory, "noted", "noted.json")
	conf := MakeDefaultConfiguration()
	err := WriteConfigurationFile(path, conf)
	if err != nil {
		log.Err(err)
		t.Fail()
	}

	assert.True(t, ConfigurationFileExists(path))
}

func TestReadConfigurationFile(t *testing.T) {
	takeDownMockDirectory := setupMockDirectory(t)
	defer takeDownMockDirectory(t)

	path := mockDirectory + string(os.PathSeparator) + "noted.json"
	conf := MakeDefaultConfiguration()
	err := WriteConfigurationFile(path, conf)
	if err != nil {
		log.Err(err)
		t.Fail()
	}

	c, err := ReadConfigurationFile(path)
	if err != nil {
		t.Errorf("error reading config file from disk: %v", err)
	}
	cwd, _ := os.Getwd()
	assert.Equal(t, cwd, c.NotesPath)
	assert.True(t, strings.HasSuffix(c.DatabasePath, "noted.sqlite3"))
	assert.True(t, c.Initialized)
	assert.False(t, c.UseGui)
	assert.True(t, c.Autosave)
}

func TestProcessEnvironment(t *testing.T) {
	debug := false
	path := ""
	err := os.Setenv("NOTED_CONFIG", "a/b/c")
	if err != nil {
		log.Err(err).Msg("unable to set environment variable")
	} else {
		log.Info().Msg("able to set NOTED_CONFIG")
	}
	defer os.Unsetenv("NOTED_CONFIG")
	_ = os.Setenv("NOTED_DEBUG", "true")
	defer os.Unsetenv("NOTED_DEBUG")
	ProcessEnvironment(&debug, &path)
	assert.Equal(t, true, debug)
	assert.Equal(t, "a/b/b", path)
}

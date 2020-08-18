package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mjslabs/ggft/internal/testhelpers"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var testConfig = "test.yml"
var testDir = "testDir"
var configContents = `
cache: ` + testDir + `
`

func TestPkgConfig(t *testing.T) {
	t.Run("InitializeConfig", testInitializeConfig)
	t.Run("ReadConfig", testReadConfig)
	t.Run("ReadConfigNoExist", testReadConfigNoExist)
	t.Run("ValidateConfig", testValidateConfigValid)
	t.Run("ValidateConfigInvalid", testValidateConfigInvalid)
	t.Run("ValidateConfigNoPerms", testValidateConfigValidNoPerms)
}

func testInitializeConfig(t *testing.T) {
	assert.NoError(t, initializeConfig(testDir, testDir, filepath.Join(testDir, testConfig)))
	_, err := os.Stat(filepath.Join(testDir, testConfig))
	assert.NoError(t, err)
	assert.NoError(t, os.RemoveAll(testDir))
}

func testReadConfig(t *testing.T) {
	assert.NoError(t, testhelpers.CreateFileWithContents(testConfig, configContents))
	assert.NoError(t, readConfig(testConfig))
	assert.NoError(t, os.Remove(testConfig))
}

func testReadConfigNoExist(t *testing.T) {
	assert.Error(t, readConfig(testConfig))
}

func testValidateConfigInvalid(t *testing.T) {
	assert.NoError(t, testhelpers.CreateFileWithContents(testConfig, configContents))
	viper.SetConfigFile(testConfig)
	assert.NoError(t, viper.ReadInConfig())
	assert.Error(t, validateConfig(viper.GetViper()))
	assert.NoError(t, os.Remove(testConfig))
}

func testValidateConfigValid(t *testing.T) {
	assert.NoError(t, testhelpers.CreateFileWithContents(testConfig, configContents))
	assert.NoError(t, os.Mkdir(testDir, 0755))
	viper.SetConfigFile(testConfig)
	assert.NoError(t, viper.ReadInConfig())
	assert.NoError(t, validateConfig(viper.GetViper()))
	assert.NoError(t, os.Remove(testConfig))
	assert.NoError(t, os.RemoveAll(testDir))
}

func testValidateConfigValidNoPerms(t *testing.T) {
	assert.NoError(t, testhelpers.CreateFileWithContents(testConfig, configContents))
	assert.NoError(t, os.Mkdir(testDir, 0400))
	assert.NoError(t, os.Chmod(testDir, 0400)) // for Windows
	viper.SetConfigFile(testConfig)
	assert.NoError(t, viper.ReadInConfig())
	assert.Error(t, validateConfig(viper.GetViper()))
	assert.NoError(t, os.Remove(testConfig))
	assert.NoError(t, os.RemoveAll(testDir))
}

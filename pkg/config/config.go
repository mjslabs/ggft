package config

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"golang.org/x/sys/unix"
)

type configFileDefault struct {
	Cache string // path to template cache
}

var defaultConfigDir = filepath.Join("~", ".ggft")
var defaultConfigCacheDir = "templates"
var defaultConfigFile = filepath.Join(defaultConfigDir, "config.yml")
var defaults = configFileDefault{
	Cache: filepath.Join(defaultConfigDir, defaultConfigFile),
}

func readConfig(config string) error {
	cfgFile, err := homedir.Expand(config)
	if err != nil {
		return err
	}

	// If a config file is found, read it in.
	viper.SetConfigFile(cfgFile)
	if err := viper.ReadInConfig(); err != nil {
		return err
		// TODO handle config not found differently from other errors
		// if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		// }
	}

	return nil
}

func validateConfig(v *viper.Viper) error {
	cache := v.GetString("cache")
	if _, err := os.Stat(cache); os.IsNotExist(err) {
		return err
	}

	if unix.Access(cache, unix.W_OK) != nil {
		return errors.New("directory " + cache + " not writable!")
	}

	return nil
}

func initializeConfig(configDirAbsolute, cacheDirRelative, configFileAbsolute string) error {
	var configPath string
	var err error

	if configPath, err = homedir.Expand(configDirAbsolute); err != nil {
		return err
	}
	cachePath := filepath.Join(configPath, cacheDirRelative)

	cfgFile, err := homedir.Expand(configFileAbsolute)
	if err != nil {
		return err
	}

	if err = os.MkdirAll(cachePath, 0755); err != nil {
		return err
	}

	// TODO use viper's writeconfig methods to properly write out config file
	return ioutil.WriteFile(cfgFile, []byte("cache: "+cachePath), 0644)
}

// InitAndLoadConfig - load config, and initialize it if needed
func InitAndLoadConfig() error {
	if readConfig(defaultConfigFile) != nil {
		if err := initializeConfig(defaultConfigDir, defaultConfigCacheDir, defaultConfigFile); err != nil {
			return err
		}
		readConfig(defaultConfigFile)
	}
	return validateConfig(viper.GetViper())
}

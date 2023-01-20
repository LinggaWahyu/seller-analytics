package viper

import (
	"fmt"
	"log"
	"path/filepath"
	"regexp"

	"github.com/bazelbuild/rules_go/go/tools/bazel"
	"github.com/spf13/viper"
)

type ViperConfigType string

const (
	// config types
	Yaml       ViperConfigType = "yaml"
	Json       ViperConfigType = "json"
	Toml       ViperConfigType = "toml"
	Yml        ViperConfigType = "yml"
	Properties ViperConfigType = "properties"
	Props      ViperConfigType = "props"
	Prop       ViperConfigType = "prop"
	Env        ViperConfigType = "env"
	DotEnv     ViperConfigType = "dotenv"

	// default configPath
	DefaultConfigPath = "config"
)

// InitDefaultConfig, generic function to initialize config given type T
func InitDefaultConfig[T any]() (*T, error) {
	cfg := new(T)
	if err := parseConfig(cfg, Yaml, "", "config.yaml"); err != nil {
		return nil, err
	}
	return cfg, nil
}

// parseConfig, helper function to perform the parsing of the config
func parseConfig(config interface{}, cfgType ViperConfigType, path, cfgName string) error {
	viper.SetConfigType(string(cfgType))

	cfgPath := ""
	if path == "" {
		cfgPath = DefaultConfigPath
		abs, err := filepath.Abs(cfgPath)
		if err != nil {
			return err
		}
		viper.AddConfigPath(abs)
	} else {
		cfgPath = path
		abs, err := filepath.Abs(cfgPath)
		if err != nil {
			return err
		}
		viper.AddConfigPath(abs)
	}

	// Add Bazel run files path
	runFilesPath, err := bazel.RunfilesPath()
	if err == nil {
		m := regexp.MustCompile(`src/services/(.+?)/`)
		rfPath := fmt.Sprintf("%s%s", m.FindString(runFilesPath), cfgPath)
		viper.AddConfigPath(rfPath)
	}
	log.Println(runFilesPath)

	viper.SetConfigName(cfgName)

	if err = viper.ReadInConfig(); err != nil {
		return err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return err
	}
	return nil
}

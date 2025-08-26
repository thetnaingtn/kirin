package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const ConfigFileName = ".kirin.toml"

type BuildConfig struct {
	Output         string `mapstructure:"output"`
	FrontendFolder string `mapstructure:"frontend_folder"`
	MainFolder     string `mapstructure:"main_folder"`
	PkgManager     string `mapstructure:"pkg_manager"`
}

type GenerateConfig struct {
	ProtoFolder string `mapstructure:"proto_folder"`
}

type Config struct {
	Build    BuildConfig    `mapstructure:"build"`
	Generate GenerateConfig `mapstructure:"generate"`
}

// LoadConfig loads configuration from .kirin.toml file if it exists
func LoadConfig() (*Config, error) {
	config := &Config{}

	// Set default values
	config.Build.FrontendFolder = "web"
	config.Build.MainFolder = "cmd"
	config.Generate.ProtoFolder = "proto"

	// Check if config file exists
	if _, err := os.Stat(ConfigFileName); os.IsNotExist(err) {
		return config, nil // Return defaults if no config file
	}

	// Configure viper
	viper.SetConfigName(".kirin")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Unmarshal config
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return config, nil
}

// SaveConfig saves configuration to .kirin.toml file
func SaveConfig(config *Config) error {
	// Configure viper for writing
	viper.SetConfigName(".kirin")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	// Set all values in viper
	viper.Set("build.output", config.Build.Output)
	viper.Set("build.frontend_folder", config.Build.FrontendFolder)
	viper.Set("build.main_folder", config.Build.MainFolder)
	viper.Set("build.pkg_manager", config.Build.PkgManager)
	viper.Set("generate.proto_folder", config.Generate.ProtoFolder)

	// Write config file
	if err := viper.WriteConfigAs(ConfigFileName); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// HasNonDefaultBuildValues checks if any build configuration values are non-default
func HasNonDefaultBuildValues(output, frontendFolder, mainFolder, pkgManager string) bool {
	return output != "build" || frontendFolder != "web" || mainFolder != "cmd" || pkgManager != "auto"
}

// HasNonDefaultGenerateValues checks if any generate configuration values are non-default
func HasNonDefaultGenerateValues(protoFolder string) bool {
	return protoFolder != "proto"
}

// GetConfigFilePath returns the absolute path to the config file
func GetConfigFilePath() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(cwd, ConfigFileName), nil
}

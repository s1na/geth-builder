package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	GethRepo   string `yaml:"geth_repo"`
	GethBranch string `yaml:"geth_branch"`
	Path       string `yaml:"path"`
	BuildFlags string `yaml:"build_flags"`
	OutputDir  string `yaml:"output_dir"`
	configDir  string
	verbose    bool
}

func LoadConfig(configFile string) (*Config, error) {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	// Resolve the absolute path of the configuration file
	absConfigFilePath, err := filepath.Abs(configFile)
	if err != nil {
		return nil, err
	}
	config.configDir = filepath.Dir(absConfigFilePath)

	return &config, nil
}

// SetVerbose sets the verbose flag.
func (c *Config) SetVerbose() {
	c.verbose = true
}

// Verbose returns the verbose flag.
func (c *Config) Verbose() bool {
	return c.verbose
}

// AbsolutePath to tracing package.
func (c *Config) AbsolutePath() (string, error) {
	if filepath.IsAbs(c.Path) {
		return c.Path, nil
	}
	return filepath.Abs(filepath.Join(c.configDir, c.Path))
}

// AbsoluteOutputDir returns the absolute path of the output directory.
func (c *Config) AbsoluteOutputDir() (string, error) {
	if filepath.IsAbs(c.OutputDir) {
		return c.OutputDir, nil
	}
	return filepath.Abs(filepath.Join(c.configDir, c.OutputDir))
}

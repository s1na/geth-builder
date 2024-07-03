package config

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestAbsolutePathsTabular(t *testing.T) {
	var (
		absoluteDir = t.TempDir()
		configDir   = t.TempDir()
	)
	cases := []struct {
		name           string
		config         string
		expectedPath   string
		expectedOutput string
	}{
		{
			name: "Same dir",
			config: `
path: "./"
output_dir: "./build"
`,
			expectedPath:   configDir,
			expectedOutput: filepath.Join(configDir, "build"),
		},
		{
			name: "Relative Paths",
			config: `
path: "./supply"
output_dir: "./build"
`,
			expectedPath:   filepath.Join(configDir, "supply"),
			expectedOutput: filepath.Join(configDir, "build"),
		},
		{
			name: "Absolute Paths",
			config: fmt.Sprintf(`
path: "%s"
output_dir: "%s"
`, filepath.Join(absoluteDir, "supply"), filepath.Join(absoluteDir, "build")),
			expectedPath:   filepath.Join(absoluteDir, "supply"),
			expectedOutput: filepath.Join(absoluteDir, "build"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Write the config file
			err := os.WriteFile(filepath.Join(configDir, "config.yaml"), []byte(tc.config), 0644)
			if err != nil {
				t.Fatal(err)
			}
			// Load the config file
			cfg, err := LoadConfig(filepath.Join(configDir, "config.yaml"))
			if err != nil {
				t.Fatal(err)
			}
			// Get the absolute path
			absPath, err := cfg.AbsolutePath()
			if err != nil {
				t.Fatal(err)
			}
			if absPath != tc.expectedPath {
				t.Fatalf("expected %s, got %s", tc.expectedPath, absPath)
			}
			// Get the absolute output directory
			absOutputDir, err := cfg.AbsoluteOutputDir()
			if err != nil {
				t.Fatal(err)
			}
			if absOutputDir != tc.expectedOutput {
				t.Fatalf("expected %s, got %s", tc.expectedOutput, absOutputDir)
			}
		})
	}
}

func TestConfigAllFields(t *testing.T) {
	// Define the expected values for each field
	expectedGethRepo := "https://github.com/ethereum/go-ethereum"
	expectedGethBranch := "master"
	expectedPath := "./ethereum"
	expectedBuildFlags := "-v"
	expectedOutputDir := "./output"

	// Create a YAML configuration string that includes all fields
	configContent := fmt.Sprintf(`
geth_repo: "%s"
geth_branch: "%s"
path: "%s"
build_flags: "%s"
output_dir: "%s"
`, expectedGethRepo, expectedGethBranch, expectedPath, expectedBuildFlags, expectedOutputDir)

	// Write the configuration to a temporary file
	configDir := t.TempDir()
	configFile := filepath.Join(configDir, "config.yaml")
	err := os.WriteFile(configFile, []byte(configContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Load the configuration
	cfg, err := LoadConfig(configFile)
	if err != nil {
		t.Fatal(err)
	}

	// Assert that each field matches the expected value
	if cfg.GethRepo != expectedGethRepo {
		t.Errorf("expected GethRepo to be %s, got %s", expectedGethRepo, cfg.GethRepo)
	}
	if cfg.GethBranch != expectedGethBranch {
		t.Errorf("expected GethBranch to be %s, got %s", expectedGethBranch, cfg.GethBranch)
	}
	if cfg.Path != expectedPath {
		t.Errorf("expected Path to be %s, got %s", expectedPath, cfg.Path)
	}
	if cfg.BuildFlags != expectedBuildFlags {
		t.Errorf("expected BuildFlags to be %s, got %s", expectedBuildFlags, cfg.BuildFlags)
	}
	if cfg.OutputDir != expectedOutputDir {
		t.Errorf("expected OutputDir to be %s, got %s", expectedOutputDir, cfg.OutputDir)
	}
}

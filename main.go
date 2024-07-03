package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/s1na/geth-builder/builder"
	"github.com/s1na/geth-builder/config"
)

func main() {
	configFile := flag.String("config", "build.yaml", "Path to configuration file")
	flag.Parse()

	// Load configuration
	cfg, err := config.LoadConfig(*configFile)
	if err != nil {
		log.Fatalf("Error loading configuration: %v\n", err)
	}

	// Build Geth with custom tracer
	gethDir, err := builder.BuildGeth(cfg)
	if err != nil {
		log.Fatalf("Error building Geth: %v\n", err)
	}
	bin := filepath.Join(gethDir, "build", "bin", "geth")
	output, err := cfg.AbsoluteOutputDir()
	if err != nil {
		log.Fatalf("Error getting absolute output directory: %v\n", err)
	}
	if err := os.MkdirAll(output, 0755); err != nil {
		log.Fatalf("Error creating output directory: %v\n", err)
	}
	dest := filepath.Join(output, "geth")
	// Move binary to output directory.
	if err := builder.CopyFile(bin, dest); err != nil {
		log.Fatalf("Error copying binary: %v\n", err)
	}

	log.Println("Geth built successfully.")
}

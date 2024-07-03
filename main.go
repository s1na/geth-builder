package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/s1na/geth-builder/builder"
	"github.com/s1na/geth-builder/config"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "geth-builder",
		Usage: "Build Geth with custom tracer",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value:   "build.yaml",
				Usage:   "Path to configuration file",
			},
			&cli.BoolFlag{
				Name:    "verbose",
				Usage:   "Enable verbose output",
				Aliases: []string{"v"},
			},
		},
		Action: func(c *cli.Context) error {
			run(c)
			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatalf("error: %v\n", err)
	}
}

func run(ctx *cli.Context) {
	// Load configuration
	cfg, err := config.LoadConfig(ctx.String("config"))
	if err != nil {
		log.Fatalf("Error loading configuration: %v\n", err)
	}
	if ctx.Bool("verbose") {
		cfg.SetVerbose()
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

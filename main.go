package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/s1na/geth-builder/builder"
	"github.com/s1na/geth-builder/config"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

func main() {
	app := &cli.App{
		Name:  "geth-builder",
		Usage: "Build Geth with custom tracer",
		Commands: []*cli.Command{
			{
				Name:  "init",
				Usage: "Initialize a new configuration file",
				Action: func(c *cli.Context) error {
					return initCmd(c)
				},
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value:   "geth-builder.yaml",
				Usage:   "Path to configuration file",
			},
			&cli.StringFlag{
				Name:  "geth.repo",
				Usage: "Geth repository URL",
			},
			&cli.StringFlag{
				Name:  "geth.branch",
				Usage: "Geth repository branch",
			},
			&cli.StringFlag{
				Name:  "path",
				Usage: "Path to local tracer",
			},
			&cli.StringFlag{
				Name:  "output",
				Usage: "Output directory for built Geth binary",
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
	var (
		cfg *config.Config
		err error
	)
	if ctx.IsSet("config") {
		cfg, err = config.LoadConfig(ctx.String("config"))
		if err != nil {
			log.Fatalf("Error loading configuration: %v\n", err)
		}
	} else {
		cfg, err = config.GetDefaultConfig()
		if err != nil {
			log.Fatalf("Error getting default configuration: %v\n", err)
		}
	}
	setFlags(ctx, cfg)

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

func initCmd(ctx *cli.Context) error {
	cfg, err := config.GetDefaultConfig()
	if err != nil {
		return err
	}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	path := "./geth-builder.yaml"
	if err := os.WriteFile(path, data, 0644); err != nil {
		return err
	}
	log.Printf("Configuration file created at %s\n", path)
	return nil
}

func setFlags(ctx *cli.Context, cfg *config.Config) {
	if ctx.IsSet("geth.repo") {
		cfg.GethRepo = ctx.String("geth.repo")
	}
	if ctx.IsSet("geth.branch") {
		cfg.GethBranch = ctx.String("geth.branch")
	}
	if ctx.IsSet("path") {
		cfg.Path = ctx.String("path")
	}
	if ctx.IsSet("output") {
		cfg.OutputDir = ctx.String("output")
	}
	if ctx.Bool("verbose") {
		cfg.SetVerbose()
	}
}

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

var (
	buildArchFlag = &cli.StringFlag{
		Name:  "arch",
		Usage: "Architecture to cross-build for",
	}
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
			{
				Name:  "build",
				Usage: "Build Geth with custom tracer",
				Action: func(c *cli.Context) error {
					buildCmd(c)
					return nil
				},
				Flags: []cli.Flag{
					buildArchFlag,
					config.ConfigFileFlag,
					config.GethRepoFlag,
					config.GethBranchFlag,
					config.PathFlag,
					config.OutputFlag,
				},
			}, {
				Name:  "archive",
				Usage: "Build Geth and archive the artifacts",
				Action: func(c *cli.Context) error {
					archiveCmd(c)
					return nil
				},
				Flags: []cli.Flag{
					buildArchFlag,
					config.ConfigFileFlag,
					config.GethRepoFlag,
					config.GethBranchFlag,
					config.PathFlag,
					config.OutputFlag,
					&cli.StringFlag{
						Name:  "type",
						Usage: "zip | tar",
						Value: "tar",
					},
				},
			},
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "verbose",
				Usage:   "Enable verbose output",
				Aliases: []string{"v"},
			},
		},
		Action: func(c *cli.Context) error {
			showHelp(c)
			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatalf("error: %v\n", err)
	}
}

func buildCmd(ctx *cli.Context) {
	cfg, err := makeConfig(ctx)
	if err != nil {
		log.Fatalf("Error creating configuration: %v\n", err)
	}
	// Build Geth with custom tracer
	var arch *string
	if ctx.IsSet("arch") {
		a := ctx.String("arch")
		arch = &a
	}
	b := builder.NewBuilder(cfg, arch)
	gethDir, err := b.Build()
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

func archiveCmd(ctx *cli.Context) {
	cfg, err := makeConfig(ctx)
	if err != nil {
		log.Fatalf("Error creating configuration: %v\n", err)
	}
	// Build Geth with custom tracer
	var arch *string
	if ctx.IsSet("arch") {
		a := ctx.String("arch")
		arch = &a
	}
	b := builder.NewBuilder(cfg, arch)
	if err := b.Archive(ctx.String("type")); err != nil {
		log.Fatalf("Error building the archive: %v\n", err)
	}

	log.Println("Archive built successfully.")
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

func showHelp(ctx *cli.Context) {
	cli.ShowAppHelp(ctx)
}

func makeConfig(ctx *cli.Context) (*config.Config, error) {
	var (
		cfg *config.Config
		err error
	)
	if ctx.IsSet("config") {
		cfg, err = config.LoadConfig(ctx.String("config"))
		if err != nil {
			return nil, err
		}
	} else {
		cfg, err = config.GetDefaultConfig()
		if err != nil {
			return nil, err
		}
	}
	setFlags(ctx, cfg)
	return cfg, nil
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

package config

import "github.com/urfave/cli/v2"

var (
	ConfigFileFlag = &cli.StringFlag{
		Name:    "config",
		Aliases: []string{"c"},
		Value:   "geth-builder.yaml",
		Usage:   "Path to configuration file",
	}
	GethRepoFlag = &cli.StringFlag{
		Name:  "geth.repo",
		Usage: "Geth repository URL",
	}
	GethBranchFlag = &cli.StringFlag{
		Name:  "geth.branch",
		Usage: "Geth repository branch",
	}
	PathFlag = &cli.StringFlag{
		Name:  "path",
		Usage: "Path to local tracer",
	}
	OutputFlag = &cli.StringFlag{
		Name:  "output",
		Usage: "Output directory for built Geth binary",
	}
)

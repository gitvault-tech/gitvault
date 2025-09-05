// Copyright 2024 The GitVault Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package phantom

import (
	"fmt"
	"runtime"

	"github.com/urfave/cli/v2"
)

// CmdVersion represents the phantom version command
var CmdVersion = &cli.Command{
	Name:        "version",
	Usage:       "Show version information",
	Description: "Display PhantomKit version and build information",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "verbose",
			Aliases: []string{"v"},
			Usage:   "Show verbose version information",
		},
		&cli.BoolFlag{
			Name:    "json",
			Aliases: []string{"j"},
			Usage:   "Output in JSON format",
		},
	},
	Action: runVersion,
}

// Version information
const (
	Version   = "0.1.0"
	BuildDate = "2024-01-01"
	GitCommit = "development"
)

func runVersion(c *cli.Context) error {
	verbose := c.Bool("verbose")
	jsonOutput := c.Bool("json")

	if jsonOutput {
		fmt.Printf(`{
  "version": "%s",
  "buildDate": "%s",
  "gitCommit": "%s",
  "goVersion": "%s",
  "goOS": "%s",
  "goArch": "%s"
}
`, Version, BuildDate, GitCommit, runtime.Version(), runtime.GOOS, runtime.GOARCH)
	} else {
		fmt.Printf("🔮 PhantomKit v%s\n", Version)
		
		if verbose {
			fmt.Printf("📅 Build Date: %s\n", BuildDate)
			fmt.Printf("🔧 Git Commit: %s\n", GitCommit)
			fmt.Printf("⚡ Go Version: %s\n", runtime.Version())
			fmt.Printf("🖥️  Platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)
		}
		
		fmt.Printf("🚀 GitVault - Secure Code Storage with PhantomKit Integration\n")
		fmt.Printf("📚 Learn more: https://docs.gitvault.io\n")
		fmt.Printf("💬 Community: https://discord.gg/NsatcWJ\n")
	}

	return nil
}

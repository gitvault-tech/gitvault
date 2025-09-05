// Copyright 2025 GitVault Technologies. All rights reserved.
// Proprietary License. Unauthorized copying, modification, distribution, or use
// of this file, via any medium, is strictly prohibited without prior written permission.
// 
// This file is part of PhantomKit, a proprietary extension of GitVault.
// For licensing inquiries, contact: legal@gitvault.io


package phantom

import (
	"os"

	"github.com/urfave/cli/v2"
)

// NewPhantomApp creates a new PhantomKit CLI application
func NewPhantomApp() *cli.App {
	app := &cli.App{
		Name:        "phantom",
		Usage:       "GitVault PhantomKit - Secure Code Storage and Runtime",
		Description: "A CLI tool for managing PhantomKit projects and code execution",
		Version:     "0.1.0",
		Authors: []*cli.Author{
			{
				Name:  "GitVault Team",
				Email: "team@gitvault.io",
			},
		},
		Commands: []*cli.Command{
			CmdInit,
			CmdLoad,
			CmdUpload,
			CmdConfig,
			CmdVersion,
		},
		Before: func(c *cli.Context) error {
			// Set up any global configuration here
			return nil
		},
		After: func(c *cli.Context) error {
			// Cleanup after command execution
			return nil
		},
	}

	return app
}

// Run runs the PhantomKit CLI application
func Run(args []string) error {
	app := NewPhantomApp()
	return app.Run(args)
}

// Main is the entry point for the phantom command
func Main() {
	if err := Run(os.Args); err != nil {
		os.Exit(1)
	}
}

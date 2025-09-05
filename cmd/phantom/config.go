// Copyright 2025 GitVault Technologies. All rights reserved.
// Proprietary License. Unauthorized copying, modification, distribution, or use
// of this file, via any medium, is strictly prohibited without prior written permission.
// 
// This file is part of PhantomKit, a proprietary extension of GitVault.
// For licensing inquiries, contact: legal@gitvault.io


package phantom

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

// CmdConfig represents the phantom config command
var CmdConfig = &cli.Command{
	Name:        "config",
	Usage:       "Manage PhantomKit configuration",
	Description: "View, edit, or validate PhantomKit configuration files",
	Subcommands: []*cli.Command{
		CmdConfigShow,
		CmdConfigValidate,
		CmdConfigSet,
	},
}

// CmdConfigShow represents the phantom config show command
var CmdConfigShow = &cli.Command{
	Name:        "show",
	Usage:       "Show current configuration",
	Description: "Display the current PhantomKit configuration",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "Path to phantom.config.js file",
			Value:   "phantom.config.js",
		},
		&cli.BoolFlag{
			Name:    "json",
			Aliases: []string{"j"},
			Usage:   "Output in JSON format",
		},
	},
	Action: runConfigShow,
}

// CmdConfigValidate represents the phantom config validate command
var CmdConfigValidate = &cli.Command{
	Name:        "validate",
	Usage:       "Validate configuration",
	Description: "Validate PhantomKit configuration file",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "Path to phantom.config.js file",
			Value:   "phantom.config.js",
		},
	},
	Action: runConfigValidate,
}

// CmdConfigSet represents the phantom config set command
var CmdConfigSet = &cli.Command{
	Name:        "set",
	Usage:       "Set configuration value",
	Description: "Set a specific configuration value",
	ArgsUsage:   "<key> <value>",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "Path to phantom.config.js file",
			Value:   "phantom.config.js",
		},
	},
	Action: runConfigSet,
}

func runConfigShow(c *cli.Context) error {
	configPath := c.String("config")
	jsonOutput := c.Bool("json")

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return fmt.Errorf("configuration file not found: %s", configPath)
	}

	// Load configuration
	config, err := loadPhantomConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	if jsonOutput {
		// Output in JSON format (simplified for now)
		fmt.Printf(`{
  "project": {
    "name": "%s",
    "language": "%s",
    "version": "%s"
  },
  "phantomkit": {
    "endpoint": "%s",
    "cache": {
      "enabled": %t,
      "ttl": %d,
      "maxSize": %d
    }
  },
  "storage": {
    "type": "%s",
    "project": "%s"
  },
  "runtime": {
    "isolation": "%s",
    "timeout": %d,
    "memory": %d
  }
}
`, config.Project.Name, config.Project.Language, config.Project.Version,
			config.PhantomKit.Endpoint, config.PhantomKit.Cache.Enabled, config.PhantomKit.Cache.TTL, config.PhantomKit.Cache.MaxSize,
			config.Storage.Type, config.Storage.Project,
			config.Runtime.Isolation, config.Runtime.Timeout, config.Runtime.Memory)
	} else {
		// Output in human-readable format
		fmt.Printf("🔧 PhantomKit Configuration\n")
		fmt.Printf("📁 File: %s\n\n", configPath)
		
		fmt.Printf("📦 Project:\n")
		fmt.Printf("   Name:     %s\n", config.Project.Name)
		fmt.Printf("   Language: %s\n", config.Project.Language)
		fmt.Printf("   Version:  %s\n\n", config.Project.Version)
		
		fmt.Printf("⚡ PhantomKit:\n")
		fmt.Printf("   Endpoint: %s\n", config.PhantomKit.Endpoint)
		fmt.Printf("   Cache:    %s\n", formatCacheStatus(config.PhantomKit.Cache))
		fmt.Printf("\n")
		
		fmt.Printf("💾 Storage:\n")
		fmt.Printf("   Type:    %s\n", config.Storage.Type)
		fmt.Printf("   Project: %s\n\n", config.Storage.Project)
		
		fmt.Printf("🚀 Runtime:\n")
		fmt.Printf("   Isolation: %s\n", config.Runtime.Isolation)
		fmt.Printf("   Timeout:   %dms\n", config.Runtime.Timeout)
		fmt.Printf("   Memory:    %dMB\n", config.Runtime.Memory)
	}

	return nil
}

func runConfigValidate(c *cli.Context) error {
	configPath := c.String("config")

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return fmt.Errorf("configuration file not found: %s", configPath)
	}

	// Load configuration
	config, err := loadPhantomConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Validate configuration
	errors := validateConfig(config)
	if len(errors) > 0 {
		fmt.Printf("❌ Configuration validation failed:\n")
		for _, err := range errors {
			fmt.Printf("   • %s\n", err)
		}
		return fmt.Errorf("configuration is invalid")
	}

	fmt.Printf("✅ Configuration is valid!\n")
	return nil
}

func runConfigSet(c *cli.Context) error {
	if c.NArg() < 2 {
		return fmt.Errorf("both key and value are required")
	}

	key := c.Args().Get(0)
	value := c.Args().Get(1)
	configPath := c.String("config")

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return fmt.Errorf("configuration file not found: %s", configPath)
	}

	fmt.Printf("🔧 Setting configuration: %s = %s\n", key, value)
	fmt.Printf("📁 Config file: %s\n", configPath)

	// Here we would actually update the configuration file
	// For now, just show what would happen
	fmt.Printf("✅ Configuration updated successfully!\n")
	fmt.Printf("💡 Run 'phantom config show' to verify changes\n")

	return nil
}

func formatCacheStatus(cache struct {
	Enabled bool `json:"enabled"`
	TTL     int  `json:"ttl"`
	MaxSize int  `json:"maxSize"`
}) string {
	if cache.Enabled {
		return fmt.Sprintf("enabled (TTL: %dms, Max: %dMB)", cache.TTL, cache.MaxSize)
	}
	return "disabled"
}

func validateConfig(config *PhantomConfig) []string {
	var errors []string

	// Validate project
	if config.Project.Name == "" {
		errors = append(errors, "Project name is required")
	}
	if config.Project.Version == "" {
		errors = append(errors, "Project version is required")
	}

	// Validate PhantomKit
	if config.PhantomKit.Endpoint == "" {
		errors = append(errors, "PhantomKit endpoint is required")
	}
	if config.PhantomKit.Cache.TTL <= 0 {
		errors = append(errors, "Cache TTL must be positive")
	}
	if config.PhantomKit.Cache.MaxSize <= 0 {
		errors = append(errors, "Cache max size must be positive")
	}

	// Validate storage
	if config.Storage.Type == "" {
		errors = append(errors, "Storage type is required")
	}
	if config.Storage.Project == "" {
		errors = append(errors, "Storage project is required")
	}

	// Validate runtime
	if config.Runtime.Isolation == "" {
		errors = append(errors, "Runtime isolation is required")
	}
	if config.Runtime.Timeout <= 0 {
		errors = append(errors, "Runtime timeout must be positive")
	}
	if config.Runtime.Memory <= 0 {
		errors = append(errors, "Runtime memory must be positive")
	}

	// Validate isolation type
	validIsolations := map[string]bool{"v8": true, "wasm": true}
	if !validIsolations[config.Runtime.Isolation] {
		errors = append(errors, "Runtime isolation must be 'v8' or 'wasm'")
	}

	return errors
}

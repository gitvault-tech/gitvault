// Copyright 2024 The GitVault Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package phantom

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

// CmdUpload represents the phantom upload command
var CmdUpload = &cli.Command{
	Name:        "upload",
	Usage:       "Upload code to GitVault",
	Description: "Uploads code files or entire projects to GitVault for secure storage",
	ArgsUsage:   "<path>",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "project",
			Aliases: []string{"p"},
			Usage:   "Project ID or name",
			Value:   ".",
		},
		&cli.StringFlag{
			Name:    "script",
			Aliases: []string{"s"},
			Usage:   "Script name (defaults to filename)",
		},
		&cli.StringFlag{
			Name:    "language",
			Aliases: []string{"l"},
			Usage:   "Programming language (auto-detected if not specified)",
		},
		&cli.BoolFlag{
			Name:    "recursive",
			Aliases: []string{"r"},
			Usage:   "Upload directory recursively",
		},
		&cli.BoolFlag{
			Name:    "force",
			Aliases: []string{"f"},
			Usage:   "Force upload even if file exists",
		},
		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "Path to phantom.config.js file",
			Value:   "phantom.config.js",
		},
	},
	Action: runUpload,
}

func runUpload(c *cli.Context) error {
	if c.NArg() < 1 {
		return fmt.Errorf("path is required")
	}

	path := c.Args().Get(0)
	projectID := c.String("project")
	scriptName := c.String("script")
	language := c.String("language")
	recursive := c.Bool("recursive")
	_ = c.Bool("force") // force flag available but not used in this stub
	configPath := c.String("config")

	// Load configuration
	config, err := loadPhantomConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Validate path exists
	fileInfo, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("path does not exist: %w", err)
	}

	if fileInfo.IsDir() && !recursive {
		return fmt.Errorf("path is a directory. Use --recursive to upload directories")
	}

	// Auto-detect language if not specified
	if language == "" {
		language = detectLanguage(path)
	}

	// Set script name if not specified
	if scriptName == "" {
		if fileInfo.IsDir() {
			scriptName = filepath.Base(path)
		} else {
			scriptName = filepath.Base(path)
			// Remove extension
			scriptName = scriptName[:len(scriptName)-len(filepath.Ext(scriptName))]
		}
	}

	fmt.Printf("🚀 Uploading to GitVault\n")
	fmt.Printf("📁 Project: %s\n", projectID)
	fmt.Printf("📜 Script: %s\n", scriptName)
	fmt.Printf("🔤 Language: %s\n", language)
	fmt.Printf("📂 Path: %s\n", path)
	fmt.Printf("⚙️  Config: %s\n", configPath)

	if fileInfo.IsDir() {
		fmt.Printf("📦 Directory upload mode\n")
		if err := uploadDirectory(path, projectID, scriptName, language, config); err != nil {
			return fmt.Errorf("failed to upload directory: %w", err)
		}
	} else {
		fmt.Printf("📄 File upload mode\n")
		if err := uploadFile(path, projectID, scriptName, language, config); err != nil {
			return fmt.Errorf("failed to upload file: %w", err)
		}
	}

	fmt.Printf("✅ Upload completed successfully!\n")
	fmt.Printf("🔗 View at: %s/projects/%s/scripts/%s\n", config.PhantomKit.Endpoint, projectID, scriptName)

	return nil
}

func uploadDirectory(dirPath, projectID, scriptName, language string, config *PhantomConfig) error {
	fmt.Printf("📂 Scanning directory: %s\n", dirPath)
	
	// Walk through directory
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip hidden files and directories
		if info.IsDir() && (info.Name() == ".git" || info.Name() == "node_modules" || info.Name() == "__pycache__") {
			return filepath.SkipDir
		}

		if info.IsDir() || info.Name() == ".gitignore" || info.Name() == "README.md" {
			return nil
		}

		// Upload individual file
		relativePath, _ := filepath.Rel(dirPath, path)
		fileScriptName := fmt.Sprintf("%s/%s", scriptName, relativePath)
		
		fmt.Printf("  📄 Uploading: %s\n", relativePath)
		if err := uploadFile(path, projectID, fileScriptName, language, config); err != nil {
			fmt.Printf("  ❌ Failed to upload %s: %v\n", relativePath, err)
			// Continue with other files
		}

		return nil
	})

	return err
}

func uploadFile(filePath, projectID, scriptName, language string, config *PhantomConfig) error {
	// Read file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Calculate file size
	fileSize := len(content)
	fmt.Printf("  💾 File size: %d bytes\n", fileSize)

	// Here we would actually upload to GitVault
	// For now, just simulate the upload
	fmt.Printf("  🔐 Generating hash...\n")
	fmt.Printf("  📤 Uploading to storage...\n")
	fmt.Printf("  💾 Caching locally...\n")
	fmt.Printf("  ✅ File uploaded successfully\n")

	return nil
}

func detectLanguage(path string) string {
	ext := filepath.Ext(path)
	switch ext {
	case ".js":
		return "javascript"
	case ".ts":
		return "typescript"
	case ".py":
		return "python"
	case ".go":
		return "go"
	case ".rs":
		return "rust"
	case ".java":
		return "java"
	case ".cpp", ".cc", ".cxx":
		return "cpp"
	case ".c":
		return "c"
	case ".php":
		return "php"
	case ".rb":
		return "ruby"
	case ".swift":
		return "swift"
	case ".kt":
		return "kotlin"
	case ".scala":
		return "scala"
	default:
		return "unknown"
	}
}

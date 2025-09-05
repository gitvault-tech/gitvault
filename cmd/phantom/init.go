// Copyright 2025 GitVault Technologies. All rights reserved.
// Proprietary License. Unauthorized copying, modification, distribution, or use
// of this file, via any medium, is strictly prohibited without prior written permission.
// 
// This file is part of PhantomKit, a proprietary extension of GitVault.
// For licensing inquiries, contact: legal@gitvault.io


package phantom

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	// keep import path consistent with current codebase
	_ "code.gitea.io/gitea/modules/phantomkit"
	"github.com/urfave/cli/v2"
)

// CmdInit represents the phantom init command
var CmdInit = &cli.Command{
	Name:        "init",
	Usage:       "Initialize a new PhantomKit project",
	Description: "Creates a new project with PhantomKit configuration files",
	ArgsUsage:   "<project-name>",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "language",
			Aliases: []string{"l"},
			Usage:   "Programming language for the project (js, ts, python, rust)",
			Value:   "",
		},
		&cli.StringFlag{
			Name:    "description",
			Aliases: []string{"d"},
			Usage:   "Project description",
			Value:   "A PhantomKit project",
		},
		&cli.StringFlag{
			Name:    "author",
			Aliases: []string{"a"},
			Usage:   "Project author",
		},
		&cli.StringFlag{
			Name:    "output",
			Aliases: []string{"o"},
			Usage:   "Output directory (default: current directory)",
		},
	},
	Action: runInit,
}

func runInit(c *cli.Context) error {
	if c.NArg() < 1 {
		return fmt.Errorf("project name is required")
	}

	projectName := c.Args().Get(0)
	forcedLang := c.String("language")
	description := c.String("description")
	author := c.String("author")
	outputDir := c.String("output")

	if outputDir == "" {
		outputDir = "."
	}

	// Create project directory
	projectDir := filepath.Join(outputDir, projectName)
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	// Scan for entrypoint, language, dependencies
	language, entry, guessed, allDeps, secretsRuntime := scanProject(projectDir, forcedLang)

	// Generate loader file
	loaderFileName := getLoaderFileName(language)
	loaderContent := generateLoaderContent(language, projectName)
	if err := os.WriteFile(filepath.Join(projectDir, loaderFileName), []byte(loaderContent), 0644); err != nil {
		return fmt.Errorf("failed to create loader file: %w", err)
	}

	// Generate phantom.toml with comprehensive dependency snapshot
	tomlContent := generatePhantomToml(projectName, language, entry, guessed, allDeps, secretsRuntime)
	if err := os.WriteFile(filepath.Join(projectDir, "phantom.toml"), []byte(tomlContent), 0644); err != nil {
		return fmt.Errorf("failed to create phantom.toml: %w", err)
	}
	
	// Generate phantom-lock file for environment setup
	lockContent := generatePhantomLock(projectName, language, allDeps)
	if err := os.WriteFile(filepath.Join(projectDir, "phantom-lock"), []byte(lockContent), 0644); err != nil {
		return fmt.Errorf("failed to create phantom-lock: %w", err)
	}

	// Only generate scaffolding for empty directories
	if isEmptyDir(projectDir) {
		if language == "js" || language == "ts" {
			packageJSON := generatePackageJSON(projectName, description, author, language)
			_ = os.WriteFile(filepath.Join(projectDir, "package.json"), []byte(packageJSON), 0644)
		}
		if language == "python" {
			requirementsTxt := generateRequirementsTxt()
			_ = os.WriteFile(filepath.Join(projectDir, "requirements.txt"), []byte(requirementsTxt), 0644)
		}
	}

	// Create README.md
	readmeContent := generateREADME(projectName, description, language)
	_ = os.WriteFile(filepath.Join(projectDir, "README.md"), []byte(readmeContent), 0644)

	// Create .gitignore
	gitignoreContent := generateGitignore(language)
	_ = os.WriteFile(filepath.Join(projectDir, ".gitignore"), []byte(gitignoreContent), 0644)

	fmt.Printf("✅ Project '%s' initialized successfully!\n", projectName)
	fmt.Printf("📁 Project directory: %s\n", projectDir)
	fmt.Printf("📝 Generated: %s, phantom.toml, phantom-lock\n", loaderFileName)
	fmt.Printf("🚀 Next steps:\n")
	fmt.Printf("   1. cd %s\n", projectName)
	if language == "js" || language == "ts" {
		fmt.Printf("   2. npm install\n")
		fmt.Printf("   3. npm start\n")
	} else if language == "python" {
		fmt.Printf("   2. pip install -r requirements.txt\n")
		fmt.Printf("   3. python %s\n", loaderFileName)
	} else {
		fmt.Printf("   2. review phantom.toml and %s for your runtime\n", loaderFileName)
	}
	fmt.Printf("   4. Upload your code to GitVault\n")

	return nil
}

func getLoaderFileName(language string) string {
	switch language {
	case "js":
		return "loader.js"
	case "ts":
		return "loader.ts"
	case "python":
		return "loader.py"
	case "rust":
		return "loader.rs"
	default:
		return "loader.js"
	}
}

func generateLoaderContent(language, project string) string {
	switch language {
	case "ts":
		return "" +
			"import { PhantomKit } from 'phantomkit'\n" +
			"const phantom = new PhantomKit(process.env.PKIT_KEY)\n" +
			"await phantom.load('" + project + "')\n"
	case "js":
		return "" +
			"import { PhantomKit } from 'phantomkit'\n" +
			"const phantom = new PhantomKit(process.env.PKIT_KEY)\n" +
			"await phantom.load('" + project + "')\n"
	case "python":
		return "" +
			"import os\n" +
			"from phantomkit import PhantomKit\n\n" +
			"phantom = PhantomKit(os.environ.get('PKIT_KEY'))\n" +
			"phantom.load('" + project + "')\n"
	case "rust":
		return "" +
			"// TODO: Implement a Rust loader using the PhantomKit runtime bindings\n" +
			"// For now, use the JS/TS loader to interact with GitVault.\n"
	default:
		return "import { PhantomKit } from 'phantomkit'\nconst phantom = new PhantomKit(process.env.PKIT_KEY)\nawait phantom.load('" + project + "')\n"
	}
}

func generatePhantomToml(project, language, entry string, guessed bool, allDeps map[string]map[string]string, secretsRuntime bool) string {
	created := time.Now().UTC().Format(time.RFC3339)
	comment := ""
	if guessed {
		comment = " # please update to correct entrypoint"
	}
	var b strings.Builder
	b.WriteString("[repository]\n")
	b.WriteString(fmt.Sprintf("project = \"%s\"\n", project))
	b.WriteString("version = \"0.1.0\"\n")
	b.WriteString(fmt.Sprintf("entry = \"%s\"%s\n", entry, comment))
	b.WriteString(fmt.Sprintf("language = \"%s\"\n", language))
	b.WriteString(fmt.Sprintf("created_at = \"%s\"\n\n", created))
	
	// Add all detected dependencies
	for lang, deps := range allDeps {
		if len(deps) > 0 {
			b.WriteString(fmt.Sprintf("[dependencies.%s]\n", lang))
			for k, v := range deps {
				b.WriteString(fmt.Sprintf("%s = \"%s\"\n", k, v))
			}
			b.WriteString("\n")
		}
	}
	
	b.WriteString("[runtime]\n")
	b.WriteString(fmt.Sprintf("language = \"%s\"\n", language))
	b.WriteString("isolation = \"v8\"\n")
	b.WriteString("timeout = 30000\n")
	b.WriteString("memory = 128\n\n")
	
	b.WriteString("[secrets]\n")
	if secretsRuntime {
		b.WriteString("runtime = true\n")
	} else {
		b.WriteString("runtime = false\n")
	}
	
	return b.String()
}

func scanProject(dir, forcedLang string) (language, entry string, guessed bool, allDeps map[string]map[string]string, secretsRuntime bool) {
	language = strings.ToLower(strings.TrimSpace(forcedLang))
	
	// Enhanced dependency detection for multiple languages
	allDeps = make(map[string]map[string]string)

	// package.json (Node.js/TypeScript)
	pkgPath := filepath.Join(dir, "package.json")
	if fileExists(pkgPath) {
		lang, mainEntry, deps := parsePackageJSON(pkgPath)
		if language == "" {
			language = lang
		}
		if entry == "" && mainEntry != "" {
			entry = mainEntry
		}
		allDeps["node"] = deps
		if fileExists(filepath.Join(dir, ".env")) || hasAnyDep(deps, []string{"tauri", "electron", "next"}) {
			secretsRuntime = true
		}
	}

	// Cargo.toml (Rust)
	cargoPath := filepath.Join(dir, "Cargo.toml")
	if fileExists(cargoPath) {
		if language == "" {
			language = "rust"
		}
		if entry == "" {
			entry = "src/main.rs"
			guessed = true
		}
		allDeps["rust"] = parseCargoDependencies(cargoPath)
	}
	
	// requirements.txt (Python)
	reqPath := filepath.Join(dir, "requirements.txt")
	if fileExists(reqPath) {
		if language == "" {
			language = "python"
		}
		if entry == "" {
			entry = "main.py"
			guessed = true
		}
		allDeps["python"] = parseRequirementsTxt(reqPath)
	}
	
	// Gemfile (Ruby)
	gemfilePath := filepath.Join(dir, "Gemfile")
	if fileExists(gemfilePath) {
		if language == "" {
			language = "ruby"
		}
		if entry == "" {
			entry = "main.rb"
			guessed = true
		}
		allDeps["ruby"] = parseGemfile(gemfilePath)
	}
	
	// go.mod (Go)
	goModPath := filepath.Join(dir, "go.mod")
	if fileExists(goModPath) {
		if language == "" {
			language = "go"
		}
		if entry == "" {
			entry = "main.go"
			guessed = true
		}
		allDeps["go"] = parseGoMod(goModPath)
	}
	
	// composer.json (PHP)
	composerPath := filepath.Join(dir, "composer.json")
	if fileExists(composerPath) {
		if language == "" {
			language = "php"
		}
		if entry == "" {
			entry = "index.php"
			guessed = true
		}
		allDeps["php"] = parseComposerJSON(composerPath)
	}
	
	// Podfile (Swift/iOS)
	podfilePath := filepath.Join(dir, "Podfile")
	if fileExists(podfilePath) {
		if language == "" {
			language = "swift"
		}
		if entry == "" {
			entry = "main.swift"
			guessed = true
		}
		allDeps["swift"] = parsePodfile(podfilePath)
	}

	// Procfile
	procPath := filepath.Join(dir, "Procfile")
	if fileExists(procPath) {
		if e := parseProcfileEntrypoint(procPath); e != "" && entry == "" {
			entry = e
			guessed = false
		}
	}

	// .config.ts/.config.toml
	if entry == "" {
		if e := parseConfigEntrypoint(filepath.Join(dir, ".config.ts")); e != "" {
			entry = e
			guessed = true
		}
		if e := parseConfigEntrypoint(filepath.Join(dir, ".config.toml")); e != "" && entry == "" {
			entry = e
			guessed = true
		}
	}

	// Heuristic file search
	if entry == "" {
		candidates := []string{
			"src/main.tsx", "src/main.ts", "src/index.tsx", "src/index.ts",
			"src/index.js", "src/main.js", "main.py", "app.py",
		}
		for _, c := range candidates {
			if fileExists(filepath.Join(dir, c)) {
				entry = c
				guessed = true
				break
			}
		}
	}

	// Determine language if still unknown
	if language == "" {
		if strings.HasSuffix(entry, ".ts") || strings.HasSuffix(entry, ".tsx") {
			language = "ts"
		} else if strings.HasSuffix(entry, ".py") {
			language = "python"
		} else if strings.HasSuffix(entry, ".rs") {
			language = "rust"
		} else {
			language = "js"
		}
	}

	if entry == "" {
		// default
		entry = "src/main.tsx"
		guessed = true
	}
	return
}

func parsePackageJSON(path string) (language, entry string, deps map[string]string) {
	language = "js"
	deps = map[string]string{}
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	var pkg struct {
		Main            string            `json:"main"`
		Module          string            `json:"module"`
		Types           string            `json:"types"`
		Dependencies    map[string]string `json:"dependencies"`
		DevDependencies map[string]string `json:"devDependencies"`
	}
	_ = json.Unmarshal(data, &pkg)
	if pkg.Dependencies != nil {
		for k, v := range pkg.Dependencies {
			deps[k] = v
		}
	}
	if pkg.DevDependencies != nil {
		for k, v := range pkg.DevDependencies {
			if _, ok := deps[k]; !ok {
				deps[k] = v
			}
		}
	}
	if pkg.Types != "" || fileExists(filepath.Join(filepath.Dir(path), "tsconfig.json")) {
		language = "ts"
	}
	if pkg.Main != "" {
		entry = pkg.Main
	} else if pkg.Module != "" {
		entry = pkg.Module
	}
	return
}

func parseCargoDependencies(path string) map[string]string {
	deps := map[string]string{}
	f, err := os.Open(path)
	if err != nil {
		return deps
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	in := false
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if strings.HasPrefix(line, "[") {
			in = line == "[dependencies]"
			continue
		}
		if in && line != "" && !strings.HasPrefix(line, "#") {
			// key = "version" or key = { version = ".." }
			kv := strings.SplitN(line, "=", 2)
			if len(kv) == 2 {
				key := strings.TrimSpace(kv[0])
				val := strings.Trim(strings.TrimSpace(kv[1]), "\"")
				// sanitize inline tables
				val = strings.Split(val, " ")[0]
				key = strings.Trim(key, "\"")
				deps[key] = val
			}
		}
	}
	return deps
}

func parseProcfileEntrypoint(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	line := strings.TrimSpace(strings.SplitN(string(data), "\n", 2)[0])
	// e.g., web: node src/index.js or web: python app.py
	parts := strings.Fields(line)
	if len(parts) >= 2 {
		return parts[len(parts)-1]
	}
	return ""
}

func parseConfigEntrypoint(path string) string {
	if !fileExists(path) {
		return ""
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	re := regexp.MustCompile(`(?i)entry\s*[:=]\s*\"([^\"]+)\"`)
	m := re.FindStringSubmatch(string(data))
	if len(m) == 2 {
		return m[1]
	}
	return ""
}

func hasAnyDep(deps map[string]string, keys []string) bool {
	for _, k := range keys {
		if _, ok := deps[k]; ok {
			return true
		}
	}
	return false
}

func fileExists(p string) bool {
	st, err := os.Stat(p)
	return err == nil && !st.IsDir()
}

func generatePackageJSON(projectName, description, author, language string) string {
	ext := "js"
	if language == "ts" {
		ext = "ts"
	}

	return fmt.Sprintf(`{
  "name": "%s",
  "version": "1.0.0",
  "description": "%s",
  "main": "loader.%s",
  "type": "module",
  "scripts": {
    "start": "node --loader ts-node/esm loader.%s"
  },
  "dependencies": {
    "phantomkit": "^0.1.0"
  },
  "devDependencies": {
    "ts-node": "^10.9.2",
    "@types/node": "^20.0.0"
  },
  "keywords": ["phantomkit", "gitvault"],
  "author": "%s",
  "license": "MIT"
}
`, projectName, description, ext, ext, author)
}

func generateRequirementsTxt() string {
	return `phantomkit>=0.1.0
asyncio
`
}

func generateREADME(projectName, description, language string) string {
	ext := "js"
	if language == "ts" {
		ext = "ts"
	} else if language == "python" {
		ext = "py"
	}

	return fmt.Sprintf("# %s\n\n%s\n\n## Getting Started\n\nThis project was initialized with GitVault PhantomKit.\n\n### Prerequisites\n\n- GitVault account\n- PhantomKit CLI installed\n\n### Installation\n\n1. Install dependencies:\n   %s\n\n2. Set your PhantomKit API key:\n   ```bash\n   export PKIT_KEY=\"your-api-key-here\"\n   ```\n\n3. Run the project:\n   ```bash\n   %s\n   ```\n\n### Project Structure\n\n- `loader.%s` - Main entry point generated by PhantomKit\n- `phantom.toml` - PhantomKit configuration\n- `README.md` - This file\n\n## Uploading to GitVault\n\n1. Initialize Git repository:\n   ```bash\n   git init\n   git add .\n   git commit -m \"Initial commit\"\n   ```\n\n2. Upload to GitVault:\n   ```bash\n   phantom upload %s\n   ```\n\n## Learn More\n\n- [GitVault Documentation](https://docs.gitvault.io)\n- [PhantomKit Guide](https://docs.gitvault.io/phantomkit)\n- [Discord Community](https://discord.gg/NsatcWJ)\n", projectName, description, getInstallCommand(language), getRunCommand(language), ext, projectName)
}

func getInstallCommand(language string) string {
	switch language {
	case "js", "ts":
		return "```bash\nnpm install\n```"
	case "python":
		return "```bash\npip install -r requirements.txt\n```"
	default:
		return "```bash\n# install dependencies for your runtime\n```"
	}
}

func getRunCommand(language string) string {
	switch language {
	case "js":
		return "```bash\nnode loader.js\n```"
	case "ts":
		return "```bash\nnode --loader ts-node/esm loader.ts\n```"
	case "python":
		return "```bash\npython loader.py\n```"
	default:
		return "```bash\n# run your project\n```"
	}
}

func generateGitignore(language string) string {
	base := `# PhantomKit
phantom.lock
phantom.cache

# Environment variables
.env
.env.local

# IDE
.vscode/
.idea/
*.swp
*.swo

# OS
.DS_Store
Thumbs.db
`

	switch language {
	case "js", "ts":
		base += `
# Node.js
node_modules/
npm-debug.log*
yarn-debug.log*
yarn-error.log*
.npm
.yarn-integrity

# Build outputs
dist/
build/
*.tsbuildinfo
`
	case "python":
		base += `
# Python
__pycache__/
*.py[cod]
*$py.class
*.so
.Python
env/
venv/
ENV/
env.bak/
venv.bak/
.pytest_cache/
.coverage
`
	}

	return base
}

// Parse requirements.txt for Python dependencies
func parseRequirementsTxt(path string) map[string]string {
	deps := map[string]string{}
	f, err := os.Open(path)
	if err != nil {
		return deps
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			// Handle version specifiers like "package==1.0.0" or "package>=1.0.0"
			parts := strings.FieldsFunc(line, func(c rune) bool {
				return c == '=' || c == '>' || c == '<' || c == '~' || c == '!'
			})
			if len(parts) >= 1 {
				packageName := strings.TrimSpace(parts[0])
				version := "latest"
				if len(parts) >= 2 {
					version = strings.TrimSpace(parts[1])
				}
				deps[packageName] = version
			}
		}
	}
	return deps
}

// Parse Gemfile for Ruby dependencies
func parseGemfile(path string) map[string]string {
	deps := map[string]string{}
	f, err := os.Open(path)
	if err != nil {
		return deps
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if strings.HasPrefix(line, "gem ") {
			// Extract gem name and version from "gem 'name', 'version'"
			re := regexp.MustCompile(`gem\s+['"]([^'"]+)['"](?:,\s*['"]([^'"]+)['"])?`)
			matches := re.FindStringSubmatch(line)
			if len(matches) >= 2 {
				gemName := matches[1]
				version := "latest"
				if len(matches) >= 3 && matches[2] != "" {
					version = matches[2]
				}
				deps[gemName] = version
			}
		}
	}
	return deps
}

// Parse go.mod for Go dependencies
func parseGoMod(path string) map[string]string {
	deps := map[string]string{}
	f, err := os.Open(path)
	if err != nil {
		return deps
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if strings.HasPrefix(line, "require ") {
			// Extract module and version from "require module version"
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				module := parts[1]
				version := parts[2]
				deps[module] = version
			}
		}
	}
	return deps
}

// Parse composer.json for PHP dependencies
func parseComposerJSON(path string) map[string]string {
	deps := map[string]string{}
	data, err := os.ReadFile(path)
	if err != nil {
		return deps
	}
	var composer struct {
		Require map[string]string `json:"require"`
	}
	_ = json.Unmarshal(data, &composer)
	if composer.Require != nil {
		for k, v := range composer.Require {
			deps[k] = v
		}
	}
	return deps
}

// Parse Podfile for Swift dependencies
func parsePodfile(path string) map[string]string {
	deps := map[string]string{}
	f, err := os.Open(path)
	if err != nil {
		return deps
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if strings.HasPrefix(line, "pod ") {
			// Extract pod name and version from "pod 'Name', 'version'"
			re := regexp.MustCompile(`pod\s+['"]([^'"]+)['"](?:,\s*['"]([^'"]+)['"])?`)
			matches := re.FindStringSubmatch(line)
			if len(matches) >= 2 {
				podName := matches[1]
				version := "latest"
				if len(matches) >= 3 && matches[2] != "" {
					version = matches[2]
				}
				deps[podName] = version
			}
		}
	}
	return deps
}

// Generate phantom-lock file for environment setup
func generatePhantomLock(project, language string, allDeps map[string]map[string]string) string {
	var b strings.Builder
	b.WriteString("# Phantom Lock File\n")
	b.WriteString("# Generated for environment setup and dependency injection\n\n")
	b.WriteString(fmt.Sprintf("project = \"%s\"\n", project))
	b.WriteString(fmt.Sprintf("language = \"%s\"\n", language))
	b.WriteString(fmt.Sprintf("generated_at = \"%s\"\n\n", time.Now().UTC().Format(time.RFC3339)))
	
	// Environment setup instructions
	b.WriteString("[environment]\n")
	switch language {
	case "js", "ts":
		b.WriteString("setup = \"npm install\"\n")
		b.WriteString("runtime = \"node\"\n")
	case "python":
		b.WriteString("setup = \"pip install -r requirements.txt\"\n")
		b.WriteString("runtime = \"python\"\n")
	case "rust":
		b.WriteString("setup = \"cargo build\"\n")
		b.WriteString("runtime = \"cargo\"\n")
	case "ruby":
		b.WriteString("setup = \"bundle install\"\n")
		b.WriteString("runtime = \"ruby\"\n")
	case "go":
		b.WriteString("setup = \"go mod download\"\n")
		b.WriteString("runtime = \"go\"\n")
	case "php":
		b.WriteString("setup = \"composer install\"\n")
		b.WriteString("runtime = \"php\"\n")
	case "swift":
		b.WriteString("setup = \"pod install\"\n")
		b.WriteString("runtime = \"swift\"\n")
	default:
		b.WriteString("setup = \"# install dependencies for your runtime\"\n")
		b.WriteString("runtime = \"custom\"\n")
	}
	b.WriteString("\n")
	
	// Dependency snapshot
	b.WriteString("[dependencies]\n")
	for lang, deps := range allDeps {
		if len(deps) > 0 {
			b.WriteString(fmt.Sprintf("[dependencies.%s]\n", lang))
			for k, v := range deps {
				b.WriteString(fmt.Sprintf("%s = \"%s\"\n", k, v))
			}
			b.WriteString("\n")
		}
	}
	
	return b.String()
}

// Check if directory is empty
func isEmptyDir(dir string) bool {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return true
	}
	return len(entries) == 0
}

func generatePhantomConfig(projectName, language string) string {
	return fmt.Sprintf(`// PhantomKit Configuration
// Generated for project: %s

module.exports = {
  project: {
    name: '%s',
    language: '%s',
    version: '1.0.0'
  },
  
  phantomkit: {
    apiKey: process.env.PKIT_KEY,
    endpoint: process.env.PHANTOM_ENDPOINT || 'https://api.gitvault.io',
    cache: {
      enabled: true,
      ttl: 300000, // 5 minutes
      maxSize: 100 // MB
    }
  },
  
  storage: {
    type: 'gitvault',
    project: '%s'
  },
  
  runtime: {
    isolation: 'v8',
    timeout: 30000, // 30 seconds
    memory: 128 // MB
  }
};
`, projectName, projectName, language, projectName)
}

package main

import (
	"bufio"
	"embed"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var version = "1.0.0"

//go:embed templates/*
var templateFS embed.FS

type Config struct {
	Root       string
	ModuleName string
	Port       string
	Gitkeep    bool
	Clean      bool
}

var dirs = []string{
	"cmd",
	"commons/constants",
	"commons/error",
	"commons/utils",
	"config/constants",
	"config/env",
	"config/init",
	"recievers",
	"services/serviceName/service_init",
	"services/serviceName/data",
	"services/serviceName/internal",
	"services/serviceName/routes",
	"services/serviceName/utils",
}

func main() {
	interactive := flag.Bool("i", false, "Interactive mode")
	showVersion := flag.Bool("version", false, "Show tool version")
	root := flag.String("r", ".", "Target directory")
	moduleName := flag.String("m", "", "Go module name")
	port := flag.String("p", "8080", "Server port")
	gitkeep := flag.Bool("g", false, "Add .gitkeep files")
	clean := flag.Bool("c", false, "Clean target directory")
	flag.Parse()

	// Show version
	if *showVersion {
		fmt.Println("hexagen version", version)
		return
	}

	// Build config
	cfg := Config{
		Root:       *root,
		ModuleName: *moduleName,
		Port:       *port,
		Gitkeep:    *gitkeep,
		Clean:      *clean,
	}

	// Interactive mode
	if *interactive {
		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Project directory (default: .): ")
		if input, _ := reader.ReadString('\n'); strings.TrimSpace(input) != "" {
			cfg.Root = strings.TrimSpace(input)
		}

		fmt.Print("Go module name (github.com/user/project): ")
		if input, _ := reader.ReadString('\n'); strings.TrimSpace(input) != "" {
			cfg.ModuleName = strings.TrimSpace(input)
		}

		fmt.Printf("Server port (default: %s): ", cfg.Port)
		if input, _ := reader.ReadString('\n'); strings.TrimSpace(input) != "" {
			cfg.Port = strings.TrimSpace(input)
		}

		fmt.Print("Add .gitkeep files? (y/N): ")
		if input, _ := reader.ReadString('\n'); strings.ToLower(strings.TrimSpace(input)) == "y" {
			cfg.Gitkeep = true
		}

		fmt.Print("Clean target directory first? (y/N): ")
		if input, _ := reader.ReadString('\n'); strings.ToLower(strings.TrimSpace(input)) == "y" {
			cfg.Clean = true
		}
	}

	// Default module name
	if cfg.ModuleName == "" {
		cfg.ModuleName = "service.com/service"
	}

	// Generate project
	if err := generate(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n✓ Done! Your project is ready.")
	fmt.Printf("\nNext steps:\n")
	fmt.Printf("  cd %s\n", cfg.Root)
	fmt.Printf("  make run\n")
}

func generate(cfg Config) error {

	// Create root if missing
	if err := os.MkdirAll(cfg.Root, 0755); err != nil {
		return err
	}

	rootAbs, _ := filepath.Abs(cfg.Root)

	// Clean directory
	if cfg.Clean {
		entries, _ := os.ReadDir(rootAbs)
		for _, e := range entries {
			os.RemoveAll(filepath.Join(rootAbs, e.Name()))
		}
	}

	// Create directories
	for _, dir := range dirs {
		path := filepath.Join(rootAbs, dir)
		os.MkdirAll(path, 0755)
		if cfg.Gitkeep {
			os.WriteFile(filepath.Join(path, ".gitkeep"), []byte(""), 0644)
		}
	}

	// Generate files
	if err := writeGoMod(rootAbs, cfg.ModuleName); err != nil {
		return err
	}
	if err := writeMakefile(rootAbs, cfg.Port); err != nil {
		return err
	}
	if err := writeTemplate(rootAbs, "cmd/main.go", "main.go.tmpl", cfg); err != nil {
		return err
	}
	if err := writeTemplate(rootAbs, "services/serviceName/routes/router.go", "router.go.tmpl", cfg); err != nil {
		return err
	}
	if err := writeTemplate(rootAbs, "config/init/serverConfig.go", "serverConfig.go.tmpl", cfg); err != nil {
		return err
	}
	if err := writeTemplate(rootAbs, "commons/utils/logger.go", "logger.go.tmpl", cfg); err != nil {
		return err
	}

	return nil
}

func writeGoMod(root, module string) error {
	content := fmt.Sprintf(`module %s

go 1.22.0
`, module)
	return os.WriteFile(filepath.Join(root, "go.mod"), []byte(content), 0644)
}

func writeMakefile(root, port string) error {
	content := `PORT ?= ` + port + `

run:
	go run ./cmd/main.go

build:
	go build -o bin/app ./cmd/main.go

test:
	go test ./...

setup:
	go mod tidy
`
	return os.WriteFile(filepath.Join(root, "Makefile"), []byte(content), 0644)
}

func writeTemplate(root, outputPath, templateName string, cfg Config) error {

	// Load template file from embed
	tmplBytes, err := templateFS.ReadFile("templates/" + templateName)
	if err != nil {
		return err
	}

	// Parse template
	tmpl, err := template.New(templateName).Parse(string(tmplBytes))
	if err != nil {
		return err
	}

	// Ensure directory exists
	outPath := filepath.Join(root, outputPath)
	os.MkdirAll(filepath.Dir(outPath), 0755)

	// Create file
	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Replace template vars
	data := map[string]string{
		"MODULE": cfg.ModuleName,
		"PORT":   cfg.Port,
	}

	// Execute template
	return tmpl.Execute(f, data)
}
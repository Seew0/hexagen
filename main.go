package main

import (
	"bufio"
	"embed"
	"flag"
	"fmt"
	"os"
	"os/exec"
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

	if *showVersion {
		fmt.Println("hexagen version", version)
		return
	}

	cfg := Config{
		Root:       *root,
		ModuleName: *moduleName,
		Port:       *port,
		Gitkeep:    *gitkeep,
		Clean:      *clean,
	}

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

	if cfg.ModuleName == "" {
		cfg.ModuleName = "service.com/service"
	}

	if err := generate(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n✓ Project structure created successfully!")
	fmt.Println("⏳ Installing dependencies...")

	if err := installDependencies(cfg.Root); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to install dependencies: %v\n", err)
		fmt.Println("You can manually run: go mod tidy")
	} else {
		fmt.Println("✓ Dependencies installed successfully!")
	}

	fmt.Println("\n✓ Done! Your project is ready.")
	fmt.Printf("\nNext steps:\n")
	fmt.Printf("  cd %s\n", cfg.Root)
	fmt.Printf("  make run\n")
}

func generate(cfg Config) error {
	if err := os.MkdirAll(cfg.Root, 0755); err != nil {
		return err
	}

	rootAbs, _ := filepath.Abs(cfg.Root)

	if cfg.Clean {
		entries, _ := os.ReadDir(rootAbs)
		for _, e := range entries {
			os.RemoveAll(filepath.Join(rootAbs, e.Name()))
		}
	}

	for _, dir := range dirs {
		path := filepath.Join(rootAbs, dir)
		os.MkdirAll(path, 0755)
		if cfg.Gitkeep {
			_ = os.WriteFile(filepath.Join(path, ".gitkeep"), []byte(""), 0644)
		}
	}

	if err := writeGoMod(rootAbs, cfg.ModuleName); err != nil {
		return err
	}
	if err := writeMakefile(rootAbs, cfg.Port); err != nil {
		return err
	}

	if err := writeTemplate(rootAbs, "cmd/main.go", "templates/app.go.tmpl", cfg); err != nil {
		return err
	}
	if err := writeTemplate(rootAbs, "services/serviceName/routes/router.go", "templates/router.go.tmpl", cfg); err != nil {
		return err
	}
	if err := writeTemplate(rootAbs, "config/init/serverConfig.go", "templates/serverConfig.go.tmpl", cfg); err != nil {
		return err
	}
	if err := writeTemplate(rootAbs, "commons/utils/logger.go", "templates/logger.go.tmpl", cfg); err != nil {
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

func writeTemplate(root, outputPath, templatePath string, cfg Config) error {
	tmplBytes, err := templateFS.ReadFile(templatePath)
	if err != nil {
		return err
	}

	tmpl, err := template.New(filepath.Base(templatePath)).Parse(string(tmplBytes))
	if err != nil {
		return err
	}

	outPath := filepath.Join(root, outputPath)
	_ = os.MkdirAll(filepath.Dir(outPath), 0755)

	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()

	data := map[string]string{
		"MODULE": cfg.ModuleName,
		"PORT":   cfg.Port,
	}

	return tmpl.Execute(f, data)
}

func installDependencies(root string) error {
	rootAbs, err := filepath.Abs(root)
	if err != nil {
		return err
	}

	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = rootAbs
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

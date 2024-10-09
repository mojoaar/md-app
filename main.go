package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

const (
	version = "1.0.0"
	author  = "Morten Johansen (mojoaar)"
)

type Template struct {
	Content string `yaml:"content"`
}

const defaultTemplateContent = `content: |
  # {{TITLE}}

  Date: {{DATE}}  
  Time: {{TIME}}

  ## Introduction

  ## Main Content

  ## Conclusion
`

func main() {
	// Define command-line flags
	helpFlag := flag.Bool("help", false, "Show help information")
	versionFlag := flag.Bool("version", false, "Show version information")
	versionShortFlag := flag.Bool("v", false, "Show version information (short flag)")
	operationType := flag.String("type", "", "Operation type: 'template' or 'note'")
	name := flag.String("name", "", "Name of the markdown file or template (without extension, optional for notes)")
	title := flag.String("title", "", "Title of the markdown document (required for notes)")
	templateFile := flag.String("template", "default", "Name of the template file (without extension, only for notes)")
	showTemplates := flag.Bool("show", false, "Show all available template files (only for template type)")

	// Custom usage message
	flag.Usage = func() {
		fmt.Printf("Markdown File Creator v%s\n", version)
		fmt.Printf("Author: %s\n\n", author)
		fmt.Println("Usage:")
		fmt.Println("  Create a new template:")
		fmt.Println("    md -type template -name <template_name>")
		fmt.Println("  Show all available templates:")
		fmt.Println("    md -type template -show")
		fmt.Println("  Create a new note:")
		fmt.Println("    md -type note -title <note_title> [-name <note_name>] [-template <template_name>]")
		fmt.Println("\nFlags:")
		flag.PrintDefaults()
	}

	flag.Parse()

	// Handle help and version flags
	if *helpFlag {
		flag.Usage()
		os.Exit(0)
	}

	if *versionFlag || *versionShortFlag {
		fmt.Printf("Markdown File Creator v%s\n", version)
		fmt.Printf("Author: %s\n", author)
		os.Exit(0)
	}

	// Validate input
	if *operationType != "template" && *operationType != "note" {
		fmt.Println("Error: type must be either 'template' or 'note'")
		flag.Usage()
		os.Exit(1)
	}

	// Ensure templates directory exists and create default template if needed
	err := ensureTemplatesDirectory()
	if err != nil {
		fmt.Printf("Error ensuring templates directory: %v\n", err)
		os.Exit(1)
	}

	if *operationType == "template" {
		if *showTemplates {
			err = showTemplateFiles()
		} else if *name != "" {
			// Convert spaces to dashes and make lowercase
			templateName := strings.ToLower(strings.ReplaceAll(*name, " ", "-"))
			err = createTemplate(templateName)
		} else {
			fmt.Println("Error: for template type, either -show or -name must be specified")
			flag.Usage()
			os.Exit(1)
		}
	} else {
		if *title == "" {
			fmt.Println("Error: title is required for notes")
			flag.Usage()
			os.Exit(1)
		}
		noteName := *name
		if noteName == "" {
			noteName = *title
		}
		noteName = strings.ToLower(strings.ReplaceAll(noteName, " ", "-"))
		err = createNote(noteName, *title, *templateFile)
	}

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func ensureTemplatesDirectory() error {
	templatesDir := "templates"

	// Check if templates directory exists
	if _, err := os.Stat(templatesDir); os.IsNotExist(err) {
		// Create templates directory
		err := os.Mkdir(templatesDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create templates directory: %v", err)
		}
		fmt.Println("Created templates directory.")

		// Create default.yaml file
		defaultFilePath := filepath.Join(templatesDir, "default.yaml")
		err = os.WriteFile(defaultFilePath, []byte(defaultTemplateContent), 0644)
		if err != nil {
			return fmt.Errorf("failed to create default template file: %v", err)
		}
		fmt.Println("Created default template file.")
	}

	return nil
}

func showTemplateFiles() error {
	templatesDir := "templates"
	files, err := os.ReadDir(templatesDir)
	if err != nil {
		return fmt.Errorf("failed to read templates directory: %v", err)
	}

	fmt.Println("Available template files:")
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".yaml") {
			fmt.Println("-", strings.TrimSuffix(file.Name(), ".yaml"))
		}
	}

	return nil
}

func createTemplate(name string) error {
	filename := filepath.Join("templates", name+".yaml")
	err := os.WriteFile(filename, []byte(defaultTemplateContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to create template file: %v", err)
	}

	fmt.Printf("Template file '%s.yaml' created successfully.\n", name)
	return nil
}

func createNote(name, title, templateFile string) error {
	// Load template
	template, err := loadTemplate(templateFile)
	if err != nil {
		return fmt.Errorf("error loading template: %v", err)
	}

	// Generate content
	content := generateContent(template, title)

	// Save markdown file
	err = saveMarkdownFile(name, content)
	if err != nil {
		return fmt.Errorf("error saving markdown file: %v", err)
	}

	fmt.Printf("Markdown note '%s.md' created successfully.\n", name)
	return nil
}

func loadTemplate(templateName string) (string, error) {
	filename := filepath.Join("templates", templateName+".yaml")
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}

	var template Template
	err = yaml.Unmarshal(data, &template)
	if err != nil {
		return "", err
	}

	return template.Content, nil
}

func generateContent(template, title string) string {
	now := time.Now()
	content := strings.ReplaceAll(template, "{{TITLE}}", title)
	content = strings.ReplaceAll(content, "{{DATE}}", now.Format("2006-01-02"))
	content = strings.ReplaceAll(content, "{{TIME}}", now.Format("15:04:05"))
	return content
}

func saveMarkdownFile(name, content string) error {
	filename := name + ".md"
	return os.WriteFile(filename, []byte(content), 0644)
}

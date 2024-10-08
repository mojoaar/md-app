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

type Template struct {
	Content string `yaml:"content"`
}

const defaultTemplateContent = `content: |
  # {{TITLE}}

  Date: {{DATE}}

  ## Introduction

  ## Main Content

  ## Conclusion
`

func main() {
	// Define command-line flags
	operationType := flag.String("type", "", "Operation type: 'template' or 'post'")
	name := flag.String("name", "", "Name of the markdown file or template (without extension)")
	title := flag.String("title", "", "Title of the markdown document (only for posts)")
	templateFile := flag.String("template", "default", "Name of the template file (without extension, only for posts)")
	flag.Parse()

	// Validate input
	if *operationType != "template" && *operationType != "post" {
		fmt.Println("Error: type must be either 'template' or 'post'")
		flag.Usage()
		os.Exit(1)
	}

	if *name == "" {
		fmt.Println("Error: name is required")
		flag.Usage()
		os.Exit(1)
	}

	if *operationType == "post" && *title == "" {
		fmt.Println("Error: title is required for posts")
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
		err = createTemplate(*name)
	} else {
		err = createPost(*name, *title, *templateFile)
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

func createTemplate(name string) error {
	filename := filepath.Join("templates", name+".yaml")
	err := os.WriteFile(filename, []byte(defaultTemplateContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to create template file: %v", err)
	}

	fmt.Printf("Template file '%s.yaml' created successfully.\n", name)
	return nil
}

func createPost(name, title, templateFile string) error {
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

	fmt.Printf("Markdown file '%s.md' created successfully.\n", name)
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
	content := strings.ReplaceAll(template, "{{TITLE}}", title)
	content = strings.ReplaceAll(content, "{{DATE}}", time.Now().Format("2006-01-02"))
	return content
}

func saveMarkdownFile(name, content string) error {
	filename := name + ".md"
	return os.WriteFile(filename, []byte(content), 0644)
}

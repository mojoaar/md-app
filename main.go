package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
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

// Custom error types
type ValidationError struct {
	Field string
	Msg   string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("Validation error for %s: %s", e.Field, e.Msg)
}

type FileError struct {
	Op   string
	Path string
	Err  error
}

func (e *FileError) Error() string {
	return fmt.Sprintf("%s error for file %s: %v", e.Op, e.Path, e.Err)
}

func main() {
	// Define command-line flags
	helpFlag := flag.Bool("help", false, "Show help information")
	helpShortFlag := flag.Bool("h", false, "Show help information (short flag)")
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
	if *helpFlag || *helpShortFlag {
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
		handleError(err)
	}

	if *operationType == "template" {
		if *showTemplates {
			err = showTemplateFiles()
		} else if *name != "" {
			err = createTemplate(*name)
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
		err = createNote(noteName, *title, *templateFile)
	}

	if err != nil {
		handleError(err)
	}
}

func handleError(err error) {
	switch e := err.(type) {
	case *ValidationError:
		fmt.Printf("Validation error: %v\n", e)
	case *FileError:
		fmt.Printf("File operation error: %v\n", e)
	default:
		fmt.Printf("An error occurred: %v\n", e)
	}
	os.Exit(1)
}

func ensureTemplatesDirectory() error {
	templatesDir := "templates"

	if _, err := os.Stat(templatesDir); os.IsNotExist(err) {
		err := os.Mkdir(templatesDir, 0755)
		if err != nil {
			return &FileError{Op: "create", Path: templatesDir, Err: err}
		}
		fmt.Println("Created templates directory.")

		defaultFilePath := filepath.Join(templatesDir, "default.yaml")
		err = os.WriteFile(defaultFilePath, []byte(defaultTemplateContent), 0644)
		if err != nil {
			return &FileError{Op: "write", Path: defaultFilePath, Err: err}
		}
		fmt.Println("Created default template file.")
	}

	return nil
}

func showTemplateFiles() error {
	templatesDir := "templates"
	files, err := os.ReadDir(templatesDir)
	if err != nil {
		return &FileError{Op: "read", Path: templatesDir, Err: err}
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
	if err := validateFileName(name); err != nil {
		return err
	}

	filename := filepath.Join("templates", sanitizeFileName(name)+".yaml")
	err := os.WriteFile(filename, []byte(defaultTemplateContent), 0644)
	if err != nil {
		return &FileError{Op: "write", Path: filename, Err: err}
	}

	fmt.Printf("Template file '%s.yaml' created successfully.\n", name)
	return nil
}

func createNote(name, title, templateFile string) error {
	if err := validateFileName(name); err != nil {
		return err
	}
	if err := validateTitle(title); err != nil {
		return err
	}

	template, err := loadTemplate(templateFile)
	if err != nil {
		return err
	}

	content := generateContent(template, title)

	err = saveMarkdownFile(name, content)
	if err != nil {
		return err
	}

	fmt.Printf("Markdown note '%s.md' created successfully.\n", name)
	return nil
}

func loadTemplate(templateName string) (string, error) {
	filename := filepath.Join("templates", sanitizeFileName(templateName)+".yaml")
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", &FileError{Op: "read", Path: filename, Err: err}
	}

	var template Template
	err = yaml.Unmarshal(data, &template)
	if err != nil {
		return "", &FileError{Op: "parse", Path: filename, Err: err}
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
	filename := sanitizeFileName(name) + ".md"
	err := os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		return &FileError{Op: "write", Path: filename, Err: err}
	}
	return nil
}

func sanitizeFileName(name string) string {
	// Replace spaces with underscores
	name = strings.ReplaceAll(name, " ", "_")

	// Remove any character that is not alphanumeric, underscore, or hyphen
	reg := regexp.MustCompile(`[^a-zA-Z0-9_-]`)
	name = reg.ReplaceAllString(name, "")

	// Convert to lowercase
	name = strings.ToLower(name)

	// Trim to a reasonable length (e.g., 255 characters)
	if len(name) > 255 {
		name = name[:255]
	}

	return name
}

func validateFileName(name string) error {
	if name == "" {
		return &ValidationError{Field: "file name", Msg: "cannot be empty"}
	}
	if len(name) > 255 {
		return &ValidationError{Field: "file name", Msg: "too long (max 255 characters)"}
	}
	if strings.ContainsAny(name, "/\\:*?\"<>|") {
		return &ValidationError{Field: "file name", Msg: "contains invalid characters"}
	}
	return nil
}

func validateTitle(title string) error {
	if title == "" {
		return &ValidationError{Field: "title", Msg: "cannot be empty"}
	}
	if len(title) > 100 {
		return &ValidationError{Field: "title", Msg: "too long (max 100 characters)"}
	}
	if strings.TrimSpace(title) == "" {
		return &ValidationError{Field: "title", Msg: "cannot be only whitespace"}
	}
	return nil
}

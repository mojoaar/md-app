package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

const (
	version = "1.0.0"
	author  = "Morten Johansen (mojoaar)"
)

type Template struct {
	Content string   `yaml:"content"`
	Tags    []string `yaml:"tags,omitempty"`
}

const defaultTemplateContent = `content: |
  # {{TITLE}}

  Date: {{DATE}}
  Time: {{TIME}}
  Tags: {{TAGS}}

  ## Introduction

  ## Main Content

  ## Conclusion

tags:
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

var rootCmd = &cobra.Command{
	Use:     "md",
	Short:   "Markdown File Creator",
	Version: version,
	Long: fmt.Sprintf(`Markdown File Creator v%s
Author: %s

A tool for creating markdown files and managing templates.

Usage Examples:
  Create a new template:
    md template create my-template

  List all templates:
    md template list

  Create a new note using the default template:
    md note -t "My Note Title"

  Create a new note with a custom name and template:
    md note -t "My Note Title" -n my-custom-note -m my-template

  Create a new note with tags:
    md note -t "My Note Title" -g tag1,tag2,tag3

  List all notes with their tags:
    md list`, version, author),
}

var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Manage templates",
	Long: `Create or list template files for markdown notes.

Usage Examples:
  Create a new template:
    md template create my-template

  List all templates:
    md template list`,
}

var createTemplateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new template",
	Long: `Create a new template file for markdown notes.

Usage Example:
  md template create my-template`,
	Args: cobra.ExactArgs(1),
	RunE: createTemplate,
}

var listTemplatesCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available templates",
	Long: `List all available template files for markdown notes.

Usage Example:
  md template list`,
	RunE: listTemplates,
}

var noteCmd = &cobra.Command{
	Use:   "note",
	Short: "Create a new note",
	Long: `Create a new markdown note using a specified template.

Usage Examples:
  Create a note with the default template:
    md note -t "My Note Title"

  Create a note with a custom name and template:
    md note -t "My Note Title" -n my-custom-note -m my-template

  Create a note with tags:
    md note -t "My Note Title" -g tag1,tag2,tag3`,
	RunE: createNote,
}

var listNotesCmd = &cobra.Command{
	Use:   "list",
	Short: "List all notes with their tags",
	Long:  `List all markdown notes in the notes directory with their associated tags.`,
	RunE:  listNotes,
}

var (
	noteTitle    string
	noteName     string
	templateName string
	noteTags     []string
)

func init() {
	rootCmd.AddCommand(templateCmd, noteCmd, listNotesCmd)
	templateCmd.AddCommand(createTemplateCmd, listTemplatesCmd)

	noteCmd.Flags().StringVarP(&noteTitle, "title", "t", "", "Title of the markdown document (required)")
	noteCmd.Flags().StringVarP(&noteName, "name", "n", "", "Name of the markdown file (optional, defaults to title)")
	noteCmd.Flags().StringVarP(&templateName, "template", "m", "default", "Name of the template to use")
	noteCmd.Flags().StringSliceVarP(&noteTags, "tags", "g", []string{}, "Tags for the note (comma-separated)")
	noteCmd.MarkFlagRequired("title")
}

func main() {
	if err := LoadConfig(); err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func ensureTemplatesDirectory() error {
	if _, err := os.Stat(AppConfig.TemplatesDir); os.IsNotExist(err) {
		err := os.Mkdir(AppConfig.TemplatesDir, 0755)
		if err != nil {
			return &FileError{Op: "create", Path: AppConfig.TemplatesDir, Err: err}
		}
		fmt.Printf("Created templates directory: %s\n", AppConfig.TemplatesDir)

		// Create default template
		defaultTemplatePath := filepath.Join(AppConfig.TemplatesDir, "default.yaml")
		err = os.WriteFile(defaultTemplatePath, []byte(defaultTemplateContent), 0644)
		if err != nil {
			return &FileError{Op: "write", Path: defaultTemplatePath, Err: err}
		}
		fmt.Println("Created default template file: default.yaml")
	}
	return nil
}

func ensureNotesDirectory() error {
	if AppConfig.NotesDir == "" {
		return nil // No notes directory specified, we'll use the current directory
	}

	if _, err := os.Stat(AppConfig.NotesDir); os.IsNotExist(err) {
		err := os.MkdirAll(AppConfig.NotesDir, 0755)
		if err != nil {
			return &FileError{Op: "create", Path: AppConfig.NotesDir, Err: err}
		}
		fmt.Printf("Created notes directory: %s\n", AppConfig.NotesDir)
	}
	return nil
}

func createTemplate(cmd *cobra.Command, args []string) error {
	name := args[0]
	if err := validateFileName(name); err != nil {
		return err
	}

	if err := ensureTemplatesDirectory(); err != nil {
		return err
	}

	filename := filepath.Join(AppConfig.TemplatesDir, sanitizeFileName(name)+".yaml")
	err := os.WriteFile(filename, []byte(defaultTemplateContent), 0644)
	if err != nil {
		return &FileError{Op: "write", Path: filename, Err: err}
	}

	fmt.Printf("Template file '%s.yaml' created successfully.\n", name)
	return nil
}

func listTemplates(cmd *cobra.Command, args []string) error {
	if err := ensureTemplatesDirectory(); err != nil {
		return err
	}

	files, err := os.ReadDir(AppConfig.TemplatesDir)
	if err != nil {
		return &FileError{Op: "read", Path: AppConfig.TemplatesDir, Err: err}
	}

	templateFiles := []string{}
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".yaml") {
			templateFiles = append(templateFiles, strings.TrimSuffix(file.Name(), ".yaml"))
		}
	}

	if len(templateFiles) == 0 {
		fmt.Println("No template files found.")
	} else {
		fmt.Println("Available template files:")
		for _, name := range templateFiles {
			fmt.Println("-", name)
		}
	}

	return nil
}

func createNote(cmd *cobra.Command, args []string) error {
	if noteName == "" {
		noteName = noteTitle
	}

	if err := validateFileName(noteName); err != nil {
		return err
	}
	if err := validateTitle(noteTitle); err != nil {
		return err
	}

	if err := ensureNotesDirectory(); err != nil {
		return err
	}

	template, err := loadTemplate(templateName)
	if err != nil {
		return err
	}

	content := generateContent(template, noteTitle, noteTags)

	err = saveMarkdownFile(noteName, content)
	if err != nil {
		return err
	}

	if len(noteTags) > 0 {
		fmt.Printf("Markdown note '%s.md' created successfully with tags: %s\n", noteName, strings.Join(noteTags, ", "))
	} else {
		fmt.Printf("Markdown note '%s.md' created successfully\n", noteName)
	}
	return nil
}

func loadTemplate(templateName string) (Template, error) {
	if err := ensureTemplatesDirectory(); err != nil {
		return Template{}, err
	}

	filename := filepath.Join(AppConfig.TemplatesDir, sanitizeFileName(templateName)+".yaml")
	data, err := os.ReadFile(filename)
	if err != nil {
		return Template{}, &FileError{Op: "read", Path: filename, Err: err}
	}

	var template Template
	err = yaml.Unmarshal(data, &template)
	if err != nil {
		return Template{}, &FileError{Op: "parse", Path: filename, Err: err}
	}

	return template, nil
}

func generateContent(template Template, title string, tags []string) string {
	now := time.Now()
	content := strings.ReplaceAll(template.Content, "{{TITLE}}", title)
	content = strings.ReplaceAll(content, "{{DATE}}", now.Format("2006-01-02"))
	content = strings.ReplaceAll(content, "{{TIME}}", now.Format("15:04:05"))

	tagsString := strings.Join(tags, ", ")
	content = strings.ReplaceAll(content, "{{TAGS}}", tagsString)

	return content
}

func saveMarkdownFile(name, content string) error {
	filename := sanitizeFileName(name) + ".md"
	if AppConfig.NotesDir != "" {
		filename = filepath.Join(AppConfig.NotesDir, filename)
	}
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

func listNotes(cmd *cobra.Command, args []string) error {
	if err := ensureNotesDirectory(); err != nil {
		return err
	}

	notesDir := AppConfig.NotesDir
	if notesDir == "" {
		notesDir = "."
	}

	files, err := os.ReadDir(notesDir)
	if err != nil {
		return &FileError{Op: "read", Path: notesDir, Err: err}
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".md") {
			content, err := os.ReadFile(filepath.Join(notesDir, file.Name()))
			if err != nil {
				fmt.Printf("Error reading %s: %v\n", file.Name(), err)
				continue
			}

			tags := extractTags(string(content))
			fmt.Printf("%s: %v\n", file.Name(), tags)
		}
	}

	return nil
}

func extractTags(content string) []string {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Tags:") {
			tagsPart := strings.TrimPrefix(line, "Tags:")
			tags := strings.Split(tagsPart, ",")
			for i, tag := range tags {
				tags[i] = strings.TrimSpace(tag)
			}
			return tags
		}
	}
	return []string{}
}

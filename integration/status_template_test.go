package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/je-sidestuff/orgonization/templates"
)

func TestStatusTemplateProcessor(t *testing.T) {
	t.Parallel()

	// Create a new instance of FilesystemTestInput
	fti := NewFilesystemTestInput()

	// Create a temporary folder for our test
	tempDir, err := fti.CreateAndMapTempFolder("status_template_test", "status_template_test_")
	if err != nil {
		t.Fatalf("Failed to create temp folder: %v", err)
	}

	// Delete the temporary directory to clean up when we are done
	defer RunTestStage(t, "cleanup", func() { fti.CleanupFilesystemTestInput() }, func() {
		t.Logf("Skipping cleanup of directory %s", tempDir)
	})

	// Create test files with PING tokens and misc text
	testFiles := map[string]string{
		"status.txt": `This is a status file with some misc text.
It contains a <PING:> token that should be replaced.
Here's another line with more content.
And another <PING:> token for good measure.
End of file.`,
		"notes.md": `# Notes File

This markdown file has some content and a <PING:> token.

## Section 2

More content here with no tokens.

## Section 3

Final section with a <PING:> at the end.`,
		"config.yaml": `# Configuration file
setting1: value1
ping_token: "<PING:>"
setting2: value2
# End of config`,
	}

	// Write the test files
	for filename, content := range testFiles {
		filePath := filepath.Join(tempDir, filename)
		err := os.WriteFile(filePath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to write test file %s: %v", filename, err)
		}
	}

	// Create a filesystem configuration pointing to our temp directory
	// Since newFilesystemConfiguration is not exported, we'll use a different approach
	// by creating a simple implementation for testing
	fsConfig := &testFilesystemConfiguration{
		name:                  "test_config",
		inputFileFolderPaths:  []string{},
		outputFileFolderPaths: []string{tempDir},
		configPath:            filepath.Join(tempDir, "config.yaml"),
		orgoDirectoryPath:     tempDir,
	}

	// Create a template processor
	processor := templates.NewTemplateProcessor("status_template_test_processor")

	// Initialize the filesystem
	err = processor.InitializeFilesystem(fsConfig)
	if err != nil {
		t.Fatalf("Failed to initialize filesystem: %v", err)
	}

	// Create directive configuration to match <PING:> tokens and replace them
	directive := templates.DirectiveConfiguration{
		Name: "ping_replacer",
		DirectiveTriggers: []templates.DirectiveTrigger{
			{
				Name: "ping_token_match",
				Class: templates.DirectiveTriggerClass{
					Name: "ReqRepTokenMatch",
					AllowedArgs: []templates.DirectiveTriggerArg{
						templates.Pattern,
					},
				},
				Args: map[templates.DirectiveTriggerArg]string{
					templates.Pattern: `<PING:>`,
				},
			},
		},
		DirectiveActions: []templates.DirectiveAction{
			{
				Name: "replace_ping_tokens",
				Class: templates.DirectiveActionClass{
					Name: "ReplaceTokenAction",
					Targets: []templates.DirectiveTriggerTarget{
						templates.MatchedToken,
					},
					Effects: []templates.DirectiveActionEffect{
						templates.ReplaceToken,
					},
				},
			},
		},
	}

	// Initialize directives
	err = processor.InitializeDirectives([]templates.DirectiveConfiguration{directive})
	if err != nil {
		t.Fatalf("Failed to initialize directives: %v", err)
	}

	// Complete initialization
	err = processor.CompleteInitialization()
	if err != nil {
		t.Fatalf("Failed to complete initialization: %v", err)
	}

	// Verify initial state - files should contain <PING:> tokens
	for filename := range testFiles {
		filePath := filepath.Join(tempDir, filename)
		content, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("Failed to read file %s before processing: %v", filename, err)
		}

		// Count occurrences of <PING:> in original content
		originalContent := string(content)
		pingCount := countSubstring(originalContent, "<PING:>")
		if pingCount == 0 {
			t.Errorf("File %s should contain <PING:> tokens before processing", filename)
		}
		t.Logf("File %s contains %d <PING:> tokens before processing", filename, pingCount)
	}

	// Tick the template processor to trigger replacements
	err = processor.Tick()
	if err != nil {
		t.Fatalf("Failed to tick template processor: %v", err)
	}

	// Verify that <PING:> tokens have been replaced with <PING:PONG>
	for filename := range testFiles {
		filePath := filepath.Join(tempDir, filename)
		content, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("Failed to read file %s after processing: %v", filename, err)
		}

		processedContent := string(content)

		// Verify no <PING:> tokens remain
		pingCount := countSubstring(processedContent, "<PING:>")
		if pingCount > 0 {
			t.Errorf("File %s still contains %d unreplaced <PING:> tokens after processing", filename, pingCount)
		}

		// Verify <PING:PONG> tokens exist
		pongCount := countSubstring(processedContent, "<PING:PONG>")
		if pongCount == 0 {
			t.Errorf("File %s should contain <PING:PONG> tokens after processing", filename)
		}

		t.Logf("File %s now contains %d <PING:PONG> tokens after processing", filename, pongCount)
	}
}

// Helper function to count occurrences of a substring
func countSubstring(text, substring string) int {
	count := 0
	start := 0
	for {
		pos := strings.Index(text[start:], substring)
		if pos == -1 {
			break
		}
		count++
		start += pos + len(substring)
	}
	return count
}

// testFilesystemConfiguration is a simple implementation of FilesystemConfiguration for testing
type testFilesystemConfiguration struct {
	name                  string
	inputFileFolderPaths  []string
	outputFileFolderPaths []string
	configPath            string
	orgoDirectoryPath     string
}

func (t *testFilesystemConfiguration) GetName() string {
	return t.name
}

func (t *testFilesystemConfiguration) GetInputFileFolderPaths() []string {
	return t.inputFileFolderPaths
}

func (t *testFilesystemConfiguration) GetOutputFileFolderPaths() []string {
	return t.outputFileFolderPaths
}

func (t *testFilesystemConfiguration) GetConfigPath() string {
	return t.configPath
}

func (t *testFilesystemConfiguration) GetOrgoDirectoryPath() string {
	return t.orgoDirectoryPath
}
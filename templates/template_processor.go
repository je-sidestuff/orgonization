package templates

import (
	"log/slog"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
)

// TemplateProcessor interface - the public interface of the template processor.
//
// Used to coordinate the management of directives and processing of files. The template
// processor is the central management object for directives and file processing. It will
// be acted on externally by the agent and server to orchestrate the "passage of time".
type TemplateProcessor interface {
	GetName() string
	InitializeFilesystem(FilesystemConfiguration FilesystemConfiguration) error
	InitializeDirectives(DirectiveConfigurations []DirectiveConfiguration) error
	CompleteInitialization() error
	Tick() error
	internalTemplateProcessorConstructor()
}

// Private struct that implements the interface
type templateProcessor struct {
	name                            string
	directives                      map[string]DirectiveConfiguration
	inputFileFolderPaths            []string
	outputFileFolderPaths           []string
	fileProcessors                  []FileProcessor
	allPathsToFiles                 map[string]string
	activePathsToFileContentStrings map[string]string
	matchRegexes                    map[string]*regexp.Regexp
	logger                          *slog.Logger
}

type MappedToken struct {
	TokenString string
	StartIndex  int
	EndIndex    int
}

// CompleteInitialization implements TemplateProcessor.
func (tp *templateProcessor) CompleteInitialization() error {

	// Traverse our config and input files to load initial directives (soon)

	return nil
}

// Tick implements TemplateProcessor.
func (tp *templateProcessor) Tick() error {

	// Log with structured fields
	tp.logger.Debug("Processing tick.")

	activeFilePathsModified := []string{}

	// Identify the set of files with changes and pull them into memory
	// (files that have not yet been examined yet, creations, and deletions - also keep the full set in case we need a refresh)
	for _, fileProcessor := range tp.fileProcessors {
		fileProcessor.TraverseFiles()
	}

	// Keep a map of directive names to token patterns so that we can
	// match the tokens within the context of the directives they belong to and come back to them to process
	directivesToSearchedTokens := make(map[string]string)
	// directivesToSearchedTokens := make(map[string][]string)

	// Read directives to gather set of tokens to match
	for directiveKey, directive := range tp.directives {

		tp.logger.Debug("Processing directive", slog.String("name", directiveKey))

		// Read our triggers to detect tokens to look for
		for _, trigger := range directive.DirectiveTriggers {

			// If the trigger is type ReqRepTokenMatch then we will
			// remember to look for this token when we search input files
			if trigger.Class.Name == "ReqRepTokenMatch" {

				// Look for the pattern argument
				// Iterate through trigger.Args and if the key is 'pattern' then store the value against this directive
				for key, value := range trigger.Args {
					if key == Pattern {

						tp.logger.Debug("Processing pattern", slog.String("pattern", value))

						directivesToSearchedTokens[directive.Name] = value
					}
				}
			}
		}
	}

	// Keep a map of directive names to token patterns so that we can
	// match the tokens within the context of the directives they belong to and come back to them to process
	// We need a map of directives to a map of map-per-file
	// The inner map contains strings as the key of startChar:endChar
	directivesToFilesToMatchedTokens := make(map[string]map[string]map[string]string)

	// Traverse all of the active files and check for token matches
	for filePath, fileContent := range tp.activePathsToFileContentStrings {

		// Go through our map of directives to searched tokens
		for directiveName, searchedToken := range directivesToSearchedTokens {

			// Get the regex for the token
			var regex *regexp.Regexp

			// If we have this regex in matchRegexes, then get it
			// Otherwise add it
			if _, ok := tp.matchRegexes[searchedToken]; ok {
				regex = tp.matchRegexes[searchedToken]
			} else {
				var err error
				regex, err = regexp.Compile(searchedToken)
				if err != nil {
					return err
				}
				tp.matchRegexes[searchedToken] = regex
			}

			// Find all matches of the token type in the input string, as well as the strings
			matches := regex.FindAllStringIndex(fileContent, -1)
			strings := regex.FindAllString(fileContent, -1)

			// If we have one or more matches then populate this filename key
			if len(matches) > 0 {

				// If we don't have a map for this directive yet, create it
				if _, ok := directivesToFilesToMatchedTokens[directiveName]; !ok {
					directivesToFilesToMatchedTokens[directiveName] = make(map[string]map[string]string)
				}

				directivesToFilesToMatchedTokens[directiveName][filePath] = make(map[string]string)
			}

			tp.logger.Debug("Finding matches for directive",
				slog.String("directive", directiveName),
				slog.String("file", filePath),
				slog.Int("matches", len(matches)))

			// Add each match to the MappedTokens list, increment a counter as we go to collect the corresponding string
			idx := 0
			for _, match := range matches {

				// The key is startIndex:endIndex
				key := strconv.Itoa(match[0]) + ":" + strconv.Itoa(match[1])

				directivesToFilesToMatchedTokens[directiveName][filePath][key] = strings[idx]
				idx++
			}
		}
	}

	// Process directives
	for _, directive := range tp.directives {

		// Mark a boolean indicating that this dirtective has not yet been triggered to start
		directiveTriggered := false

		// Check the trigger set to see if the action should occur
		// (For now we'll do it in a dumb way and just check if any have triggered)
		for _, trigger := range directive.DirectiveTriggers {

			// If the trigger is type always, then always execute
			// (For now to check if we have an always trigger we can just check if the
			// DirectiveTriggerClass has the name 'AlwaysExecute')
			if trigger.Class.Name == "AlwaysExecute" {
				directiveTriggered = true
			}

			// If the trigger is type ReqRepTokenMatch then we will consider the
			// trigger to be triggered if we have a match in the map
			if trigger.Class.Name == "ReqRepTokenMatch" {

				// Look in directivesToFilesToMatchedTokens to see if we have a match
				// (If the outer keys contain this directive name then we know we have 1 or more matches)
				if _, ok := directivesToFilesToMatchedTokens[directive.Name]; ok {
					directiveTriggered = true

					tp.logger.Debug("Directive triggered", slog.String("name", directive.Name))
				}
			}
		}

		// If the trigger has triggered, then execute the action
		if directiveTriggered {
			for _, action := range directive.DirectiveActions {

				// Find the directive trigger target in the DirectiveActionClass for now
				// (For now start by looking to see if the DirectiveTriggerTarget has the type
				// 'MatchedFile')

				// Iterate the DirectiveTriggerTargets in the DirectiveActionClass
				for _, target := range action.Class.Targets {

					if target == MatchedFile {

						// So what will happen in the longer term is that we'll find a specific matched file;
						// In this iteration we know it's the special system.out file so we'll just print it to the console
						PrintWeekdays(time.Now())
					}
					if target == MatchedToken {

						// If we have a 'ReplaceToken' effect then we will traverse our map of files to matched tokens
						// and replace the matched tokens with the new tokens
						if slices.Contains(action.Class.Effects, ReplaceToken) {

							for filePath, mappedTokens := range directivesToFilesToMatchedTokens[directive.Name] {

								// Grab the file content for this filepath
								fileContent := tp.activePathsToFileContentStrings[filePath]

								forwardStartAndEndIndex := make([]string, 0)
								reverseStartAndEndIndex := make([]string, 0)

								// Iterate the MappedTokens in reverse to get the start and end indices
								for startAndEndIndex := range mappedTokens {
									forwardStartAndEndIndex = append(forwardStartAndEndIndex, startAndEndIndex)
								}
								for i := len(forwardStartAndEndIndex) - 1; i >= 0; i-- {
									reverseStartAndEndIndex = append(reverseStartAndEndIndex, forwardStartAndEndIndex[i])
								}

								// Iterate the reversed MappedTokens and replace the matched tokens with the new tokens
								for _, startAndEndIndex := range reverseStartAndEndIndex {

									mappedToken := mappedTokens[startAndEndIndex]

									//for startAndEndIndex, mappedToken := range mappedTokens {

									// Extract the start and end ints from the <start>:<end> formatted string
									start, _ := strconv.Atoi(strings.Split((startAndEndIndex), ":")[0])
									end, _ := strconv.Atoi(strings.Split(startAndEndIndex, ":")[1])

									tp.logger.Info("Replacing token",
										slog.String("file", filePath),
										slog.String("token", mappedToken),
										slog.Int("start", start),
										slog.Int("end", end))

									// Replace the string with "<PING:PONG>" using the numeric indices
									// (Grab the substring before the match and the substring after)
									fileContent = fileContent[:start] + "<PING:PONG>" + fileContent[end:]

									// We can extract some logging with good visibility in the future for this section
									// It probably makes sense to show less than the full files though
									// tp.logger.Info("Replaced content:",
									// 	slog.String("OldFileContent", tp.activePathsToFileContentStrings[filePath]),
									// 	slog.String("NewFileContent", fileContent),
									// )

									// Replace the file content with the updated string
									tp.activePathsToFileContentStrings[filePath] = fileContent

									// If activeFilePathsModified doesn't contain the file path, then add it
									if !slices.Contains(activeFilePathsModified, filePath) {
										activeFilePathsModified = append(activeFilePathsModified, filePath)
									}
								}
							}
						}
					}
				}
			}
		}
	}

	// Overwrite our files with the new content if they have been modified
	for filePath, fileContent := range tp.activePathsToFileContentStrings {
		if !slices.Contains(activeFilePathsModified, filePath) {
			continue
		}
		err := os.WriteFile(filePath, []byte(fileContent), 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetName returns the name of the template processor.
func (tp *templateProcessor) GetName() string {
	return tp.name
}

func (tp *templateProcessor) InitializeDirectives(DrectiveConfigurations []DirectiveConfiguration) error {
	for _, directive := range DrectiveConfigurations {

		tp.logger.Info("Adding directive", slog.String("name", directive.Name))

		tp.directives[directive.Name] = directive
	}
	return nil
}

// InitializeFilesystem implements TemplateProcessor.
func (tp *templateProcessor) InitializeFilesystem(FilesystemConfiguration FilesystemConfiguration) error {

	// Read the config file immediately for top level config

	// Add the IO file locations to our maps for later traversal
	tp.inputFileFolderPaths = append(tp.inputFileFolderPaths, FilesystemConfiguration.GetInputFileFolderPaths()...)
	tp.outputFileFolderPaths = append(tp.outputFileFolderPaths, FilesystemConfiguration.GetOutputFileFolderPaths()...)

	outputFileProcessor := NewFileProcessor("output")
	outputFileProcessor.UpdateFileFolderPaths(tp.outputFileFolderPaths)
	outputFileProcessor.UpdateProcessorFuncs([]func(string, ChangeType, bool) error{
		tp.templateProcessorFunc,
	})
	tp.fileProcessors = append(tp.fileProcessors, outputFileProcessor)

	return nil
}

// This satisfies the unexported interface method
func (tp *templateProcessor) internalTemplateProcessorConstructor() {
	tp.directives = make(map[string]DirectiveConfiguration)
	tp.inputFileFolderPaths = make([]string, 0)
	tp.outputFileFolderPaths = make([]string, 0)
	tp.fileProcessors = make([]FileProcessor, 0)
	tp.allPathsToFiles = make(map[string]string)
	tp.activePathsToFileContentStrings = make(map[string]string)
	tp.matchRegexes = make(map[string]*regexp.Regexp)

	// Create a JSON logger
	tp.logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

// NewTemplateProcessor creates a new TemplateProcessor with the given name.
func NewTemplateProcessor(name string) TemplateProcessor {
	tp := &templateProcessor{name: name}
	tp.internalTemplateProcessorConstructor()
	return tp
}

func (tp *templateProcessor) templateProcessorFunc(path string, changeType ChangeType, isDir bool) error {

	tp.logger.Debug("Processing file",
		slog.String("path", path),
		slog.String("changeType", convertChangeTypeToString(changeType)),
		slog.Bool("isDir", isDir))

	// Make sure the file is not a directory then read the content into a string
	if !isDir {

		fileContentBytes, err := os.ReadFile(path)
		if err != nil {
			tp.logger.Error("Error reading file",
				slog.String("path", path),
				slog.String("error", err.Error()))
			return err
		}

		tp.logger.Debug("Adding file",
			slog.String("path", path))

		// Load the file if it is not yet loaded, or update it if it is active
		tp.activePathsToFileContentStrings[path] = string(fileContentBytes)

		// In the future we will need to deal with deleted files... different map?
		// Or extend this map?
	}

	return nil
}

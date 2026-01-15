package templates

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

type FileProcessor interface {
	GetName() string
	UpdateFileFolderPaths([]string) error
	UpdateProcessorFuncs([]func(string, ChangeType, bool) error) error
	TraverseFiles() error
	InjectSysout(input string) error
}

func (fp *fileProcessor) InjectSysout(input string) error {
	fp.ephemeralInput = input

	return nil
}

type fileProcessor struct {
	name               string
	fileFolderPaths    []string
	fileProcessorFuncs []func(string, ChangeType, bool) error
	initialFiles       map[string]FileInfo
	finalFiles         map[string]FileInfo
	ephemeral          bool
	firstScan          bool
	logger             *slog.Logger
	ephemeralInput     string
}

type FileInfo struct {
	ModTime time.Time
	IsDir   bool
}

type ChangeType int

const (
	Create ChangeType = iota
	Update
	Delete
	Initialize
)

func convertChangeTypeToString(changeType ChangeType) string {
	switch changeType {
	case Create:
		return "Create"
	case Update:
		return "Update"
	case Delete:
		return "Delete"
	case Initialize:
		return "Initialize"
	default:
		return "Unknown"
	}
}

// GetName implements FileProcessor.
func (fp *fileProcessor) GetName() string {
	return fp.name
}

func NewFileProcessor(name string, ephemeral bool, logLevel slog.Level) FileProcessor {
	fp := &fileProcessor{name: name}
	fp.internalFileProcessorConstructor(logLevel, ephemeral)
	return fp
}

func (fp *fileProcessor) UpdateFileFolderPaths(fileFolderPaths []string) error {
	fp.fileFolderPaths = fileFolderPaths

	return nil
}

func (fp *fileProcessor) UpdateProcessorFuncs(processorFuncs []func(string, ChangeType, bool) error) error {

	fp.fileProcessorFuncs = make([]func(string, ChangeType, bool) error, 0)

	fp.fileProcessorFuncs = append(fp.fileProcessorFuncs, processorFuncs...)

	return nil
}

func (fp *fileProcessor) TraverseFiles() error {

	// First process standard files
	fp.finalFiles = make(map[string]FileInfo)
	for _, path := range fp.fileFolderPaths {
		err := filepath.Walk(path, fp.walkDirCallback)
		if err != nil {
			panic(err)
		}
	}

	// Compare initial and final files - if any keys exist in initial but not in final, remove the fiile from initial
	for initialFile := range fp.initialFiles {
		_, ok := fp.finalFiles[initialFile]
		if !ok {
			delete(fp.initialFiles, initialFile)
			for _, processorFunc := range fp.fileProcessorFuncs {
				err := processorFunc(initialFile, Delete, fp.initialFiles[initialFile].IsDir)
				if err != nil {
					return err
				}
			}
		}
	}

	// Copy all final files to initial files
	for finalFile := range fp.finalFiles {
		fp.initialFiles[finalFile] = fp.finalFiles[finalFile]
	}

	fp.firstScan = false

	fp.logger.Debug("Completed first scan.")

	// For ephemeral files (sysout so far) we always consider them new if there is content
	for _, processorFunc := range fp.fileProcessorFuncs {
		fp.logger.Debug("Processing ephemeral files for function:",
			slog.String("function", fmt.Sprintf("%p", processorFunc)))
		if fp.ephemeral && fp.ephemeralInput != "" {
			fp.logger.Debug("Injecting ephemeral file",
				slog.String("content", fp.ephemeralInput))
			err := processorFunc("\x00sysout", Create, false)
			if err != nil {
				return err
			}
		}
	}
	// This doesn't quite make sense yet - we don't need the content in the file processor
	fp.ephemeralInput = ""

	return nil
}

// This satisfies the unexported interface method
func (fp *fileProcessor) internalFileProcessorConstructor(logLevel slog.Level, ephemeral bool) {
	fp.initialFiles = make(map[string]FileInfo)
	fp.finalFiles = make(map[string]FileInfo)
	fp.fileProcessorFuncs = make([]func(string, ChangeType, bool) error, 0)
	fp.firstScan = true
	fp.ephemeral = ephemeral
	fp.ephemeralInput = ""

	// Create a JSON logger with the specified log level
	fp.logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
}

func (fp *fileProcessor) walkDirCallback(path string, info os.FileInfo, err error) error {
	if err != nil {
		fp.logger.Error("Error accessing file",
			slog.String("path", path),
			slog.String("error", err.Error()))
		return nil
	}

	fileInfo, ok := fp.initialFiles[path]

	if !ok { // New file detected
		fp.finalFiles[path] = FileInfo{
			ModTime: info.ModTime(),
			IsDir:   info.IsDir(),
		}
		if fp.firstScan {
			for _, processorFunc := range fp.fileProcessorFuncs {
				err := processorFunc(path, Initialize, info.IsDir())
				if err != nil {
					fp.logger.Error("Error processing file",
						slog.String("path", path),
						slog.String("error", err.Error()))
					return err
				}
			}
		} else {
			for _, processorFunc := range fp.fileProcessorFuncs {
				err := processorFunc(path, Create, info.IsDir())
				if err != nil {
					fp.logger.Error("Error processing file",
						slog.String("path", path),
						slog.String("error", err.Error()))
					return err
				}
			}
		}
	} else if info.ModTime().After(fileInfo.ModTime) {
		fp.finalFiles[path] = FileInfo{
			ModTime: info.ModTime(),
			IsDir:   info.IsDir(),
		}

		for _, processorFunc := range fp.fileProcessorFuncs {
			fmt.Println("Updated file:", path)
			err := processorFunc(path, Update, info.IsDir())
			if err != nil {
				fp.logger.Error("Error processing file",
					slog.String("path", path),
					slog.String("error", err.Error()))
				return err
			}
		}
	} else {
		fp.finalFiles[path] = FileInfo{
			ModTime: info.ModTime(),
			IsDir:   info.IsDir(),
		}
	}

	return nil
}

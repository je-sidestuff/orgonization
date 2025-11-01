package templates

import (
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
}

type fileProcessor struct {
	name               string
	fileFolderPaths    []string
	fileProcessorFuncs []func(string, ChangeType, bool) error
	initialFiles       map[string]FileInfo
	finalFiles         map[string]FileInfo
	firstDetection     bool
	logger             *slog.Logger
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
	panic("unimplemented")
}

func NewFileProcessor(name string) FileProcessor {
	fp := &fileProcessor{name: name}
	fp.internalFileProcessorConstructor()
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

	fp.firstDetection = false

	return nil
}

// This satisfies the unexported interface method
func (fp *fileProcessor) internalFileProcessorConstructor() {
	fp.initialFiles = make(map[string]FileInfo)
	fp.finalFiles = make(map[string]FileInfo)
	fp.fileProcessorFuncs = make([]func(string, ChangeType, bool) error, 0)
	fp.firstDetection = true

	// Create a JSON logger
	fp.logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
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
		if fp.firstDetection {
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

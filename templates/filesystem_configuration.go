package templates

import (
	"os/user"
	"path/filepath"
	"strings"
)

type FilesystemConfiguration interface {
	GetName() string
	GetInputFileFolderPaths() []string
	GetOutputFileFolderPaths() []string
	GetConfigPath() string
	GetOrgoDirectoryPath() string
}

type filesystemConfiguration struct {
	name                  string
	inputFileFolderPaths  []string
	outputFileFolderPaths []string
	configPath            string
	orgoDirectoryPath     string
}

func newFilesystemConfiguration(name string, inputPaths, outputPaths []string, configPath, orgoPath string) (FilesystemConfiguration, error) {
	expandedInputPaths := make([]string, len(inputPaths))
	for i, path := range inputPaths {
		expanded, err := expandPath(path)
		if err != nil {
			return nil, err
		}
		expandedInputPaths[i] = expanded
	}

	expandedOutputPaths := make([]string, len(outputPaths))
	for i, path := range outputPaths {
		expanded, err := expandPath(path)
		if err != nil {
			return nil, err
		}
		expandedOutputPaths[i] = expanded
	}

	expandedConfigPath, err := expandPath(configPath)
	if err != nil {
		return nil, err
	}

	expandedOrgoPath, err := expandPath(orgoPath)
	if err != nil {
		return nil, err
	}

	return &filesystemConfiguration{
		name:                  name,
		inputFileFolderPaths:  expandedInputPaths,
		outputFileFolderPaths: expandedOutputPaths,
		configPath:            expandedConfigPath,
		orgoDirectoryPath:     expandedOrgoPath,
	}, nil
}

func (f *filesystemConfiguration) GetName() string {
	return f.name
}

func (f *filesystemConfiguration) GetInputFileFolderPaths() []string {
	return f.inputFileFolderPaths
}

func (f *filesystemConfiguration) GetOutputFileFolderPaths() []string {
	return f.outputFileFolderPaths
}

func (f *filesystemConfiguration) GetConfigPath() string {
	return f.configPath
}

func (f *filesystemConfiguration) GetOrgoDirectoryPath() string {
	return f.orgoDirectoryPath
}

func GetDefaultFilesystemConfiguration() (FilesystemConfiguration, error) {
	return newFilesystemConfiguration(
		"default",
		[]string{"~/.orgo/infiles", "~/notes"},
		[]string{"~/.orgo/outfiles", "~/notes"},
		"~/.orgo/config.yaml",
		"~/.orgo/",
	)
}

func expandPath(path string) (string, error) {
	if strings.HasPrefix(path, "~/") {
		usr, err := user.Current()
		if err != nil {
			return "", err
		}
		return filepath.Join(usr.HomeDir, path[2:]), nil
	}
	return path, nil
}

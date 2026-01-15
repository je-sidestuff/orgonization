package integration

import (
	"errors"
	"os"
	"strings"

	cp "github.com/otiai10/copy"
)

// FilesystemTestInput contains a map of the created temporary directory handles to their paths.
type FilesystemTestInput struct {
	TempFolderHandleToPath map[string]string
}

// NewFilesystemTestInput creates a new instance of FilesystemTestInput with an initialized map
// that associates temporary folder handles with their corresponding paths.
func NewFilesystemTestInput() FilesystemTestInput {
	return FilesystemTestInput{
		TempFolderHandleToPath: make(map[string]string),
	}
}

// CreateAndMapTempFolder creates a temporary directory and add it to the map of created directories.
func (fti *FilesystemTestInput) CreateAndMapTempFolder(folderHandle string, tempFolderPrefix string) (string, error) {

	// Verify that tempFolderPrefix is not empty and does not contain more than one asterix
	if tempFolderPrefix == "" || strings.Count(tempFolderPrefix, "*") > 1 {
		return "", errors.New("tempFolderPrefix is empty or contains more than one asterix")
	}

	folderPath, err := os.MkdirTemp("", tempFolderPrefix)

	if err != nil {
		return "", err
	}

	fti.TempFolderHandleToPath[folderHandle] = folderPath

	return folderPath, nil
}

// GetTempFolderPath retrieves the path of the temporary folder associated with the given folder handle.
func (fti *FilesystemTestInput) GetTempFolderPath(folderHandle string) string {
	return fti.TempFolderHandleToPath[folderHandle]
}

// CloneDirectoryTreeToNewTempFolder clones the directory tree rooted at sourcePath to a new temporary directory and maps it to a new temp folder.
// The new temporary directory is created with a name that starts with "orgo_test_clone_".
func (fti *FilesystemTestInput) CloneDirectoryTreeToNewTempFolder(sourcePath string, folderHandle string) (string, error) {

	folderPath, err := fti.CreateAndMapTempFolder(folderHandle, "orgo_test_clone_")

	if err != nil {
		return "", err
	}

	err = cp.Copy(sourcePath, folderPath)

	if err != nil {
		return "", err
	}

	return folderPath, nil
}

// DeleteAndUnmapTempFolder deletes a temporary directory and removes it from the map of created directories
func (fti *FilesystemTestInput) DeleteAndUnmapTempFolder(folderHandle string) error {
	err := os.RemoveAll(fti.TempFolderHandleToPath[folderHandle])

	if err != nil {
		return err
	}

	delete(fti.TempFolderHandleToPath, folderHandle)

	return nil
}

// CleanupFilesystemTestInput deletes all temporary directories that were created by FilesystemTestInput.
func (fti *FilesystemTestInput) CleanupFilesystemTestInput() {
	for folderHandle := range fti.TempFolderHandleToPath {
		_ = fti.DeleteAndUnmapTempFolder(folderHandle)
	}
}

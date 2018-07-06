package FileSystem

import "time"
import Models "../../shared/models"
import "errors"

// Note: This package doesn't deal with any of the server config properties.
// This responsibility is left to the callers of these methods.

// Server configuration
var GlobalServerConfig *Models.ServerConfig

// Hold information for each file in the vfs
var fileEntries []Models.FileEntry

// Hold actual contents for each file
var fileContentsMap map[int] []byte

// Hold next fileID to assign
var nextFileIDToAssign int

func init() {
  fileEntries = []Models.FileEntry{}
  fileContentsMap = make(map[int] []byte)
  nextFileIDToAssign = 0
}

func GetFileManifest() ([]Models.FileEntry) {
  return fileEntries
}

func GetFileContentsMap() (map[int] []byte) {
  return fileContentsMap
}

// Add a file to the file system
func AddFile(fileName string, owner string, contents []byte) (Models.FileEntry, error) {
  fileID := nextFileIDToAssign; nextFileIDToAssign++
  curTime := time.Now()

  // Add entry to file contents map
  fileContentsMap[fileID] = make([]byte,len(contents))
  copy(fileContentsMap[fileID], contents)

  // Update manifest
  newEntry := Models.FileEntry{fileName,fileID,len(contents),owner,curTime,0}
  fileEntries = append(fileEntries,newEntry)
  return newEntry, nil
}

// Remove the specified manifest entry from the file system.
// manifestIdx is the position of the entry within fileEntires
// fileID is the unique ID assigned to the file (used as the key in fileContentsMap
// and is a filed in the FileEntry structure)
func RemoveFile(manifestIdx int, fileID int) error {
  // Sanity checks: Make sure manifestIdx in bounds and that it and 
  // fileID correspond to the same file
  if manifestIdx < 0 || manifestIdx >= len(fileEntries) {
    return errors.New("Manifest index out of bounds")
  }

  if fileEntries[manifestIdx].FileID != fileID {
    return errors.New("Mismatched file IDs")
  }

  // Remove the file from the manifest
  // Replace the entry to delete with the last entry (order doesn't matter)
  fileEntries[manifestIdx] = fileEntries[len(fileEntries)-1]
  fileEntries = fileEntries[:len(fileEntries)-1]

  // Remove file contents mapping (the actual "file")
  delete(fileContentsMap, fileID)
  return nil
}

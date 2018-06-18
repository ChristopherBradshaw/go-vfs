package FileSystem

import "time"
import Models "../../shared/models"

// Note: This package doesn't deal with any of the server config properties.
// This responsibility is left to the callers of these methods.

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

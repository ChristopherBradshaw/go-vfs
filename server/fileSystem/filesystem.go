package FileSystem

import "time"

// Note: This package doesn't deal with any of the server config properties.
// This responsibility is left to the callers of these methods.

type FileEntry struct {
  FileName string
  FileID int
  FileSize int
  Owner string
  LastModified time.Time
  NumDownloads int
}

// Hold information for each file in the vfs
var fileEntries []FileEntry

// Hold actual contents for each file
var fileContentsMap map[int] []byte

// Hold next fileID to assign
var nextFileIDToAssign int

func init() {
  fileEntries = []FileEntry{}
  fileContentsMap = make(map[int] []byte)
  nextFileIDToAssign = 0
}

func GetFileManifest() ([]FileEntry) {
  return fileEntries
}

func GetFileContentsMap() (map[int] []byte) {
  return fileContentsMap
}

func AddFile(fileName string, owner string, contents []byte) (FileEntry, error) {
  fileID := nextFileIDToAssign; nextFileIDToAssign++
  curTime := time.Now()

  // Add entry to file contents map
  fileContentsMap[fileID] = make([]byte,len(contents))
  copy(fileContentsMap[fileID], contents)

  newEntry := FileEntry{fileName,fileID,len(contents),owner,curTime,0}
  // Update manifest
  fileEntries = append(fileEntries,newEntry)
  return newEntry, nil
}

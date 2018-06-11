package FileSystem

import "time"

type FileEntry struct {
  FileName string
  FileID int
  Owner string
  LastModified time.Time
  NumDownloads int
}


// Hold information for each file in the vfs
var fileEntries []FileEntry

// Hold actual contents for each file
var fileContentsMap map[int] []byte

func init() {
  fileEntries = []FileEntry{}
  fileContentsMap = make(map[int] []byte)
}

func GetFileManifest() ([]FileEntry) {
  return fileEntries
}

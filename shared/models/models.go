package Models
import "time"
import "fmt"

type GetConfigResponse struct {
  ServerConfig ServerConfig
}

type GetManifestResponse struct {
  FileEntries []FileEntry
}

type UploadFileRequest struct {
  Owner string
}

type UploadFileResponse struct {
  Filename string
}

type DownloadFileRequest struct {
  FileID int
}

type ClientConfig struct {
  Hostname string
  Port string
  Username string
}

type ServerConfig struct {
  Port int
  MaxNumFiles int
  MaxFilesizeMB int
  MaxUserUploadsPerMinute int
}

func (config ServerConfig) String() string {
  return fmt.Sprintf("Port: %d\nMax amount of files: %d\nMax filesize: %dMB\n"+
  "Max uploads per user/minute: %d",config.Port, config.MaxNumFiles,
  config.MaxFilesizeMB, config.MaxUserUploadsPerMinute)
}

type FileEntry struct {
  FileName string
  FileID int
  FileSize int
  Owner string
  LastModified time.Time
  NumDownloads int
}

func (entry FileEntry) String() string {
  return fmt.Sprintf("ID: %d | Name: %s | Owner: %s | Size: %d bytes | Last Modified: %s",
    entry.FileID, entry.FileName, entry.Owner, entry.FileSize,
    entry.LastModified.Format(time.Stamp))
}

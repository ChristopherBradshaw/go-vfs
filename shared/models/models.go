package Models
import "time"
import "fmt"

type GetConfigResponse struct {
  ServerConfig ServerConfig `json:"server_config"`
}

type GetManifestResponse struct {
  FileEntries []FileEntry `json:"files"`
}

type ClientConfig struct {
  Hostname string
  Port string
  Username string
}

type ServerConfig struct {
  Port int `json:"port"`
  MaxNumFiles int `json:"max_num_files"`
  MaxFilesizeMB int `json:"max_filesize_mb"`
  MaxUserUploadsPerMinute int `json:"max_user_uploads_per_minute"`
}

func (config ServerConfig) String() string {
  return fmt.Sprintf("Port: %d\nMax amount of files: %d\nMax filesize: %dMB\n"+
  "Max uploads per user/minute: %d",config.Port, config.MaxNumFiles,
  config.MaxFilesizeMB, config.MaxUserUploadsPerMinute)
}

type FileEntry struct {
  FileName string `json:"file_name"`
  FileID int `json:"file_id"`
  FileSize int `json:"file_size"`
  Owner string `json:"owner"`
  LastModified time.Time `json:"last_modified"`
  NumDownloads int `json:"num_downloads"`
}

func (entry FileEntry) String() string {
  return fmt.Sprintf("ID: %d | Name: %s | Owner: %s | Size: %d bytes | Last Modified: %s",
    entry.FileID, entry.FileName, entry.Owner, entry.FileSize,
    entry.LastModified.Format(time.Stamp))
}

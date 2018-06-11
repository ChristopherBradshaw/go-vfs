package Handlers
import "../fileSystem"
import "net/http"
import "encoding/json"

type GetManifestResponse struct {
  Manifest []FileSystem.FileEntry
}

func GetManifestHandler(w http.ResponseWriter, r *http.Request) {
  encoder := json.NewEncoder(w)
  encoder.Encode(FileSystem.GetFileManifest())
}


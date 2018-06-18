package Handlers
import "../fileSystem"
import "net/http"
import "encoding/json"

// GET request for server file listing
func GetManifestHandler(w http.ResponseWriter, r *http.Request) {
  if !EnforceHttpMethods(w,r,[]string{"GET"}) {
    return
  }

  encoder := json.NewEncoder(w)
  encoder.Encode(FileSystem.GetFileManifest())
}

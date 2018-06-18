package Handlers
import Models "../../shared/models"
import "../fileSystem"
import "net/http"
import "encoding/json"

// GET request for server file listing
func GetManifestHandler(w http.ResponseWriter, r *http.Request) {
  if !EnforceHttpMethods(w,r,[]string{"GET"}) {
    return
  }

  encoder := json.NewEncoder(w)
  response := Models.GetManifestResponse{FileSystem.GetFileManifest()}
  encoder.Encode(response)
}

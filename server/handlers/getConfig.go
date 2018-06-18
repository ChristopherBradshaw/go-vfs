package Handlers
import Config "../config"
import Models "../../shared/models"
import "net/http"
import "encoding/json"

// GET request for server configuration
func GetConfigHandler(w http.ResponseWriter, r *http.Request) {
  if !EnforceHttpMethods(w,r,[]string{"GET"}) {
    return
  }

  encoder := json.NewEncoder(w)
  response := Models.GetConfigResponse{*Config.GlobalServerConfig}
  encoder.Encode(response)
}


package Handlers
import Config "../config"
import "net/http"
import "encoding/json"

// GET request for server configuration
func GetConfigHandler(w http.ResponseWriter, r *http.Request) {
  if !EnforceHttpMethods(w,r,[]string{"GET"}) {
    return
  }

  encoder := json.NewEncoder(w)
  encoder.Encode(Config.GlobalServerConfig)
}


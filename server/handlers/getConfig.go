package Handlers
import Config "../config"
import "net/http"
import "encoding/json"

func GetConfigHandler(w http.ResponseWriter, r *http.Request) {
  encoder := json.NewEncoder(w)
  encoder.Encode(Config.GlobalServerConfig)
}


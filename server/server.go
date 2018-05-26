package main

import "fmt"
import "log"
import "os"
import "io/ioutil"
import "net/http"
import Handlers "./handlers"
import Config "./config"

func main() {
  args := os.Args

  if len(args) < 2 {
    log.Fatal("No server config provided")
  }

  // Read in server configuration
  file, err := os.Open(args[1])
  defer file.Close()

  if err != nil {
    log.Fatal(err)
    return
  }

  // Parse server config
  fileContents, err := ioutil.ReadAll(file)
  Config.GlobalServerConfig, _ = Config.ParseServerConfig(string(fileContents))

  // Register endpoints
  http.HandleFunc("/getConfig", Handlers.GetConfigHandler)

  // Start server
  log.Printf("Starting server on port %v\n", Config.GlobalServerConfig.Port)
  log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", Config.GlobalServerConfig.Port), nil))
}

package main

import "fmt"
import "log"
import "os"
import "io/ioutil"
import "net/http"
import Handlers "./handlers"
import Config "../shared/config"
import FileSystem "./fileSystem"


func main() {
  args := os.Args

  if len(args) < 2 {
    log.Fatal("No server config provided")
    return
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
  FileSystem.GlobalServerConfig, _ = Config.ParseServerConfig(string(fileContents))

  // Register endpoints
  http.HandleFunc("/getConfig", Handlers.GetConfigHandler)
  http.HandleFunc("/getManifest", Handlers.GetManifestHandler)
  http.HandleFunc("/uploadFile", Handlers.UploadFileHandler)
  http.HandleFunc("/getFile/", Handlers.GetFileHandler)

  // Start server
  log.Printf("Starting server on port %v\n", FileSystem.GlobalServerConfig.Port)
  log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", FileSystem.GlobalServerConfig.Port), nil))
}

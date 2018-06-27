package main

import Config "../shared/config"
import Models "../shared/models"
import "errors"
import "encoding/json"
import "bufio"
import "fmt"
import "os"
import "strconv"
import "io/ioutil"
import "net/http"
import "strings"
import "io"
import "mime/multipart"
import "bytes"

var url string
var username string

func init() {
  url = ""
  username = ""
}

func main() {
  fmt.Println("\n--------- Virtual Filesystem ---------")
  reader := bufio.NewReader(os.Stdin)
  // If they provided a config file, pull fields from there
  args := os.Args
  if len(args) == 2 {
    parseServerInfo(args[1])
  } else  {
    // Otherwise ask them for the fields
    getServerInfo(reader)
  }
  testConnection()
  fmt.Printf("Hello, %s. Connected to: %s\n",username,url)
  doPromptLoop(reader)
}

func saveConfig(config Models.ClientConfig) {
  url = "http://" + config.Hostname + ":" + config.Port
  username = config.Username
}

// Read server/username information from a file instead of prompting the user
func parseServerInfo(fileName string) {
  file, err := os.Open(fileName)
  if err != nil {
    fmt.Println("Failed to open specified config file")
    os.Exit(1)
  }

  defer file.Close()

  // Parse the file input into the fields we need
  fileContents, err := ioutil.ReadAll(file)
  conf, err := Config.ParseClientConfig(string(fileContents))
  if err != nil {
    fmt.Println("Failed to parse specified config file")
    os.Exit(1)
  }

  saveConfig(*conf)
}

// Ask the user for the specified fields
// (call if they don't provide a config file)
func getServerInfo(reader *bufio.Reader) {
  fmt.Print("Hostname: ")
  host, _ := reader.ReadString('\n')
  host = strings.TrimSpace(host)

  fmt.Print("Port: ")
  port, _ := reader.ReadString('\n')
  port = strings.TrimSpace(port)

  fmt.Print("Your name: ")
  name, _ := reader.ReadString('\n')
  name = strings.TrimSpace(name)

  if len(name) == 0 {
    fmt.Println("Username can't be empty")
    os.Exit(1)
  }

  conf := Models.ClientConfig{host,port,name}
  saveConfig(conf)
}

// Make sure we can reach the server they provided
func testConnection() {
  resp, err := http.Get(url)
  if err != nil {
    fmt.Printf("%s", err)
    os.Exit(1)
  } else {
    defer resp.Body.Close()
    _, err := ioutil.ReadAll(resp.Body)
    if err != nil {
      fmt.Printf("%s",err)
      os.Exit(1)
    }
  }
}

// Main prompt loop
func doPromptLoop(reader *bufio.Reader) {
  for ;; {
    fmt.Println("\n----- What would you like to do? -----")
    fmt.Println("1. View server manifest")
    fmt.Println("2. Upload a file")
    fmt.Println("3. Download a file")
    fmt.Println("4. Remove a file")
    fmt.Println("5. View server settings")
    fmt.Println("6. Exit")
    fmt.Print("\nSelection (1-6): ")

    selectionStr, _ := reader.ReadString('\n')
    selectionInt , err := strconv.Atoi(strings.TrimSpace(selectionStr))

    if err != nil {
      fmt.Errorf("%s\n",err)
      continue
    }

    switch selectionInt {
      case 1:
        printManifest()
        break
      case 2:
        uploadFile(reader)
        break
      case 3:
        downloadFile(reader)
        break
      case 4:
        removeFile()
        break
      case 5:
        printConfig()
        break
      case 6:
        os.Exit(0)
        break
      default:
        fmt.Println("Invalid selection")
        continue
    }
  }
}

// Retrieve the manifest data from the server
func getManifest() (*Models.GetManifestResponse, error) {
  resp, err := http.Get(url + "/getManifest")
  if err != nil {
    fmt.Printf("%s", err)
    return nil, errors.New("Could not retrieve manifest")
  }
  defer resp.Body.Close()

  // Decode response data
  var decodedResponse Models.GetManifestResponse
  decoder := json.NewDecoder(resp.Body)
  err = decoder.Decode(&decodedResponse)
  if err != nil {
    fmt.Printf("%s",err)
    return nil, errors.New("Could not decode manifest response")
  }

  return &decodedResponse, nil
}

// Request the manifest and print it
func printManifest() {
  fmt.Println("\n----- File listing -----")
  manifest, err := getManifest()

  // Could not get manifest
  if err != nil {
    fmt.Println(err)
    return
  }

  // Got the manifest, display them nicely
  for _, elm := range manifest.FileEntries {
    fmt.Println(elm)
  }
}

func mustOpen(f string) *os.File {
    r, err := os.Open(f)
    if err != nil {
        panic(err)
    }
    return r
}

func uploadFile(reader *bufio.Reader) {
  fmt.Println("\n----- Upload File -----")
  fmt.Print("Enter file path: ")
  filePath, _ := reader.ReadString('\n')
  filePath = strings.TrimSpace(filePath)

  httpClient := http.Client{}

	// Prepare a form
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

  // Add file to upload to the form
	file := mustOpen(filePath)
	defer file.Close()
	formWriter, _ := w.CreateFormFile("file",file.Name())
  io.Copy(formWriter, file)

  // Add owner/username to the form
  formWriter, _ = w.CreateFormField("owner")
  io.Copy(formWriter, strings.NewReader(username))

	w.Close()

	// Now that we have a form, submit it to the handler.
	req, err := http.NewRequest("POST", url+"/uploadFile", &b)
	if err != nil {
			return
	}

	req.Header.Set("Content-Type", w.FormDataContentType())

	// Submit the request
	res, err := httpClient.Do(req)
	if err != nil {
			return
	}

	// Check the response
	if res.StatusCode != http.StatusOK {
			err = fmt.Errorf("Bad status: %s", res.Status)
	}

  // Decode response
  defer res.Body.Close()
  var decodedResponse Models.UploadFileResponse
  decoder := json.NewDecoder(res.Body)
  err = decoder.Decode(&decodedResponse)
  if err != nil {
    fmt.Printf("%s", err)
    return
  }

  // If the server saved the file with a different file name than the
  // one we provided, let the user know
  if decodedResponse.Filename != file.Name() {
    fmt.Printf("File was renamed: %s->%s\n",file.Name(),decodedResponse.Filename)
  }

  fmt.Println("Successfully uploaded " + decodedResponse.Filename)
}

func downloadFile(reader *bufio.Reader) {
  fmt.Println("\n----- Download File -----")
  fmt.Print("Enter file ID: ")
  fileID, _ := reader.ReadString('\n')
  fileID = strings.TrimSpace(fileID)

  resp, err := http.Get(url + "/getFile/" + fileID)
  if err != nil {
    fmt.Printf("%s", err)
    return
  }
  defer resp.Body.Close()

  contents, err := ioutil.ReadAll(resp.Body)
  fmt.Println(string(contents))
  // TODO actually save the file
}

func removeFile() {

}

// Retrieve the server config data from the server
func getConfig() (*Models.GetConfigResponse, error) {
  resp, err := http.Get(url + "/getConfig")
  if err != nil {
    fmt.Printf("%s", err)
    return nil, errors.New("Could not retrieve config")
  }
  defer resp.Body.Close()

  // Decode response data
  var decodedResponse Models.GetConfigResponse
  decoder := json.NewDecoder(resp.Body)
  err = decoder.Decode(&decodedResponse)
  if err != nil {
    fmt.Printf("%s",err)
    return nil, errors.New("Could not decode manifest response")
  }

  return &decodedResponse, nil
}

// Request and display the server config data
func printConfig() {
  fmt.Println("\n----- Server information -----")
  fmt.Printf("Connected to: %s as %s\n\n",url,username)
  config, err := getConfig()

  // Could not get manifest
  if err != nil {
    fmt.Println(err)
    return
  }

  // Got the config, display it nicely
  fmt.Println(config.ServerConfig)
}

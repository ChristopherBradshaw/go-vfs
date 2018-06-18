package main

import "bufio"
import "fmt"
import "os"
import "strconv"
import "io/ioutil"
import "net/http"
import "strings"

var url string
var owner string

func init() {
  url = ""
}

func main() {
  fmt.Println("--------- Virtual Filesystem ---------")
  reader := bufio.NewReader(os.Stdin)
  getServerInfo(reader)
  testConnection()
  doPromptLoop(reader)
}

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

  url = "http://" + host + ":" + port
  owner = name
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

func doPromptLoop(reader *bufio.Reader) {
  for ;; {
    fmt.Println("----- What would you like to do? -----")
    fmt.Println("1. View server manifest")
    fmt.Println("2. Upload a file")
    fmt.Println("3. Download a file")
    fmt.Println("4. Remove a file")
    fmt.Println("5. View server configuration")
    fmt.Println("6. Exit")
    fmt.Print("Selection (1-6): ")

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
        uploadFile()
        break
      case 3:
        downloadFile()
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

func printManifest() {

}

func uploadFile() {

}

func downloadFile() {

}

func removeFile() {

}

func printConfig() {

}

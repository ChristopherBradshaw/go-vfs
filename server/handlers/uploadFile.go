package Handlers
import FileSystem "../fileSystem"
import Models "../../shared/models"
import "net/http"
import "bytes"
import "strings"
import "fmt"
import "io"
import "encoding/json"

// POST request for uploading files to the server
func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
  if !EnforceHttpMethods(w,r,[]string{"POST"}) {
    return
  }

  // Read request body
  owner := r.FormValue("owner")

  // They didn't include the owner
  if owner == "" {
    http.Error(w, "Missing file owner field", http.StatusBadRequest)
    return
  }

  // Find the file they want to upload
	var buffer bytes.Buffer
  file, header, err := r.FormFile("file")
  if err != nil {
    http.Error(w, "Upload file not provided", http.StatusBadRequest)
    return
  }

  defer file.Close()

  // Tokenize filename
  name := strings.Split(header.Filename, ".")

  // We want a name and an extension (Ex: a.txt ["a","txt"])
  if len(name) != 2 {
    http.Error(w, "File name must have a name and a type", http.StatusBadRequest)
    return
  }

  // See if a file with this name already exists
  // If it does, generate a similar unique file name
  filenameAlreadyExists := existsFilename(header.Filename)
  filename := header.Filename
  if filenameAlreadyExists {
    // Keep trying to append a number to the end of the filename until the 
    // filename becomes unique
    newFilename := header.Filename
    for renameAttemptNumber := 1; existsFilename(newFilename); renameAttemptNumber++ {
      newFilename = fmt.Sprintf("%s_%d.%s",name[0],renameAttemptNumber,name[1])
    }

    filename = newFilename
  }

  // Read in the file to our filesystem
  // TODO: We're copying the entire file twice (once here and once in AddFile)
  //       find a better way to do this
  io.Copy(&buffer, file)
  FileSystem.AddFile(filename, strings.ToLower(owner), buffer.Bytes())

  encoder := json.NewEncoder(w)
  encoder.Encode(Models.UploadFileResponse{filename})
}

// Check the given filename against all existing names in the manifest
func existsFilename(filename string) bool {
  var exists bool
  exists = false
  for _, elm := range FileSystem.GetFileManifest() {
    if elm.FileName == filename {
      exists = true
      break
    }
  }

  return exists
}

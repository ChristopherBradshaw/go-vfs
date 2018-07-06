package Handlers
import FileSystem "../fileSystem/"
import Models "../../shared/models/"
import "net/http"
import "path"
import "strconv"
import "errors"
import "encoding/json"
import "strings"

// DELETE request for removing files from this server (/removeFile/<fileID>)
func RemoveFileHandler(w http.ResponseWriter, r *http.Request) {
  if !EnforceHttpMethods(w,r,[]string{"DELETE"}) {
    return
  }

  // Parse URL to find which file they want to remove
  fileID := path.Base(r.URL.Path)
  if fileID == "removeFile" {
    // They made a request to /removeFile and not /removeFile/<fileID>
    http.Error(w, "File ID not provided", http.StatusBadRequest)
    return
  }

  // Convert file path string to integer
  fileIDInt, err := strconv.Atoi(fileID)
  if err != nil {
    // File ID they gave wasn't a number
    http.Error(w, "File ID not a number", http.StatusBadRequest)
    return
  }

  // Check if this file ID exists
  _, containsKey := FileSystem.GetFileContentsMap()[fileIDInt]
  if !containsKey {
    // File doesn't exist
    http.Error(w, "File ID not found", http.StatusNotFound)
    return
  }

  // Parse request body
  var request Models.RemoveFileRequest
  decoder := json.NewDecoder(r.Body)
  err = decoder.Decode(&request)
  if err != nil {
    http.Error(w, "Bad request", http.StatusBadRequest)
    return
  }

  // Get the manifest index of file entry associated with fileIDInt
  manifestIdx, err := getFileIDIndex(fileIDInt)

  // An entry with this file ID doesn't exist in the manifest
  if err != nil {
    http.Error(w, "File ID not found", http.StatusNotFound)
    return
  }

  // Make sure the person making the request is allowed to delete this file
  // TODO: Use some form of authentication instead of this dumb username checking
  if FileSystem.GetFileManifest()[manifestIdx].Owner != strings.ToLower(request.Username) {
    http.Error(w, "Not authorized", http.StatusUnauthorized)
    return
  }

  // Remove it from the file system
  fileEntry, err := FileSystem.RemoveFile(manifestIdx, fileIDInt)
  if err != nil {
    http.Error(w, "Failed to remove file", http.StatusInternalServerError)
    return
  }


  encoder := json.NewEncoder(w)
  encoder.Encode(Models.RemoveFileResponse{FileInfo: fileEntry})
}

func getFileIDIndex(fileID int) (int,error) {
  for idx, elm := range FileSystem.GetFileManifest() {
    if elm.FileID == fileID {
      return idx,nil
    }
  }

  return -1, errors.New("File ID not found")
}

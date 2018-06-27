package Handlers
import FileSystem "../fileSystem/"
import "fmt"
import "net/http"
import "io"
import "bytes"
import "path"
import "strconv"

// GET request for retrieving files from this server (/getFile/<fileID>)
func GetFileHandler(w http.ResponseWriter, r *http.Request) {
  if !EnforceHttpMethods(w,r,[]string{"GET"}) {
    return
  }

  fileID := path.Base(r.URL.Path)
  if fileID == "getFile" {
    // They made a request to /getFile and not /getFile/<fileID>
    http.Error(w, "File ID not provided", http.StatusBadRequest)
    return
  }

  fileIDInt, err := strconv.Atoi(fileID)
  if err != nil {
    // File ID they gave wasn't a number
    http.Error(w, "File ID not a number", http.StatusBadRequest)
    return
  }


  _, containsKey := FileSystem.GetFileContentsMap()[fileIDInt]
  if !containsKey {
    // File doesn't exist
    http.Error(w, "File ID not found", http.StatusNotFound)
    return
  }

  // Find the filename (to set it in the response)
  // Would probably want to use a bimap or something for faster lookup, but just
  // loop through all the entires for now
  var fName *string
  fName = nil
  for _, elm := range FileSystem.GetFileManifest() {
    if elm.FileID == fileIDInt {
      fName = &elm.FileName
    }
  }

  // This should not happen. If it does, that means we're not keeping
  // fileEntires and fileContentsMap in sync (file IDs are shared across them)
  if fName == nil {
    http.Error(w,"File entry structures out of sync", http.StatusInternalServerError)
    return
  }

  // Send the file back
  w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s",*fName))
  w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
  io.Copy(w,bytes.NewReader(FileSystem.GetFileContentsMap()[fileIDInt]))
}

package Config

import Models "../models"
import "strings"
import "errors"
import "fmt"
import "strconv"

// Parse the entire server config file. Return a struct
// containing the data
func ParseServerConfig(fileContents string) (*Models.ServerConfig, error) {
  // Parse each property, fail if any are missing
  port, err := readConfigProperty(fileContents,"port")
  if err != nil {
    return nil,errors.New("Port not specified")
  }

  maxNumFiles, err := readConfigProperty(fileContents,"max_num_files")
  if err != nil {
    return nil, errors.New("max_num_files not specified")
  }

  maxFilesizeMb, err := readConfigProperty(fileContents,"max_filesize_mb")
  if err != nil {
    return nil, errors.New("max_filesiz_mb not specified")
  }

  maxUserUploadsPerMinute, err :=
    readConfigProperty(fileContents,"max_user_uploads_per_minute")
  if err != nil {
    return nil, errors.New("max_user_uploads_per_minute not specified")
  }

  // Convert to ints
  portN,_ := strconv.Atoi(port)
  maxNumFilesN,_ := strconv.Atoi(maxNumFiles)
  maxFilesizeMbN,_ := strconv.Atoi(maxFilesizeMb)
  maxUserUploadsPerMinuteN,_ := strconv.Atoi(maxUserUploadsPerMinute)

  conf := Models.ServerConfig{portN,maxNumFilesN,maxFilesizeMbN,maxUserUploadsPerMinuteN}
  return &conf,nil
}

// Parse the entire client config file. Return a struct
// containing the data
func ParseClientConfig(fileContents string) (*Models.ClientConfig, error) {
  // Parse each property, fail if any are missing
  hostname, err := readConfigProperty(fileContents,"hostname")
  if err != nil {
    return nil, errors.New("Hostname not specified")
  }

  port, err := readConfigProperty(fileContents,"port")
  if err != nil {
    return nil,errors.New("Port not specified")
  }

  username, err := readConfigProperty(fileContents,"username")
  if err != nil {
    return nil, errors.New("Username not specified")
  }

  conf := Models.ClientConfig{hostname,port,username}
  return &conf,nil
}

// Look for one specific key inside of the server config. Return value or
// an error if not found
func readConfigProperty(fileContents string, key string) (string,error) {
  key = strings.ToLower(key)
  lines := strings.Split(fileContents,"\n")

  // Loop through each line in the config file
  for _,line := range lines {
    line = strings.ToLower(line)
    // This line is a comment, skip it
    if len(strings.TrimSpace(line)) > 0 && []rune(strings.TrimSpace(line))[0] == '#' {
      continue
    }

    // This line has a key assignment
    if strings.Contains(line,"=") {
      components := strings.Split(line,"=")
      lhs := strings.TrimSpace(components[0])
      rhs := strings.TrimSpace(components[1])
      // Check if it's the key assignment we're interested in
      if lhs == key {
        return rhs,nil
      }
    }
  }

  return "",errors.New(fmt.Sprintf("Key %s not found", key))
}

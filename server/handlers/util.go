package Handlers
import "net/http"
import "strings"

// Common functions for handlers

func EnforceHttpMethods(w http.ResponseWriter, r *http.Request, allowedMethods []string) bool {

  // r has a valid request type/method?
  var valid bool
  valid = false

  // Make sure the given request type (r.Method) matches an allowed type
  for _, elm := range allowedMethods {
    if strings.ToLower(elm) == strings.ToLower(r.Method) {
      valid = true
      break
    }
  }

  // Didn't match any of the allowed types
  if !valid {
    http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
  }

  return valid
}

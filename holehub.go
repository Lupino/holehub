package main

import (
  "github.com/codegangsta/negroni"
  "github.com/tylerb/graceful"
  "github.com/xyproto/permissions2"
  "net/http"
  "fmt"
  "time"
)

func main() {
  mux := http.NewServeMux()

  // New permissions middleware
  perm := permissions.New()

  // Get the userstate, used in the handlers below
  // userstate := perm.UserState()

  mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
    fmt.Fprintf(w, "Hello HoleHub.")
  })

   // Custom handler for when permissions are denied
   perm.SetDenyFunction(func(w http.ResponseWriter, req *http.Request) {
     http.Error(w, "Permission denied!", http.StatusForbidden)
   })

  n := negroni.Classic()

  n.Use(perm)
  n.UseHandler(mux)

  //n.Run(":3000")
  graceful.Run(":3000", 10*time.Second, n)
}

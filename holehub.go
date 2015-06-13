package main

import (
  "github.com/codegangsta/negroni"
  "github.com/tylerb/graceful"
  "net/http"
  "fmt"
  "time"
)

func main() {
  mux := http.NewServeMux()
  mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
    fmt.Fprintf(w, "Hello HoleHub.")
  })

  n := negroni.Classic()
  n.UseHandler(mux)
  //n.Run(":3000")
  graceful.Run(":3000", 10*time.Second, n)
}

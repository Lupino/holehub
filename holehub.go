package main

import (
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/mholt/binding"
	"github.com/tylerb/graceful"
	"github.com/xyproto/permissions2"
	"net/http"
	"time"
)

type UserForm struct {
	Name     string
	Email    string
	Password string
}

func (uf *UserForm) FieldMap() binding.FieldMap {
	return binding.FieldMap{
		&uf.Name:     "username",
		&uf.Email:    "email",
		&uf.Password: "password",
	}
}

func checkMethod(w http.ResponseWriter, req *http.Request, method string) bool {
	if req.Method != method {
		http.Error(w, "404 page not found.", http.StatusNotFound)
		return false
	}
	return true
}

func main() {
	router := mux.NewRouter()

	// New permissions middleware
	perm := permissions.New()

	// Get the userstate, used in the handlers below
	// userstate := perm.UserState()

	router.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Hello HoleHub.")
	})

	router.HandleFunc("/api/users/", func(w http.ResponseWriter, req *http.Request) {
		if !checkMethod(w, req, "POST") {
			return
		}
		userForm := new(UserForm)
		errs := binding.Bind(req, userForm)
		if errs.Handle(w) {
			return
		}
		fmt.Fprintf(w, "register>> userName: %s, method: %s\n", userForm.Name, req.Method)
	})

	// Custom handler for when permissions are denied
	perm.SetDenyFunction(func(w http.ResponseWriter, req *http.Request) {
		http.Error(w, "Permission denied!", http.StatusForbidden)
	})

	n := negroni.Classic()

	n.Use(perm)
	n.UseHandler(router)

	//n.Run(":3000")
	graceful.Run(":3000", 10*time.Second, n)
}

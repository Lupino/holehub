package main

import (
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/mholt/binding"
	e "github.com/pjebs/jsonerror"
	"github.com/tylerb/graceful"
	"github.com/unrolled/render"
	"github.com/xyproto/permissions2"
	"net/http"
	"time"
)

var ErrorMessages = map[int]map[string]string{
	0: e.New(0, "", "Success").Render(),
	1: e.New(1, "User is already exists.", "Please try a new one.").Render(),
	2: e.New(2, "Email is already exists.", "Please try a new one or reset the password.").Render(),
}

type NewUserForm struct {
	Name     string
	Email    string
	Password string
}

func (uf *NewUserForm) FieldMap() binding.FieldMap {
	return binding.FieldMap{
		&uf.Name: binding.Field{
			Form:     "username",
			Required: true,
		},
		&uf.Email: binding.Field{
			Form:     "email",
			Required: true,
		},
		&uf.Password: binding.Field{
			Form:     "password",
			Required: true,
		},
	}
}

func main() {
	router := mux.NewRouter()

	r := render.New()

	// New permissions middleware
	perm := permissions.New()

	// Get the userstate, used in the handlers below
	userstate := perm.UserState()

	creator := userstate.Creator()
	emails, _ := creator.NewKeyValue("emails")

	router.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Hello HoleHub.")
	})

	router.HandleFunc("/api/users/", func(w http.ResponseWriter, req *http.Request) {
		userForm := new(NewUserForm)
		errs := binding.Bind(req, userForm)
		if errs.Handle(w) {
			return
		}
		if userstate.HasUser(userForm.Name) {
			r.JSON(w, http.StatusOK, ErrorMessages[1])
			return
		}
		if name, _ := emails.Get(userForm.Email); name != "" {
			r.JSON(w, http.StatusOK, ErrorMessages[2])
			return
		}
		userstate.AddUser(userForm.Name, userForm.Password, userForm.Email)
		emails.Set(userForm.Email, userForm.Name)
		r.JSON(w, http.StatusOK, ErrorMessages[0])
	}).Methods("POST")

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

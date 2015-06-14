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
	"github.com/xyproto/pinterface"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const HOLE_SERVER = "hole-server"

var defaultMinPort = 10000

var ErrorMessages = map[int]map[string]string{
	0: e.New(0, "", "Success").Render(),
	1: e.New(1, "User is already exists.", "Please try a new one.").Render(),
	2: e.New(2, "Email is already exists.", "Please try a new one or reset the password.").Render(),
	3: e.New(3, "Email format error", "Please type a valid email.").Render(),
	4: e.New(4, "User name or password invalid.", "").Render(),
}

var reEmail, _ = regexp.Compile("(\\w[-._\\w]*\\w@\\w[-._\\w]*\\w\\.\\w{2,3})")

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

type AuthForm struct {
	NameOrEmail string
	Password    string
}

func (af *AuthForm) FieldMap() binding.FieldMap {
	return binding.FieldMap{
		&af.NameOrEmail: binding.Field{
			Form:     "username",
			Required: true,
		},
		&af.Password: binding.Field{
			Form:     "password",
			Required: true,
		},
	}
}

func isEmail(email string) bool {
	return reEmail.MatchString(email)
}

type HoleServer struct {
	addr  string
	ca    string
	cakey string
	cmd   *exec.Cmd
}

func NewHoleServer(addr, ca, cakey string) *HoleServer {
	return &HoleServer{
		addr:  addr,
		ca:    ca,
		cakey: cakey,
	}
}

func (h *HoleServer) Run() error {
	h.cmd = exec.Command(HOLE_SERVER, "-addr", h.addr, "-use-tls", "-ca", h.ca, "-key", h.cakey)
	h.cmd.Stdout = os.Stdout
	h.cmd.Stderr = os.Stderr
	return h.cmd.Run()

}

func (h *HoleServer) Kill() error {
	if h.cmd != nil && h.cmd.Process != nil {
		return h.cmd.Process.Kill()
	}
	return nil
}

func (h *HoleServer) Exited() bool {
	if h.cmd != nil && h.cmd.ProcessState != nil {
		return h.cmd.ProcessState.Exited()
	}
	return true
}

type UsersHoleServer struct {
	state   pinterface.IUserState
	holes   pinterface.IHashMap
	seq     pinterface.IKeyValue
	servers map[string]*HoleServer
}

func NewUsersHoleServer(state pinterface.IUserState) *UsersHoleServer {
	uhs := new(UsersHoleServer)
	creator := state.Creator()
	uhs.state = state
	uhs.holes, _ = creator.NewHashMap("holes")
	uhs.seq, _ = creator.NewKeyValue("seq")
	uhs.servers = make(map[string]*HoleServer)
	return uhs
}

func (h *UsersHoleServer) New(username string) *HoleServer {
	if !h.state.HasUser(username) {
		return nil
	}
	users := h.state.Users()
	port := strconv.Itoa(h.GetLastPort())
	ca := username + "-ca.pem"
	cakey := username + "-ca.key"
	addr := "tcp://:" + port
	users.Set(username, "ca", ca)
	users.Set(username, "cakey", cakey)
	h.holes.Set(port, "ca", ca)
	h.holes.Set(port, "cakey", cakey)
	h.holes.Set(port, "addr", addr)
	userholes, _ := users.Get(username, "holes")
	users.Set(username, "holes", userholes+port+",")
	hs := NewHoleServer(addr, ca, cakey)
	h.servers[port] = hs
	return hs
}

func (h *UsersHoleServer) GetAll(username string) []*HoleServer {
	if !h.state.HasUser(username) {
		return nil
	}
	users := h.state.Users()
	userholes, _ := users.Get(username, "holes")
	ports := strings.Split(userholes, ",")
	servers := make([]*HoleServer, 0)
	var ok bool
	var server *HoleServer
	for _, port := range ports {
		if port == "" {
			continue
		}
		if server, ok = h.servers[port]; !ok {
			addr, _ := h.holes.Get(port, "addr")
			ca, _ := h.holes.Get(port, "ca")
			cakey, _ := h.holes.Get(port, "cakey")
			server = NewHoleServer(addr, ca, cakey)
		}
		servers = append(servers, server)
	}
	return servers
}

func (h *UsersHoleServer) GetLastPort() int {
	lastport, _ := h.seq.Inc("holeserverport")
	port, _ := strconv.Atoi(lastport)
	if port < defaultMinPort {
		port = defaultMinPort
		h.seq.Set("holeserverport", strconv.Itoa(port))
	}
	return port
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

	router.HandleFunc("/api/signup/", func(w http.ResponseWriter, req *http.Request) {
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
		if !isEmail(userForm.Email) {
			r.JSON(w, http.StatusOK, ErrorMessages[3])
			return
		}
		userstate.AddUser(userForm.Name, userForm.Password, userForm.Email)
		emails.Set(userForm.Email, userForm.Name)
		r.JSON(w, http.StatusOK, ErrorMessages[0])
	}).Methods("POST")

	router.HandleFunc("/api/signin/", func(w http.ResponseWriter, req *http.Request) {
		authForm := new(AuthForm)
		errs := binding.Bind(req, authForm)
		if errs.Handle(w) {
			return
		}
		name := authForm.NameOrEmail
		if isEmail(authForm.NameOrEmail) {
			name, _ = emails.Get(authForm.NameOrEmail)
		}
		if !userstate.CorrectPassword(name, authForm.Password) {
			r.JSON(w, http.StatusOK, ErrorMessages[4])
			return
		}
		userstate.Login(w, name)
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

package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/mholt/binding"
	e "github.com/pjebs/jsonerror"
	"github.com/satori/go.uuid"
	"github.com/tylerb/graceful"
	"github.com/unrolled/render"
	"github.com/xyproto/permissions2"
	"github.com/xyproto/pinterface"
	"io/ioutil"
	"log"
	"math/big"
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
var defaultHost = "127.0.0.1"
var defaultCaPath = "certs/"

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

func GenerateUserCa(username string) {
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(1653),
		Subject: pkix.Name{
			Country:            []string{"China"},
			Organization:       []string{"HoleHUB"},
			OrganizationalUnit: []string{"HoleHUB"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		SubjectKeyId:          []byte{1, 2, 3, 4, 5},
		BasicConstraintsValid: true,
		IsCA:        true,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}

	priv, _ := rsa.GenerateKey(rand.Reader, 1024)
	pub := &priv.PublicKey
	ca_b, err := x509.CreateCertificate(rand.Reader, ca, ca, pub, priv)
	if err != nil {
		log.Println("create ca failed", err)
		return
	}
	ca_f := defaultCaPath + username + "-ca.pem"
	log.Println("write to", ca_f)
	ioutil.WriteFile(ca_f, ca_b, 0777)

	priv_f := defaultCaPath + username + "-ca.key"
	priv_b := x509.MarshalPKCS1PrivateKey(priv)
	log.Println("write to", priv_f)
	ioutil.WriteFile(priv_f, priv_b, 0777)
}

type HoleServer struct {
	ID    string
	Addr  string
	Ca    string
	Cakey string
	Cmd   *exec.Cmd
}

func NewHoleServer(ID, addr, ca, cakey string) *HoleServer {
	return &HoleServer{
		ID:    ID,
		Addr:  addr,
		Ca:    ca,
		Cakey: cakey,
	}
}

func (h *HoleServer) Start() error {
	h.Cmd = exec.Command(HOLE_SERVER, "-addr", h.Addr, "-use-tls", "-ca", defaultCaPath+h.Ca, "-key", defaultCaPath+h.Cakey)
	h.Cmd.Stdout = os.Stdout
	h.Cmd.Stderr = os.Stderr
	return h.Cmd.Start()

}

func (h *HoleServer) Kill() error {
	if h.Cmd != nil && h.Cmd.Process != nil {
		return h.Cmd.Process.Kill()
	}
	return nil
}

func (h *HoleServer) Exited() bool {
	if h.Cmd != nil && h.Cmd.ProcessState != nil {
		return h.Cmd.ProcessState.Exited()
	}
	return true
}

type UsersHole struct {
	state   pinterface.IUserState
	holes   pinterface.IHashMap
	seq     pinterface.IKeyValue
	servers map[string]*HoleServer
}

func NewUsersHole(state pinterface.IUserState) *UsersHole {
	uh := new(UsersHole)
	creator := state.Creator()
	uh.state = state
	uh.holes, _ = creator.NewHashMap("holes")
	uh.seq, _ = creator.NewKeyValue("seq")
	uh.servers = make(map[string]*HoleServer)
	return uh
}

func (h *UsersHole) NewHoleServer(username string) *HoleServer {
	if !h.state.HasUser(username) {
		return nil
	}
	users := h.state.Users()
	port := strconv.Itoa(h.GetLastPort())
	ca := username + "-ca.pem"
	cakey := username + "-ca.key"
	addr := "tcp://" + defaultHost + ":" + port
	holeID := uuid.NewV4().String()
	h.holes.Set(holeID, "ca", ca)
	h.holes.Set(holeID, "cakey", cakey)
	h.holes.Set(holeID, "addr", addr)
	userholes, _ := users.Get(username, "holes")
	users.Set(username, "holes", userholes+holeID+",")
	hs := NewHoleServer(holeID, addr, ca, cakey)
	h.servers[holeID] = hs
	return hs
}

func (h *UsersHole) GetAll(username string) []*HoleServer {
	if !h.state.HasUser(username) {
		return nil
	}
	users := h.state.Users()
	userholes, _ := users.Get(username, "holes")
	holeIDs := strings.Split(userholes, ",")
	servers := make([]*HoleServer, 0)
	var ok bool
	var server *HoleServer
	for _, holeID := range holeIDs {
		if holeID == "" {
			continue
		}
		if server, ok = h.servers[holeID]; !ok {
			addr, _ := h.holes.Get(holeID, "addr")
			ca, _ := h.holes.Get(holeID, "ca")
			cakey, _ := h.holes.Get(holeID, "cakey")
			server = NewHoleServer(holeID, addr, ca, cakey)
			h.servers[holeID] = server
		}
		servers = append(servers, server)
	}
	return servers
}

func (h *UsersHole) GetOne(username, holeID string) *HoleServer {
	if !h.state.HasUser(username) {
		return nil
	}
	users := h.state.Users()
	userholes, _ := users.Get(username, "holes")
	if !strings.Contains(userholes, holeID) {
		return nil
	}
	hs, ok := h.servers[holeID]
	if !ok {
		addr, _ := h.holes.Get(holeID, "addr")
		ca, _ := h.holes.Get(holeID, "ca")
		cakey, _ := h.holes.Get(holeID, "cakey")
		hs = NewHoleServer(holeID, addr, ca, cakey)
		h.servers[holeID] = hs
	}
	return hs
}

func (h *UsersHole) GetLastPort() int {
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

	perm.AddUserPath("/api/holes/")

	// Get the userstate, used in the handlers below
	userstate := perm.UserState()

	creator := userstate.Creator()
	emails, _ := creator.NewKeyValue("emails")
	usershole := NewUsersHole(userstate)

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
		GenerateUserCa(userForm.Name)
		users := userstate.Users()
		ca := userForm.Name + "-ca.pem"
		cakey := userForm.Name + "-ca.key"
		users.Set(userForm.Name, "ca", ca)
		users.Set(userForm.Name, "cakey", cakey)
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

	router.HandleFunc("/api/holes/create/", func(w http.ResponseWriter, req *http.Request) {
		username := userstate.Username(req)
		hs := usershole.NewHoleServer(username)
		out := ErrorMessages[0]
		out["ID"] = hs.ID
		r.JSON(w, http.StatusOK, out)
	}).Methods("POST")

	router.HandleFunc("/api/holes/{holeID}/start/", func(w http.ResponseWriter, req *http.Request) {
		holeID := mux.Vars(req)["holeID"]
		username := userstate.Username(req)
		hs := usershole.GetOne(username, holeID)
		hs.Start()
		r.JSON(w, http.StatusOK, ErrorMessages[0])
	}).Methods("POST")

	router.HandleFunc("/api/holes/{holeID}/kill/", func(w http.ResponseWriter, req *http.Request) {
		holeID := mux.Vars(req)["holeID"]
		username := userstate.Username(req)
		hs := usershole.GetOne(username, holeID)
		hs.Kill()
		r.JSON(w, http.StatusOK, ErrorMessages[0])
	}).Methods("POST")

	router.HandleFunc("/api/holes/", func(w http.ResponseWriter, req *http.Request) {
		username := userstate.Username(req)
		holes := usershole.GetAll(username)
		r.JSON(w, http.StatusOK, map[string][]*HoleServer{"holes": holes})
	}).Methods("GET")

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

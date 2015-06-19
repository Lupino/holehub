package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"flag"
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/mholt/binding"
	e "github.com/pjebs/jsonerror"
	"github.com/satori/go.uuid"
	"github.com/sendgrid/sendgrid-go"
	"github.com/tylerb/graceful"
	"github.com/unrolled/render"
	"github.com/xyproto/permissions2"
	"github.com/xyproto/pinterface"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"
)

const HOLE_SERVER = "hole-server"

var defaultMinPort int
var defaultHost string
var configPath string
var port int
var command string
var sg *sendgrid.SGClient

var ErrorMessages = map[int]map[string]string{
	0: e.New(0, "", "Success").Render(),
	1: e.New(1, "User is already exists.", "Please try a new one.").Render(),
	2: e.New(2, "Email is already exists.", "Please try a new one or reset the password.").Render(),
	3: e.New(3, "Email format error", "Please type a valid email.").Render(),
	4: e.New(4, "User name or password invalid.", "").Render(),
	5: e.New(5, "User is confimd or ConfirmationCode is expired.", "Resend a new confirmation code?").Render(),
	6: e.New(6, "User is confimd.", "No need resend twice.").Render(),
}

var reEmail, _ = regexp.Compile("(\\w[-._\\w]*\\w@\\w[-._\\w]*\\w\\.\\w{2,3})")

const runsitConfig = `{
  "cwd": "{{.Cwd}}",
  "standardEnv": true,
  "env": {
    "WANT_USER": ["_env", "want-${USER}"]
  },
  "binary": "{{.Command}}",
  "args": [
    "--addr", "{{.Addr}}",
    "--ca", "{{.Ca}}",
    "--key", "{{.Cakey}}",
    "--use-tls"
  ]
}`

var tpl = template.Must(template.New("runsit").Parse(runsitConfig))

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
	ca_f := configPath + "certs/" + username + "-ca.pem"
	log.Println("write to", ca_f)
	ioutil.WriteFile(ca_f, ca_b, 0777)

	priv_f := configPath + "certs/" + username + "-ca.key"
	priv_b := x509.MarshalPKCS1PrivateKey(priv)
	log.Println("write to", priv_f)
	ioutil.WriteFile(priv_f, priv_b, 0777)
}

func GenerateUserCert(username string) {
	caFile := configPath + "certs/" + username + "-ca.pem"
	privFile := configPath + "certs/" + username + "-ca.key"

	ca_b, _ := ioutil.ReadFile(caFile)
	ca, _ := x509.ParseCertificate(ca_b)
	priv_b, _ := ioutil.ReadFile(privFile)
	priv, _ := x509.ParsePKCS1PrivateKey(priv_b)

	cert2 := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			Country:            []string{"China"},
			Organization:       []string{"HoleHUB"},
			OrganizationalUnit: []string{"HoleHUB"},
		},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}
	priv2, _ := rsa.GenerateKey(rand.Reader, 1024)
	pub2 := &priv2.PublicKey
	cert2_b, err2 := x509.CreateCertificate(rand.Reader, cert2, ca, pub2, priv)
	if err2 != nil {
		log.Println("create cert2 failed", err2)
		return
	}

	cert2_f := configPath + "certs/" + username + "-cert.pem"
	log.Println("write to", cert2_f)
	ioutil.WriteFile(cert2_f, cert2_b, 0777)

	priv2_f := configPath + "certs/" + username + "-cert.key"
	priv2_b := x509.MarshalPKCS1PrivateKey(priv2)
	log.Println("write to", priv2_f)
	ioutil.WriteFile(priv2_f, priv2_b, 0777)
}

func SendConfirmationCode(username, email, confirmationCode string) bool {
	message := sendgrid.NewMail()
	message.AddTo(email)
	message.AddToName(username)
	message.SetSubject("欢迎注册 HoleHUB")
	message.SetText("Hi，欢迎加入HoleHUB！\n\n在这里您可以方便地穿透路由器。\n\n为了保障该帐号可以正常使用，请于24小时内点击以下链接验证您的账号:\nhttp://holehub.com/api/confirm/" + confirmationCode)
	message.SetFrom("support@holehub.com")
	message.SetFromName("HoleHUB Support")
	if r := sg.Send(message); r == nil {
		fmt.Println("Email sent!")
		return true
	} else {
		fmt.Println(r)
		return false
	}
}

type HoleServer struct {
	ID      string
	Cwd     string
	Addr    string
	Ca      string
	Cakey   string
	IsAlive bool
	Command string
}

func NewHoleServer(ID, addr, ca, cakey string) *HoleServer {
	hs := &HoleServer{
		ID:      ID,
		Addr:    addr,
		Ca:      ca,
		Cakey:   cakey,
		Cwd:     configPath + "certs",
		Command: command,
	}
	hs.IsAlive = hs.Alive()
	return hs
}

func (h *HoleServer) Start() error {
	fp, err := os.Create(configPath + h.ID + ".json")
	if err != nil {
		return err
	}
	err = tpl.Execute(fp, h)
	h.IsAlive = true
	return err
}

func (h *HoleServer) Kill() error {
	h.IsAlive = false
	return os.Remove(configPath + h.ID + ".json")
}

func (h *HoleServer) Alive() bool {
	_, err := os.Stat(configPath + h.ID + ".json")
	if err == nil || os.IsExist(err) {
		return true
	}
	return false
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

func init() {
	flag.StringVar(&defaultHost, "host", "127.0.0.1", "The server host.")
	flag.IntVar(&port, "port", 3000, "The server port.")
	flag.StringVar(&configPath, "config_dir", "config/", "The config path.")
	flag.IntVar(&defaultMinPort, "min_port", 10000, "The min hole server port.")
	var sgUser = flag.String("sendgrid_user", "", "The SendGrid username.")
	var sgKey = flag.String("sendgrid_key", "", "The SendGrid password.")
	gopath := os.Getenv("GOPATH")
	flag.StringVar(&command, "cmd", gopath+"/bin/hole-server", "The hole server binary path.")
	flag.Parse()
	sg = sendgrid.NewSendGridClient(*sgUser, *sgKey)
}

func main() {
	router := mux.NewRouter()

	r := render.New()

	// New permissions middleware
	perm := permissions.New()

	perm.AddUserPath("/api/holes/")
	perm.AddUserPath("/api/new_ca/")
	perm.AddUserPath("/api/new_cert/")
	perm.AddUserPath("/api/ca.pem")
	perm.AddUserPath("/api/ca.key")
	perm.AddUserPath("/api/cert.pem")
	perm.AddUserPath("/api/cert.key")

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
		GenerateUserCert(userForm.Name)
		users := userstate.Users()
		ca := userForm.Name + "-ca.pem"
		cakey := userForm.Name + "-ca.key"
		cert := userForm.Name + "-cert.pem"
		certkey := userForm.Name + "-cert.key"
		users.Set(userForm.Name, "ca", ca)
		users.Set(userForm.Name, "cakey", cakey)
		users.Set(userForm.Name, "cert", cert)
		users.Set(userForm.Name, "certkey", certkey)

		code, _ := userstate.GenerateUniqueConfirmationCode()
		userstate.AddUnconfirmed(userForm.Name, code)
		SendConfirmationCode(userForm.Name, userForm.Email, code)

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

	router.HandleFunc("/api/holes/{holeID}/", func(w http.ResponseWriter, req *http.Request) {
		holeID := mux.Vars(req)["holeID"]
		username := userstate.Username(req)
		hs := usershole.GetOne(username, holeID)
		r.JSON(w, http.StatusOK, hs)
	}).Methods("GET")

	router.HandleFunc("/api/holes/", func(w http.ResponseWriter, req *http.Request) {
		username := userstate.Username(req)
		holes := usershole.GetAll(username)
		r.JSON(w, http.StatusOK, map[string][]*HoleServer{"holes": holes})
	}).Methods("GET")

	router.HandleFunc("/api/new_ca/", func(w http.ResponseWriter, req *http.Request) {
		username := userstate.Username(req)
		GenerateUserCa(username)
		r.JSON(w, http.StatusOK, ErrorMessages[0])
	}).Methods("POST")

	router.HandleFunc("/api/ca.pem", func(w http.ResponseWriter, req *http.Request) {
		username := userstate.Username(req)
		data, _ := ioutil.ReadFile(username + "-ca.pem")
		r.Data(w, http.StatusOK, data)
	}).Methods("GET")

	router.HandleFunc("/api/ca.key", func(w http.ResponseWriter, req *http.Request) {
		username := userstate.Username(req)
		data, _ := ioutil.ReadFile(username + "-ca.key")
		r.Data(w, http.StatusOK, data)
	}).Methods("GET")

	router.HandleFunc("/api/new_cert/", func(w http.ResponseWriter, req *http.Request) {
		username := userstate.Username(req)
		GenerateUserCert(username)
		r.JSON(w, http.StatusOK, ErrorMessages[0])
	}).Methods("POST")

	router.HandleFunc("/api/cert.pem", func(w http.ResponseWriter, req *http.Request) {
		username := userstate.Username(req)
		data, _ := ioutil.ReadFile(username + "-cert.pem")
		r.Data(w, http.StatusOK, data)
	}).Methods("GET")

	router.HandleFunc("/api/cert.key", func(w http.ResponseWriter, req *http.Request) {
		username := userstate.Username(req)
		data, _ := ioutil.ReadFile(username + "-cert.key")
		r.Data(w, http.StatusOK, data)
	}).Methods("GET")

	router.HandleFunc("/api/confirm/{confirmationCode}", func(w http.ResponseWriter, req *http.Request) {
		code := mux.Vars(req)["confirmationCode"]
		if err := userstate.ConfirmUserByConfirmationCode(code); err != nil {
			r.JSON(w, http.StatusOK, ErrorMessages[5])
			return
		}
		msg := ErrorMessages[0]
		r.JSON(w, http.StatusOK, msg)
	}).Methods("GET")

	router.HandleFunc("/api/resend/confirmationcode", func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		email := req.Form.Get("email")
		username, _ := emails.Get(email)

		if userstate.IsConfirmed(username) {
			r.JSON(w, http.StatusOK, ErrorMessages[6])
			return
		}

		code, _ := userstate.GenerateUniqueConfirmationCode()
		userstate.AddUnconfirmed(username, code)

		SendConfirmationCode(username, email, code)
		msg := ErrorMessages[0]
		r.JSON(w, http.StatusOK, msg)
	}).Methods("POST")

	// Custom handler for when permissions are denied
	perm.SetDenyFunction(func(w http.ResponseWriter, req *http.Request) {
		http.Error(w, "Permission denied!", http.StatusForbidden)
	})

	n := negroni.Classic()

	n.Use(perm)
	n.UseHandler(router)

	//n.Run(":3000")
	graceful.Run(fmt.Sprintf("%s:%d", defaultHost, port), 10*time.Second, n)
}

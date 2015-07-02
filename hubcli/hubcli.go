package main

import (
	"fmt"
	"github.com/Lupino/hole"
	"github.com/codegangsta/cli"
	"github.com/levigross/grequests"
	"github.com/xyproto/simplebolt"
	"io"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type JE struct {
	Code    string `json:"code"`
	Domain  string `json:"domain"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

var defaultReTryTime = 1000
var reTryTimes = defaultReTryTime

var boltFile = os.Getenv("HOME") + "/.holehub.db"
var certFile = "/tmp/cert.pem"
var privFile = "/tmp/cert.key"

var db *simplebolt.Database
var config *simplebolt.KeyValue
var holes *simplebolt.HashMap
var apps *simplebolt.Set
var appNames *simplebolt.KeyValue

func init() {
	initDB()
}

func initDB() {
	var err error
	db, err = simplebolt.New(boltFile)
	if err != nil {
		log.Fatalf("Could not create database! %s", err)
	}
	config, _ = simplebolt.NewKeyValue(db, "config")
	holes, _ = simplebolt.NewHashMap(db, "holes")
	apps, _ = simplebolt.NewSet(db, "apps")
	appNames, _ = simplebolt.NewKeyValue(db, "appnames")
}

func Login(host string) {
	name, _ := config.Get("email")
	passwd, _ := config.Get("password")
	if name == "" || passwd == "" {
		log.Fatalf("Error: email or password is not config\n")
	}

	ro := &grequests.RequestOptions{
		Data: map[string]string{"username": name, "password": passwd},
	}

	rsp, err := grequests.Post(host+"/api/signin/", ro)
	if err != nil {
		log.Fatal(err)
	}
	defer rsp.Close()

	if !rsp.Ok {
		log.Fatalf("Error: %s\n", rsp.String())
	}

	var msg JE
	err = rsp.JSON(&msg)

	if err != nil {
		log.Fatal(err)
	}

	if msg.Code != "0" {
		fmt.Printf("Error: %s\n", msg.Error)
		os.Exit(1)
	} else {
		fmt.Printf("Login HoleHUB %s\n", msg.Message)
		cookie := rsp.Header.Get("Set-Cookie")
		config.Set("cookie", cookie)
	}
}

func Ping(host string) bool {
	cookie, _ := config.Get("cookie")
	var ro = &grequests.RequestOptions{
		Headers: map[string]string{"Cookie": cookie},
	}

	rsp, err := grequests.Get(host+"/api/ping/", ro)
	if err != nil {
		log.Fatal(err)
	}
	defer rsp.Close()

	if !rsp.Ok {
		log.Fatalf("Error: %s\n", rsp.String())
	}

	var pong = true
	if rsp.String() == "false" {
		pong = false
	}
	return pong
}

type HoleApp struct {
	ID      string
	Name    string
	Port    string
	Scheme  string
	Lscheme string
	Lport   string
	Status  string
	Pid     int
}

func NewHoleApp(ID string) (holeApp HoleApp, err error) {
	if ok, _ := apps.Has(ID); !ok {
		err = fmt.Errorf("hole app: not exists.")
		return
	}
	holeApp = HoleApp{ID: ID}
	holeApp.Name, _ = holes.Get(ID, "name")
	holeApp.Port, _ = holes.Get(ID, "port")
	holeApp.Scheme, _ = holes.Get(ID, "scheme")
	holeApp.Lport, _ = holes.Get(ID, "local-port")
	holeApp.Lscheme, _ = holes.Get(ID, "local-scheme")
	holeApp.Status, _ = holes.Get(ID, "status")

	if holeApp.Status == "started" {
		Pid, _ := holes.Get(ID, "pid")
		holeApp.Pid, _ = strconv.Atoi(Pid)
	}

	return holeApp, nil
}

func NewHoleAppByName(name string) (holeApp HoleApp, err error) {
	holeID, _ := appNames.Get(name)
	if holeID == "" {
		err = fmt.Errorf("hole app: not exists.")
		return
	}
	holeApp, err = NewHoleApp(holeID)
	return
}

func (hole HoleApp) run(host, command string) {
	cookie, _ := config.Get("cookie")
	var ro = &grequests.RequestOptions{
		Headers: map[string]string{"Cookie": cookie},
	}

	rsp, err := grequests.Post(host+"/api/holes/"+hole.ID+"/"+command+"/", ro)
	if err != nil {
		log.Fatal(err)
	}
	defer rsp.Close()

	if !rsp.Ok {
		log.Fatalf("Error: %s\n", rsp.String())
	}

	var msg JE
	err = rsp.JSON(&msg)

	if err != nil {
		log.Fatal(err)
	}

	if msg.Code != "0" {
		fmt.Printf("Error: %s\n", msg.Error)
		os.Exit(1)
	}
}

func (hole HoleApp) Start(host string) {
	hole.run(host, "start")
	holes.Set(hole.ID, "status", "started")
	holes.Set(hole.ID, "pid", strconv.Itoa(os.Getpid()))
}

func (hole HoleApp) Kill(host string) {
	if db == nil {
		initDB()
	}
	hole.run(host, "kill")
	holes.Set(hole.ID, "status", "stoped")
}

func createHoleApp(host, scheme, name string) HoleApp {
	cookie, _ := config.Get("cookie")
	var ro = &grequests.RequestOptions{
		Headers: map[string]string{"Cookie": cookie},
		Data:    map[string]string{"scheme": scheme, "name": name},
	}

	rsp, err := grequests.Post(host+"/api/holes/create/", ro)
	if err != nil {
		log.Fatal(err)
	}
	defer rsp.Close()

	if !rsp.Ok {
		log.Fatalf("Error: %s\n", rsp.String())
	}

	var msg map[string]HoleApp
	err = rsp.JSON(&msg)

	if err != nil {
		log.Fatal(err)
	}

	hole := msg["hole"]
	holes.Set(hole.ID, "name", hole.Name)
	holes.Set(hole.ID, "scheme", hole.Scheme)
	holes.Set(hole.ID, "port", hole.Port)
	holes.Set(hole.ID, "status", "stoped")
	apps.Add(hole.ID)
	if hole.Name != "" {
		appNames.Set(hole.Name, hole.ID)
	}

	return hole
}

func getCert(host, name, outName string) {
	cookie, _ := config.Get("cookie")
	var ro = &grequests.RequestOptions{
		Headers: map[string]string{"Cookie": cookie},
	}

	rsp, err := grequests.Get(host+"/api/"+name, ro)
	if err != nil {
		log.Fatal(err)
	}
	defer rsp.Close()

	if !rsp.Ok {
		log.Fatalf("Error: %s\n", rsp.String())
	}

	fp, _ := os.Create(outName)
	io.Copy(fp, rsp)
	fp.Close()
}

func processHoleClient(host string, holeApp HoleApp) {
	getCert(host, "cert.pem", certFile)
	getCert(host, "cert.key", privFile)
	db.Close()
	db = nil

	var realAddr = holeApp.Lscheme + "://127.0.0.1:" + holeApp.Lport
	var hostPort = strings.Split(host, "://")[1]
	var parts = strings.Split(hostPort, ":")
	var serverAddr = holeApp.Scheme + "://" + parts[0] + ":" + holeApp.Port
	var client = hole.NewClient(realAddr)
	client.ConfigTLS(certFile, privFile)

	for {
		if err := client.Connect(serverAddr); err == nil {
			break
		}
		reTryTimes = reTryTimes - 1
		if reTryTimes == 0 {
			log.Fatal("Error: unable to connect %s\n", serverAddr)
		}
		log.Printf("Retry after 2 second...")
		time.Sleep(2 * time.Second)
	}

	fmt.Printf("Publish: %s\n", serverAddr)
	client.Process()
}

func Run(host, scheme, name, port string) {
	if !Ping(host) {
		Login(host)
	}

	holeApp := createHoleApp(host, scheme, name)
	holes.Set(holeApp.ID, "local-port", port)
	holes.Set(holeApp.ID, "local-scheme", scheme)
	holeApp.Lport = port
	holeApp.Lscheme = scheme

	holeApp.Start(host)
	defer holeApp.Kill(host)
	go processHoleClient(host, holeApp)
	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt, os.Kill)
	<-s
}

func ListApp(host string) {
	holeIDs, _ := apps.GetAll()
	var hostPort = strings.Split(host, "://")[1]
	host = strings.Split(hostPort, ":")[0]
	fmt.Println("ID\t\t\t\t\tName\t\tPort\t\t\t\t\tStatus")
	for _, holeID := range holeIDs {
		holeApp, err := NewHoleApp(holeID)
		if err != nil {
			continue
		}
		fmt.Printf("%s\t%s\t\t127.0.0.1:%s/%s->%s:%s/%s\t%s\n", holeApp.ID,
			holeApp.Name, holeApp.Lport, holeApp.Lscheme, host, holeApp.Port, holeApp.Scheme, holeApp.Status)
	}
}

func StartApp(host, nameOrID string) {
	var holeApp HoleApp
	var err error
	if holeApp, err = NewHoleAppByName(nameOrID); err != nil {
		if holeApp, err = NewHoleApp(nameOrID); err != nil {
			log.Fatal(err)
		}
	}
	if holeApp.Status == "started" {
		log.Fatalf("HoleApp: %s is already started.", nameOrID)
	}

	holeApp.Start(host)
	defer holeApp.Kill(host)
	go processHoleClient(host, holeApp)
	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt, os.Kill)
	<-s
}

func StopApp(nameOrID string) {
	var holeApp HoleApp
	var err error
	if holeApp, err = NewHoleAppByName(nameOrID); err != nil {
		if holeApp, err = NewHoleApp(nameOrID); err != nil {
			log.Fatal(err)
		}
	}
	if holeApp.Status == "stoped" {
		log.Fatalf("HoleApp: %s is already stoped.", nameOrID)
	}

	syscall.Kill(holeApp.Pid, syscall.SIGINT)
}

func main() {
	app := cli.NewApp()
	app.Name = "hubcli"
	app.Usage = "HoleHUB command line."
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "host, H",
			Value:  "http://holehub.com",
			Usage:  "The HoleHUB Host",
			EnvVar: "HOLEHUB_HOST",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:  "login",
			Usage: "Login HoleHUB",
			Action: func(c *cli.Context) {
				Login(c.GlobalString("host"))
			},
		},
		{
			Name:        "config",
			Usage:       "Config HoleHUB cli",
			Description: "config set key value\n   config get key",
			Action: func(c *cli.Context) {
				var args = c.Args()
				switch args.First() {
				case "get":
					if len(args) != 2 {
						fmt.Printf("Not enough arguments.\n\n")
						cli.ShowCommandHelp(c, "config")
						os.Exit(1)
					}
					var value, _ = config.Get(args[1])
					fmt.Printf("%s\n", value)
					return
				case "set":
					if len(args) != 3 {
						fmt.Printf("Not enough arguments.\n\n")
						cli.ShowCommandHelp(c, "config")
						os.Exit(1)
					}
					config.Set(args[1], args[2])
				default:
					cli.ShowCommandHelp(c, "config")
				}
			},
		},
		{
			Name:  "run",
			Usage: "Create and run a new holeapp",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "scheme, s",
					Value: "tcp",
					Usage: "The scheme. tcp udp tcp6 udp6",
				},
				cli.StringFlag{
					Name:  "name, n",
					Value: "",
					Usage: "The app name.",
				},
				cli.StringFlag{
					Name:  "port, p",
					Value: "8080",
					Usage: "The source server port.",
				},
			},
			Action: func(c *cli.Context) {
				var scheme = c.String("scheme")
				var name = c.String("name")
				var port = c.String("port")
				Run(c.GlobalString("host"), scheme, name, port)
			},
		},
		{
			Name:  "ls",
			Usage: "List HoleApps",
			Action: func(c *cli.Context) {
				ListApp(c.GlobalString("host"))
			},
		},
		{
			Name:        "start",
			Usage:       "Start a HoleApp",
			Description: "start name\n   start ID",
			Action: func(c *cli.Context) {
				if len(c.Args()) == 0 {
					fmt.Printf("Not enough arguments.\n\n")
					cli.ShowCommandHelp(c, "start")
					os.Exit(1)
				}
				StartApp(c.GlobalString("host"), c.Args().First())
			},
		},
		{
			Name:        "stop",
			Usage:       "Stop a started HoleApp",
			Description: "stop name\n   stop ID",
			Action: func(c *cli.Context) {
				if len(c.Args()) == 0 {
					fmt.Printf("Not enough arguments.\n\n")
					cli.ShowCommandHelp(c, "stop")
					os.Exit(1)
				}
				StopApp(c.Args().First())
			},
		},
	}

	app.Action = func(c *cli.Context) {
		cli.ShowAppHelp(c)
	}

	app.Run(os.Args)
}

package main

import (
	"fmt"
	"github.com/Lupino/hole"
	"github.com/codegangsta/cli"
	"github.com/levigross/grequests"
	"github.com/xyproto/simplebolt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strings"
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

func init() {
	var err error
	db, err = simplebolt.New(boltFile)
	if err != nil {
		log.Fatalf("Could not create database! %s", err)
	}
	config, _ = simplebolt.NewKeyValue(db, "config")
	holes, _ = simplebolt.NewHashMap(db, "holes")
	apps, _ = simplebolt.NewSet(db, "apps")
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
	ID     string
	Name   string
	Port   string
	Scheme string
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
}

func (hole HoleApp) Kill(host string) {
	hole.run(host, "kill")
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
	apps.Add(hole.ID)

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

	ioutil.WriteFile(outName, rsp.Bytes(), 0444)
}

func Run(host, scheme, name, port string) {
	if !Ping(host) {
		Login(host)
	}

	holeApp := createHoleApp(host, scheme, name)
	holes.Set(holeApp.ID, "local-port", port)
	holes.Set(holeApp.ID, "local-scheme", scheme)

	holeApp.Start(host)
	defer holeApp.Kill(host)

	getCert(host, "cert.pem", certFile)
	getCert(host, "cert.key", privFile)

	var realAddr = scheme + "://127.0.0.1:" + port
	var hostPort = strings.Split(host, "://")[1]
	var parts = strings.Split(hostPort, ":")
	var serverAddr = scheme + "://" + parts[0] + ":" + holeApp.Port
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
	go client.Process()
	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt, os.Kill)
	<-s
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
	}

	app.Action = func(c *cli.Context) {
		cli.ShowAppHelp(c)
	}

	app.Run(os.Args)
}

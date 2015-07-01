package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/levigross/grequests"
	"github.com/xyproto/simplebolt"
	"log"
	"os"
)

type JE struct {
	Code    string `json:"code"`
	Domain  string `json:"domain"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

var boltFile = "/tmp/bolt.db"

var db *simplebolt.Database
var config *simplebolt.KeyValue

func init() {
	var err error
	db, err = simplebolt.New(boltFile)
	if err != nil {
		log.Fatalf("Could not create database! %s", err)
	}
	config, _ = simplebolt.NewKeyValue(db, "config")
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
	}

	app.Action = func(c *cli.Context) {
		cli.ShowAppHelp(c)
	}

	app.Run(os.Args)
}

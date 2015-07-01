package main

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/xyproto/simplebolt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
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

func GetDB() *simplebolt.Database {
	var err error
	if db == nil {
		db, err = simplebolt.New(boltFile)
		if err != nil {
			log.Fatalf("Could not create database! %s", err)
		}
	}
	return db
}

func GetConfig() *simplebolt.KeyValue {
	var err error
	if config == nil {
		db := GetDB()
		config, err = simplebolt.NewKeyValue(db, "config")
		if err != nil {
			log.Fatalf("Could not create database! %s", err)
		}
	}
	return config
}

func Login(host string) {
	var config = GetConfig()
	name, _ := config.Get("email")
	passwd, _ := config.Get("password")
	if name == "" || passwd == "" {
		log.Fatalf("Error: email or password is not config\n")
	}

	var data = url.Values{}
	data.Set("username", name)
	data.Set("password", passwd)
	res, err := http.PostForm(host+"/api/signin/", data)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	b, _ := ioutil.ReadAll(res.Body)

	if res.StatusCode != 200 {
		fmt.Printf("%s\n", b)
		return
	}

	var msg JE
	err = json.Unmarshal(b, &msg)
	if err != nil {
		log.Fatal(err)
	}

	if msg.Code != "0" {
		fmt.Printf("%s\n", msg.Error)
	} else {
		fmt.Printf("%s\n", msg.Message)
		cookie := res.Header.Get("Set-Cookie")
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
				var config = GetConfig()
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
					if len(args) != 2 {
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

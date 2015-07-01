package main

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
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

var rcFile = ".hubclirc"

func Login(host, name, passwd string) {
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
		ioutil.WriteFile(rcFile, []byte(res.Header.Get("Set-Cookie")), 0444)
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "hubcli"
	app.Usage = "HoleHUB command line."
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "H",
			Value:  "http://holehub.com",
			Usage:  "The HoleHUB Host",
			EnvVar: "HOLEHUB_HOST",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:  "login",
			Usage: "Login HoleHUB",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "name",
					Value: "",
					Usage: "Username or email",
				},
				cli.StringFlag{
					Name:  "password",
					Value: "",
					Usage: "Password",
				},
			},
			Action: func(c *cli.Context) {
				Login(c.GlobalString("H"), c.String("name"), c.String("password"))
			},
		},
	}

	app.Action = func(c *cli.Context) {
		cli.ShowAppHelp(c)
	}

	app.Run(os.Args)
}

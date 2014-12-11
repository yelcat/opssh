package main

import "fmt"

import (
	"bytes"
	"code.google.com/p/go.crypto/ssh"
	"flag"
	"github.com/fatih/color"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type Config struct {
	Hosts []string
	Auth  Auth
}

type Auth struct {
	Name string
	Pass string
}

func genCommand(task string) string {
	switch task {
	case "uptime":
		return "uptime"
	default:
		return task
	}
}

func execute(config *ssh.ClientConfig, hostname string, task string) {
	client, err := ssh.Dial("tcp", hostname+":22", config)
	if err != nil {
		panic("unable to connect: %s" + err.Error())
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		panic("Failed to create session: " + err.Error())
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	cmd := genCommand(task)
	if err := session.Run(cmd); err != nil {
		color.Red("[%s]", hostname)
		fmt.Print(err)
	} else {
		color.Green("[%s]", hostname)
		fmt.Print(b.String())
	}
}

func main() {
	bytes, err = ioutil.ReadFile("opssh.yaml")
	if err != nil {
		fmt.Print("read config file opssh.yaml error: " + err)
	}

	var conf Conf
	err = yaml.Unmarshal(bytes, &conf)
	if err != nil {
		panic(err)
	}

	flag.Parse()

	config := &ssh.ClientConfig{
		User: conf.auth.user,
		Auth: []ssh.AuthMethod{
			ssh.Password(conf.auth.pass),
		},
	}

	for _, hostname := range conf.hosts {
		execute(config, hostname, flag.Arg(0))
	}
}

package main

import (
	"github.com/GeertJohan/go.linenoise"
	"github.com/Sirupsen/logrus"
	"os"
	"strings"
)

func main() {
	logrus.Println("---- Ulule CLI ----")

	username := ""
	apiKey := ""

	// get username from args
	if len(os.Args) > 1 {
		username = os.Args[1]
	}

	// get apiKey from args
	if len(os.Args) > 2 {
		apiKey = os.Args[2]
	}

	var err error

	for username == "" {
		username, err = linenoise.Line("username> ")
		if err != nil {
			username = ""
		}
	}

	for apiKey == "" {
		apiKey, err = linenoise.Line("apiKey> ")
		if err != nil {
			apiKey = ""
		}
	}

	var completionHandler = func(input string) []string {
		return []string{"projects"}
	}

	linenoise.SetCompletionHandler(completionHandler)

	for {
		cmd, err := linenoise.Line("> ")
		if err != nil {
			logrus.Fatal(err)
		}

		args := strings.Split(cmd, " ")

		if len(args) > 0 {
			switch args[0] {
			case "projects":

			}
		}
	}
}

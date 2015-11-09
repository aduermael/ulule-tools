package main

import (
	"fmt"
	"github.com/GeertJohan/go.linenoise"
	"github.com/Sirupsen/logrus"
	"os"
	"strconv"
	"strings"
	"ulule/clientapi"
)

func main() {
	logrus.Println("---- Ulule CLI ----")

	username := ""
	apikey := ""

	// get username from args
	if len(os.Args) > 1 {
		username = os.Args[1]
	}

	// get apikey from args
	if len(os.Args) > 2 {
		apikey = os.Args[2]
	}

	var err error

	for username == "" {
		username, err = linenoise.Line("username> ")
		if err != nil {
			logrus.Fatal(err)
		}
	}

	for apikey == "" {
		apikey, err = linenoise.Line("apikey> ")
		if err != nil {
			logrus.Fatal(err)
		}
	}

	ululeClient := clientapi.New(username, apikey)

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
				projects, err := ululeClient.GetProjects("created")
				if err != nil {
					logrus.Fatal(err)
				}
				// logrus.Printf("projects: %#v", projects[0])
				for _, project := range projects {
					percentage := int(float32(project.AmountRaised) / float32(project.Goal) * 100.0)
					percentageStr := strconv.Itoa(percentage)
					fmt.Println("-", project.Slug, "|", project.AmountRaised, project.CurrencyDisplay, "|", percentageStr+"%")
				}
			}
		}
	}
}

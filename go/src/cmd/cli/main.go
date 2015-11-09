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

var (
	selectedProject *clientapi.Project
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

	linenoise.SetCompletionHandler(completionHandler)

	for {
		cmd, err := linenoise.Line("> ")
		if err != nil {
			logrus.Fatal(err)
		}

		args := strings.Split(cmd, " ")

		if len(args) > 0 {
			switch args[0] {
			case "project":
				if len(args) > 1 {
					switch args[1] {
					case "list":
						projects, err := ululeClient.GetProjects("created")
						if err != nil {
							logrus.Fatal(err)
						}
						// logrus.Printf("projects: %#v", projects[0])
						for _, project := range projects {
							percentage := int(float32(project.AmountRaised) / float32(project.Goal) * 100.0)
							percentageStr := strconv.Itoa(percentage)
							fmt.Println(project.Id, "|", project.Slug, "|", project.AmountRaised, project.CurrencyDisplay, "|", percentageStr+"%")
						}
					case "select":
						if len(args) > 2 {
							selectedProject, err = ululeClient.GetProject(args[2])
							if err != nil {
								fmt.Println(err.Error())
							} else {
								fmt.Println("project selected:", selectedProject.Id, "|", selectedProject.Slug)
							}
						} else {
							fmt.Println("error: `project select` expects a project id or slug argument")
						}
					}
				}
			}
		}
	}
}

func completionHandler(input string) []string {

	input = strings.Trim(input, " ")

	commands := []string{
		"project list",
		"project select",
	}

	autocomplete := []string{}

	for _, cmd := range commands {
		if strings.HasPrefix(cmd, input) {
			autocomplete = append(autocomplete, cmd)
		}
	}

	return autocomplete
}

package main

import (
	"crypto/tls"
	"encoding/json"
	"github.com/Sirupsen/logrus"
	"io"
	"net/http"
	"os"
)

type projects struct {
	Projects []project `json:"projects"`
}

type project struct {
	Id       int    `json:"id"`
	Url      string `json:"absolute_url"`
	Goal     int    `json:"goal"`
	Commited int    `json:"committed"`
}

func main() {

	username := ""
	apiKey := ""

	args := os.Args
	if len(args) > 2 {
		username = args[1]
		apiKey = args[2]
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("GET", "https://api.ulule.com/v1/users/bloglaurel/projects?filter=created", nil)

	// project id: 31458
	// req, err := http.NewRequest("GET", "https://api.ulule.com/v1/projects/31458/orders", nil)

	req.Header.Add("Authorization", "ApiKey "+username+":"+apiKey)

	resp, err := client.Do(req)
	if err != nil {
		logrus.Fatal(err.Error())
	}

	decoder := json.NewDecoder(resp.Body)

	var projs projects

	for {
		err = decoder.Decode(&projs)
		if err != nil && err != io.EOF {
			logrus.Fatal(err.Error())
		}

		if err != nil && err == io.EOF {
			logrus.Printf("%#v", projs)
			break
		}
	}

	//p := make([]byte, 64)
	// for {
	// n, err := resp.Body.Read(p)
	// if err != nil && err != io.EOF {
	// 	logrus.Fatal(err)
	// }
	// fmt.Printf("%s", string(p[:n]))

	// if err != nil && err == io.EOF {
	// 	fmt.Printf("\n")
	// 	break
	// }
	// }
}

package clientapi

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Sirupsen/logrus"
	"io"
	"net/http"
)

// ClientAPI is a structure that can
// be used to get data from Ulule's API
type ClientAPI struct {
	username   string
	apikey     string
	httpClient *http.Client
}

// New returns a ClientAPI structure
// initialized with given credentials
func New(username, apikey string) *ClientAPI {

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	httpClient := &http.Client{Transport: transport}

	clientAPI := &ClientAPI{
		username:   username,
		apikey:     apikey,
		httpClient: httpClient,
	}

	return clientAPI
}

// GetProjects returns ClientAPI user's projects.
// Supported string filters: created, followed, supported
func (c *ClientAPI) GetProjects(filter string) ([]*Project, error) {
	if filter != "created" && filter != "followed" && filter != "supported" {
		return nil, errors.New("ClientAPI GetProjects error: string filter not supported (" + filter + ")")
	}

	req, err := http.NewRequest("GET", "https://api.ulule.com/v1/users/"+c.username+"/projects?filter="+filter, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "ApiKey "+c.username+":"+c.apikey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	listProjectResp := &ListProjectResponse{}
	decodeHTMLBody(resp, listProjectResp)

	return listProjectResp.Projects, nil
}

// HTML utils

func decodeHTMLBody(response *http.Response, i interface{}) error {
	decoder := json.NewDecoder(response.Body)
	for {
		err := decoder.Decode(i)
		if err != nil && err != io.EOF {
			return err
		}
		if err != nil && err == io.EOF {
			break
		}
	}
	return nil
}

func logHTMLBody(response *http.Response) {
	p := make([]byte, 64)
	for {
		n, err := response.Body.Read(p)
		if err != nil && err != io.EOF {
			logrus.Fatal(err)
		}
		fmt.Printf("%s", string(p[:n]))
		if err != nil && err == io.EOF {
			fmt.Printf("\n")
			break
		}
	}
}

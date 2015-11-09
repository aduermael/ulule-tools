package clientapi

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Sirupsen/logrus"
	"io"
	"net/http"
	"strconv"
	"strings"
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
// Supported string filters: "created", "followed", "supported", "" (no filter)
func (c *ClientAPI) GetProjects(filter string) ([]*Project, error) {
	if filter != "created" && filter != "followed" && filter != "supported" && filter != "" {
		return nil, errors.New("error: string filter not supported (" + filter + ")")
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

// GetProject returns one specific ClientAPI user's
// project identified by its Id or Slug.
func (c *ClientAPI) GetProject(identifier string) (*Project, error) {
	identifier = strings.Trim(identifier, " ")
	projects, err := c.GetProjects("")
	if err != nil {
		return nil, err
	}
	for _, project := range projects {
		if identifier == project.Slug || identifier == strconv.Itoa(project.Id) {
			return project, nil
		}
	}
	return nil, errors.New("error: project not found (" + identifier + ")")
}

// GetProjectSupporters lists supporters for a project
// limit and offset stand for pagination
// the boolean returns indicates if it was the last
// page of supporters or not.
func (c *ClientAPI) GetProjectSupporters(projectID, limit, offset int) ([]*Supporter, error, bool) {

	projectIDStr := strconv.Itoa(projectID)

	req, err := http.NewRequest("GET", "https://api.ulule.com/v1/projects/"+projectIDStr+"/supporters?limit="+strconv.Itoa(limit)+"&offset="+strconv.Itoa(offset), nil)
	if err != nil {
		return nil, err, false
	}

	req.Header.Add("Authorization", "ApiKey "+c.username+":"+c.apikey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err, false
	}

	listSupporterResp := &ListSupporterResponse{}
	decodeHTMLBody(resp, listSupporterResp)

	lastPage := listSupporterResp.Meta.Next == ""

	return listSupporterResp.Supporters, nil, lastPage
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

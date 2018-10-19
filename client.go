package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	baseURL string
	token   string
	c       *http.Client
}

func NewClient(baseURL, token string) *Client {
	httpClient := http.DefaultClient
	httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	return &Client{
		baseURL: baseURL,
		token:   token,
		c:       httpClient,
	}
}

func (c *Client) getPath(path string) string {
	separator := "?"
	if strings.Contains(path, separator) {
		separator = "&"
	}
	return fmt.Sprintf("%s/%s%sprivate_token=%s&per_page=1000", c.baseURL, path, separator, c.token)
}

func (c *Client) doRequest(method, path string, body interface{}) (*http.Response, error) {
	var buf io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		buf = bytes.NewBuffer(b)
	}
	req, err := http.NewRequest(method, c.getPath(path), buf)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Close = true
	return c.c.Do(req)
}

func (c *Client) doFormRequest(method, path string, values map[string]interface{}) (*http.Response, error) {
	form := url.Values{}
	for key, val := range values {
		form.Add(key, fmt.Sprintf("%v", val))
	}
	req, err := http.NewRequest(method, c.getPath(path), strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "multipart/form-data")
	req.PostForm = form
	return c.c.Do(req)
}

type Project map[string]interface{}

func (p Project) Get(key string) interface{} {
	return p[key]
}

// ref: https://docs.gitlab.com/ee/api/groups.html#list-a-group-s-projects
func (c *Client) GetGroupProjects(id int) ([]*Project, error) {
	resp, err := c.doRequest("GET", fmt.Sprintf("groups/%d/projects", id), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	projects := []*Project{}
	if err := json.Unmarshal(bytes, &projects); err != nil {
		return nil, err
	}
	return projects, nil
}

// ref: https://docs.gitlab.com/ee/api/projects.html#edit-project
func (c *Client) UpdateProject(id float64, settings map[string]interface{}) error {
	resp, err := c.doFormRequest("PUT", fmt.Sprintf("projects/%v", id), settings)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Unexpected response code: %s", resp.Status)
	}
	return nil
}

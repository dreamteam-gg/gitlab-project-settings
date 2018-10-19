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
	return fmt.Sprintf("%s/%s%sprivate_token=%s&per_page=100", c.baseURL, path, separator, c.token)
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
	var projects []*Project
	err := c.readPaginatedProjects(
		fmt.Sprintf("groups/%d/projects", id),
		func(obj interface{}) {
			resp := obj.([]*Project)
			for _, project := range resp {
				projects = append(projects, project)
			}
		},
	)
	return projects, err
}

// ref: https://docs.gitlab.com/ee/api/projects.html#edit-project
func (c *Client) UpdateProject(id float64, settings map[string]interface{}) error {
	resp, err := c.doFormRequest(http.MethodPut, fmt.Sprintf("projects/%v", id), settings)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Unexpected response code: %s", resp.Status)
	}
	return nil
}

func (c *Client) readPaginatedProjects(path string, accumulate func(interface{})) error {
	pagedURL, err := url.Parse(path)
	if err != nil {
		return err
	}
	values := url.Values{
		"page":     []string{"1"},
		"per_page": []string{"50"},
	}
	pagedURL.RawQuery = values.Encode()
	for {
		resp, err := c.doRequest(http.MethodGet, pagedURL.String(), nil)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode < 200 || resp.StatusCode > 299 {
			return fmt.Errorf("return code not 2XX: %s", resp.Status)
		}
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		obj := []*Project{}
		if err := json.Unmarshal(b, &obj); err != nil {
			return err
		}
		accumulate(obj)
		page := resp.Header.Get("x-next-page")
		if len(page) == 0 {
			break
		}
		q := pagedURL.Query()
		q.Set("page", page)
		pagedURL.RawQuery = q.Encode()
	}
	return nil
}

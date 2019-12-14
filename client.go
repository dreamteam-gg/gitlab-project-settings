package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
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
	return fmt.Sprintf("%s/%s%sprivate_token=%s", c.baseURL, path, separator, c.token)
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
		return nil, fmt.Errorf("Error querying %s", path)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Close = true
	return c.c.Do(req)
}

func (c *Client) doFormRequest(method, path string, values map[string]interface{}) ([]byte, error) {
	s, err := json.Marshal(values)
	if err != nil {
		return nil, err
	}

	b := bytes.NewBuffer(s)
	req, err := http.NewRequest(method, c.getPath(path), b)
	if err != nil {
		return nil, fmt.Errorf("Error querying %s", path)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	r, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("return code not 2XX: %s, message: %s", resp.Status, string(r))
	}

	if err != nil {
		return nil, err
	}

	return r, nil
}

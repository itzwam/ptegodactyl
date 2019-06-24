package ptegodactyl

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
)

// Client manage communication with API
type Client struct {
	APIToken   string
	BaseURL    *url.URL
	UserAgent  string
	httpClient *http.Client
}

// NewClient returns a ready to use client
func NewClient(apiURL string, token string) (c *Client) {
	u, err := url.Parse(apiURL)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return &Client{
		APIToken:   token,
		BaseURL:    u,
		UserAgent:  "Ptegodactyl",
		httpClient: http.DefaultClient,
	}
}

func (c *Client) newRequest(method, path string, body interface{}) (*http.Request, error) {
	u, err := url.Parse(c.BaseURL.String() + path)
	if err != nil {
		return nil, err
	}
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "Application/vnd.pterodactyl.v1+json")
	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("Authorization", "Bearer "+c.APIToken)
	return req, nil
}

func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&v)
	return resp, err
}

// AnswerList is the format of a List answer
type AnswerList struct {
	Data json.RawMessage `json:"data"`
	Meta struct {
		Pagination struct {
			Count       int64         `json:"count"`
			CurrentPage int64         `json:"current_page"`
			Links       []interface{} `json:"links"`
			PerPage     int64         `json:"per_page"`
			Total       int64         `json:"total"`
			TotalPages  int64         `json:"total_pages"`
		} `json:"pagination"`
	} `json:"meta"`
	Object string `json:"object"`
}

// List return a list of things
func (c *Client) list(path string, v interface{}) error {
	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return err
	}
	var a AnswerList
	_, err = c.do(req, &a)
	if err != nil {
		return err
	}
	err = json.Unmarshal(a.Data, &v)

	return nil
}

// Get return things
func (c *Client) get(path string, v interface{}) error {
	req, err := c.newRequest("GET", path, nil)
	if err != nil {
		return err
	}
	_, err = c.do(req, &v)
	if err != nil {
		return err
	}

	return nil
}

// Send sends infos to API
func (c *Client) send(path string, body interface{}, v interface{}) error {
	req, err := c.newRequest("POST", path, body)
	if err != nil {
		return err
	}

	_, err = c.do(req, &v)
	if err != nil {
		return err
	}
	return nil
}

package ptegodactyl

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
)

// AppClient manage communication with API
type AppClient struct {
	APIToken   string
	BaseURL    *url.URL
	UserAgent  string
	httpClient *http.Client
}

// NewApp returns a ready to use appClient
func NewApp(apiURL string, token string) (c *AppClient) {
	u, err := url.Parse(apiURL)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return &AppClient{
		APIToken:   token,
		BaseURL:    u,
		UserAgent:  "Ptegodactyl",
		httpClient: http.DefaultClient,
	}
}

func (c *AppClient) newRequest(method, path string, body interface{}) (*http.Request, error) {
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

func (c *AppClient) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&v)
	return resp, err
}

// List return a list of things
func (c *AppClient) list(path string, v interface{}) error {
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
func (c *AppClient) get(path string, v interface{}) error {
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
func (c *AppClient) send(path string, body interface{}) error {
	req, err := c.newRequest("POST", path, body)
	if err != nil {
		return err
	}
	_, err = c.do(req, nil)
	if err != nil {
		return err
	}
	return nil
}

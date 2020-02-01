package statuspage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

const STATUSPAGE_HOSTNAME = "https://api.statuspage.io/v1"
const DEFAULT_MAX_RETRIES = 2
const DEFAULT_RETRY_INTERVAL = 10 * time.Second

type Client struct {
	token         string
	HttpClient    *http.Client
	MaxRetries    int
	RetryInterval time.Duration
}

func NewClient(token string) *Client {
	return &Client{
		token:         token,
		HttpClient:    &http.Client{},
		MaxRetries:    DEFAULT_MAX_RETRIES,
		RetryInterval: DEFAULT_RETRY_INTERVAL,
	}
}

func (client *Client) do(method, endpoint string, bodyObject interface{}) (rsp *http.Response, err error) {
	var body io.Reader
	if bodyObject != nil {
		data, err := json.Marshal(bodyObject)
		if err != nil {
			return nil, err
		}

		body = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, STATUSPAGE_HOSTNAME+endpoint, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "OAuth "+client.token)

	rsp, err = client.HttpClient.Do(req)

	var retries int
	for retries = 0; rsp != nil && (rsp.StatusCode == 420 || rsp.StatusCode == 429) && retries <= client.MaxRetries; retries++ {
		time.Sleep(client.RetryInterval)

		rsp, err = client.HttpClient.Do(req)
	}

	if err == nil && rsp.StatusCode > 299 {
		defer rsp.Body.Close()

		body, _ := ioutil.ReadAll(rsp.Body)

		return nil, fmt.Errorf("StatusPage request %s failed(%d): %s", endpoint, rsp.StatusCode, body)
	}

	return rsp, err
}

func (client *Client) doGET(endpoint string, bodyObject interface{}) (*http.Response, error) {
	return client.do("GET", endpoint, bodyObject)
}

func (client *Client) doPATCH(endpoint string, bodyObject interface{}) (err error) {
	_, err = client.do("PATCH", endpoint, bodyObject)

	return err
}

func (client *Client) unmarshal(rsp *http.Response, result interface{}) (err error) {
	defer rsp.Body.Close()

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, result)
}

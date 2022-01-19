package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type PageResponse struct {
	Content  string `json:"content"`
	Children []int  `json:"children"`
}

type LoginResponse struct {
	Token   string    `json:"Token"`
	Expires time.Time `json:"Expires"`
}

type APIClient struct {
	mutex          sync.RWMutex
	client         *http.Client
	token          string
	tokenExpiresAt time.Time
}

func NewApiClient() *APIClient {
	c := APIClient{}
	c.client = NewHttpClient() // This would have to be injected as a parameter for testing
	c.GetToken()
	return &c
}

func (c *APIClient) GetToken() string {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	timeTokenShouldBeRenewed := c.tokenExpiresAt.Add(time.Second * -1 * SECONDS_TO_RENEW_TOKEN)
	if time.Now().Before(timeTokenShouldBeRenewed) {
		// Token is still valid
		return c.token
	}

	url := fmt.Sprintf("%s/api/login", BASE_URL)
	payload := []byte(fmt.Sprintf(`{"username":"%s", "password":"%s"}`, API_USERNAME, API_PASSWORD))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		panic(err)
	}
	res, err := c.client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	var result LoginResponse
	if err := json.Unmarshal(body, &result); err != nil {
		panic(err)
	}

	c.token = result.Token
	c.tokenExpiresAt = result.Expires

	return c.token
}

func (c *APIClient) GetPage(id int) (*PageResponse, error) {
	token := c.GetToken()
	url := fmt.Sprintf("%s/api/page?id=%d", BASE_URL, id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Status is %d", res.Status))
	}
	var result PageResponse
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &result, nil
}

func NewHttpClient() *http.Client {
	// Ref: https://www.loginradius.com/blog/async/tune-the-go-http-client-for-high-performance/
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxConnsPerHost = MAX_CONNECTIONS_PER_HOST
	t.MaxIdleConnsPerHost = MAX_CONNECTIONS

	return &http.Client{
		Timeout:   ACTIVE_CONNECTION_TIMEOUT * time.Second,
		Transport: t,
	}
}

package networkService

import (
	"net/http"
	"time"
)

// Notably missing is a thread pool or something equivalent to prevent the number of concurrent goroutines
// from increasing exponentially. In production I would consider having a thread executor + workers pattern if
// the size of the data is going to mean that there might be a lot of goroutines (e.g. 100k) running
// at the same time.
//
// The number of connections is limited by the http client instance, so that's not an issue in
// this case.

// AuthMonitor keeps the login and renews it on demand.
func InitAuthMonitor() {
	// The http client is passed as an argument so that it is easy to test this module by passing a test instance
	// E.g. https://www.thegreatcodeadventure.com/mocking-http-requests-in-golang/

	if isRunning {
		return
	}
	login := authenticate()
	for {
		select {
		case <-getTokenRequest:
			getTokenResponse <- login
		case <-renewTokenRequest:
			// In production I would implement some way to prevent multiple consecutive
			// requests triggering a renew, e.g. keep a timestamp of last fetch or a flag
			login = getLogin()
			getTokenResponse <- login
		}
	}
}

func GetToken() {
	getTokenRequest <- true
	result := <-getTokenResponse
	return result
}

func RenewToken() {
	renewTokenRequest <- true
	result := <-getTokenResponse
	return result
}

// Private

// authenticate logs in and fetches the token
func authenticate() {

}

func newHttpClient() {
	// Ref: https://www.loginradius.com/blog/async/tune-the-go-http-client-for-high-performance/
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100

	return &http.Client{
		Timeout:   10 * time.Second,
		Transport: t,
	}
}

// Channels
var getTokenRequest = make(chan bool)
var getTokenResponse = make(chan string)
var renewTokenRequest = make(chan bool)

var isRunning = false

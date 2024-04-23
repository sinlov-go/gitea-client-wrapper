package gitea_token_client

import (
	"code.gitea.io/sdk/gitea"
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"sync"
	"time"
)

// NewClientWithHttpTimeout creates a new Gitea client with a timeout
// for the http client
// The timeout is in seconds, and the minimum is 30 seconds
// If insecure is true, the client will skip SSL verification
// This function is a wrapper around gitea.NewClient
func (g *GiteaTokenClient) NewClientWithHttpTimeout(url, accessToken string, timeoutSecond uint, insecure bool) error {
	rwMutex := &sync.RWMutex{}
	rwMutex.Lock()

	if timeoutSecond < 30 {
		timeoutSecond = 30
	}

	httpClient := &http.Client{
		Timeout: time.Duration(timeoutSecond) * time.Second,
	}
	if insecure {
		cookieJar, _ := cookiejar.New(nil)
		httpClient = &http.Client{
			Jar: cookieJar,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				DialContext: (&net.Dialer{
					Timeout:   time.Duration(timeoutSecond*3) * time.Second,
					KeepAlive: time.Duration(timeoutSecond*3) * time.Second,
				}).DialContext,
				TLSHandshakeTimeout:   time.Duration(timeoutSecond) * time.Second,
				ResponseHeaderTimeout: time.Duration(timeoutSecond) * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			},
			Timeout: time.Duration(timeoutSecond) * time.Second,
		}
	}
	baseUrl := strings.TrimSuffix(url, "/")
	client, errNewClient := gitea.NewClient(baseUrl,
		gitea.SetToken(accessToken),
		gitea.SetHTTPClient(httpClient),
	)
	rwMutex.Unlock()
	if errNewClient != nil {
		return fmt.Errorf("failed to create gitea client with NewClientWithHttpTimeout: %s", errNewClient)
	}
	g.client = client
	g.debug = false
	g.baseUrl = baseUrl
	g.ctx = context.Background()
	g.mutex = rwMutex
	g.httpClient = httpClient
	g.accessToken = accessToken

	return nil
}

// NewClient creates a new Gitea client
// This function is a wrapper around gitea.NewClient
func (g *GiteaTokenClient) NewClient(url, accessToken string, httpClient *http.Client) error {
	rwMutex := &sync.RWMutex{}
	rwMutex.Lock()
	baseUrl := strings.TrimSuffix(url, "/")
	client, errNewClient := gitea.NewClient(baseUrl,
		gitea.SetToken(accessToken),
		gitea.SetHTTPClient(httpClient),
	)
	rwMutex.Unlock()
	if errNewClient != nil {
		return fmt.Errorf("failed to create gitea client: %s", errNewClient)
	}
	g.client = client
	g.debug = false
	g.baseUrl = baseUrl
	g.ctx = context.Background()
	g.mutex = rwMutex
	g.httpClient = httpClient
	g.accessToken = accessToken

	return nil
}

// SetDebug sets the debug status for the client from gitea.Client
func (g *GiteaTokenClient) SetDebug(debug bool) {
	g.mutex.Lock()
	g.debug = debug
	g.mutex.Unlock()
}

// IsDebug returns the debug status
func (g *GiteaTokenClient) IsDebug() bool {
	return g.debug
}

func (g *GiteaTokenClient) GetContext() context.Context {
	return g.ctx
}

// GiteaClient returns the gitea.Client
func (g *GiteaTokenClient) GiteaClient() *gitea.Client {
	return g.client
}

// GetBaseUrl returns the base URL of the Gitea instance
func (g *GiteaTokenClient) GetBaseUrl() string {
	return g.baseUrl
}

// GetUsername returns the username by SetBasicAuth
func (g *GiteaTokenClient) GetUsername() string {
	return g.username
}

// SetOTP sets the otp for the client from gitea.Client
func (g *GiteaTokenClient) SetOTP(otp string) {
	g.mutex.Lock()
	g.otp = otp
	g.client.SetOTP(otp)
	g.mutex.Unlock()
}

// SetSudo sets the sudo for the client from gitea.Client
func (g *GiteaTokenClient) SetSudo(sudo string) {
	g.mutex.Lock()
	g.sudo = sudo
	g.client.SetSudo(sudo)
	g.mutex.Unlock()
}

// SetBasicAuth sets the basic auth for the client from gitea.Client
func (g *GiteaTokenClient) SetBasicAuth(username, password string) {
	g.mutex.Lock()
	g.username, g.password = username, password
	g.client.SetBasicAuth(username, password)
	g.mutex.Unlock()
}

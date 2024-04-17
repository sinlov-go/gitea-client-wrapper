package gitea_token_client

import (
	"code.gitea.io/sdk/gitea"
	"context"
	"net/http"
	"sync"
)

type GiteaTokenClient struct {
	GiteaTokenClientFunc GiteaTokenClientFunc `json:"-"`

	client *gitea.Client
	debug  bool

	baseUrl    string
	ctx        context.Context
	mutex      *sync.RWMutex
	httpClient *http.Client

	accessToken string
	username    string
	password    string
	otp         string
	sudo        string
}

type GiteaTokenClientFunc interface {
	NewClientWithHttpTimeout(url, accessToken string, timeoutSecond uint, insecure bool) error

	NewClient(url, accessToken string, httpClient *http.Client) error

	SetDebug(debug bool)

	IsDebug() bool

	SetOTP(otp string)

	SetSudo(sudo string)

	SetBasicAuth(username, password string)

	GiteaClient() *gitea.Client

	GetBaseUrl() string

	GetUsername() string
}

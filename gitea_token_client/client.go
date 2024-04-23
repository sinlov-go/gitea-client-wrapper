package gitea_token_client

import (
	"code.gitea.io/sdk/gitea"
	"context"
	"io"
	"net/http"
	"sync"
)

type GiteaTokenClient struct {
	GiteaTokenClientFunc GiteaTokenClientFunc `json:"-"`
	GiteaApiFunc         GiteaApiFunc         `json:"-"`

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

	GetContext() context.Context

	SetOTP(otp string)

	SetSudo(sudo string)

	SetBasicAuth(username, password string)

	GiteaClient() *gitea.Client

	GetBaseUrl() string

	GetUsername() string
}

type GiteaApiFunc interface {
	ApiGiteaGet(httpPath string, header http.Header, response interface{}) (*GiteaApiResponse, error)

	ApiGiteaPostJson(httpPath string, header http.Header, body interface{}, response interface{}) (*GiteaApiResponse, error)

	ApiGiteaPost(httpPath string, header http.Header, body io.Reader, response interface{}) (*GiteaApiResponse, error)

	ApiGiteaPatchJson(httpPath string, header http.Header, body interface{}, response interface{}) (*GiteaApiResponse, error)

	ApiGiteaPatch(httpPath string, header http.Header, body io.Reader, response interface{}) (*GiteaApiResponse, error)

	ApiGiteaPutJson(httpPath string, header http.Header, body io.Reader, response interface{}) (*GiteaApiResponse, error)

	ApiGiteaPut(httpPath string, header http.Header, body io.Reader, response interface{}) (*GiteaApiResponse, error)

	ApiGiteaDelete(httpPath string, header http.Header, response interface{}) (*GiteaApiResponse, error)
}

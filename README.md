[![ci](https://github.com/sinlov-go/gitea-client-wrapper/actions/workflows/ci.yml/badge.svg)](https://github.com/sinlov-go/gitea-client-wrapper/actions/workflows/ci.yml)

[![go mod version](https://img.shields.io/github/go-mod/go-version/sinlov-go/gitea-client-wrapper?label=go.mod)](https://github.com/sinlov-go/gitea-client-wrapper)
[![GoDoc](https://godoc.org/github.com/sinlov-go/gitea-client-wrapper?status.png)](https://godoc.org/github.com/sinlov-go/gitea-client-wrapper)
[![goreportcard](https://goreportcard.com/badge/github.com/sinlov-go/gitea-client-wrapper)](https://goreportcard.com/report/github.com/sinlov-go/gitea-client-wrapper)

[![GitHub license](https://img.shields.io/github/license/sinlov-go/gitea-client-wrapper)](https://github.com/sinlov-go/gitea-client-wrapper)
[![codecov](https://codecov.io/gh/sinlov-go/gitea-client-wrapper/branch/main/graph/badge.svg)](https://codecov.io/gh/sinlov-go/gitea-client-wrapper)
[![GitHub latest SemVer tag)](https://img.shields.io/github/v/tag/sinlov-go/gitea-client-wrapper)](https://github.com/sinlov-go/gitea-client-wrapper/tags)
[![GitHub release)](https://img.shields.io/github/v/release/sinlov-go/gitea-client-wrapper)](https://github.com/sinlov-go/gitea-client-wrapper/releases)

## for what

- this project used to github golang lib project

## Contributing

[![Contributor Covenant](https://img.shields.io/badge/contributor%20covenant-v1.4-ff69b4.svg)](.github/CONTRIBUTING_DOC/CODE_OF_CONDUCT.md)
[![GitHub contributors](https://img.shields.io/github/contributors/sinlov-go/gitea-client-wrapper)](https://github.com/sinlov-go/gitea-client-wrapper/graphs/contributors)

We welcome community contributions to this project.

Please read [Contributor Guide](.github/CONTRIBUTING_DOC/CONTRIBUTING.md) for more information on how to get started.

请阅读有关 [贡献者指南](.github/CONTRIBUTING_DOC/zh-CN/CONTRIBUTING.md) 以获取更多如何入门的信息

## depends

in go mod project

```bash
# warning use private git host must set
# global set for once
# add private git host like github.com to evn GOPRIVATE
$ go env -w GOPRIVATE='github.com'
# use ssh proxy
# set ssh-key to use ssh as http
$ git config --global url."git@github.com:".insteadOf "https://github.com/"
# or use PRIVATE-TOKEN
# set PRIVATE-TOKEN as gitlab or gitea
$ git config --global http.extraheader "PRIVATE-TOKEN: {PRIVATE-TOKEN}"
# set this rep to download ssh as https use PRIVATE-TOKEN
$ git config --global url."ssh://github.com/".insteadOf "https://github.com/"

# before above global settings
# test version info
$ git ls-remote -q https://github.com/sinlov-go/gitea-client-wrapper.git

# test depends see full version
$ go list -mod readonly -v -m -versions github.com/sinlov-go/gitea-client-wrapper
# or use last version add go.mod by script
$ echo "go mod edit -require=$(go list -mod=readonly -m -versions github.com/sinlov-go/gitea-client-wrapper | awk '{print $1 "@" $NF}')"
$ echo "go mod vendor"
```

## Features

- [x] `gitea_token_client.GiteaTokenClient` to let api client easy to use
- [ ] more perfect test case coverage
- [ ] more perfect benchmark case

## env

- minimum go version: go 1.19
- change `go 1.19`, `^1.19`, `1.19.12-bullseye`, `1.19.12` to new go version

### libs

| lib                                 | version |
|:------------------------------------|:--------|
| https://github.com/stretchr/testify | v1.8.4  |
| https://github.com/sebdah/goldie    | v2.5.3  |

- more libs see [go.mod](https://github.com/sinlov-go/gitea-client-wrapper/blob/main/go.mod)

## usage

### gitea_token_client.GiteaTokenClient

```go
package foo

import (
  "github.com/sinlov-go/gitea-client-wrapper/gitea_token_client"
)
type GiteaApiClient struct {
	gitea_token_client.GiteaTokenClient
}

func NewGiteaApiClient(url, accessToken string, timeoutSecond uint, insecure bool) (*GiteaApiClient, error) {
	rc := GiteaApiClient{}
	err := rc.NewClientWithHttpTimeout(url, accessToken, timeoutSecond, insecure)
	if err != nil {
		return nil, err
	}
	return &rc, nil
}
```

# dev

```bash
# It needs to be executed after the first use or update of dependencies.
$ make init dep
```

- test code

```bash
$ make test testBenchmark
```

add main.go file and run

```bash
# run at env dev use cmd/main.go
$ make dev
```

- ci to fast check

```bash
# check style at local
$ make style

# run ci at local
$ make ci
```

## docker

```bash
# then test build as test/Dockerfile
$ make dockerTestRestartLatest
# clean test build
$ make dockerTestPruneLatest

# more info see
$ make helpDocker
```

package gitea_token_client

import (
	"fmt"
	"net/url"
)

const (
	BaseGiteaApi = "/api/v1"
)

var giteaBaseApi = BaseGiteaApi

func SetBaseApi(api string) {
	giteaBaseApi = api
}

func GetBaseApi() string {
	return giteaBaseApi
}

// GiteaApiParsef is a help function to parse api path
func GiteaApiParsef(format string, a ...any) string {
	return fmt.Sprintf("%s%s", giteaBaseApi, fmt.Sprintf(format, a...))
}

// EscapeValidatePathSegments is a help function to validate and encode url path segments
// use as
//
//	pkgName := "foo.one.two"
//	errEsCapePkgName := EscapeValidatePathSegments(&pkgName)
func EscapeValidatePathSegments(seg ...*string) error {
	for i := range seg {
		if seg[i] == nil || len(*seg[i]) == 0 {
			return fmt.Errorf("path segment [%d] is empty", i)
		}
		*seg[i] = url.PathEscape(*seg[i])
	}
	return nil
}

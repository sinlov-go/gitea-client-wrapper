package gitea_token_client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// GiteaApiResponse represents the gitea response
type GiteaApiResponse struct {
	*http.Response

	FirstPage int
	PrevPage  int
	NextPage  int
	LastPage  int
}

// ApiGiteaStatusCode sends a request to the Gitea API and returns the status code
func (g *GiteaTokenClient) ApiGiteaStatusCode(method, path string, header http.Header, body io.Reader) (int, error) {
	resp, err := g.doApiRequest(method, path, header, body)
	if err != nil {
		return -1, err
	}
	return resp.StatusCode, nil
}

// ApiGiteaGet sends a GET request to the Gitea API and returns the response
func (g *GiteaTokenClient) ApiGiteaGet(httpPath string, header http.Header, response interface{}) (*GiteaApiResponse, error) {
	return g.getApiParsedResponse(http.MethodGet, httpPath, header, nil, response)
}

// ApiGiteaPostJson sends a POST request to the Gitea API with a JSON body and returns the response
func (g *GiteaTokenClient) ApiGiteaPostJson(httpPath string, header http.Header, body interface{}, response interface{}) (*GiteaApiResponse, error) {
	bodyData, errJson := json.Marshal(body)
	if errJson != nil {
		return nil, errJson
	}
	return g.ApiGiteaPost(httpPath, header, bytes.NewBuffer(bodyData), response)
}

// ApiGiteaPost sends a POST request to the Gitea API and returns the response
func (g *GiteaTokenClient) ApiGiteaPost(httpPath string, header http.Header, body io.Reader, response interface{}) (*GiteaApiResponse, error) {
	return g.getApiParsedResponse(http.MethodPost, httpPath, header, body, response)
}

// ApiGiteaPatchJson sends a PATCH request to the Gitea API with a JSON body and returns the response
func (g *GiteaTokenClient) ApiGiteaPatchJson(httpPath string, header http.Header, body interface{}, response interface{}) (*GiteaApiResponse, error) {
	bodyData, errJson := json.Marshal(body)
	if errJson != nil {
		return nil, errJson
	}
	return g.ApiGiteaPatch(httpPath, header, bytes.NewBuffer(bodyData), response)
}

// ApiGiteaPatch sends a PATCH request to the Gitea API and returns the response
func (g *GiteaTokenClient) ApiGiteaPatch(httpPath string, header http.Header, body io.Reader, response interface{}) (*GiteaApiResponse, error) {
	return g.getApiParsedResponse(http.MethodPatch, httpPath, header, body, response)
}

// ApiGiteaPutJson sends a PUT request to the Gitea API with a JSON body and returns the response
func (g *GiteaTokenClient) ApiGiteaPutJson(httpPath string, header http.Header, body io.Reader, response interface{}) (*GiteaApiResponse, error) {
	bodyData, errJson := json.Marshal(body)
	if errJson != nil {
		return nil, errJson
	}
	return g.ApiGiteaPut(httpPath, header, bytes.NewBuffer(bodyData), response)
}

// ApiGiteaPut sends a PUT request to the Gitea API and returns the response
func (g *GiteaTokenClient) ApiGiteaPut(httpPath string, header http.Header, body io.Reader, response interface{}) (*GiteaApiResponse, error) {
	return g.getApiParsedResponse(http.MethodPut, httpPath, header, body, response)
}

// ApiGiteaDelete sends a DELETE request to the Gitea API and returns the response
func (g *GiteaTokenClient) ApiGiteaDelete(httpPath string, header http.Header, response interface{}) (*GiteaApiResponse, error) {
	return g.getApiParsedResponse(http.MethodDelete, httpPath, header, nil, response)
}

func (g *GiteaTokenClient) getApiParsedResponse(method, path string, header http.Header, body io.Reader, obj interface{}) (*GiteaApiResponse, error) {
	data, resp, err := g.getApiResponse(method, path, header, body)
	if err != nil {
		if g.debug {
			if resp != nil {
				fmt.Printf("getApiParsedResponse code %d err: %s\n", resp.StatusCode, err)
			} else {
				fmt.Printf("getApiParsedResponse err: %s\n", err)
			}
		}
		return resp, err
	}
	return resp, json.Unmarshal(data, obj)
}

func (g *GiteaTokenClient) getApiResponse(method, path string, header http.Header, body io.Reader) ([]byte, *GiteaApiResponse, error) {
	resp, err := g.doApiRequest(method, path, header, body)
	if err != nil {
		return nil, nil, err
	}
	defer func(Body io.ReadCloser) {
		errReadClose := Body.Close()
		if errReadClose != nil {
			fmt.Printf("Error getApiResponse closing response body: %v\n", errReadClose)
		}
	}(resp.Body)

	// check for errors
	data, err := statusCodeToErr(resp)
	if err != nil {
		return data, resp, err
	}
	// success (2XX), read body
	data, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp, err
	}
	return data, resp, nil
}

func (g *GiteaTokenClient) doApiRequest(method, path string, header http.Header, body io.Reader) (*GiteaApiResponse, error) {
	if g.client == nil {
		return nil, fmt.Errorf("gitea client is nil")
	}
	g.mutex.Lock()
	debug := g.debug
	urlFullPath := g.baseUrl + path
	if debug {
		fmt.Printf("%s: %s\nHeader: %v\nBody: %s\n", method, urlFullPath, header, body)
	}
	req, err := http.NewRequestWithContext(g.ctx, method, urlFullPath, body)
	if err != nil {
		g.mutex.RUnlock()
		return nil, err
	}

	if len(g.accessToken) != 0 {
		req.Header.Set("Authorization", "token "+g.accessToken)
	}
	if len(g.otp) != 0 {
		req.Header.Set("X-GITEA-OTP", g.otp)
	}
	if len(g.username) != 0 {
		req.SetBasicAuth(g.username, g.password)
	}
	if len(g.sudo) != 0 {
		req.Header.Set("Sudo", g.sudo)
	}

	for k, v := range header {
		req.Header[k] = v
	}

	g.mutex.Unlock()
	httpClient := g.httpClient
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if debug {
		fmt.Printf("Response: %v\n\n", resp)
	}

	return newResponse(resp), nil
}

func newResponse(r *http.Response) *GiteaApiResponse {
	response := &GiteaApiResponse{Response: r}
	response.parseLinkHeader()
	return response
}

func (r *GiteaApiResponse) parseLinkHeader() {
	link := r.Header.Get("Link")
	if link == "" {
		return
	}

	links := strings.Split(link, ",")
	for _, l := range links {
		u, param, ok := strings.Cut(l, ";")
		if !ok {
			continue
		}
		u = strings.Trim(u, " <>")

		key, value, ok := strings.Cut(strings.TrimSpace(param), "=")
		if !ok || key != "rel" {
			continue
		}

		value = strings.Trim(value, "\"")

		parsed, err := url.Parse(u)
		if err != nil {
			continue
		}

		page := parsed.Query().Get("page")
		if page == "" {
			continue
		}

		switch value {
		case "first":
			r.FirstPage, _ = strconv.Atoi(page)
		case "prev":
			r.PrevPage, _ = strconv.Atoi(page)
		case "next":
			r.NextPage, _ = strconv.Atoi(page)
		case "last":
			r.LastPage, _ = strconv.Atoi(page)
		}
	}
}

// Converts a response for a HTTP status code indicating an error condition
// (non-2XX) to a well-known error value and response body. For non-problematic
// (2XX) status codes nil will be returned. Note that on a non-2XX response, the
// response body stream will have been read and, hence, is closed on return.
func statusCodeToErr(resp *GiteaApiResponse) (body []byte, err error) {
	// no error
	if resp.StatusCode/100 == 2 {
		return nil, nil
	}

	//
	// error: body will be read for details
	//
	defer func(Body io.ReadCloser) {
		errReadClose := Body.Close()
		if errReadClose != nil {
			fmt.Printf("Error statusCodeToErr closing response body: %v\n", errReadClose)
		}
	}(resp.Body)
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("body read on HTTP error %d: %v", resp.StatusCode, err)
	}

	switch resp.StatusCode {
	case http.StatusForbidden:
		return data, errors.New("403 Forbidden")
	case http.StatusNotFound:
		return data, errors.New("404 Not Found")
	case http.StatusConflict:
		return data, errors.New("409 Conflict")
	case http.StatusUnprocessableEntity:
		return data, fmt.Errorf("422 Unprocessable Entity: %s", string(data))
	}

	urlPath := resp.Request.URL.Path
	method := resp.Request.Method
	header := resp.Request.Header
	errMap := make(map[string]interface{})
	if err = json.Unmarshal(data, &errMap); err != nil {
		// when the JSON can't be parsed, data was probably empty or a
		// plain string, so we try to return a helpful error anyway
		return data, fmt.Errorf("Unknown API Error: %d\nRequest: '%s' with '%s' method '%s' header and '%s' body", resp.StatusCode, urlPath, method, header, string(data))
	}
	return data, errors.New(errMap["message"].(string))
}

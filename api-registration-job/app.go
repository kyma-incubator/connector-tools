package main

import (
	"net/http"
)

type endpointInfo struct {
	Path        string
	Name        string
	Description string
}

type app interface {
	setCredentials(request *http.Request) *http.Request
	getAPIUrl(systemURL string, path string) string
	verifyActiveResponse(resp *http.Response) (bool, error)
	readEndpoints(apis []API, r registrationApp)
}

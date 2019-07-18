package main

import (
	"fmt"
	"net/http"
	"strings"
)

const format = "json"

type restWithAPIKey struct {
	apikey string
	source string
}

func (l *restWithAPIKey) generateMetadata(endpoint endpointInfo, r registrationApp) []byte {

	metadata := fmt.Sprintf(`
			{
				"provider" : "%s",
				"name": "%s - %s",
				"description":"%s",
				"api": {
					"targetUrl": "%s",
					"queryParameters": {
						"format": ["%s"],
						"source": ["%s"]
					},
					"headers": {
						"apikey": ["%s"]
					}
				}
			}
	`, r.ProviderName, r.ProductName, endpoint.Name, endpoint.Description, r.SystemURL, format, l.source, l.apikey)
	return []byte(metadata)
}

func (l *restWithAPIKey) setCredentials(request *http.Request) *http.Request {
	request.Header.Set("apikey", l.apikey)
	request.URL.Query().Set("format", format)
	request.URL.Query().Set("source", l.source)
	return request
}

func (l *restWithAPIKey) getAPIUrl(systemURL string, path string) string {
	if strings.HasSuffix(systemURL, "/") && strings.HasPrefix(path, "/") {
		runes := []rune(path)
		return systemURL + string(runes[1:])
	} else if strings.HasSuffix(systemURL, "/") || strings.HasPrefix(path, "/") {
		return systemURL + path
	} else {
		return systemURL + "/" + path
	}
}

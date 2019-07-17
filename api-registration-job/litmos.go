package main

import (
	"fmt"
	"net/http"
	"net/url"
)

const format = "json"

type litmos struct {
	apikey string
	source string
}

func (l *litmos) generateMetadata(endpoint endpointInfo, r registrationApp) []byte {
	specificationsURL, err := url.Parse(r.SystemURL)
	check(err)
	specificationsURL.Path = endpoint.Path + "/$metadata"

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

func (l *litmos) setCredentials(request *http.Request) *http.Request {
	request.Header.Set("apikey", l.apikey)
	request.URL.Query().Set("format", format)
	request.URL.Query().Set("source", l.source)
	return request
}

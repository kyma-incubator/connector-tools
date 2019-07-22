package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const format = "json"

type restWithAPIKey struct {
	apikey      string
	source      string
	description string
}

func (l *restWithAPIKey) generateMetadata(r registrationApp) []byte {
	apiName := l.apiName(r)
	metadata := fmt.Sprintf(`
			{
				"provider" : "%s",
				"name": "%s",
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
	`, r.ProviderName, apiName, l.description, r.SystemURL, format, l.source, l.apikey)
	return []byte(metadata)
}

func (l *restWithAPIKey) setCredentials(request *http.Request) *http.Request {
	request.Header.Set("apikey", l.apikey)
	q := request.URL.Query()
	q.Add("format", format)
	q.Add("source", l.source)

	request.URL.RawQuery = q.Encode()

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

func (l *restWithAPIKey) verifyActiveResponse(resp *http.Response) (bool, error) {
	jsonResponse := make([]map[string]interface{}, 0)
	err := json.NewDecoder(resp.Body).Decode(&jsonResponse)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (l *restWithAPIKey) readEndpoints(apis []API, r registrationApp) () {
	apiName := l.apiName(r)
	contains, id := containsAPI(apis, apiName)
	if contains {
		fmt.Printf("API %s is already registered at kyma application\n", apiName)
		err := r.updateSingleAPI(id, l.generateMetadata(r))

		if err != nil {
			fmt.Printf("error while updating API %s\n", err)
		}
	} else {
		fmt.Printf("API %s is not registered yet at kyma application\n", apiName)
		err := r.registerSingleAPI(l.generateMetadata(r))

		if err != nil {
			fmt.Printf("error while registering API %s\n", err)
		}
	}
}

func (l *restWithAPIKey) apiName(r registrationApp) string {
	return fmt.Sprintf("%s-API", r.ProductName)
}

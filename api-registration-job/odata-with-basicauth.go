package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type oDataWithBasicAuth struct {
	BasicUser     string
	BasicPassword string
}

func (a *oDataWithBasicAuth) generateMetadata(endpoint endpointInfo, r registrationApp) []byte {
	specificationsURL, err := url.Parse(r.SystemURL)
	check(err)
	specificationsURL.User = url.UserPassword(a.BasicUser, a.BasicPassword)
	specificationsURL.Path = endpoint.Path + "/$metadata"

	tokenEndpointURL := fmt.Sprintf("%s/%s/", r.SystemURL, endpoint.Path)
	metadata := fmt.Sprintf(`
			{
				"provider" : "%s",
				"name": "%s - %s",
				"description":"%s",
				"api": {
					"targetUrl": "%s",
					"SpecificationUrl":"%s",
					"ApiType": "oDataWithBasicAuth",
					"credentials": {
						"basic": {
							"username":"%s",
							"password":"%s",
							"csrfInfo":{
								"tokenEndpointURL":"%s"
							}
						}
					}
				}
			}
	`, r.ProviderName, r.ProductName, endpoint.Name, endpoint.Description, r.SystemURL, specificationsURL.String(), a.BasicUser, a.BasicPassword, tokenEndpointURL)
	return []byte(metadata)
}

func (a *oDataWithBasicAuth) setCredentials(request *http.Request) *http.Request {
	request.SetBasicAuth(a.BasicUser, a.BasicPassword)
	return request
}

func (a *oDataWithBasicAuth) getAPIUrl(systemURL string, path string) string {
	return systemURL + "/" + path + "/"
}

func (a *oDataWithBasicAuth) verifyActiveResponse(resp *http.Response) (bool, error) {
	jsonResponse := make(map[string]map[string][]string)
	err := json.NewDecoder(resp.Body).Decode(&jsonResponse)
	if err != nil {
		return false, err
	}
	if len(jsonResponse["d"]["EntitySets"]) > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

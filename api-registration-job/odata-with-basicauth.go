package main

import (
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

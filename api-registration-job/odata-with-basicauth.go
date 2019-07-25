package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
					"ApiType": "OData",
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

func (a *oDataWithBasicAuth) readEndpoints(apis []API, r registrationApp) {
	configString, err := ioutil.ReadFile("files/apis.json")
	if err != nil {
		fmt.Println("new_config.json not found... Moving on.")
		return
	}

	fmt.Println("Registering new APIs")
	var endpoints []endpointInfo
	err = json.Unmarshal(configString, &endpoints)
	check(err)
	var errors = ""
	for _, e := range endpoints {
		fmt.Printf("Processing API %s\n", e.Name)
		active, err := r.isAPIActive(e.Path)
		if err != nil {
			errors = errors + err.Error() + "\n"
			fmt.Println(err)
			continue
		}
		if active {
			fmt.Printf("API %s is enabled in remote system\n", e.Name)
			contains, id := containsAPI(apis, fmt.Sprintf("%s - %s", r.ProductName, e.Name))
			if contains {
				fmt.Printf("API %s is already registered at kyma application\n", e.Name)
				err = r.updateSingleAPI(id, a.generateMetadata(e, r))
				if err != nil {
					errors = errors + err.Error() + "\n"
					fmt.Printf("Error while update: %s", err)
					continue
				}
			} else {
				fmt.Printf("API %s is not registered yet at kyma application\n", e.Name)
				err = r.registerSingleAPI(a.generateMetadata(e, r))
				if err != nil {
					errors = errors + err.Error() + "\n"
					fmt.Printf("Error while registration: %s", err)
					continue
				}
			}
		} else {
			fmt.Printf("Skipping API %s as it is not enabled in remote system\n", e.Name)
		}
	}
	if errors != "" {
		panic(fmt.Errorf("There were errors while API registration:\n%s", errors))
	}
}

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

type registrationApp struct {
	hostname     string
	authUsername string
	authPassword string
	apiURL       string
	provider     string
	product      string
}

type endpointInfo struct {
	API          string
	Name         string
	Description  string
	HelpDoc      string
	Scenario     string
	ScenarioName string
}

func main() {
	applicationName := os.Getenv("APPLICATION_NAME")
	apiURL := fmt.Sprintf("http://application-registry-external-api.kyma-integration.svc.cluster.local:8081/%s/v1/metadata/services", applicationName)
	r := registrationApp{
		hostname:     os.Getenv("SOURCE_ID"),
		authUsername: os.Getenv("auth_username"),
		authPassword: os.Getenv("auth_password"),
		apiURL:       apiURL,
		provider:     os.Getenv("provider"),
		product:      os.Getenv("product"),
	}
	r.registerStaticEvents()
	r.readEndpoints()
}

func (r registrationApp) registerStaticEvents() {
	eventsString, err := ioutil.ReadFile("files/events.json")
	if err != nil {
		fmt.Println("events.json not found... Moving on.")
		return
	}
	req, err := http.NewRequest("POST", r.apiURL, bytes.NewBuffer(eventsString))
	check(err)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println("response Status:", resp.StatusCode)

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func (r registrationApp) readEndpoints() {
	configString, err := ioutil.ReadFile("files/new_config.json")
	if err != nil {
		fmt.Println("new_config.json not found... Moving on.")
		return
	}
	var endpoints []endpointInfo
	err = json.Unmarshal(configString, &endpoints)
	check(err)
	for _, e := range endpoints {
		fmt.Println(e.API)
		if r.isAPIActive(e.API) {
			r.registerSingleAPI(r.generateMetadata(e))
		}
	}

}

func (r registrationApp) isAPIActive(api string) bool {
	url := r.hostname + "/" + api
	req, err := http.NewRequest("GET", url, nil)
	check(err)
	req.SetBasicAuth(r.authUsername, r.authPassword)
	req.Header.Set("accept", "application/json")
	req.Header.Set("responseType", "json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	jsonResponse := make(map[string]map[string][]string)
	if resp.StatusCode == 200 {
		json.NewDecoder(resp.Body).Decode(&jsonResponse)

		if len(jsonResponse["d"]["EntitySets"]) > 0 {
			return true
		}
	}
	return false
}

func (r registrationApp) registerSingleAPI(apiMetadata []byte) {
	req, err := http.NewRequest("POST", r.apiURL, bytes.NewBuffer(apiMetadata))
	check(err)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

func (r registrationApp) generateMetadata(endpoint endpointInfo) []byte {
	name := endpoint.Name
	targetURL := r.hostname
	specificationsURL, err := url.Parse(r.hostname)
	check(err)
	specificationsURL.User = url.UserPassword(r.authUsername, r.authPassword)
	specificationsURL.Path = endpoint.API + "/$metadata"

	tokenEndpointURL := r.hostname + endpoint.API
	metadata := fmt.Sprintf(`
			{
				"provider" : "%s",
				"name": "%s - %s",
				"identifier":"%s",
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
	`, r.provider, r.product, name, endpoint.API, endpoint.Description, targetURL, specificationsURL.String(), r.authUsername, r.authPassword, tokenEndpointURL)

	fmt.Println(metadata)
	return []byte(metadata)

}

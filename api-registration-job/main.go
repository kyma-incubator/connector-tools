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
	ApplicationName string
	SystemURL       string
	BasicUser       string
	BasicPassword   string
	RegistrationURL string
	ProviderName    string
	ProductName     string
}

type endpointInfo struct {
	Path        string
	Name        string
	Description string
}

func main() {
	fmt.Println("Started registration job")
	r := registrationApp{
		ApplicationName: os.Getenv("APPLICATION_NAME"),
		SystemURL:       os.Getenv("SYSTEM_URL"),
		BasicUser:       os.Getenv("BASIC_USER"),
		BasicPassword:   os.Getenv("BASIC_PASSWORD"),
		RegistrationURL: os.Getenv("REGISTRATION_URL"),
		ProviderName:    os.Getenv("PROVIDER_NAME"),
		ProductName:     os.Getenv("PRODUCT_NAME"),
	}

	if r.RegistrationURL == "" {
		r.RegistrationURL = fmt.Sprintf("http://application-registry-external-api.kyma-integration.svc.cluster.local:8081/%s/v1/metadata/services", r.ApplicationName)
	}
	r.validateSystemURL()
	r.registerStaticEvents()
	r.readEndpoints()
	fmt.Println("Finished registration job")
}

func (r registrationApp) validateSystemURL() {
	u, err := url.Parse(r.SystemURL)
	check(err)
	if u.Scheme == "" {
		u.Scheme = "https"
	}
	r.SystemURL = u.String()
}

func (r registrationApp) registerStaticEvents() {
	eventsString, err := ioutil.ReadFile("files/events.json")
	if err != nil {
		fmt.Println("events.json not found... Moving on.")
		return
	}
	fmt.Println("Registering events")
	req, err := http.NewRequest("POST", r.RegistrationURL, bytes.NewBuffer(eventsString))
	check(err)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	check(err)
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		fmt.Println("Events registered with success")
	} else {
		check(fmt.Errorf("Registration of events failed with status code %d", resp.StatusCode))
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func (r registrationApp) readEndpoints() {
	configString, err := ioutil.ReadFile("files/apis.json")
	if err != nil {
		fmt.Println("new_config.json not found... Moving on.")
		return
	}
	fmt.Println("Registering APIs")
	var endpoints []endpointInfo
	err = json.Unmarshal(configString, &endpoints)
	check(err)
	for _, e := range endpoints {
		fmt.Printf("Processing API %s\n", e.Name)
		if r.isAPIActive(e.Path) {
			fmt.Printf("API %s is enabled, continue with registration\n", e.Name)
			r.registerSingleAPI(r.generateMetadata(e))
		} else {
			fmt.Printf("Skipping API %s as it is not enabled\n", e.Name)
		}
	}
}

func (r registrationApp) isAPIActive(path string) bool {
	url := r.SystemURL + "/" + path
	req, err := http.NewRequest("GET", url, nil)
	check(err)
	req.SetBasicAuth(r.BasicUser, r.BasicPassword)
	req.Header.Set("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	jsonResponse := make(map[string]map[string][]string)
	if resp.StatusCode == 200 {
		err := json.NewDecoder(resp.Body).Decode(&jsonResponse)
		check(err)
		if len(jsonResponse["d"]["EntitySets"]) > 0 {
			return true
		}
	}
	return false
}

func (r registrationApp) registerSingleAPI(apiMetadata []byte) {
	fmt.Printf("Registering API with payload %s", apiMetadata)
	req, err := http.NewRequest("POST", r.RegistrationURL, bytes.NewBuffer(apiMetadata))
	check(err)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	check(err)
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		fmt.Println("API registered with success")
	} else {
		check(fmt.Errorf("Registration of API failed with status code %d", resp.StatusCode))
	}
}

func (r registrationApp) generateMetadata(endpoint endpointInfo) []byte {
	specificationsURL, err := url.Parse(r.SystemURL)
	check(err)
	specificationsURL.User = url.UserPassword(r.BasicUser, r.BasicPassword)
	specificationsURL.Path = endpoint.Path + "/$metadata"

	tokenEndpointURL := r.SystemURL + endpoint.Path
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
	`, r.ProviderName, r.ProductName, endpoint.Name, endpoint.Description, r.SystemURL, specificationsURL.String(), r.BasicUser, r.BasicPassword, tokenEndpointURL)
	return []byte(metadata)
}

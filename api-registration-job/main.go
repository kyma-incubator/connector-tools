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
	EventAPIName    string
}

type endpointInfo struct {
	Path        string
	Name        string
	Description string
}

type API struct {
	Name string `json:"name"`
	Id   string `json:"id"`
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
		EventAPIName:    os.Getenv("EVENT_API_NAME"),
	}

	if r.RegistrationURL == "" {
		r.RegistrationURL = fmt.Sprintf("http://application-registry-external-api.kyma-integration.svc.cluster.local:8081/%s/v1/metadata/services", r.ApplicationName)
	}
	r.validateSystemURL()
	fmt.Println("Retrieving already registered APIs")
	apis := r.getRegisteredAPIs()
	r.registerStaticEvents(apis)
	r.readEndpoints(apis)
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

func (r registrationApp) registerStaticEvents(apis []API) {
	eventsString, err := ioutil.ReadFile("files/events.json")
	if err != nil {
		fmt.Println("events.json not found... Moving on.")
		return
	}

	contains := false
	id := ""
	if r.EventAPIName != "" {
		contains, id = containsAPI(apis, r.EventAPIName)
	}

	var req *http.Request
	if contains {
		fmt.Println("Updating events")
		req, err = http.NewRequest("PUT", fmt.Sprintf("%s/%s", r.RegistrationURL, id), bytes.NewBuffer(eventsString))
	} else {
		fmt.Println("Registering events")
		req, err = http.NewRequest("POST", r.RegistrationURL, bytes.NewBuffer(eventsString))
	}
	check(err)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	check(err)
	defer resp.Body.Close()
	if resp.StatusCode < 300 {
		fmt.Println("Events registered with success")
	} else {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		check(err)
		bodyString := string(bodyBytes)
		check(fmt.Errorf("Registration of events failed with status code %d and response body %s", resp.StatusCode, bodyString))
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func (r registrationApp) readEndpoints(apis []API) {
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
				err = r.updateSingleAPI(id, r.generateMetadata(e))
				if err != nil {
					errors = errors + err.Error() + "\n"
					fmt.Printf("Error while update: %s", err)
					continue
				}
			} else {
				fmt.Printf("API %s is not registered yet at kyma application\n", e.Name)
				err = r.registerSingleAPI(r.generateMetadata(e))
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

func containsAPI(apis []API, name string) (bool, string) {
	for _, v := range apis {
		if v.Name == name {
			return true, v.Id
		}
	}
	return false, ""
}

func (r registrationApp) getRegisteredAPIs() []API {
	req, err := http.NewRequest("GET", r.RegistrationURL, nil)
	check(err)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var result []API
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		panic(err)
	}
	return result
}

func (r registrationApp) isAPIActive(path string) (bool, error) {
	url := r.SystemURL + "/" + path
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}
	req.SetBasicAuth(r.BasicUser, r.BasicPassword)
	req.Header.Set("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	jsonResponse := make(map[string]map[string][]string)
	if resp.StatusCode == 200 {
		err := json.NewDecoder(resp.Body).Decode(&jsonResponse)
		if err != nil {
			return false, err
		}
		if len(jsonResponse["d"]["EntitySets"]) > 0 {
			return true, nil
		}
	}
	return false, nil
}

func (r registrationApp) registerSingleAPI(apiMetadata []byte) error {
	fmt.Println("Registering API")
	req, err := http.NewRequest("POST", r.RegistrationURL, bytes.NewBuffer(apiMetadata))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 300 {
		fmt.Println("API registered with success")
	} else {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		bodyString := string(bodyBytes)
		return fmt.Errorf("Registration of API failed with status code %d and response body %s", resp.StatusCode, bodyString)
	}
	return nil
}

func (r registrationApp) updateSingleAPI(id string, apiMetadata []byte) error {
	fmt.Println("Updating API")
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/%s", r.RegistrationURL, id), bytes.NewBuffer(apiMetadata))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 300 {
		fmt.Println("API updated with success")
	} else {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		bodyString := string(bodyBytes)
		return fmt.Errorf("Update of API failed with status code %d and response body %s", resp.StatusCode, bodyString)
	}
	return nil
}

func (r registrationApp) generateMetadata(endpoint endpointInfo) []byte {
	specificationsURL, err := url.Parse(r.SystemURL)
	check(err)
	specificationsURL.User = url.UserPassword(r.BasicUser, r.BasicPassword)
	specificationsURL.Path = endpoint.Path + "/$metadata"

	tokenEndpointURL := fmt.Sprintf("%s/%s", r.SystemURL, endpoint.Path)
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

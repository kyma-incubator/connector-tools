package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type registrationApp struct {
	ApplicationName string
	SystemURL       string
	RegistrationURL string
	ProviderName    string
	ProductName     string
	EventAPIName    string
	app             app
}

type API struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

func main() {
	fmt.Println("Started registration job")

	registrable := getRegistrableApp()

	r := registrationApp{
		ApplicationName: os.Getenv("APPLICATION_NAME"),
		SystemURL:       os.Getenv("SYSTEM_URL"),
		RegistrationURL: os.Getenv("REGISTRATION_URL"),
		ProviderName:    os.Getenv("PROVIDER_NAME"),
		ProductName:     os.Getenv("PRODUCT_NAME"),
		EventAPIName:    os.Getenv("EVENT_API_NAME"),
		app:             registrable,
	}

	if r.RegistrationURL == "" {
		r.RegistrationURL = fmt.Sprintf("http://application-registry-external-api.kyma-integration.svc.cluster.local:8081/%s/v1/metadata/services", r.ApplicationName)
	}
	r.validateSystemURL()
	fmt.Println("Retrieving already registered APIs")
	apis := r.getRegisteredAPIs()
	r.registerStaticEvents(apis)
	r.app.readEndpoints(apis, r)
	fmt.Println("Finished registration job")
}

func getRegistrableApp() app {
	appKind := os.Getenv("APP_KIND")

	defaultAppKind := "odata-with-basic-auth"
	if appKind == "" {
		appKind = defaultAppKind
	}

	switch strings.ToLower(appKind) {
	case defaultAppKind:
		return &oDataWithBasicAuth{
			BasicUser:     os.Getenv("BASIC_USER"),
			BasicPassword: os.Getenv("BASIC_PASSWORD"),
		}
	case "rest-with-apikey":
		return &restWithAPIKey{
			apikey:      os.Getenv("API_KEY"),
			source:      os.Getenv("SOURCE"),
			description: os.Getenv("API_DESCRIPTION"),
		}
	default:
		panic("app kind: " + appKind + "not implemented yet")
	}
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
		check(fmt.Errorf("registration of events failed with status code %d and response body %s", resp.StatusCode, bodyString))
	}
}

func check(e error) {
	if e != nil {
		panic(e)
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
	url := r.app.getAPIUrl(r.SystemURL, path)
	fmt.Printf("url to test active path %s\n", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}
	r.app.setCredentials(req)

	req.Header.Set("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return r.app.verifyActiveResponse(resp)
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
		return fmt.Errorf("registration of API failed with status code %d and response body %s", resp.StatusCode, bodyString)
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
		return fmt.Errorf("update of API failed with status code %d and response body %s", resp.StatusCode, bodyString)
	}
	return nil
}

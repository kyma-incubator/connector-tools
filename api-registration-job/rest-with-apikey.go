package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const format = "json"

type restWithAPIKey struct {
	apikey      string
	source      string
	description string
}

type ApiDefinition struct {
	TargetUrl         string
	SpecificationUrl  string
	QueryParameters   *map[string][]string         `json:",omitempty"`
	Headers           *map[string][]string         `json:",omitempty"`
	Credentials       *map[string]OAuthCredentials `json:",omitempty"`
	RequestParameters *RequestParams               `json:",omitempty"`
}

type RequestParams struct {
	QueryParameters *map[string][]string `json:",omitempty"`
	Headers         *map[string][]string `json:",omitempty"`
}
type ApiMetadata struct {
	Provider    string
	Name        string
	Description string
	Api         ApiDefinition
}

type OAuthCredentials struct {
	Url          string `json:",omitempty"`
	ClientId     string `json:",omitempty"`
	ClientSecret string `json:",omitempty"`
}

func (l *restWithAPIKey) generateMetadata(r registrationApp) []byte {

	targetUrl := r.SystemURL

	metadata := ApiMetadata{
		Provider:    r.ProviderName,
		Name:        l.apiName(r),
		Description: l.description,
		Api: ApiDefinition{
			TargetUrl: targetUrl,
		},
	}

	headers := getParams("headers.json")
	queryParams := getParams("params.json")

	if len(headers) != 0 || len(queryParams) != 0 {

		// The schema has changed for kyma 1.4 so we need to support 1.4 and
		// prior releases. To do this we are setting the header and query params
		// in 2 different ways. While there is no validation in the app registry
		// this works

		// create struct for 1.4
		requestParams := new(RequestParams)
		metadata.Api.RequestParameters = requestParams

		if len(headers) != 0 {
			headersMap := convertMap(headers)
			metadata.Api.Headers = headersMap
			metadata.Api.RequestParameters.Headers = headersMap
		}

		if len(queryParams) != 0 {
			queryParamsMap := convertMap(queryParams)
			metadata.Api.QueryParameters = queryParamsMap
			metadata.Api.RequestParameters.QueryParameters = queryParamsMap
		}
	}

	specURL := getSpecificationUrl(r)

	if specURL != "" {
		metadata.Api.SpecificationUrl = specURL
	}

	if os.Getenv("CLIENT_ID") != "" && os.Getenv("CLIENT_SECRET") != "" && os.Getenv("OAUTH_URL") != "" {
		fmt.Printf("Configuring oauth credentials")
		credentialsMap := make(map[string]OAuthCredentials, 0)
		credentialsMap["oauth"] = OAuthCredentials{
			ClientId:     os.Getenv("CLIENT_ID"),
			ClientSecret: os.Getenv("CLIENT_SECRET"),
			Url:          os.Getenv("OAUTH_URL"),
		}
		metadata.Api.Credentials = &credentialsMap
	}

	var metadataData []byte
	metadataData, err := json.Marshal(metadata)
	check(err)

	fmt.Printf("Metadata = %s\n", string(metadataData))

	return metadataData
}

func (l *restWithAPIKey) setCredentials(request *http.Request) *http.Request {
	fmt.Println("setCredentials not required for rest-with-apikey registration")
	return nil
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

func (l *restWithAPIKey) readEndpoints(apis []API, r registrationApp) {
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
	return fmt.Sprintf("%s - APIs", r.ProductName)
}

/**
* Convert the headers and query parameters from the ConfigMap to json
 */
func getParams(file string) map[string]string {
	var params map[string]string
	fmt.Printf("Loading file %s \n", file)
	paramsString, err := ioutil.ReadFile("files/" + file)
	if err != nil {
		fmt.Printf("%s not found... Moving on.\n", file)
		return nil // TODO error handling
	}

	err = json.Unmarshal(paramsString, &params)
	check(err)
	return params
}

/**
* Convert from map[string]string to map[string][]string, as required by the connector service
 */
func convertMap(origMap map[string]string) *map[string][]string {

	newMap := make(map[string][]string)
	for key, value := range origMap {
		newMap[key] = []string{value}
	}
	return &newMap
}

func getSpecificationUrl(r registrationApp) string {

	specUrlValue := os.Getenv("API_SPECIFICATION_URL")
	fmt.Printf("API_SPECIFICATION_URL = %s\n", specUrlValue)
	if specUrlValue != "" {
		// if fully qualified then use as is
		if !strings.HasPrefix(strings.ToLower(specUrlValue), "http") {
			specUrl := r.SystemURL
			if !strings.HasSuffix(r.SystemURL, "/") {
				specUrl += "/"
			}
			specUrl += specUrlValue
			return specUrl
		}
	}
	return specUrlValue
}

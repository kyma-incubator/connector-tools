package open_connectors

import (
	"context"
	"encoding/json"
	"fmt"
	errorWrap "github.com/kyma-incubator/connector-tools/open-connectors-job/pkg/error"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Connector interface {
	GetConnectorInstances(ctx context.Context, tags []string) ([]Instance, error)
	GetOpenAPISpec(ctx context.Context, ID string, version string) (string, error)
	GetOpenConnectorsContext(ctx context.Context)(string, error)
	GetOpenConnectorsAPIURL(ctx context.Context) string
	CreateAPIAuthorizationHeader(ctx context.Context, instance *Instance) string
}

type Client struct {
	apiBaseUrl *url.URL
	organizationID string
	userID string
	httpClient *http.Client
}

func (c *Client) createAuthorizationHeader() string {
	return fmt.Sprintf("User %s, Organization %s", c.userID, c.organizationID)
}


func NewWithClient(ctx context.Context, organizationID string, userID string, apiHost string,
	httpClient *http.Client) (*Client, error) {


	log.Debugf("creating new open connectors client with custom client")
	url, err := url.Parse(fmt.Sprintf("https://%s/elements/api-v2", apiHost))

	if err != nil {
		log.Errorf("Invalid API host: %s. Value should be like " +
			"(my.openconnectors.ext.hanatrial.ondemand.com)", apiHost)
		return nil, errorWrap.WrapError(err, "Invalid API host: %s. Value should be like " +
			"(my.openconnectors.ext.hanatrial.ondemand.com)", apiHost)
	}

	return &Client{
		apiBaseUrl:     url,
		organizationID: organizationID,
		userID:         userID,
		httpClient:     httpClient,
	}, nil
}

func NewWithTimeout(ctx context.Context, organizationID string, userID string, apiHost string,
	timeoutMills int) (*Client, error) {

	log.Debugf("creating new open connectors client with timeout")

	return NewWithClient(ctx, organizationID, userID, apiHost,
		&http.Client{Timeout: time.Duration(timeoutMills)*time.Millisecond})
}

func New(ctx context.Context, organizationID string, userID string, apiHost string) (*Client, error) {

	log.Debugf("creating new open connectors client with timeout")

	return NewWithClient(ctx, organizationID, userID, apiHost,
		&http.Client{})
}


// GetConnectorInstances retrieves a list of registered Connector Instances for an SAP CP Open Connectors
// Tenant
func (c *Client) GetConnectorInstances(ctx context.Context, tags []string) ([]Instance, error) {

	log.Debugf("retrieving open connectors instances")

	var requestUrl string

	if len(tags) > 0 {
		tempRequestUrl := *c.apiBaseUrl
		tempRequestUrl.Path = fmt.Sprintf("%s/instances", tempRequestUrl.Path)

		var query string
		for i := range tags {
			query += fmt.Sprintf("&tags[]=%s", tags[i])
		}

		if len(tempRequestUrl.RawQuery) > 0 {
			tempRequestUrl.RawQuery= fmt.Sprintf("%s&%s",query[1:], tempRequestUrl.RawQuery)
		} else {
			tempRequestUrl.RawQuery = query[1:]
		}

		requestUrl = tempRequestUrl.String()
		fmt.Println(requestUrl)

	} else {
		requestUrl = fmt.Sprintf("%s/instances", c.apiBaseUrl.String())
	}

	req, err := http.NewRequest(http.MethodGet, requestUrl, nil)
	if err != nil {
		log.Errorf("error requesting instances from open connectors: %s", err.Error())
		return nil, errorWrap.WrapError(err, "error requesting instances from open connectors")
	}

	req.Header.Set("Authorization", c.createAuthorizationHeader())
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Errorf("error requesting instances from open connectors: %s", err.Error())
		return nil, errorWrap.WrapError(err, "error requesting instances from open connectors")
	}

	if resp.StatusCode != http.StatusOK {
		log.Errorf("error requesting instances form open connectors (status %d)", resp.StatusCode)
		return nil, errorWrap.WrapError(fmt.Errorf("invalid http status received:  %d", resp.StatusCode),
			"error requesting instances from open connectors")
	}

	dec := json.NewDecoder(resp.Body)
	defer resp.Body.Close()

	var connectorInstanceResponse []connectorInstanceResponse
	err = dec.Decode(&connectorInstanceResponse)

	if err != nil {
		log.Errorf("error decoding response from open connectors: %s", err.Error())
		return nil, errorWrap.WrapError(err, "error requesting instances from open connectors")
	}

	instances := make([]Instance, len(connectorInstanceResponse))

	for i := range connectorInstanceResponse {
		instances[i].ID = strconv.FormatInt(connectorInstanceResponse[i].ID, 10)
		instances[i].Name = connectorInstanceResponse[i].Name
		instances[i].APIKey = connectorInstanceResponse[i].Token
		instances[i].ConnectorName = connectorInstanceResponse[i].Element.Name
		instances[i].ConnectorKey = connectorInstanceResponse[i].Element.Key
		instances[i].ConnectorID =strconv.FormatInt(connectorInstanceResponse[i].Element.ID, 10)
	}

	return instances, nil
}

// GetOpenAPISpec retrieves the Open API spec of the respective Connector Instance, if a version is
// provided, that version is fetched
func (c *Client) GetOpenAPISpec(ctx context.Context, ID string, version string) (string, error){

	if log.GetLevel() == log.TraceLevel {
		log.Tracef("retrieving Open API spec for instance %q, version %q", ID, version)
	} else {
		log.Debugf("retrieving Open API spec")
	}



	//set version to latest if not provided
	if version == "" {
		version = "-1"
	}

	requestUrl := fmt.Sprintf("%s/instances/%s/docs?version=%s", c.apiBaseUrl.String(), ID, version)

	log.Tracef("requesting Open API spec from %q", requestUrl)

	req, err := http.NewRequest(http.MethodGet, requestUrl, nil)
	if err != nil {
		log.Errorf("error requesting Open API spec from open connectors: %s", err.Error())
		return "", errorWrap.WrapError(err, "error requesting Open API spec from open connectors")
	}

	req.Header.Set("Authorization", c.createAuthorizationHeader())
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Errorf("error requesting Open API spec from open connectors: %s", err.Error())
		return "", errorWrap.WrapError(err, "error requesting Open API spec from open connectors")
	}

	if resp.StatusCode != http.StatusOK {
		log.Errorf("error requesting Open API spec from open connectors (status %d)", resp.StatusCode)
		return "", errorWrap.WrapError(fmt.Errorf("invalid http status received:  %d", resp.StatusCode),
			"error requesting Open API spec from open connectors")
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Errorf("error reading Open API spec response from open connectors: %s", err.Error())
		return "", errorWrap.WrapError(err, "error reading Open API spec response from open connectors")
	}

	return string(respBytes), nil
}

func (c *Client) GetOpenConnectorsContext(ctx context.Context)(string, error) {
	return c.apiBaseUrl.Host, nil
}

func (c *Client) GetOpenConnectorsAPIURL(ctx context.Context) string {
	return c.apiBaseUrl.String()
}

func (c *Client) CreateAPIAuthorizationHeader(ctx context.Context, instance *Instance) string {
	return fmt.Sprintf("%s, Element %s", c.createAuthorizationHeader(),instance.APIKey)
}



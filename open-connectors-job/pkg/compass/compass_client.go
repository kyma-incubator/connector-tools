package compass

import (
	"context"
	errorWrap "github.com/kyma-incubator/connector-tools/open-connectors-job/pkg/error"
	"github.com/machinebox/graphql"
	log "github.com/sirupsen/logrus"
	"net/http"
)

const (
	openConnectorsLabelKey = "open_connectors"
)

type Connector interface {
	GetApplications(ctx context.Context, connectorInstanceContext string) ([]Application, error)
	CreateApplication(ctx context.Context, name string, description string, connectorInstanceContext string,
		connectorInstanceID string) (string, error)
	CreateAPIForApplication(ctx context.Context, applicationId string, apiName string,
		version string, targetUrl string, authorizationHeader string, openAPISpecJsonString string) (string, error)
	DeleteApplication(ctx context.Context, applicationId string) (string, error)
}



//Client to access Director in a specific tenant context
type Client struct {
	graphQlClient *graphql.Client
	tenantId      string
}


//Creates new client to access Director in a specific tenant context
func New(httpClient *http.Client, directorURL string, tenantId string) (*Client, error) {

	log.Debugf("creating new compass client for tenant %q", tenantId)
	graphQlClient := graphql.NewClient(directorURL, graphql.WithHTTPClient(httpClient))

	return &Client{
		graphQlClient: graphQlClient,
		tenantId:      tenantId,
	}, nil
}

func (c *Client) GetApplications(ctx context.Context, connectorInstanceContext string) ([]Application, error) {
	log.Debug("reading applications from compass")

	graphQlRequest := graphql.NewRequest(`
    query ($key: String!) {
  		applications (filter: {
			key: $key
  		}){
    		data {
      			id
      			name
				labels
				apis {
        			data {
          				id
          				name
						version {
            				value
          				}
					}
      			}
    		}
  		}
	}`)

	graphQlRequest.Header.Set("tenant", c.tenantId)
	graphQlRequest.Var("key", openConnectorsLabelKey)

	var respData compassApplicationResponseApplications

	log.Debug("making graphql request to compass")
	err := c.graphQlClient.Run(ctx, graphQlRequest, &respData)

	if err != nil {
		log.Errorf("Error establishing connection to Compass: Requesting list of "+
			"Applications failed. Original error was %q", err.Error())
		return nil, errorWrap.WrapError(err, "Error establishing connection to Compass: Requesting list of "+
			"Applications failed")
	}

	log.Debugf("received response with %d applications", len(respData.Applications.Data))

	//ToDo make length as big as response size once API is adapted
	applications := make([]Application,0)
	for i := range respData.Applications.Data {
		//ToDo remove once filtering is applicable to API
		if respData.Applications.Data[i].Labels.OpenConnectors.ConnectorInstanceContext != connectorInstanceContext {
			log.Debugf("skipping application %q due to connectorinstancecontext %q",
				respData.Applications.Data[i].Name,
				respData.Applications.Data[i].Labels.OpenConnectors.ConnectorInstanceContext)
			continue
		}


		log.Debugf("adding application %q due to connectorinstancecontext %q",
			respData.Applications.Data[i].Name,
			respData.Applications.Data[i].Labels.OpenConnectors.ConnectorInstanceContext)

		apis := make([]API, len(respData.Applications.Data[i].APIs.Data))
		for j := range respData.Applications.Data[i].APIs.Data {
			apis[j] = API{
				ID: respData.Applications.Data[i].APIs.Data[j].ID,
				Version: respData.Applications.Data[i].APIs.Data[j].Version.Value,
				Name: respData.Applications.Data[i].APIs.Data[j].Name,
			}
		}

		applications = append(applications, Application{
			ID:                       respData.Applications.Data[i].ID,
			Name:                     respData.Applications.Data[i].Name,
			ConnectorInstanceID:      respData.Applications.Data[i].Labels.OpenConnectors.ConnectorInstanceID,
			ConnectorInstanceContext: respData.Applications.Data[i].Labels.OpenConnectors.ConnectorInstanceContext,
			APIs:                     &apis,
		})

	}

	log.Debugf("successfully returned %d applications for context %q",
		len(applications),
		connectorInstanceContext)

	return applications, nil
}

func (c *Client) CreateApplication(ctx context.Context, name string, description string,
	connectorInstanceContext string, connectorInstanceID string) (string, error) {

	if log.GetLevel() == log.TraceLevel {
		log.Tracef("creating new compass application with name %q, connectorInstanceContext %q and " +
			"connectorInstanceID %q", name, connectorInstanceContext, connectorInstanceID)
	} else {
		log.Debug("creating new compass application")
	}

	graphQlRequest := graphql.NewRequest(`
    mutation (
  		$name:String!,
		$description:String!,
  		$connectorInstanceContext:String!,
  		$connectorInstanceID:String!,
	) {
  		createApplication (in: {
    		name: $name, 
			description: $description,
			labels: {
      			open_connectors: {
        			connectorInstanceContext: $connectorInstanceContext,
        			connectorInstanceID:$connectorInstanceID
      			}
    		}
		})
  		{
    		id
  		}
	}`)

	graphQlRequest.Header.Set("tenant", c.tenantId)
	graphQlRequest.Var("name", name)
	graphQlRequest.Var("description", description)
	graphQlRequest.Var("connectorInstanceContext", connectorInstanceContext)
	graphQlRequest.Var("connectorInstanceID", connectorInstanceID)

	var respData compassCreateApplicationResponse

	err := c.graphQlClient.Run(ctx, graphQlRequest, &respData)

	if err != nil {
		log.Errorf("Error establishing connection to Compass: Creation of "+
			"Application failed. Original error was %q", err.Error())
		return "", errorWrap.WrapError(err, "Error establishing connection to Compass: Creation of "+
			"Application failed.")
	}


	return respData.CreateApplication.ID, nil
}

func (c *Client) CreateAPIForApplication(ctx context.Context, applicationId string, apiName string,
	version string, targetUrl string, authorizationHeader string, openAPISpecJsonString string) (string, error) {

	if log.GetLevel() == log.TraceLevel {
		log.Tracef("creating new compass API for application %q with apiName %q, version %q and " +
			"targetUrl %q", applicationId, apiName, version, targetUrl)
	} else {
		log.Debugf("creating new compass API for application %q", applicationId)
	}

	graphQlRequest := graphql.NewRequest(`
    mutation (
  		$openAPISpecJsonString:CLOB!,
  		$applicationID:ID!,
  		$name:String!,
  		$targetURL:String!,
  		$version:String!,
  		$authorization:String!
  	){
   	addAPI( 
    	applicationID: $applicationID,
    	in: {
      		name: $name
      		targetURL: $targetURL,
      		version: {
        		value: $version
      		}
      		spec: {
        		type: OPEN_API
        		format: JSON
        		data:  $openAPISpecJsonString
      		}
      		defaultAuth: {
        		credential: {
          			basic: {
            			username: "karl"
            			password: "küma"
          			}
        		}
        		additionalHeaders: {
          			Authorization: [$authorization]
        		}
      		}
  		})
  		{
    		id
  		}
	}`)

	graphQlRequest.Header.Set("tenant", c.tenantId)
	graphQlRequest.Var("openAPISpecJsonString", openAPISpecJsonString)
	graphQlRequest.Var("applicationID", applicationId)
	graphQlRequest.Var("name", apiName)
	graphQlRequest.Var("targetURL", targetUrl)
	graphQlRequest.Var("version", version)
	graphQlRequest.Var("authorization", authorizationHeader)

	var respData compassCreateAPIResponse

	err := c.graphQlClient.Run(ctx, graphQlRequest, &respData)

	if err != nil {
		log.Errorf("Error establishing connection to Compass: Creation of "+
			"API failed. Original error was %q", err.Error())
		return "", errorWrap.WrapError(err, "Error establishing connection to Compass: Creation of "+
			"API failed.")
	}


	return respData.AddAPI.ID, nil
}

func (c *Client) DeleteApplication(ctx context.Context, applicationId string) (string, error) {

	log.Debugf("deleting Compass application %q", applicationId)

	graphQlRequest := graphql.NewRequest(`
    mutation ($appId:ID!) {
  		deleteApplication(id:$appId){
    		id
  		}
	}`)

	graphQlRequest.Header.Set("tenant", c.tenantId)
	graphQlRequest.Var("appId", applicationId)

	var respData compassDeleteApplicationResponse

	err := c.graphQlClient.Run(ctx, graphQlRequest, &respData)

	if err != nil {
		log.Errorf("Error establishing connection to Compass: Deletion of "+
			"Application failed. Original error was %q", err.Error())
		return "", errorWrap.WrapError(err, "Error establishing connection to Compass: Deletion of "+
			"Application failed")
	}


	return respData.DeleteApplication.ID, nil

}




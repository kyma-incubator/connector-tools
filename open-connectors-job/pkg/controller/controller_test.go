package controller

import (
	"context"
	"fmt"
	"github.com/kyma-incubator/connector-tools/open-connectors-job/pkg/compass"
	errorWrap "github.com/kyma-incubator/connector-tools/open-connectors-job/pkg/error"
	"github.com/kyma-incubator/connector-tools/open-connectors-job/pkg/open_connectors"
	"strings"
	"sync"
	"testing"
)

type compassClient struct {
	applicationsCreatedMutex *sync.Mutex
	apisCreatedMutex         *sync.Mutex
	applicationsDeletedMutex *sync.Mutex

	applicationsCreated map[string]struct{}
	apisCreated         map[string]struct{}
	applicationsDeleted map[string]struct{}
}

func (c *compassClient) GetApplications(ctx context.Context, connectorInstanceContext string) ([]compass.Application, error) {
	if connectorInstanceContext == "error" {
		return nil, errorWrap.WrapError(fmt.Errorf("dummy"), "Something went wrong")
	}

	return []compass.Application{
		compass.Application{
			ID:                       "d7282c23-332f-43cf-8ca4-c2fd81f24af3",
			Name:                     "App1",
			ConnectorInstanceID:      "11",
			ConnectorInstanceContext: connectorInstanceContext,
			APIs: &[]compass.API{
				compass.API{
					ID:      "5a3b43d3-4a1c-4e47-90f7-0cbcb5fa657f",
					Name:    "Test API",
					Version: "1",
				},
			},
		},
		compass.Application{
			ID:                       "0cd1d74b-08ac-4155-b461-624c54ce4920",
			Name:                     "App2",
			ConnectorInstanceID:      "12",
			ConnectorInstanceContext: connectorInstanceContext,
			APIs: &[]compass.API{
				compass.API{
					ID:      "4f4515a8-8911-4aad-94b1-a2b2cae6033f",
					Name:    "Test API",
					Version: "1",
				},
				compass.API{
					ID:      "569f1d35-0e8d-4b9a-aaa9-0cb5f86c78f7",
					Name:    "Test API",
					Version: "1",
				},
			},
		},
		compass.Application{
			ID:                       "2e3470cf-7e26-41a6-93d3-d245eb98975b",
			Name:                     "App3",
			ConnectorInstanceID:      "13",
			ConnectorInstanceContext: connectorInstanceContext,
			APIs: &[]compass.API{
				compass.API{
					ID:      "b6ac8228-514d-4dd0-a071-fc5749ba6ac9",
					Name:    "Test API",
					Version: "1",
				},
			},
		},
	}, nil
}
func (c *compassClient) CreateApplication(ctx context.Context, name string, description string,
	connectorInstanceContext string, connectorInstanceID string) (string, error) {

	if connectorInstanceID == "0815-error" {
		return "", errorWrap.WrapError(fmt.Errorf("dummy"), "Something went wrong")
	}

	c.applicationsCreatedMutex.Lock()
	c.applicationsCreated[connectorInstanceID] = struct{}{}
	c.applicationsCreatedMutex.Unlock()

	return name, nil
}
func (c *compassClient) CreateAPIForApplication(ctx context.Context, applicationId string, apiDescription string, apiID string,
	version string, targetUrl string, authorizationHeader string, openAPISpecJsonString string) (string, error) {

	if apiDescription == "Error: 4711-error" {
		return "", errorWrap.WrapError(fmt.Errorf("dummy"), "Something went wrong")
	}

	c.applicationsCreatedMutex.Lock()
	c.apisCreated[fmt.Sprintf("%s-%s", applicationId, apiDescription)] = struct{}{}
	c.applicationsCreatedMutex.Unlock()

	return fmt.Sprintf("%s-%s", applicationId, apiDescription), nil
}
func (c *compassClient) DeleteApplication(ctx context.Context, applicationId string) (string, error) {

	if applicationId == "error" {
		return "", errorWrap.WrapError(fmt.Errorf("dummy"), "Something went wrong")
	}

	c.applicationsDeletedMutex.Lock()
	c.applicationsDeleted[applicationId] = struct{}{}
	c.applicationsDeletedMutex.Unlock()

	return applicationId, nil
}

type openConnectorsClient struct {
	apiHost string
}

func (c *openConnectorsClient) GetConnectorInstances(ctx context.Context, tags []string) ([]open_connectors.Instance, error) {
	if len(tags) > 0 && tags[0] == "error" {
		return nil, errorWrap.WrapError(fmt.Errorf("dummy"), "Something went wrong")
	}
	return []open_connectors.Instance{
		open_connectors.Instance{
			ID:            "11",
			Name:          "Karl Küma",
			APIKey:        "1a7ff55f-0468-4bd5-a0eb-9410836fbb0a",
			ConnectorName: "Twitter",
			ConnectorID:   "1347",
		},
		open_connectors.Instance{
			ID:            "14",
			Name:          "New App",
			APIKey:        "16007909-6486-4ee4-832c-d95c453cacd1",
			ConnectorName: "Facebook",
			ConnectorID:   "1348",
		},
	}, nil
}
func (c *openConnectorsClient) GetOpenAPISpec(ctx context.Context, ID string, version string) (string, error) {
	if ID == "error" {
		return "", errorWrap.WrapError(fmt.Errorf("dummy"), "Something went wrong")
	}
	return `{"swagger": "2.0"}`, nil
}
func (c *openConnectorsClient) GetOpenConnectorsContext(ctx context.Context) (string, error) {
	return c.apiHost, nil
}
func (c *openConnectorsClient) GetOpenConnectorsAPIURL(ctx context.Context) string {
	return fmt.Sprintf("https://%s/elements/api-v2/", c.apiHost)
}
func (c *openConnectorsClient) CreateAPIAuthorizationHeader(ctx context.Context, instance *open_connectors.Instance) string {
	return "User user, Organization org, Element 0815"
}

func TestNew(t *testing.T) {

	compass := &compassClient{}
	openConnectors := &openConnectorsClient{}

	ctx := context.Background()

	controller, err := New(ctx, compass, openConnectors, nil, "prefix")
	if err != nil {
		t.Fatalf("error creating controller instance: %s", err.Error())
	}

	if controller.compass != compass || controller.openConnectors != openConnectors {
		t.Error("inputs not passed correctly")
	}

	if controller.tags != nil {
		t.Error("expected controller.tags to be nil, but found non-nil")
	}

	tags := []string{"test1", "test2"}
	controller, err = New(ctx, compass, openConnectors, tags, "prefix")
	if err != nil {
		t.Fatalf("error creating controller instance: %s", err.Error())
	}
	if controller.tags[0] != tags[0] || controller.tags[1] != tags[1] {
		t.Error("expected controller.tags to match inputs but it didn't")
	}

}

func Test_DetermineStatus(t *testing.T) {

	ctx := context.Background()
	controller, _ := New(ctx, &compassClient{}, &openConnectorsClient{apiHost: "hello.world.com"}, []string{},
		"prefix")

	instancesToAdd, existingApps, appsToDelete, err := controller.DetermineStatus(ctx)
	if err != nil {
		t.Fatalf("error classifying work: %s", err.Error())
	}

	if len(instancesToAdd) != 1 {
		t.Fatalf("expected 1 application to add, received %d", len(instancesToAdd))
	}

	if instancesToAdd[0].ID != "14" {
		t.Errorf("expected ID 14 to add, received %s", instancesToAdd[0].ID)
	}

	if len(existingApps) != 1 {
		t.Fatalf("expected 1 application to be existing, received %d", len(instancesToAdd))
	}

	if existingApps[0].ID != "d7282c23-332f-43cf-8ca4-c2fd81f24af3" {
		t.Errorf("expected ID \"d7282c23-332f-43cf-8ca4-c2fd81f24af3\" to be existing, received %s", instancesToAdd[0].ID)
	}

	if len(appsToDelete) != 2 {
		t.Fatalf("expected 2 applications to delete, received %d", len(appsToDelete))
	}
	if !(appsToDelete[0].ID == "0cd1d74b-08ac-4155-b461-624c54ce4920" ||
		appsToDelete[1].ID != "0cd1d74b-08ac-4155-b461-624c54ce4920") {
		t.Errorf("expected application \"0cd1d74b-08ac-4155-b461-624c54ce4920\" to be on the apps" +
			" to delete list, but did not find it")
	}

	if !(appsToDelete[0].ID != "2e3470cf-7e26-41a6-93d3-d245eb98975b" ||
		appsToDelete[1].ID != "2e3470cf-7e26-41a6-93d3-d245eb98975b") {
		t.Errorf("expected application \"2e3470cf-7e26-41a6-93d3-d245eb98975b\" to be on the apps" +
			" to delete list, but did not find it")
	}
}

func Test_createNewApplications(t *testing.T) {
	ctx := context.Background()
	openConnectorsClient := &openConnectorsClient{apiHost: "hello.world.com"}

	instances, _ := openConnectorsClient.GetConnectorInstances(ctx, nil)
	compassClnt := &compassClient{
		applicationsCreatedMutex: &sync.Mutex{},
		apisCreatedMutex:         &sync.Mutex{},
		applicationsDeletedMutex: &sync.Mutex{},
		applicationsCreated:      make(map[string]struct{}),
		apisCreated:              make(map[string]struct{}),
		applicationsDeleted:      make(map[string]struct{}),
	}

	controller, _ := New(ctx, compassClnt, openConnectorsClient, []string{}, "prefix")

	if err := controller.createNewApplications(ctx, instances); err != nil {
		t.Errorf("error creating applications: %s", err.Error())
	}

	for i := range instances {
		if _, ok := compassClnt.applicationsCreated[instances[i].ID]; ok {
			delete(compassClnt.applicationsCreated, instances[i].ID)
		} else {
			t.Errorf("expected to find connector instance id %s, but didn't", instances[i].ID)
		}

		apiName := fmt.Sprintf("prefix-%s-%s: %s", instances[i].ID, instances[i].ConnectorName, instances[i].Name)
		if _, ok := compassClnt.apisCreated[apiName]; ok {
			delete(compassClnt.apisCreated, apiName)
		} else {
			t.Errorf("expected to find api %s, but didn't", apiName)
		}
	}

	if len(compassClnt.applicationsCreated) != 0 {
		t.Errorf("too many applications were created")
	}

	if len(compassClnt.apisCreated) != 0 {
		t.Errorf("too many apis were created")
	}
}

func Test_createNewApplicationsError(t *testing.T) {
	ctx := context.Background()
	openConnectorsClient := &openConnectorsClient{apiHost: "hello.world.com"}

	instances := []open_connectors.Instance{
		open_connectors.Instance{
			ID:            "error",
			Name:          "Karl Küma",
			APIKey:        "1a7ff55f-0468-4bd5-a0eb-9410836fbb0a",
			ConnectorName: "Twitter",
			ConnectorKey:  "twitter",
			ConnectorID:   "1347",
		},
		open_connectors.Instance{
			ID:            "0815-error",
			Name:          "Karl Küma",
			APIKey:        "1a7ff55f-0468-4bd5-a0eb-9410836fbb0a",
			ConnectorName: "Twitter",
			ConnectorKey:  "twitter",
			ConnectorID:   "1347",
		},
		open_connectors.Instance{
			ID:            "15",
			Name:          "4711-error",
			APIKey:        "1a7ff55f-0468-4bd5-a0eb-9410836fbb0a",
			ConnectorName: "Error",
			ConnectorKey:  "error",
			ConnectorID:   "1347",
		},
		open_connectors.Instance{
			ID:            "11",
			Name:          "Karl Küma",
			APIKey:        "1a7ff55f-0468-4bd5-a0eb-9410836fbb0a",
			ConnectorName: "Twitter",
			ConnectorKey:  "twitter",
			ConnectorID:   "1347",
		},
	}
	compassClnt := &compassClient{
		applicationsCreatedMutex: &sync.Mutex{},
		apisCreatedMutex:         &sync.Mutex{},
		applicationsDeletedMutex: &sync.Mutex{},
		applicationsCreated:      make(map[string]struct{}),
		apisCreated:              make(map[string]struct{}),
		applicationsDeleted:      make(map[string]struct{}),
	}

	controller, _ := New(ctx, compassClnt, openConnectorsClient, []string{}, "prefix")

	if err := controller.createNewApplications(ctx, instances); err == nil {
		t.Error("error creating applications expected, but not received")
	}

	if len(compassClnt.applicationsCreated) != 2 {
		t.Errorf("expected number of application creations is 1, received %d", len(compassClnt.applicationsCreated))
	}

	if len(compassClnt.apisCreated) != 1 {
		t.Errorf("expected number of api creations is 1, received %d", len(compassClnt.apisCreated))
	}
}

func Test_deleteApplications(t *testing.T) {

	ctx := context.Background()

	compassClnt := &compassClient{
		applicationsCreatedMutex: &sync.Mutex{},
		apisCreatedMutex:         &sync.Mutex{},
		applicationsDeletedMutex: &sync.Mutex{},
		applicationsCreated:      make(map[string]struct{}),
		apisCreated:              make(map[string]struct{}),
		applicationsDeleted:      make(map[string]struct{}),
	}

	openConnectorsClient := &openConnectorsClient{apiHost: "hello.world.com"}

	controller, _ := New(ctx, compassClnt, openConnectorsClient, []string{}, "prefix")

	apps, _ := compassClnt.GetApplications(ctx, "dummy")

	if err := controller.deleteApplications(ctx, apps); err != nil {
		t.Errorf("error deleting applications: %s", err.Error())
	}

	for i := range apps {
		if _, ok := compassClnt.applicationsDeleted[apps[i].ID]; ok {
			delete(compassClnt.applicationsDeleted, apps[i].ID)
		} else {
			t.Errorf("expected application %s to be deleted but it wasn't", apps[i].ID)
		}
	}

	if len(compassClnt.applicationsDeleted) > 0 {
		t.Errorf("expected all apps to be delted but %d weren't", len(compassClnt.applicationsDeleted))
	}
}

func Test_deleteApplicationsError(t *testing.T) {

	ctx := context.Background()

	compassClnt := &compassClient{
		applicationsCreatedMutex: &sync.Mutex{},
		apisCreatedMutex:         &sync.Mutex{},
		applicationsDeletedMutex: &sync.Mutex{},
		applicationsCreated:      make(map[string]struct{}),
		apisCreated:              make(map[string]struct{}),
		applicationsDeleted:      make(map[string]struct{}),
	}

	openConnectorsClient := &openConnectorsClient{apiHost: "hello.world.com"}

	controller, _ := New(ctx, compassClnt, openConnectorsClient, []string{}, "prefix")

	apps := []compass.Application{
		compass.Application{
			ID:                       "error",
			Name:                     "App1",
			ConnectorInstanceID:      "11",
			ConnectorInstanceContext: "hello.world.com",
			APIs: &[]compass.API{
				compass.API{
					ID:      "5a3b43d3-4a1c-4e47-90f7-0cbcb5fa657f",
					Name:    "Test API",
					Version: "1",
				},
			},
		},
		compass.Application{
			ID:                       "no-error",
			Name:                     "App1",
			ConnectorInstanceID:      "11",
			ConnectorInstanceContext: "hello.world.com",
			APIs: &[]compass.API{
				compass.API{
					ID:      "5a3b43d3-4a1c-4e47-90f7-0cbcb5fa657f",
					Name:    "Test API",
					Version: "1",
				},
			},
		},
	}

	if err := controller.deleteApplications(ctx, apps); err == nil {
		t.Error("error expected but none received")
	}

	if len(compassClnt.applicationsDeleted) != 1 {
		t.Errorf("expected 1 app to be delted but received %d ", len(compassClnt.applicationsDeleted))
	}
}

func Test_createApplicationDescription(t *testing.T) {
	instance := &open_connectors.Instance{
		Name:          "test 123",
		ID:            "4711",
		ConnectorName: "Twitter",
		ConnectorID:   "4711",
		APIKey:        "0815",
	}

	description := createApplicationDescription(instance)

	if description !=
		fmt.Sprintf("SAP Cloud Platform Open Connectors %s: %s", instance.ConnectorName, instance.Name) {
		t.Errorf("wrong application description: %s", description)
	}
}

func Test_createApplicationName(t *testing.T) {

	instance := &open_connectors.Instance{
		Name:          "test 123",
		ID:            "4711",
		ConnectorName: "Twitter",
		ConnectorID:   "4711",
		APIKey:        "0815",
	}

	name := createApplicationName("pref", instance)

	if name != fmt.Sprintf("%s-%s", "pref", instance.ID) {
		t.Errorf("wrong application description: %s", name)
	}

	if strings.Contains(name, " ") {
		t.Errorf("name %q contains a blank", name)
	}

}

func TestController_Synchronize(t *testing.T) {

	ctx := context.Background()

	compassClnt := &compassClient{
		applicationsCreatedMutex: &sync.Mutex{},
		apisCreatedMutex:         &sync.Mutex{},
		applicationsDeletedMutex: &sync.Mutex{},
		applicationsCreated:      make(map[string]struct{}),
		apisCreated:              make(map[string]struct{}),
		applicationsDeleted:      make(map[string]struct{}),
	}

	openConnectorsClient := &openConnectorsClient{apiHost: "hello.world.com"}

	controller, _ := New(ctx, compassClnt, openConnectorsClient, []string{}, "prefix")

	if err := controller.Synchronize(ctx); err != nil {
		t.Fatalf("unexpected error synchronizing: %s", err.Error())
	}

	if len(compassClnt.applicationsCreated) != 1 {
		t.Errorf("expected creation of 1 application, received %d", len(compassClnt.applicationsCreated))
	}

	if len(compassClnt.apisCreated) != 1 {
		t.Errorf("expected creation of 1 api, received %d", len(compassClnt.apisCreated))
	}

	if len(compassClnt.applicationsDeleted) != 2 {
		t.Errorf("expected deletion of 2 applications, received %d", len(compassClnt.applicationsDeleted))
	}

}

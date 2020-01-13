package open_connectors

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNew(t *testing.T) {

	client, err := NewWithTimeout(context.Background(), "org", "user",
		"api.openconnectors.ext.hanatrial.ondemand.com", 2000)

	if err != nil {
		t.Fatalf("error creating open connectors API client: %s", err.Error())
	}

	if client.apiBaseUrl.String() != "https://api.openconnectors.ext.hanatrial.ondemand.com/elements/api-v2/" {
		t.Errorf("apiHost should be \"https://api.openconnectors.ext.hanatrial.ondemand.com/elements/api-v2/\""+
			"but is %q", client.apiBaseUrl.String())
	}
}

func TestClient_GetConnectorInstances(t *testing.T) {

	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header.Get("Authorization") != "User user, Organization org" {
			t.Errorf("invalid Authorization header supplied, should be "+
				"\"User user, Organization org\" but was %q", r.Header.Get("Authorization"))
		}

		w.Write([]byte(`
[
  {
    "id": 164195,
    "name": "Karl Küma",
    "token": "test1",
    "tags": [
      "Karl Küma"
    ],
    "element": {
      "id": 2113,
      "name": "Twitter",
      "key": "twitter",
      "description": "Add a Twitter Instance to connect your existing Twitter account to the Social Hub, allowing you to manage statuses and followers across multiple Social Elements. You will need your Twitter account information to add an instance."
    }
  },
  {
    "id": 164196,
    "name": "Karl Kyma",
    "token": "test2",
    "tags": [
      "Karl Küma"
    ],
    "element": {
      "id": 2113,
      "name": "Twitter",
      "key": "twitter",
      "description": "Add a Twitter Instance to connect your existing Twitter account to the Social Hub, allowing you to manage statuses and followers across multiple Social Elements. You will need your Twitter account information to add an instance."
    }
  }
]
`))
	}))
	ctx := context.Background()
	client, _ := NewWithClient(ctx, "org", "user", server.URL[8:], server.Client())

	instances, err := client.GetConnectorInstances(ctx, []string{})

	if err != nil {
		t.Fatalf("error reading open connectors instances: %s", err.Error())
	}

	if len(instances) != 2 {
		t.Errorf("expecting to receive 2 instances, received %d", len(instances))
	}

}

func TestClient_GetOpenAPISpec(t *testing.T) {

	spec := `{
  "swagger": "2.0",
  "info": {
    "title": "Sample API",
    "description": "API description in Markdown.",
    "version": "1.0.0"
  },
  "host": "api.example.com",
  "basePath": "/v1",
  "schemes": [
    "https"
  ],
  "paths": {
    "/users": {
      "get": {
        "summary": "Returns a list of users.",
        "description": "Optional extended description in Markdown.",
        "produces": [
          "application/json"
        ],
        "responses": {
          "200": {
            "description": "OK"
          }
        }
      }
    }
  }
}`
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header.Get("Authorization") != "User user, Organization org" {
			t.Errorf("invalid Authorization header supplied, should be "+
				"\"User user, Organization org\" but was %q", r.Header.Get("Authorization"))
		}

		w.Write([]byte(spec))
	}))

	ctx := context.Background()
	client, _ := NewWithClient(ctx, "org", "user", server.URL[8:], server.Client())

	openAPISpec, err := client.GetOpenAPISpec(ctx, "12345", "")

	if err != nil {
		t.Fatalf("error reading open connectors instances: %s", err.Error())
	}

	if openAPISpec != spec {
		t.Error("received spec does is not equal to the expected result")
	}
}

func TestClient_GetOpenConnectorsContext(t *testing.T) {
	connectsContext := "api.openconnectors.ext.hanatrial.ondemand.com"
	ctx := context.Background()
	client, err := New(ctx, "org", "user", connectsContext)

	if err != nil {
		t.Fatalf("error creating open connectors API client: %s", err.Error())
	}

	receivedContext, err := client.GetOpenConnectorsContext(ctx)
	if err != nil {
		t.Fatalf("error reciving open connectors context: %s", err.Error())
	}

	if receivedContext != connectsContext {
		t.Errorf("expected to receive context \"api.openconnectors.ext.hanatrial.ondemand.com\", but received " +
			"%q", receivedContext)
	}
}

func TestClient_GetOpenConnectorsAPIURL(t *testing.T) {
	connectsContext := "api.openconnectors.ext.hanatrial.ondemand.com"
	ctx := context.Background()
	client, err := New(ctx, "org", "user", connectsContext)

	if err != nil {
		t.Fatalf("error creating open connectors API client: %s", err.Error())
	}

	targetUrl := "https://api.openconnectors.ext.hanatrial.ondemand.com/elements/api-v2"
	if receivedContext := client.GetOpenConnectorsAPIURL(ctx);
		receivedContext != targetUrl {
			t.Errorf("expected to receive %q, received %q", targetUrl, receivedContext)
	}

}

func TestClient_CreateAPIAuthorizationHeader(t *testing.T) {
	apiHost := "api.openconnectors.ext.hanatrial.ondemand.com"
	ctx := context.Background()
	client, err := New(ctx, "org", "user", apiHost)

	if err != nil {
		t.Fatalf("error creating open connectors API client: %s", err.Error())
	}

	instance := &Instance{
		Name: "test 123",
		ID: "4711",
		ConnectorName:"Twitter",
		ConnectorID: "4711",
		APIKey: "0815",
	}

	targetAuth := "User user, Organization org, Element 0815"

	if auth := client.CreateAPIAuthorizationHeader(ctx, instance); auth != targetAuth {
		t.Errorf("expected authorization header %q, received %q", auth, targetAuth)
	}

}

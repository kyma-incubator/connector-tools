package compass

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNew(t *testing.T) {

	tenantId := "tenant123"

	client, err := New(&http.Client{}, "https://myDirecxtor.com", "tenant123")

	if err != nil {
		t.Fatalf("error creating compass client: %s", err.Error())
	}

	if client.tenantId != tenantId {
		t.Errorf("tenantId should be %q, but was %q", tenantId, client.tenantId )
	}
}

func TestClient_GetApplications(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func (w http.ResponseWriter,r *http.Request) {

		w.Write([]byte(`
{
  "data": {
    "applications": {
      "data": [
        {
          "id": "1b8873a9-c468-4b71-83fa-41248469a75a",
          "name": "fifth-try",
          "labels": {
            "open_connectors": {
              "connectorInstanceContext": "api.openconnectors.ext.hanatrial.ondemand.com",
              "connectorInstanceID": "164195"
            },
            "scenarios": [
              "DEFAULT"
            ]
          },
          "apis": {
            "data": [
              {
                "id": "47318f96-25a0-48e3-a8b7-4d7bb6db1348",
                "name": "Twitter",
                "version": null
              }
            ]
          }
        },
{
          "id": "1b8873a9-c468-4b71-83fa-41248469a75a",
          "name": "sixth-try",
          "labels": {
            "open_connectors": {
              "connectorInstanceContext": "api.openconnectors.ext.hana.ondemand.com",
              "connectorInstanceID": "164196"
            },
            "scenarios": [
              "DEFAULT"
            ]
          },
          "apis": {
            "data": [
              {
                "id": "47318f96-25a0-48e3-a8b7-4d7bb6db1348",
                "name": "Twitter",
                "version": null
              }
            ]
          }
        }
      ]
    }
  }
}
		`))

	}))

	httpClient := server.Client()


	client, err := New(httpClient, server.URL, "12345")
	if err != nil {
		t.Fatalf("error creating compass client: %s", err.Error())
	}

	applications, err := client.GetApplications(context.Background(),
		"api.openconnectors.ext.hanatrial.ondemand.com")

	if err != nil {
		t.Fatalf("error creading applications: %s", err.Error())
	}

	if len(applications) != 1 {
		t.Errorf("expected to receive 1 applications, received %d", len(applications))
	}

	if applications[0].Name != "fifth-try" {
		t.Errorf("expected application name \"fifth-try\", received %q", applications[0].Name)
	}

	if len(*applications[0].APIs) != 1 {
		t.Errorf("expected to receive 1 api, but received %d", len(*applications[0].APIs))
	}
}

func TestClient_CreateApplication(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func (w http.ResponseWriter,r *http.Request) {

		w.Write([]byte(`
{
  "data": {
    "createApplication": {
      "id": "12181107-f70b-48ca-8ab7-e97eafb36018"
    }
  }
}
		`))

	}))

	httpClient := server.Client()

	client, err := New(httpClient, server.URL, "12345")
	if err != nil {
		t.Fatalf("error creating compass client: %s", err.Error())
	}

	ctx := context.Background()

	id, err := client.CreateApplication(ctx,"testapp", "tespapp description",
		"context", "12345")

	if err != nil {
		t.Fatalf("error creating compass application: %s", err.Error())
	}

	if id != "12181107-f70b-48ca-8ab7-e97eafb36018" {
		t.Errorf("expected application id to be \"12181107-f70b-48ca-8ab7-e97eafb36018\", received " +
			"%q", id)
	}
}

func TestClient_CreateAPIForApplication(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){

		w.Write([]byte(`{
  "data": {
    "addAPI": {
      "id": "b731cbfc-de12-449c-8690-0d9c6e6db5cd"
    }
  }
}`))
	}))

	httpClient := server.Client()

	client, err := New(httpClient, server.URL, "12345")
	if err != nil {
		t.Fatalf("error creating compass client: %s", err.Error())
	}
	ctx := context.Background()

	id, err := client.CreateAPIForApplication(ctx, "12181107-f70b-48ca-8ab7-e97eafb36018",
		"Test", "1", "https://api.openconnectors.ext.hanatrial.ondemand.com/elements/api-v2/",
		"User usr, Organization org, Element eYB6pKskPBa21vzK599fwX+Lw2jAg=",
		`{"swagger":"2.0"}`)

	if err != nil {
		t.Fatalf("error creating compass api spec: %s", err.Error())
	}

	if id != "b731cbfc-de12-449c-8690-0d9c6e6db5cd" {
		t.Errorf("expected application id to be \"b731cbfc-de12-449c-8690-0d9c6e6db5cd\", received " +
			"%q", id)
	}
}

func TestClient_DeleteApplication(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func (w http.ResponseWriter,r *http.Request) {

		w.Write([]byte(`
{
  "data": {
    "deleteApplication": {
      "id": "12181107-f70b-48ca-8ab7-e97eafb36018"
    }
  }
}
		`))

	}))

	httpClient := server.Client()

	client, err := New(httpClient, server.URL, "12345")
	if err != nil {
		t.Fatalf("error creating compass client: %s", err.Error())
	}

	ctx := context.Background()

	id, err := client.DeleteApplication(ctx,"12181107-f70b-48ca-8ab7-e97eafb36018")

	if err != nil {
		t.Fatalf("error deleting compass application: %s", err.Error())
	}

	if id != "12181107-f70b-48ca-8ab7-e97eafb36018" {
		t.Errorf("expected application id to be \"12181107-f70b-48ca-8ab7-e97eafb36018\", received " +
			"%q", id)
	}
}

/*func TestClient_CreateAPIForApplication2(t *testing.T) {

	ctx := context.Background()
	httpClient, _ := CreateHttpClientTimeout(ctx, "85a45c96-dbad-45b6-a028-fcc695accedd",
		"OnO~KAj5MSwH", "https://oauth2.compass-ak.cluster.extend.cx.cloud.sap/oauth2/token",
		5000)

	client, _ := New(httpClient,
		"https://compass-gateway-auth-oauth.compass-ak.cluster.extend.cx.cloud.sap/director/graphql",
		"3e64ebae-38b5-46a0-b1ed-9ccee153a0ae")

	apps, err := client.GetApplications(ctx, "api.openconnectors.ext.hanatrial.ondemand.com")

	if err != nil {
		t.Errorf("Error reading apps: %s", err.Error())
	}

	fmt.Printf("read %d applications", len(apps))

	appId, err := client.CreateApplication(ctx, "generated-app-3",
		"api.openconnectors.ext.hanatrial.ondemand.com", "3")

	if err != nil {
		t.Errorf("Error creating app: %s", err.Error())
	}

	fmt.Printf("created application %q", appId)

	api1Id, err := client.CreateAPIForApplication(ctx, appId, "test-1", "1",
		"https://petstore.swagger.io/api", "User 123", swaggerspec)

	if err != nil {
		t.Errorf("Error creating API: %s", err.Error())
	}

	fmt.Printf("created api %q", api1Id)

	api2Id, err := client.CreateAPIForApplication(ctx, appId, "test-2", "1",
		"https://petstore.swagger.io/api", "User 123", swaggerspec)

	if err != nil {
		t.Errorf("Error creating API: %s", err.Error())
	}

	fmt.Printf("created api %q", api2Id)
}

const swaggerspec = `{
  "swagger": "2.0",
  "info": {
    "version": "1.0.0",
    "title": "Swagger Petstore",
    "description": "A sample API that uses a petstore as an example to demonstrate features in the swagger-2.0 specification",
    "termsOfService": "http://swagger.io/terms/",
    "contact": {
      "name": "Swagger API Team"
    },
    "license": {
      "name": "MIT"
    }
  },
  "host": "petstore.swagger.io",
  "basePath": "/api",
  "schemes": [
    "http"
  ],
  "consumes": [
    "application/json"
  ],
  "produces": [
    "application/json"
  ],
  "paths": {
    "/pets": {
      "get": {
        "description": "Returns all pets from the system that the user has access to",
        "produces": [
          "application/json"
        ],
        "responses": {
          "200": {
            "description": "A list of pets.",
            "schema": {
              "type": "array",
              "items": {
                "$ref": "#/definitions/Pet"
              }
            }
          }
        }
      }
    }
  },
  "definitions": {
    "Pet": {
      "type": "object",
      "required": [
        "id",
        "name"
      ],
      "properties": {
        "id": {
          "type": "integer",
          "format": "int64"
        },
        "name": {
          "type": "string"
        },
        "tag": {
          "type": "string"
        }
      }
    }
  }
}`*/

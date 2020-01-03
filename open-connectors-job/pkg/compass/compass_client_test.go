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



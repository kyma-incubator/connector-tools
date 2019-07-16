package apiclient

import (
	"encoding/json"
	"github.com/kyma-incubator/connector-tools/qualtrics-webhook-registration/pkg/util"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func Test_NewQualtricsSubscription(t *testing.T) {

	inst, err := NewQualtricsSubscription("test", "http://www.apiclient.com", 2000*time.Microsecond)

	if err != nil {
		t.Errorf("instance must not fail on creation, error %q", err.Error())
	}

	if inst.URL != "http://www.apiclient.com" {
		t.Errorf("url  should be http://www.apiclient.com, but was %q", inst.URL)
	}

	if inst.APIKey != "test" {
		t.Errorf("api key should be test, but was %q", inst.APIKey)
	}

	inst, err = NewQualtricsSubscription("", "http://www.apiclient.com", 2000*time.Microsecond)

	if err == nil {
		t.Error("instance must fail on creation, but didn't")
	}

	inst, err = NewQualtricsSubscription("test", "", 2000*time.Microsecond)

	if err == nil {
		t.Error("instance must fail on creation, but didn't")
	}

}

func TestInstance_CreateSubscription(t *testing.T) {

	createSubscription := QualtricsSubscription{
		Topics:         "test.topic",
		PublicationURL: "https://testemp.com",
		SharedKey:      "kyma",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header.Get("X-API-TOKEN") != "apikey" {
			t.Errorf("X-API-TOKEN header must be apikey, but received %q instead",
				r.Header.Get("X-API-TOKEN"))
		}

		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Content-Type header must be application/json, but received %q instead",
				r.Header.Get("Content-Type"))
		}
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()

		var receivedSubscription QualtricsSubscription

		decoder.Decode(&receivedSubscription)

		if receivedSubscription != createSubscription {
			t.Errorf("expected subscription %1v, received %1v", createSubscription, receivedSubscription)
		}

		w.Header().Set("Content-Type", "application.json")
		file, _ := os.Open("../../testdata/POST_API_v3_eventsubscriptions.json")
		fileBytes, _ := ioutil.ReadAll(file)

		w.Write(fileBytes)

	},
	))

	inst, err := NewQualtricsSubscriptionWithClient("apikey", server.URL, server.Client())

	if err != nil {
		t.Errorf("test should not error out, but error %s received", err.Error())
	}

	id, err := inst.CreateSubscription(&createSubscription, &util.RequestContext{TraceHeaders:http.Header{}})

	if err != nil {
		t.Errorf("test should not error out, but error %s received", err.Error())
	}

	if id != "SUB_8pSh6Kbqjgvwynb" {
		t.Errorf("subscription ID should be SUB_8pSh6Kbqjgvwynb, %q received", id)
	}

}

func TestInstance_DeleteSubscription(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header.Get("X-API-TOKEN") != "apikey" {
			t.Errorf("X-API-TOKEN header must be apikey, but received %q instead",
				r.Header.Get("X-API-TOKEN"))
		}

		if r.URL.Path != "/API/v3/eventsubscriptions/SUB_1zaGQnkmLGYvV2J" {
			t.Errorf("expected path to be \"/API/v3/eventsubscriptions/SUB_1zaGQnkmLGYvV2J\", received" +
				" %q", r.URL.Path )
		}

		w.Header().Set("Content-Type", "application/json")
		file, _ := os.Open("../../testdata/DELETE_API_v3_eventsubscriptions.json")
		fileBytes, _ := ioutil.ReadAll(file)

		w.Write(fileBytes)

	},
	))

	inst, err := NewQualtricsSubscriptionWithClient("apikey", server.URL, server.Client())

	if err != nil {
		t.Errorf("test should not error out, but error %s received", err.Error())
	}

	err = inst.DeleteSubscription("SUB_1zaGQnkmLGYvV2J",
		&util.RequestContext{TraceHeaders:http.Header{}})

	if err != nil {
		t.Errorf("test should not error out, but error %s received", err.Error())
	}
}

func TestInstance_GetSubscriptionList(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header.Get("X-API-TOKEN") != "apikey" {
			t.Errorf("X-API-TOKEN header must be apikey, but received %q instead",
				r.Header.Get("X-API-TOKEN"))
		}

		w.Header().Set("Content-Type", "application/json")
		file, _ := os.Open("../../testdata/GET_API_v3_eventsubscriptions.json")
		fileBytes, _ := ioutil.ReadAll(file)

		w.Write(fileBytes)

	},
	))

	inst, err := NewQualtricsSubscriptionWithClient("apikey", server.URL, server.Client())

	if err != nil {
		t.Errorf("test should not error out, but error %s received", err.Error())
	}

	subs, err := inst.GetSubscriptionList(&util.RequestContext{TraceHeaders:http.Header{}})

	if err != nil {
		t.Errorf("test should not error out, but error %s received", err.Error())
	}

	if len(subs) != 5 {
		t.Errorf("expected 5 subscriptions, received %d", len(subs))
	}

	for _, sub := range subs {
		if sub.SharedKey != "" {
			t.Errorf("expected sharedkey to be empty, received %q", sub.SharedKey)
		}

		if sub.Topics == "" {
			t.Errorf("expected topics to be non-empty, received %q", sub.Topics)
		}

		if sub.PublicationURL == "" {
			t.Errorf("expected PublicationURL to be non-empty, received %q", sub.PublicationURL)
		}

		if sub.ID == "" {
			t.Errorf("expected ID to be non-empty, received %q", sub.ID)
		}
	}
}

func TestInstance_UpdateSubscription(t *testing.T) {

	updateSubscription := QualtricsSubscription{
		ID:             "4711",
		Topics:         "test.topic",
		PublicationURL: "https://testemp.com",
		SharedKey:      "kyma",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header.Get("X-API-TOKEN") != "apikey" {
			t.Errorf("X-API-TOKEN header must be apikey, but received %q instead",
				r.Header.Get("X-API-TOKEN"))
		}

		if r.Method == http.MethodPost {
			decoder := json.NewDecoder(r.Body)
			defer r.Body.Close()

			var receivedSubscription QualtricsSubscription

			decoder.Decode(&receivedSubscription)

			if receivedSubscription != updateSubscription {
				t.Errorf("expected subscription %1v, received %1v", updateSubscription, receivedSubscription)
			}

			w.Header().Set("Content-Type", "application/json")
			file, _ := os.Open("../../testdata/POST_API_v3_eventsubscriptions.json")
			fileBytes, _ := ioutil.ReadAll(file)

			w.Write(fileBytes)
		} else if r.Method == http.MethodDelete {
			w.Header().Set("Content-Type", "application.json")
			file, _ := os.Open("../../testdata/DELETE_API_v3_eventsubscriptions.json")
			fileBytes, _ := ioutil.ReadAll(file)

			w.Write(fileBytes)
		} else {
			t.Errorf("non expected HTTP Method %q received", r.Method)
		}

	},
	))

	inst, err := NewQualtricsSubscriptionWithClient("apikey", server.URL, server.Client())

	if err != nil {
		t.Errorf("test should not error out, but error %s received", err.Error())
	}

	id, err := inst.UpdateSubscription(&updateSubscription, &util.RequestContext{TraceHeaders:http.Header{}})

	if err != nil {
		t.Errorf("test should not error out, but error %s received", err.Error())
	}

	if id != "SUB_8pSh6Kbqjgvwynb" {
		t.Errorf("subscription ID should be SUB_8pSh6Kbqjgvwynb, %q received", id)
	}
}
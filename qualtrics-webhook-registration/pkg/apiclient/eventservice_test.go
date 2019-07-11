package apiclient

import (
	"github.com/kyma-incubator/connector-tools/qualtrics-webhook-registration/pkg/util"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)


func TestNewEventService(t *testing.T) {
	inst, err := NewEventService("http://www.apiclient.com", "qualtrics", 2000*time.Microsecond)

	if err != nil {
		t.Errorf("instance must not fail on creation, error %q", err.Error())
	}

	if inst.URL != "http://www.apiclient.com" {
		t.Errorf("url should be http://www.apiclient.com, but was %q", inst.URL)
	}

	if inst.ApplicationName != "qualtrics" {
		t.Errorf("ApplicationName should be \"qualtrics\", but was %q", inst.ApplicationName)
	}

	inst, err = NewEventService( "", "qualtrics", 000*time.Microsecond)

	if err == nil {
		t.Error("instance must fail on creation, but didn't")
	}

	inst, err = NewEventService( "http://www.apiclient.com", "", 000*time.Microsecond)

	if err == nil {
		t.Error("instance must fail on creation, but didn't")
	}
}

func TestEventService_GetActiveSubscriptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")
		file, _ := os.Open("../../testdata/GET_EventserviceSubscription.json")
		fileBytes, _ := ioutil.ReadAll(file)

		w.Write(fileBytes)

	},
	))

	inst, err := NewEventServiceWithClient(server.URL, "qualtrics", server.Client())

	if err != nil {
		t.Errorf("test should not error out, but error %s received", err.Error())
	}

	subs, err := inst.GetActiveSubscriptions(&util.RequestContext{TraceHeaders:http.Header{}})

	if err != nil {
		t.Errorf("test should not error out, but error %s received", err.Error())
	}

	if len(subs) != 3 {
		t.Errorf("expected response to contain 3 subscriptions, %+v received", subs)
	}

}
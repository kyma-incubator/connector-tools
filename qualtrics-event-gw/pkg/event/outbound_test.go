package event

import (
	"encoding/json"
	"github.com/kyma-incubator/connector-tools/qualtrics-event-gw/pkg/httphandler"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProcessor_ForwardEvent(t *testing.T) {

	ctx := &httphandler.RequestContext{TraceHeaders: http.Header{
		"X-Request-Id": []string{"ABCD"},
	}}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()

		var event KymaEvent
		decoder.Decode(&event)


		//check whether X-Request-Id header is included
		if r.Header.Get("X-Request-Id") != "ABCD" {
			t.Errorf("expected X-Request-Id header ABCD, but recieved %q", r.Header.Get("X-Request-Id"))
		}

		//Check whether Timestamp is included
		if event.EventTime == "" {
			t.Error("event time must be populated")
		}

		if event.Data == `{"target":"success"}` {
			w.Write([]byte(`{"status":"success"}`))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"status":"error"}`))
		}

	}))

	processor := NewOutboundProcessorWithCustomClient(server.URL,
		server.Client(),
		prometheus.NewSummary(prometheus.SummaryOpts{
			Name: "dummy",
		}))

	//Success
	resp, err := processor.ForwardEvent(&KymaEvent{
		EventType:        "test",
		SourceID:         "qualtrics",
		EventTypeVersion: "v1",
		Data:             `{"target":"success"}`,
	}, ctx)

	if err != nil {
		t.Errorf("success response expected, error received: %s", err.Error())
	}

	if respStatus, ok := resp["status"]; ok {
		statusValue := respStatus.(string)
		if statusValue != "success" {
			t.Errorf("different response expected, status should be success, but is %s", statusValue)
		}

	} else {
		t.Errorf("different response expected, it should contain status, but doesn't %+v", resp)
	}

	//Error

	resp, err = processor.ForwardEvent(&KymaEvent{
		EventType:        "test",
		SourceID:         "qualtrics",
		EventTypeVersion: "v1",
		Data:             `{"target":"error"}`,
	}, ctx)

	if err == nil {
		t.Errorf("error response expected, success received: %+v", resp)
	}

}

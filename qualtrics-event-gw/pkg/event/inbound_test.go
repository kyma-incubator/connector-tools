package event

import (
	"fmt"
	"github.com/kyma-incubator/connector-tools/qualtrics-event-gw/pkg/httphandler"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

type mockTopicMapper struct{}

func (m *mockTopicMapper) MapTopic(qualtricsTopicName string) (eventName string, eventVersion string, err error) {
	if qualtricsTopicName == "success" {
		return "success", "v1", nil
	} else if qualtricsTopicName == "ForwarderError" {
		return "ForwarderError", "v1", nil
	} else {
		return "", "", fmt.Errorf("it went south for %s", qualtricsTopicName)
	}
}

type mockEventForwarder struct{}

func (m *mockEventForwarder) ForwardEvent(evt *KymaEvent, ctx *httphandler.RequestContext) (map[string]interface{}, error) {

	if evt.EventType == "ForwarderError" {
		return nil, fmt.Errorf("something went wrong")
	} else {
		return map[string]interface{}{"result": "success"}, nil
	}
}

func Test_HandleRequest(t *testing.T) {

	topicMapper := mockTopicMapper{}
	eventForwarder := mockEventForwarder{}
	ctx := &httphandler.RequestContext{TraceHeaders: http.Header{
		"X-Request-Id": []string{"ABCD"},
	}}

	processor := InboundProcessor{
		TopicMapper:    &topicMapper,
		EventForwarder: &eventForwarder,
		SourceID:       "test",
	}

	// Test successes
	form := url.Values{}

	form.Set(topicField, "success")
	form.Set(dataField, `{"hello", "world"}`)

	req, _ := http.NewRequest(http.MethodPost, "http://www.kyma-project.io", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp := processor.HandleRequest(req, ctx)

	if !resp.IsSuccess {
		t.Errorf("success response expected, failure received")
	}

	if resp.ResponseCode != 200 {
		t.Errorf("status code 200 expected, %d received", resp.ResponseCode)
	}

	respMap := resp.Response.(map[string]interface{})

	if respMessage, ok := respMap["result"]; !ok {
		t.Error("Expected response to be map containing key \"result\", but it didn't")

	} else if respString := respMessage.(string); respString != "success" {
		t.Errorf("Expected response to be \"success\" but received %q", respString)
	}

	// Test wrong topic
	form = url.Values{}

	form.Set(topicField, "fail")
	form.Set(dataField, `{"hello", "world"}`)

	req, _ = http.NewRequest(http.MethodPost, "http://www.kyma-project.io", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp = processor.HandleRequest(req, ctx)

	if resp.IsSuccess {
		t.Errorf("error response expected, success received")
	}

	if resp.ResponseCode != 400 {
		t.Errorf("status code 400 expected, %d received", resp.ResponseCode)
	}

	// Test error from forwarder
	form = url.Values{}

	form.Set(topicField, "ForwarderError") //creates specific event type
	form.Set(dataField, `{"hello", "world"}`)

	req, _ = http.NewRequest(http.MethodPost, "http://www.kyma-project.io", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp = processor.HandleRequest(req, ctx)

	if resp.IsSuccess {
		t.Errorf("error response expected, success received")
	}

	if resp.ResponseCode != 500 {
		t.Errorf("status code 500 expected, %d received", resp.ResponseCode)
	}

}

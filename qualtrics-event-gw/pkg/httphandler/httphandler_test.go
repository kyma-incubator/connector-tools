package httphandler

import (
	"encoding/json"
	"github.com/prometheus/client_golang/prometheus"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type handlerMock struct {
	t *testing.T
}

func (h *handlerMock) HandleRequest(r *http.Request, ctx *RequestContext) *Response {

	if r.Header.Get(requestID) != "ABC" {
		h.t.Errorf("expected request to contain %s header ABC, but found %q", requestID,
			r.Header.Get(requestID))
	}

	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if string(body) == "success" {
		return &Response{
			IsSuccess: true,
			Response: map[string]string{
				"status": "success",
			},
			ResponseCode: 200,
		}
	} else {
		return &Response{
			IsSuccess:    false,
			Response:     JsonError{Message: "Oops"},
			ResponseCode: 500,
		}
	}
}

func TestHandlerContext_ServeHTTP(t *testing.T) {

	//Prepare
	var response map[string]string

	counter := prometheus.NewCounterVec(prometheus.CounterOpts{Name: "requests"}, []string{"responseCode"})
	handler := HandlerContext{
		Metrics: &Metrics{
			InFlightRequests:    prometheus.NewGauge(prometheus.GaugeOpts{Name: "in_flight"}),
			ServerResponseTimes: prometheus.NewSummary(prometheus.SummaryOpts{Name: "response_time"}),
			HTTPCalls2xx:        counter.With(prometheus.Labels{"responseCode": "2xx"}),
			HTTPCalls3xx:        counter.With(prometheus.Labels{"responseCode": "3xx"}),
			HTTPCalls4xx:        counter.With(prometheus.Labels{"responseCode": "4xx"}),
			HTTPCalls5xx:        counter.With(prometheus.Labels{"responseCode": "5xx"}),
		},
		NextHandler: &handlerMock{t: t},
	}

	//Success
	req, err := http.NewRequest(http.MethodPost, "/", strings.NewReader("success"))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set(requestID, "ABC")

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	decoder := json.NewDecoder(rr.Body)
	err = decoder.Decode(&response)

	if err != nil {
		t.Errorf("expected valid json response, error decoding %s", err.Error())
	}

	if status, ok := response["status"]; ok {
		if status != "success" {
			t.Errorf("expected status field to be success, received %q", status)
		}
	} else {
		t.Errorf("expected response to contain status field, received %1v", response)
	}

	if rr.Header().Get("Content-Type") != "application/json" {
		t.Errorf("expected Content-Type application/json, received %q", rr.Header().Get("Content-Type"))
	}

	//Error
	req, err = http.NewRequest(http.MethodPost, "/", strings.NewReader("fail"))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set(requestID, "ABC")

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

	decoder = json.NewDecoder(rr.Body)
	err = decoder.Decode(&response)

	if err != nil {
		t.Errorf("expected valid json response, error decoding %s", err.Error())
	}

	if status, ok := response["message"]; ok {
		if status != "Oops" {
			t.Errorf("expected status field to be Oops, received %q", status)
		}
	} else {
		t.Errorf("expected response to contain message field, received %1v", response)
	}

	if rr.Header().Get("Content-Type") != "application/json" {
		t.Errorf("expected Content-Type application/json, received %q", rr.Header().Get("Content-Type"))
	}

}

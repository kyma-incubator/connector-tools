package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestReady(t *testing.T) {
	req, err := http.NewRequest("GET", "/ready", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ready)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"code": success}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestHealthz(t *testing.T) {

	//Test working case

	eventURL, _ = url.Parse("http://kyma-project.io/docs/")

	req, err := http.NewRequest("GET", "/healthz", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(healthz)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"code": success}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	//Test not working case (wrong port)

	eventURL, _ = url.Parse("http://kyma-project.io:500/docs/")

	req, err = http.NewRequest("GET", "/healthz", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		temp := rr.Body.String()
		fmt.Println(temp)
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}

func TestSetLogLevel(t *testing.T) {

	if result := setLogLevel("foo"); result != "ERROR" {
		t.Errorf("wrong log level returned, expected ERROR, got %s", result)
	}

	if result := setLogLevel("ERROR"); result != "ERROR" {
		t.Errorf("wrong log level returned, expected ERROR, got %s", result)
	}

	if result := setLogLevel("eRRoR"); result != "ERROR" {
		t.Errorf("wrong log level returned, expected ERROR, got %s", result)
	}

	if result := setLogLevel("warn"); result != "WARN" {
		t.Errorf("wrong log level returned, expected WARN, got %s", result)
	}

	if result := setLogLevel("INFO"); result != "INFO" {
		t.Errorf("wrong log level returned, expected INFO, got %s", result)
	}

	if result := setLogLevel("debug"); result != "DEBUG" {
		t.Errorf("wrong log level returned, expected DEBUG, got %s", result)
	}

	if result := setLogLevel("TRACE"); result != "TRACE" {
		t.Errorf("wrong log level returned, expected TRACE, got %s", result)
	}
}

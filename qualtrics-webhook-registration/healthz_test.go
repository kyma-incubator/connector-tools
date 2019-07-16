package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHealthHandler_ServeHTTP(t *testing.T) {

	theBeginnig := time.Unix(0,0)
	handler := HealthHandler{
		LastSuccessfulSynchTime: &theBeginnig,
		RefreshIntervalSeconds: 60,
	}

	req := httptest.NewRequest(http.MethodGet, "http://www.kyma-project.io/", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code == 200 {
		t.Errorf("Expected health check to fail, as no sucessful refresh was indicated")
	}


	now := time.Now()

	handler.LastSuccessfulSynchTime = &now

	req = httptest.NewRequest(http.MethodGet, "http://www.kyma-project.io/", nil)
	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != 200 {
		t.Errorf("Expected health check to be sucessful, butr it wasn't")
	}

}

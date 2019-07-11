package main

import (
	"fmt"
	"net/http"
	"time"
)

type HealthHandler struct {
	LastSuccessfulSynchTime *time.Time
	RefreshIntervalSeconds  int64
}

func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	//assuming that one in the last 3 refreshes must have been successful
	minTime := time.Now().Add(time.Duration(-h.RefreshIntervalSeconds) * time.Second)



	if h.LastSuccessfulSynchTime.Before(minTime) {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf("{\"lastSuccessfulSynch\": \"%s\"}", h.LastSuccessfulSynchTime.String())))
}

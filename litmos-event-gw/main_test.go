package main

import (
	"bytes"
	"encoding/json"
	"github.com/kyma-incubator/connector-tools/litmos-event-gw/internal/config"
	"github.com/kyma-incubator/connector-tools/litmos-event-gw/internal/model/events"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

var httpClient *http.Client = &http.Client{}

func TestMain(m *testing.M) {
	ts := setUpEventPublishingService()
	defer ts.Close()

	setupGlobalConfig(ts)
	interrupt := make(chan os.Signal, 1)

	go func() {
		startGateway(interrupt)
	}()

	time.Sleep(100 * time.Millisecond) //waiting for server to start
	code := m.Run()

	interrupt <- os.Interrupt
	os.Exit(code)
}

func setupGlobalConfig(ts *httptest.Server) {
	config.GlobalConfig = &config.Opts{
		LogRequest:         false,
		AppName:            "litmos",
		EventPublishURL:    ts.URL,
		BaseTopic:          "litmos",
		InSecureSkipVerify: false,
	}
}

func setUpEventPublishingService() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		println("called...")
		m := make(map[string]string)
		m["x"] = "y"
		ba, _ := json.Marshal(m)
		_, _ = w.Write(ba)
	}))
}

func TestEventIngestionSuccess(t *testing.T) {
	g := NewGomegaWithT(t)
	le := &events.LitmosEvent{
		ID:      123,
		Created: "2019-05-06T01:13:19.533",
		Type:    "achievement.earned",
		Object:  "event",
		Data:    struct{ x string }{x: "y"},
	}
	b, _ := json.Marshal(le)

	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/events", bytes.NewReader(b))
	resp, err := httpClient.Do(req)

	g.Expect(err).Should(BeNil())
	g.Expect(resp).ShouldNot(BeNil())
	g.Expect(resp.StatusCode).To(Equal(200))
}

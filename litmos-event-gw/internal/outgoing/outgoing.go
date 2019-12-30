package outgoing

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/kyma-incubator/connector-tools/litmos-event-gw/internal/config"
	"github.com/kyma-incubator/connector-tools/litmos-event-gw/internal/logger"
	"github.com/kyma-incubator/connector-tools/litmos-event-gw/internal/model/events"
	"net/http"
	"time"
)

type EventForwarder struct {
	eventPublishURL string
	client          *http.Client
}

func NewEventForwarder() *EventForwarder {
	client := &http.Client{
		Transport: &http.Transport{
			DisableCompression:  false,
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 0,
			MaxConnsPerHost:     0,
			IdleConnTimeout:     30 * time.Second,
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: config.GlobalConfig.InSecureSkipVerify},
		},
	}

	return &EventForwarder{
		eventPublishURL: config.GlobalConfig.EventPublishURL,
		client:          client,
	}
}

func (e *EventForwarder) Forward(event *events.KymaEvent) (map[string]interface{}, error) {
	eventBytes, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, e.eventPublishURL, bytes.NewReader(eventBytes))
	if err != nil {
		return nil, err
	}

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	var respMap map[string]interface{}
	err = dec.Decode(&respMap)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		errMsg := fmt.Sprintf("unexpected response when publishing event %d (%s)", resp.StatusCode, resp.Status)
		logger.Logger.Error(errMsg)
		return respMap, fmt.Errorf(errMsg)
	}

	return respMap, nil

}

package events

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/kyma-incubator/connector-tools/gen-event-gw/pkg/config"

	log "github.com/sirupsen/logrus"
)

//EventForwarder -  todo
type EventForwarder struct {
	eventPublishURL *string
	client          *http.Client
}

//NewEventForwarder - todo
func NewEventForwarder() *EventForwarder {
	client := &http.Client{
		Transport: &http.Transport{
			DisableCompression:  false,
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 0,
			MaxConnsPerHost:     0,
			IdleConnTimeout:     30 * time.Second,
		},
	}

	return &EventForwarder{
		eventPublishURL: config.GlobalConfig.EventPublishURL,
		client:          client,
	}
}

//ForwardEvent - submit events to the kyma event bus
func (e *EventForwarder) ForwardEvent(event *KymaEvent) (map[string]interface{}, error) {

	eventBytes, err := json.Marshal(event)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	log.Printf("Will submit: %s", string(eventBytes))

	req, err := http.NewRequest(http.MethodPost, *e.eventPublishURL, bytes.NewReader(eventBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

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
		log.Println(errMsg)
		return respMap, fmt.Errorf(errMsg)
	}

	return respMap, nil
}

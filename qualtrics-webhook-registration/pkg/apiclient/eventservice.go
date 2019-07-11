package apiclient

import (
	"encoding/json"
	"fmt"
	"github.com/kyma-incubator/connector-tools/qualtrics-webhook-registration/pkg/util"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	kymaApiPath  = "v1/events/subscribed"
)

type EventService struct {
	URL             string
	ApplicationName string
	Client          *http.Client
}

type EventSubscription struct {
	EventType    string `json:"name"`
	EventVersion string `json:"version"`
}

type eventSubscriptionsResponse struct {
	EventsInfo []EventSubscription `json:"eventsInfo"`
}

func NewEventService(url string, applicationName string, timeout time.Duration) (*EventService, error) {
	return NewEventServiceWithClient(url, applicationName,
		&http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				MaxIdleConns:    10,
				MaxConnsPerHost: 20,
			},
		},
	)
}

func NewEventServiceWithClient(url string, applicationName string, client *http.Client) (*EventService, error) {

	if url == "" {
		return nil, fmt.Errorf("url must not be empty")
	}

	if client == nil {
		return nil, fmt.Errorf("client must not be nil")
	}

	if applicationName == "" {
		return nil, fmt.Errorf("applicationName must not be empty")
	}

	return &EventService{
		URL:             url,
		ApplicationName: applicationName,
		Client:          client,
	}, nil
}

func (e *EventService) GetActiveSubscriptions(ctx *util.RequestContext) ([]EventSubscription, error) {

	log.WithFields(ctx.GetLoggerFields()).Debug("Reading Kyma Subscriptions")

	url, err := url.Parse(fmt.Sprintf("%s/%s/%s", e.URL, e.ApplicationName, kymaApiPath))

	if err != nil {
		log.WithFields(ctx.GetLoggerFields()).Errorf("error assembling get subscription list url: %s", err.Error())
		return nil, fmt.Errorf("error assembling get subscription list url: %s", err.Error())
	}

	req, err := http.NewRequest(http.MethodGet, url.String(), nil)

	if err != nil {
		log.WithFields(ctx.GetLoggerFields()).Errorf("error assembling get subscription list request: %s",
			err.Error())
		return nil, fmt.Errorf("error assembling get subscription list request: %s", err.Error())
	}

	ctx.IncludeTraceHeaders(req.Header)

	resp, err := e.Client.Do(req)

	if err != nil {
		log.WithFields(ctx.GetLoggerFields()).Errorf("error getting subscription list: %s", err.Error())
		return nil, fmt.Errorf("error getting subscription list: %s", err.Error())
	}

	if resp.StatusCode != http.StatusOK {

		if log.GetLevel() != log.TraceLevel {
			log.WithFields(ctx.GetLoggerFields()).Errorf("error getting subscription list: %d (%s)",
				resp.StatusCode, resp.Status)
		} else {
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()

			if err != nil {
				bodyBytes = []byte(err.Error())
			}

			log.WithFields(ctx.GetLoggerFields()).Tracef("error getting subscription list: %d (%s): %s",
				resp.StatusCode, resp.Status, string(bodyBytes))
		}
		return nil, fmt.Errorf("error getting subscription list: %d (%s)",
			resp.StatusCode, resp.Status)
	}

	var respJson eventSubscriptionsResponse

	dec := json.NewDecoder(resp.Body)
	defer resp.Body.Close()

	err = dec.Decode(&respJson)

	if err != nil {
		log.WithFields(ctx.GetLoggerFields()).Errorf("error parsing subscription list: %s", err.Error())
		return nil, fmt.Errorf("error parsing subscription list: %s", err.Error())
	}
	return respJson.EventsInfo, nil
}

package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kyma-incubator/connector-tools/qualtrics-event-gw/pkg/httphandler"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const (
	contentType = "application/json"
)

//Processor serves as a client to forward events coming from various sources to Kyma
type OutboundProcessor struct {
	kymaEventURL              string
	client                    *http.Client
	responseTimeSummaryMetric prometheus.Summary
}

func NewOutboundProcessor(kymaEventURL string, responseTimeSummaryMetric prometheus.Summary, timeout time.Duration) OutboundProcessor {

	client := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			MaxConnsPerHost:     0,
			MaxIdleConnsPerHost: 0,
			MaxIdleConns:        100,
			IdleConnTimeout:     30 * time.Second,
			DisableCompression:  false,
		},
	}
	return NewOutboundProcessorWithCustomClient(kymaEventURL, client, responseTimeSummaryMetric)
}

func NewOutboundProcessorWithCustomClient(kymaEventURL string, clnt *http.Client, responseTimeSummaryMetric prometheus.Summary) OutboundProcessor {
	return OutboundProcessor{
		kymaEventURL:              kymaEventURL,
		client:                    clnt,
		responseTimeSummaryMetric: responseTimeSummaryMetric,
	}
}

//ForwardEvent sends event to configured event URL
func (p *OutboundProcessor) ForwardEvent(evt *KymaEvent, ctx *httphandler.RequestContext) (map[string]interface{}, error) {

	if evt.EventTime == "" {
		evt.EventTime = fmt.Sprintf(time.Now().Format(time.RFC3339))
	}

	evtBytes, err := json.Marshal(evt)

	if err != nil {
		return nil, err
	}

	if log.GetLevel() == log.TraceLevel {
		log.WithFields(
			ctx.GetLoggerFields(),
		).Tracef("Event sent to %s: %s)", p.kymaEventURL, string(evtBytes))
	}

	req, err := http.NewRequest(http.MethodPost, p.kymaEventURL, bytes.NewReader(evtBytes))

	if err != nil {
		log.WithFields(
			ctx.GetLoggerFields(),
		).Errorf("Error creating request for sending event %s", err.Error())
		return nil, err
	}

	ctx.IncludeTraceHeaders(req.Header)

	//Do request and measure execution
	startTime := time.Now()
	resp, err := p.client.Do(req)

	p.responseTimeSummaryMetric.Observe(float64(time.Since(startTime)))

	if err != nil {
		log.WithFields(
			ctx.GetLoggerFields(),
		).Errorf("Error sending event %s", err.Error())
		return nil, err
	}

	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)

	var respParsed map[string]interface{}

	err = dec.Decode(&respParsed)

	if err != nil {
		log.WithFields(
			ctx.GetLoggerFields(),
		).Errorf("Error parsing event response %s", err.Error())
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		log.WithFields(
			ctx.GetLoggerFields(),
		).Errorf("Error forwarding event: %d (%s)", resp.StatusCode, resp.Status)
		return respParsed, fmt.Errorf("error forwarding event: %d (%s)", resp.StatusCode, resp.Status)
	}

	return respParsed, nil
}

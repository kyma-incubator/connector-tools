package event

import (
	"fmt"
	"github.com/kyma-incubator/connector-tools/qualtrics-event-gw/pkg/httphandler"
	log "github.com/sirupsen/logrus"
	"net/http"
)

const (
	topicField = "Topic"
	dataField  = "MSG"
)

//Processor serves as a Server to receive events coming from Qualtrics relaying them to a Kyma Connector
type InboundProcessor struct {
	SourceID       string
	EventForwarder EventForwarder
	TopicMapper    TopicMapper
}

//HandleRequest takes in a qualtrics event and extracts the needed data
func (p *InboundProcessor) HandleRequest(r *http.Request, ctx *httphandler.RequestContext) *httphandler.Response {
	r.ParseForm()

	topic := r.FormValue(topicField)
	dataString := r.FormValue(dataField)

	if topic == "" || dataString == "" {

		log.WithFields(
			ctx.GetLoggerFields(),
		).Errorf("%s: %q and %s: %q is invalid", topicField, topic, dataField, dataString)

		return &httphandler.Response{
			ResponseCode: 400,
			IsSuccess:    false,
			Response: httphandler.JsonError{
				Message: fmt.Sprintf("invalid or missing %s or %s attribute",
					topicField, dataField),
			},
		}
	}

	kymaEventType, kymaEventVersion, err := p.TopicMapper.MapTopic(topic)

	if err != nil {
		log.WithFields(
			ctx.GetLoggerFields(),
		).Errorf("%s: %q is invalid: %q", topicField, topic, err.Error())

		return &httphandler.Response{
			ResponseCode: 400,
			IsSuccess:    false,
			Response: httphandler.JsonError{
				Message: fmt.Sprintf("%s: %q is invalid: %q", topicField, topic, err.Error()),
			},
		}
	}

	log.WithFields(
		ctx.GetLoggerFields(),
	).Debugf("Event received for topic: %q", topic)

	evt := KymaEvent{
		EventType:        kymaEventType,
		EventTypeVersion: kymaEventVersion,
		Data:             JSONString(dataString),
	}

	resp, err := p.EventForwarder.ForwardEvent(&evt, ctx)
	if err != nil {
		log.WithFields(
			ctx.GetLoggerFields(),
		).Errorf("Error forwarding Event: %s", err.Error())
		return &httphandler.Response{
			ResponseCode: 500,
			IsSuccess:    false,
			Response: httphandler.JsonError{
				Message: fmt.Sprintf("Error forwarding Event: %s", err.Error()),
			},
		}
	}

	return &httphandler.Response{
		ResponseCode: 200,
		IsSuccess:    true,
		Response:     resp,
	}
}

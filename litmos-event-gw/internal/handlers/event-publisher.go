package handlers

import (
	"github.com/kyma-incubator/connector-tools/litmos-event-gw/internal/incoming"
	"github.com/kyma-incubator/connector-tools/litmos-event-gw/internal/logger"
	"github.com/kyma-incubator/connector-tools/litmos-event-gw/internal/model/errors"
	"github.com/kyma-incubator/connector-tools/litmos-event-gw/internal/outgoing"
	"io/ioutil"
	"net/http"
)

type EventPublisher struct {
	eventForwarder *outgoing.EventForwarder
}

func NewEventPublisher() *EventPublisher {
	return &EventPublisher{eventForwarder: outgoing.NewEventForwarder()}
}
func (ep *EventPublisher) EventHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			logger.Logger.Errorw("error when parsing request", "error", err)
		}
		logger.Logger.Infow("event request body", "event", string(body))

		kymaEvent, err := incoming.Process(body)
		if err != nil {
			errors.HandleError(w, err, errors.InternalError)
			return
		}

		resp, err := ep.eventForwarder.Forward(kymaEvent)
		if err != nil {
			errors.HandleError(w, err, errors.InternalError)
			return
		}

		logger.Logger.Infow("Received response for event publishing", "response", resp)

		w.WriteHeader(http.StatusOK)
	})
}

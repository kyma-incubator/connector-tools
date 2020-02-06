package incoming

import (
	"encoding/json"
	"github.com/kyma-incubator/connector-tools/litmos-event-gw/internal/logger"
	"github.com/kyma-incubator/connector-tools/litmos-event-gw/internal/model/events"
)

func Process(requestBody []byte) (*events.KymaEvent, error) {
	le, err := to(requestBody)
	if err != nil {
		return nil, err
	}

	ke := events.Map(le)
	logger.Logger.Infow("kyma Event", "kyma-event", ke)

	return ke, nil
}

func to(requestBody []byte) (*events.LitmosEvent, error) {
	le := events.LitmosEvent{}
	err := json.Unmarshal(requestBody, &le)
	if err != nil {
		return nil, err
	}

	return &le, nil
}

package events

import (
	"github.com/gofrs/uuid"
	"github.com/kyma-incubator/connector-tools/litmos-event-gw/internal/config"
	"time"
)

type KymaEvent struct {
	SourceID         string     `json:"source-id"`
	EventType        string      `json:"event-type"`
	EventTypeVersion string      `json:"event-type-version"`
	EventID          string      `json:"event-id"`
	EventTime        string      `json:"event-time"`
	Data             interface{} `json:"data"`
}

type LitmosEvent struct {
	ID      int32       `json:"id"`
	Created string      `json:"created"`
	Type    string      `json:"type"`
	Object  string      `json:"object"`
	Data    interface{} `json:"data"`
}

func Map(litmosEvent *LitmosEvent) *KymaEvent {
	eventId, _ := generateEventID()
	eventType := config.GlobalConfig.BaseTopic + "." + litmosEvent.Type

	return &KymaEvent{
		SourceID:         config.GlobalConfig.BaseTopic,
		EventType:        eventType,
		EventTypeVersion: "v1",
		EventTime:        time.Now().Format(time.RFC3339),
		Data:             litmosEvent.Data,
		EventID:          eventId,
	}
}

func generateEventID() (string, error) {
	uid, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	return uid.String(), nil
}

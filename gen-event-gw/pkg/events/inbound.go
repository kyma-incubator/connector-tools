package events

import (
	"fmt"
	"time"

	"github.com/kyma-incubator/connector-tools/gen-event-gw/pkg/config"

	"github.com/gofrs/uuid"
	"github.com/tidwall/gjson"
)

//unmarshall will alter the property order - using custom type to prevent
type jsonString string

//MarshalJSON prevents string marshaling from kicking in
func (j jsonString) MarshalJSON() ([]byte, error) {
	return []byte(j), nil
}

//KymaEvent - data structure kyma supports
type KymaEvent struct {
	SourceID         *string    `json:"source-id"`
	EventType        string     `json:"event-type"`
	EventTypeVersion string     `json:"event-type-version"`
	EventID          string     `json:"event-id"`
	EventTime        string     `json:"event-time"`
	Data             jsonString `json:"data"`
}

//InBoundProcesser - consumes event and transforms to kyma event
func InBoundProcesser(reqBody []byte) (*KymaEvent, error) {

	//Fieldglass implemenation
	eventType := getEventType(string(reqBody))
	println("eventType:", eventType)

	if eventType == "" {
		return nil, fmt.Errorf("Could not determine an event type using  the query %s", *config.GlobalConfig.EventTypeQuery)
	}

	id, _ := generateEventID()

	//build kyma event
	var event = KymaEvent{
		SourceID:         config.GlobalConfig.AppName,
		EventTypeVersion: "v1",
		EventTime:        time.Now().Format(time.RFC3339),
		EventID:          id,
		EventType:        eventType,
		Data:             jsonString(reqBody),
	}

	return &event, nil
}

//GetFGEventType - determines the event type by checking the value of config.GlobalConfig.EventTypeQuery
//FieldGlass - "*.@xmlns"
func getEventType(json string) string {
	value := gjson.Get(json, *config.GlobalConfig.EventTypeQuery)
	return value.String()
}

func generateEventID() (string, error) {
	uid, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	return uid.String(), nil
}

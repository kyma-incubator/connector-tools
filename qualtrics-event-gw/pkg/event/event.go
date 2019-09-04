package event

import "github.com/kyma-incubator/connector-tools/qualtrics-event-gw/pkg/httphandler"

//JSONString represents a string containing a JSON Object that can be passed as string to avoid
//parsing overhead
type JSONString string

//MarshalJSON prevents string marhaling from kicking in
func (j JSONString) MarshalJSON() ([]byte, error) {
	return []byte(j), nil
}

func (j *JSONString) UnmarshalJSON(data []byte) error {

	*j = JSONString(data)
	return nil
}

func (j *JSONString) String() string {
	return string(*j)
}

type KymaEvent struct {
	EventType        string     `json:"event-type"`
	EventTypeVersion string     `json:"event-type-version"`
	EventTime        string     `json:"event-time"`
	Data             JSONString `json:"data"`
}

//EventForwarder represents an abstract event forwarder
type EventForwarder interface {
	ForwardEvent(evt *KymaEvent, ctx *httphandler.RequestContext) (map[string]interface{}, error)
}

package events

import (
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

//KymaEventProcesser - consumption and forwarding of events
type KymaEventProcesser struct {
	outbound *EventForwarder
}

//NewKymaEventProcesser -
func NewKymaEventProcesser() *KymaEventProcesser {
	return &KymaEventProcesser{
		outbound: NewEventForwarder(),
	}
}

//EventsHandler - route to handle events
func (k *KymaEventProcesser) EventsHandler(w http.ResponseWriter, r *http.Request) {

	// w.Header().Set("Content-Type", "application/json")

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Could not read data")
		return
	}

	kymaEvent, err := InBoundProcesser(reqBody)

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, err.Error())
		return
	}
	log.Printf("kymaEvent: %+v \n", kymaEvent)

	resp, err := k.outbound.ForwardEvent(kymaEvent)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "An error occurred during event publishing - %+v \n", err.Error())
		log.Printf("An error occurred during event publishing - %+v \n", err.Error())
	} else {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Event published - %+v \n", resp)
		log.Printf("Event published - %+v \n", resp)
	}

}

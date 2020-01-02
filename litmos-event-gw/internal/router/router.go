package router

import (
	gh "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/kyma-incubator/connector-tools/litmos-event-gw/internal/config"
	"github.com/kyma-incubator/connector-tools/litmos-event-gw/internal/handlers"
	"net/http"
	"os"
)

type Rtr struct {
	http.Handler
	*handlers.EventPublisher
}

func New() http.Handler {
	r := mux.NewRouter()
	ep := handlers.NewEventPublisher()

	r.HandleFunc("/events", ep.EventHandler()).Methods(http.MethodPost)

	return &Rtr{
		Handler:        applyLogging(r),
		EventPublisher: ep,
	}
}

func applyLogging(r *mux.Router) http.Handler {
	if !config.GlobalConfig.LogRequest {
		return r
	} else {
		return gh.LoggingHandler(os.Stdout, r)
	}
}

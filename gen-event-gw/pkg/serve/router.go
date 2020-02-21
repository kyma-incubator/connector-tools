package serve

import (
	"crypto/subtle"
	"net/http"

	"github.com/kyma-incubator/connector-tools/gen-event-gw/pkg/config"
	"github.com/kyma-incubator/connector-tools/gen-event-gw/pkg/events"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

//NewRouter - sets up routing
func NewRouter() error {
	log.Printf("Starting server with config...")
	log.Printf("AppName:%s, EventPublishURL:%s, UserName:%s", *config.GlobalConfig.AppName, *config.GlobalConfig.EventPublishURL, *config.GlobalConfig.UserName)

	router := mux.NewRouter().StrictSlash(true)

	t := events.NewKymaEventProcesser()

	eventsHandler := http.HandlerFunc(t.EventsHandler)

	router.Handle("/events", authMiddleware(eventsHandler)).Methods("POST")

	router.HandleFunc("/healthz", healthz)
	router.HandleFunc("/ready", ready)

	return http.ListenAndServe(":8080", router)

}

// middleware concept, as explained e.g. here: https://www.alexedwards.net/blog/making-and-using-middleware
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		passwordConfig := *config.GlobalConfig.Password
		usernameConfig := *config.GlobalConfig.UserName

		if passwordConfig == "" || usernameConfig == "" {
			log.Println("Basic Authentication has not been configured")
			next.ServeHTTP(w, r)
			return
		}

		user, pass, ok := r.BasicAuth()

		if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(usernameConfig)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(passwordConfig)) != 1 {
			log.Println("Unauthorized request has been attempted")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

//TODO
func healthz(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"code": success}`))

}

func ready(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"code": success}`))

}

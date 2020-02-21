package main

import (
	"github.com/kyma-incubator/connector-tools/gen-event-gw/pkg/serve"
	log "github.com/sirupsen/logrus"
)

func main() {
	err := serve.NewRouter()
	if err != nil {
		log.Fatal(err)
	}
}

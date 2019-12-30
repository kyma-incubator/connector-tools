package config

import (
	"flag"
	"log"
)

type Opts struct {
	LogRequest         bool
	AppName            string
	EventPublishURL    string
	BaseTopic          string
	InSecureSkipVerify bool
}

var GlobalConfig *Opts

func ParseFlags() {
	logRequest := flag.Bool("verbose", false, "log each incoming event request")
	appName := flag.String("app-name", "", "Application Name")
	eventPublishURL := flag.String("event-publish-url", "http://event-publish-service.kyma-system.svc.cluster.local:8080/v1/events", "URL to forward incoming events to Kyma Eventing")
	baseTopic := flag.String("base-topic", "litmos", "Base Topic defined in the Async API specification")
	insecureSkipVerify := flag.Bool("skip-tls-verify", false, "Skip TLS verify")
	flag.Parse()

	GlobalConfig = &Opts{
		LogRequest:         *logRequest,
		AppName:            *appName,
		EventPublishURL:    *eventPublishURL,
		BaseTopic:          *baseTopic,
		InSecureSkipVerify: *insecureSkipVerify,
	}

	if GlobalConfig.AppName == "" {
		log.Panic("Invalid configuration - Missing APP Name", "config", GlobalConfig)
	}

	log.Println("App config", "config", GlobalConfig)
}

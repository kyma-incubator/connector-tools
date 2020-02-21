package config

import (
	"flag"

	log "github.com/sirupsen/logrus"
)

//Config - configuration data for app
type Config struct {
	AppName         *string
	Password        *string
	UserName        *string
	EventTypeQuery  *string
	EventPublishURL *string
}

//GlobalConfig contains config
var GlobalConfig *Config

//Initialize the configuration
func init() {

	GlobalConfig = &Config{
		AppName:         flag.String("app-name", "", "Application Name"),
		Password:        flag.String("password", "", "Basic Auth Password"),
		UserName:        flag.String("username", "", "Basic Auth UserName"),
		EventTypeQuery:  flag.String("event-type-query", "", "The json query based on the get function of https://github.com/tidwall/gjson"),
		EventPublishURL: flag.String("event-publish-url", "http://event-publish-service.kyma-system.svc.cluster.local:8080/v1/events", "URL to forward incoming events to Kyma Eventing"),
	}

	flag.Parse()

	if *GlobalConfig.AppName == "" || *GlobalConfig.EventTypeQuery == "" {
		log.Fatalf("Invalid configuration - Missing APP Name or the Event Type Query")
	}

}

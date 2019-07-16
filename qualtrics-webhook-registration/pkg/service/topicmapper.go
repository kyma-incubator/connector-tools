package service

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

type mapperConfig struct {
	QualtricsTopic 		string `json:"qualtricsTopic"`
	KymaEventName       string `json:"kymaEventName"`
	KymaEventVersion    string `json:"kymaEventVersion"`
}



type Mapper struct {
	topic2EventType map[string]*mapperConfig
	eventType2Topic map[string]*mapperConfig
}

func NewTopicmapper(file string) (*Mapper, error) {
	configFile, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("error opening topic mapper config: %s", err.Error())
	}
	//noinspection GoUnhandledErrorResult
	defer configFile.Close()

	configData, err := ioutil.ReadAll(configFile)
	if err != nil {
		return nil, fmt.Errorf("error reading topic mapper config: %s", err.Error())
	}

	var mapperConfigList []mapperConfig

	if err := json.Unmarshal(configData, &mapperConfigList); err != nil {
		return nil, fmt.Errorf("error in topic mapper config json: %s", err.Error())
	}

	topic2EventType := make(map[string]*mapperConfig)
	eventType2Topic := make(map[string]*mapperConfig)


	//index configuration
	for i, _ := range mapperConfigList{
		topic2EventType[mapperConfigList[i].QualtricsTopic] = &mapperConfigList[i]
		eventType2Topic[fmt.Sprintf("%s.%s", mapperConfigList[i].KymaEventName,
			mapperConfigList[i].KymaEventVersion)] = &mapperConfigList[i]

	}


	return &Mapper{
		topic2EventType: topic2EventType,
		eventType2Topic: eventType2Topic,
	}, nil

}

func (m *Mapper) MapEventTypeVersionToTopic(eventType string, version string) (string, error) {

	if config, ok := m.eventType2Topic[fmt.Sprintf("%s.%s", eventType, version)]; ok {
		return config.QualtricsTopic, nil
	} else {
		log.Errorf("could not map eventType %q and version %q", eventType, version)
		return "", fmt.Errorf("could not map eventType %q and version %q", eventType, version)
	}
}

func (m *Mapper) MapTopicToEventTypeVersion(topic string) (string, string, error) {
	if config, ok := m.topic2EventType[topic]; ok {
		return config.KymaEventName, config.KymaEventVersion, nil
	} else {
		log.Errorf("could not map topic %q", topic)
		return "", "", fmt.Errorf("could not map topic %q", topic)
	}
}



package topicmapper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
)

type mapperConfig struct {
	QualtricsTopicRegex string `json:"qualtricsTopicRegex"`
	KymaEventName       string `json:"kymaEventName"`
	KymaEventVersion    string `json:"kymaEventVersion"`
}

type configCache struct {
	Regex            *regexp.Regexp
	KymaEventName    string
	KymaEventVersion string
}

type Mapper struct {
	cache []configCache
}

func New(file string) (*Mapper, error) {
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

	result := make([]configCache, len(mapperConfigList))

	for i, currentConfig := range mapperConfigList {

		regex, err := regexp.Compile(currentConfig.QualtricsTopicRegex)

		if err != nil {
			return nil, fmt.Errorf("error parsing \"qualtricsTopicRegex\": %q at config item %d: %s",
				currentConfig.QualtricsTopicRegex, i, err.Error())
		}
		result[i] = configCache{
			Regex:            regex,
			KymaEventName:    currentConfig.KymaEventName,
			KymaEventVersion: currentConfig.KymaEventVersion,
		}
	}

	return &Mapper{cache: result}, nil

}

func (m *Mapper) MapTopic(qualtricsTopicName string) (eventName string, eventVersion string, err error) {
	for _, cacheItem := range m.cache {

		if cacheItem.Regex.MatchString(qualtricsTopicName) {
			return cacheItem.KymaEventName, cacheItem.KymaEventVersion, nil
		}
	}

	return eventName, eventVersion, fmt.Errorf("no matching event Type found for topic %q", qualtricsTopicName)
}

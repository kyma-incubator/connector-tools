package config

import (
	"context"
	"encoding/json"
	errorWrap "github.com/kyma-incubator/connector-tools/open-connectors-job/pkg/error"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
)

// Default Timeout applied
const defaultTimeoutMills = 5000
const defaultLogLevel = "Error"

func readConfigFromReader(ctx context.Context, reader io.Reader, jobConfig *JobConfig) error {

	dec := json.NewDecoder(reader)

	if err := dec.Decode(jobConfig); err != nil {
		log.Errorf("error parsing config: %s", err.Error())
		return errorWrap.WrapError(err, "error parsing config")
	}

	if jobConfig.OpenConnectors.TimeoutMills == 0 {
		jobConfig.OpenConnectors.TimeoutMills = defaultTimeoutMills
	}

	if jobConfig.Compass.TimeoutMills == 0 {
		jobConfig.Compass.TimeoutMills = defaultTimeoutMills
	}

	if jobConfig.LogLevel == "" {
		jobConfig.LogLevel = defaultLogLevel
	}

	return nil
}


func NewFromFile(ctx context.Context, configFile *os.File) (*JobConfig, error) {
	log.Debugf("reading config from file: %s", configFile.Name())

	var jobConfig JobConfig
	if err := readConfigFromReader(ctx, configFile, &jobConfig); err != nil {
		log.Errorf("error parsing config: %s", err.Error())
		return nil, err
	}
	return &jobConfig, nil
}
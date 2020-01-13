package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/kyma-incubator/connector-tools/open-connectors-job/pkg/compass"
	"github.com/kyma-incubator/connector-tools/open-connectors-job/pkg/config"
	"github.com/kyma-incubator/connector-tools/open-connectors-job/pkg/controller"
	errorWrap "github.com/kyma-incubator/connector-tools/open-connectors-job/pkg/error"
	"github.com/kyma-incubator/connector-tools/open-connectors-job/pkg/open_connectors"
	log "github.com/sirupsen/logrus"
	"os"
)

func getConfigLocation() string {
	var configLocation string
	flag.StringVar(&configLocation, "config", "config/config.json", "location of the config file")
	flag.Parse()

	return configLocation
}

func getConfig(ctx context.Context, configLocation string) (*config.JobConfig, error) {
	file, err := os.Open(configLocation)
	if err != nil {
		return nil, errorWrap.WrapError(err, "error opening config file: %s", err.Error())
	}
	return config.NewFromFile(ctx, file)
}

func maskString(stringToMask string) string {
	if len(stringToMask) == 0 {
		return ""
	} else {
		return "**************"
	}
}

func setLogLevel(logLevel string) {
	level, err := log.ParseLevel(logLevel)
	if err != nil {
		log.Errorf("Log level %q not recognized, defaulting to \"Error\"", logLevel)
		level = log.ErrorLevel
	}
	log.SetLevel(level)
}

func printJobConfig(config *config.JobConfig) {
	fmt.Println("Starting job with the following configuration:")
	fmt.Printf("\tLog Level: %s\n", config.LogLevel)
	fmt.Printf("\tCompass Director URL: %s\n", config.Compass.DirectorURL)
	fmt.Printf("\tCompass Tenant ID: %s\n", config.Compass.TenantId)
	fmt.Printf("\tCompass OAuth2 Token URL: %s\n", config.Compass.TokenUrl)
	fmt.Printf("\tCompass Client ID: %s\n", maskString(config.Compass.ClientID))
	fmt.Printf("\tCompass Client Secret %s\n", maskString(config.Compass.ClientSecret))
	fmt.Printf("\tCompass Timeout in Milliseconds: %d\n", config.Compass.TimeoutMills)
	fmt.Printf("\tOpen Connectors Hostname: %s\n", config.OpenConnectors.Hostname)
	fmt.Printf("\tOpen Connectors Organization Secret: %s\n",
		maskString(config.OpenConnectors.OrganizationSecret))
	fmt.Printf("\tOpen Connectors User Secret: %s\n",
		maskString(config.OpenConnectors.UserSecret))
	fmt.Printf("\tOpen Connectors Timeout in Milliseconds: %d\n", config.OpenConnectors.TimeoutMills)
}

func startJob(ctx context.Context, config *config.JobConfig) error{

	httpClient, err := compass.CreateHttpClientTimeout(
		ctx,
		config.Compass.ClientID,
		config.Compass.ClientSecret,
		config.Compass.TokenUrl,
		config.Compass.TimeoutMills)

	if err != nil {
		return err
	}

	compassClient, err := compass.New(httpClient, config.Compass.DirectorURL, config.Compass.TenantId)
	if err != nil {
		return err
	}

	openConnectors, err := open_connectors.NewWithTimeout(ctx,
		config.OpenConnectors.OrganizationSecret,
		config.OpenConnectors.UserSecret,
		config.OpenConnectors.Hostname,
		config.OpenConnectors.TimeoutMills)
	if err != nil {
		return err
	}

	controller, err := controller.New(ctx,
		compassClient,
		openConnectors,
		config.OpenConnectors.Tags,
		config.Compass.ApplicationPrefix)
	if err != nil {
		return err
	}


	return controller.Synchronize(ctx)
}

func main() {
	ctx := context.Background()

	jobConfig, err := getConfig(ctx, getConfigLocation())
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	setLogLevel(jobConfig.LogLevel)

	printJobConfig(jobConfig)
	if err := startJob(ctx, jobConfig); err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	fmt.Println("Job successfully executed")
}

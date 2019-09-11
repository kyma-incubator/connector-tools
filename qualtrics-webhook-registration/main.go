package main

import (
	"flag"
	"fmt"
	"github.com/kyma-incubator/connector-tools/qualtrics-webhook-registration/pkg/apiclient"
	"github.com/kyma-incubator/connector-tools/qualtrics-webhook-registration/pkg/service"
	"github.com/kyma-incubator/connector-tools/qualtrics-webhook-registration/pkg/servicediscovery"
	"github.com/kyma-incubator/connector-tools/qualtrics-webhook-registration/pkg/util"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strings"
	"time"
)

func instantiateReconciler(kymaEventGatewayBaseURL string, applicationName string, timeout int64,
	qualtricsAPIKey string, qualtricsAPIBaseURL string, subscriptionUrl string, sharedKey string,
	configurationFileReference string) (service.ReconcilerType, error) {

	eventServiceAPIClient, err := apiclient.NewEventService(kymaEventGatewayBaseURL, applicationName,
		time.Duration(timeout)*time.Millisecond)

	if err != nil {
		log.Errorf("error creating event service API client: %s", err.Error())
		return nil, fmt.Errorf("error creating event service API client: %s", err.Error())
	}

	qualtricsAPIClient, err := apiclient.NewQualtricsSubscription(qualtricsAPIKey, qualtricsAPIBaseURL,
		time.Duration(timeout)*time.Millisecond)

	if err != nil {
		log.Errorf("error creating event qualtrics API client: %s", err.Error())
		return nil, fmt.Errorf("error creating qualtrics service API client: %s", err.Error())
	}

	topicMapper, err := service.NewTopicmapper(configurationFileReference)

	if err != nil {
		log.Errorf("error creating event topicmapper: %s", err.Error())
		return nil, fmt.Errorf("error creating event topicmapper: %s", err.Error())
	}

	reconciler, err := service.NewReconciler(qualtricsAPIClient, eventServiceAPIClient, topicMapper,
		sharedKey, subscriptionUrl)

	if err != nil {
		log.Errorf("error creating event reconciler: %s", err.Error())
		return nil, fmt.Errorf("error creating event reconciler: %s", err.Error())
	}

	return reconciler, err
}

func init() {
	log.SetOutput(os.Stdout)

}

func setLogLevel(logLevel string) string {
	switch strings.ToUpper(logLevel) {
	case "ERROR":
		log.SetLevel(log.ErrorLevel)
		return "ERROR"
	case "WARN":
		log.SetLevel(log.WarnLevel)
		return "WARN"
	case "INFO":
		log.SetLevel(log.InfoLevel)
		return "INFO"
	case "DEBUG":
		log.SetLevel(log.DebugLevel)
		return "DEBUG"
	case "TRACE":
		log.SetLevel(log.TraceLevel)
		return "TRACE"
	default:
		log.SetLevel(log.ErrorLevel)
		return "ERROR"
	}
}

func manageReconcileLoop(lastSuccessfulSynchPtr **time.Time, reconciler service.ReconcilerType, refreshInterval int64,
	refreshCycleQualtrics int64) {

	ctx := util.RequestContext{http.Header{}}
	var refreshCount int64 = 0
	for {
		err := reconciler.Reconcile(&ctx)

		if err != nil {
			log.WithFields(ctx.GetLoggerFields()).Errorf("Reconciling failed with error: %s", err.Error())
		} else {
			refreshCount++
			log.WithFields(ctx.GetLoggerFields()).Debug("Successfully reconciled")
			now := time.Now()
			**lastSuccessfulSynchPtr = now
		}
		if refreshCount >= refreshCycleQualtrics {
			err = reconciler.RefreshQualtricsState(&ctx)
			if err != nil {
				log.WithFields(ctx.GetLoggerFields()).Errorf("Error refreshing state from qualtrics: %s",
					err.Error())
			} else {
				refreshCount = 0
				log.WithFields(ctx.GetLoggerFields()).Debug("Successfully refreshed state from qualtrics")
			}
		}

		//sleep until next synch
		log.Debugf("Sleeping for %d seconds", refreshInterval)
		time.Sleep(time.Duration(refreshInterval) * time.Second)
	}
}

func main() {

	var kymaEventGatewayBaseURL string
	var labelSelector string
	var kubeConfig string
	var namespace string
	var applicationName string
	var timeout int64
	var qualtricsAPIKey string
	var qualtricsAPIBaseURL string
	var subscriptionUrl string
	var sharedKey string
	var configurationFileReference string
	var logLevel string
	var refreshInterval int64
	var refreshCycleQualtrics int64


	flag.StringVar(&labelSelector, "event-gateway-label-selector", "", "kubernetes label selector "+
		"used to identify standard event gateway service inside the kyma cluster (optional, as otherwise default will " +
		"be used)")
	flag.StringVar(&kubeConfig, "kubeconfig", "", "path pointing towards kubeconfig file to "+
		"be used for local testing")
	flag.StringVar(&namespace, "event-gateway-namespace", "kyma-integration", "namespace for "+
		"discovery of standard event gateway service inside the kyma cluster")

	flag.StringVar(&applicationName, "application-name", "qualtrics", "name of the kyma "+
		"application for qualtrics")
	flag.Int64Var(&timeout, "timeout-mil", 2000, "timeout in milliseconds used for all API Calls ")
	flag.StringVar(&qualtricsAPIKey, "qualtrics-apikey", "", "APIKey used for authenticating qualtrics "+
		"API Calls")
	flag.StringVar(&qualtricsAPIBaseURL, "qualtrics-base-url", "", "url pointing towards "+
		"qualtrics v3 API (without path)")
	flag.StringVar(&subscriptionUrl, "subscription-url", "", "url pointing towards the qualtrics gateway"+
		"which will be registered as endpoint for all qualtrics subscriptions")
	flag.StringVar(&sharedKey, "shared-key", "", "key used for authenticating qualtrics subscriptions "+
		"(HMAC)")
	flag.StringVar(&configurationFileReference, "config-file", "conf/topic-config.json", "reference to "+
		"json file containing topic to kyma event type / version mapping")
	flag.StringVar(&logLevel, "log-level", "ERROR", "log level that should be used (can be ERROR, WARN, INFO, DEBUG, TRACE). "+
		"Trace logs full events and requests ")
	flag.Int64Var(&refreshInterval, "refresh-interval", 60, "refresh interval in seconds for aligning kyma"+
		" and Qualtrics")
	flag.Int64Var(&refreshCycleQualtrics, "refresh-cycle", 0, "refresh cycle (in number of refresh intervals) "+
		"for refreshing qualtrics subscription state cache (0 means never)")
	flag.Parse()

	logLevel = setLogLevel(logLevel)


	fmt.Printf("Label Selector used for the kyma event gateway discovery (default is empty): %s\n",
		labelSelector)
	fmt.Printf("Kubeconfig file used for local testing (default is empty): %s\n", kubeConfig)
	fmt.Printf("Namespace used for the kyma event gateway discovery: %s\n", namespace)
	fmt.Printf("Kyma Application Name: %s\n", applicationName)
	fmt.Printf("Timeout in milliseconds for API calls: %d\n", timeout)
	fmt.Printf("Qualtrics API Key provided: %t\n", len(qualtricsAPIKey) > 0)
	fmt.Printf("Base URL for the Qualtrics API: %s\n", qualtricsAPIBaseURL)
	fmt.Printf("Shared Key for authentication provided: %t\n", len(sharedKey) > 0)
	fmt.Printf("Configuration file location: %s\n", configurationFileReference)
	fmt.Printf("Log Level: %s\n", logLevel)
	fmt.Printf("Refresh Interval: %d\n", refreshInterval)
	fmt.Printf("Refresh cycle Qualtrics: %d\n", refreshCycleQualtrics)

	//Discover Event Gateway based on Inputs

	var client *servicediscovery.KubernetesClient
	var err error
	//local testing
	if kubeConfig != "" {
		client, err = servicediscovery.InitOutOfCluster(kubeConfig)

		if err != nil {
			log.Fatalf("error instantiating kubernetes client: %s", err.Error())
		}
	} else {
		client, err = servicediscovery.InitInCluster()

		if err != nil {
			log.Fatalf("error instantiating kubernetes client: %s", err.Error())
		}
	}

	if labelSelector == "" {
		labelSelector = fmt.Sprintf("application=%s, heritage=Tiller-event-service", applicationName)
	}

	kymaEventGatewayBaseURL, err = client.DiscoverEventServiceURL(namespace, labelSelector)

	if err != nil {
		log.Fatalf("error discovering kyma event gateway base url: %s", err.Error())
	}

	fmt.Printf("Base URL for the kyma event Gateway: %s\n", kymaEventGatewayBaseURL)

	lastsucessfulSynch := time.Unix(0, 0)
	lastsucessfulSynchPtr := &lastsucessfulSynch


	//setup health checks
	healthHandler := HealthHandler{
		LastSuccessfulSynchTime: lastsucessfulSynchPtr,
		RefreshIntervalSeconds:  refreshInterval,
	}

	http.Handle("/healthz", &healthHandler)

	//start reconciler
	reconciler, err := instantiateReconciler(kymaEventGatewayBaseURL, applicationName, timeout, qualtricsAPIKey,
		qualtricsAPIBaseURL, subscriptionUrl, sharedKey, configurationFileReference)

	if err != nil {
		log.Fatalf("error instantiating reconciler: %s", err.Error())
	}

	go manageReconcileLoop(&lastsucessfulSynchPtr, reconciler, refreshInterval, refreshCycleQualtrics)


	//start health check
	log.Fatal(http.ListenAndServe(":8081", nil))

}

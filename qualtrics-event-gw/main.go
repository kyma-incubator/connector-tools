package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/kyma-incubator/connector-tools/qualtrics-event-gw/pkg/event"
	"github.com/kyma-incubator/connector-tools/qualtrics-event-gw/pkg/httphandler"
	"github.com/kyma-incubator/connector-tools/qualtrics-event-gw/pkg/servicediscovery"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/kyma-incubator/connector-tools/qualtrics-event-gw/pkg/hmac"
	"github.com/kyma-incubator/connector-tools/qualtrics-event-gw/pkg/topicmapper"
	log "github.com/sirupsen/logrus"
)

const (
	responseCodeLabel = "responseCode"
)

var (
	httpCalls = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "requests_processed_total",
		Help: "The total number of processed requests",
	},
		[]string{responseCodeLabel})

	clientResponseTimeSummaryMetric = promauto.NewSummary(prometheus.SummaryOpts{
		Name: "client_response_time",
		Help: "The summary of client response times",
	})

	serverResponseTimeSummaryMetric = promauto.NewSummary(prometheus.SummaryOpts{
		Name: "server_response_time",
		Help: "The summary of server response times",
	})

	inFlightCalls = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "in_flight_requests",
		Help: "The number of requests currently active",
	})

	eventURL *url.URL
)

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

func healthz(w http.ResponseWriter, r *http.Request) {
	//not much to check other then destination is reachable, otherwise forwarding does not make
	//a lot of sense
	w.Header().Set("Content-Type", "application/json")

	port := eventURL.Port()

	if port == "" {
		scheme := strings.ToUpper(eventURL.Scheme)

		if scheme == "HTTP" {
			port = "80"
		} else {
			port = "443"
		}
	}

	conn, err := net.DialTimeout("tcp", net.JoinHostPort(eventURL.Hostname(), port), 2*time.Second)
	if err != nil {
		// ignore error, we are in a readiness check
		resp, _ := json.Marshal(map[string]interface{}{
			"code":  "error",
			"error": err.Error(),
		})
		w.WriteHeader(500)
		w.Write(resp)

		return
	}
	conn.Close()

	//Success branch

	w.Write([]byte(`{"code": success}`))

}

func ready(w http.ResponseWriter, r *http.Request) {
	//always ready when server is started :-)
	w.Header().Set("Content-Type", "application/json")

	w.Write([]byte(`{"code": success}`))

}

func management() {

	serveMux := http.NewServeMux()

	server := http.Server{
		Addr:    ":8081",
		Handler: serveMux,
	}
	serveMux.HandleFunc("/healthz", healthz)
	serveMux.HandleFunc("/ready", ready)
	serveMux.Handle("/metrics", promhttp.Handler())

	server.ListenAndServe()

}

func main() {

	var internalEventURL string
	var labelSelector string
	var kubeConfig string
	var namespace string
	var applicationName string
	var hmacKey string
	var topicConfigLocation string
	var logLevel string
	var validateHMAC bool
	var timeoutMills int64


	flag.StringVar(&labelSelector, "event-gateway-label-selector", "", "kubernetes label selector "+
		"used to identify standard event gateway service inside the kyma cluster (optional, as otherwise default will " +
		"be used)")
	flag.StringVar(&kubeConfig, "kubeconfig", "", "path pointing towards kubeconfig file to "+
		"be used for local testing")
	flag.StringVar(&namespace, "event-gateway-namespace", "kyma-integration", "namespace for "+
		"discovery of standard event gateway service inside the kyma cluster")

	flag.StringVar(&applicationName, "applicationname", "qualtrics",
		"Name of the application that sends the events (in Kyma)")
	flag.StringVar(&hmacKey, "hmac-key", "", "shared key used to validate origin of incoming webhook calls (simple string)")
	flag.BoolVar(&validateHMAC, "hmac", false, "supplied hmac should be validated")
	flag.StringVar(&topicConfigLocation, "topic-conf", "conf/topic_config.json", "location of the topic mapper configuration file ")
	flag.StringVar(&logLevel, "log-level", "ERROR", "log level that should be used (can be ERROR, WARN, INFO, DEBUG, TRACE). "+
		"Trace logs full events and requests ")
	flag.Int64Var(&timeoutMills, "timeout", 2000, "timeout for forwarding requests to the event bus")

	flag.Parse()
	var err error

	//Discover Event Gateway based on Inputs
	var client *servicediscovery.KubernetesClient

	//local testing
	if kubeConfig != "" {
		client, err = servicediscovery.InitOutOfCluster(kubeConfig)

		if err != nil {
			log.Fatalf("error instantiating kubernetes client: %s", err.Error())
		}

	//In Cluster
	} else {
		client, err = servicediscovery.InitInCluster()

		if err != nil {
			log.Fatalf("error instantiating kubernetes client: %s", err.Error())
		}
	}

	if labelSelector == "" {
		labelSelector = fmt.Sprintf("application=%s, heritage=Tiller-event-service", applicationName)
	}

	internalEventURL, err = client.DiscoverEventServiceURL(namespace, labelSelector, applicationName)

	if err != nil {
		log.Fatalf("error discovering kyma event gateway base url: %s", err.Error())
	}


	eventURL, err = url.Parse(internalEventURL)

	if err != nil {
		log.Fatalf("Error parsing \"-kyma-eventurl\" with value %q: %s", internalEventURL, err.Error())
	}

	logLevel = setLogLevel(logLevel)

	server := http.Server{
		Addr: ":8080",
	}

	forwarder := event.NewOutboundProcessor(internalEventURL,
		clientResponseTimeSummaryMetric,
		time.Duration(timeoutMills)*time.Millisecond)

	topicMapper, err := topicmapper.New(topicConfigLocation)

	if err != nil {
		log.Fatalf("Setup of topic configuration failed with error: %s", err.Error())
	}

	var handler httphandler.Handler
	handler = &event.InboundProcessor{
		SourceID:       applicationName,
		EventForwarder: &forwarder,
		TopicMapper:    topicMapper,
	}

	// is hmac checking enabled?

	if validateHMAC {

		if hmacKey == "" {
			log.Fatalln("HMAC validation is turned on, but no key is supplied. Please supply key (-hmac-key)")
		}

		handler = &hmac.HMAC{
			Key:         hmacKey,
			NextHandler: handler,
		}
	}

	http.Handle("/", &httphandler.HandlerContext{
		NextHandler: handler,
		Metrics: &httphandler.Metrics{
			HTTPCalls2xx:        httpCalls.With(prometheus.Labels{responseCodeLabel: "2xx"}),
			HTTPCalls3xx:        httpCalls.With(prometheus.Labels{responseCodeLabel: "3xx"}),
			HTTPCalls4xx:        httpCalls.With(prometheus.Labels{responseCodeLabel: "4xx"}),
			HTTPCalls5xx:        httpCalls.With(prometheus.Labels{responseCodeLabel: "5xx"}),
			ServerResponseTimes: serverResponseTimeSummaryMetric,
			InFlightRequests:    inFlightCalls,
		},
	})


	fmt.Printf("Label Selector used for the kyma event gateway discovery (default is empty): %s\n",
		labelSelector)
	fmt.Printf("Kubeconfig file used for local testing (default is empty): %s\n", kubeConfig)
	fmt.Printf("Namespace used for the kyma event gateway discovery: %s\n", namespace)
	fmt.Printf("Server listening on: %q\n", server.Addr)
	fmt.Printf("Events are forwarded to: %q\n", internalEventURL)
	fmt.Printf("Events published in context of application: %q\n", applicationName)
	fmt.Printf("Validation of HMAC enabled: %t\n", validateHMAC)
	fmt.Printf("Topic Mapper Configuration Location: %s\n", topicConfigLocation)
	fmt.Printf("Log Level: %s\n", logLevel)
	fmt.Printf("Request timeout (milliseconds): %d\n", timeoutMills)
	go management()
	log.Fatal(server.ListenAndServe())

}

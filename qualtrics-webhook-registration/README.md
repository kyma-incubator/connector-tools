# qualtrics-webhook-registration

## About

This "Job" automates the registration of subscriptions via the Qualtrics API (https://api.qualtrics.com/reference#create-subscription). The determination which subscriptions need to be registered/de-registered leverages the Kyma Event API (https://kyma-project.io/assets/docs/1.3/application-connector/docs/assets/eventsapi.yaml) through the GET `/{application}/v1/events/subscribed` operation. The file `conf/topic-config.json` specifies the mapping from Kyma Event Types to Qualtrics topic subscriptions.

## Command Line Parameters

The application uses the following command line arguments to start: 


  - **application-name** (string) - name of the kyma application for qualtrics (default "qualtrics")
  - **config-file** (string) - reference to json file containing topic to kyma event type / version mapping (default "conf/topic-config.json")
  - **event-gateway-base-url** (string) - url pointing towards the service of the standard kyma event gateway (without path)
  - **log-level** (string) - log level that should be used (can be ERROR, WARN, INFO, DEBUG, TRACE). Trace logs full events and requests  (default "ERROR")
  - **qualtrics-apikey** (string) - APIKey used for authenticating qualtrics API Calls
  - **qualtrics-base-url** (string) - url pointing towards qualtrics v3 API (without path)
  - **refresh-cycle** (int) -refresh cycle (in number of refresh intervals) for refreshing qualtrics subscription state cache (0 means never)
  - **refresh-interval** (int) - refresh interval in seconds for aligning kyma and Qualtrics (default 60)
  - **shared-key** (string) - key used for authenticating qualtrics subscriptions (HMAC)
  - **subscription-url** (string) - url pointing towards the qualtrics gateway which will be registered as endpoint for all qualtrics subscriptions
  - **timeout-mil** (int) - timeout in milliseconds used for all API Calls  (default 2000)


## Build

```
docker build -t <username>/qualtrics-webhook-registration:<version> .
docker push <username>/qualtrics-webhook-registration:<version>
```

## Local Test

```
docker run -d -p 8081:8081  --rm <username>/qualtrics-webhook-registration:<version> \
-event-gateway-base-url http://qualtrics-event-service-external-api.kyma-integration.svc.cluster.local:8081 \
-application-name qualtrics -timeout-mil 2000 -qualtrics-apikey <your apikey> \
-qualtrics-base-url https://env.qualtrics.com \
-subscription-url https://<gw>.kymahost \
-shared-key <something secret> -config-file conf/topic-config.json -log-level TRACE \
-refresh-interval 60 -refresh-cycle 10
```




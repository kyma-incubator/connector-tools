# qualtrics-event-gw

## About

This is an HTTP gateway that takes in event/webhook notifications from Qualtrics (https://api.qualtrics.com/docs/webhooks) into Kyma. It is a Qualtrics specific alternative to the default Kyma Event Gateway component. Subscriptions in Qualtrics can be managed using the qualtrics-webhook-registtration component or directly using the Qualtrics API as described in [Event Registration on Qualtrics](#Event-Registration-on-Qualtrics).

## Command Line Parameters

The application uses the following command line arguments to start: 

  - **applicationname** (string) - Name of the application that sends the events (in Kyma) (default "qualtrics")
  - **hmac** - supplied hmac should be validated
  - **hmac-key** (string) - shared key used to validate origin of incoming webhook calls (simple string)
  - **kyma-eventurl** (string) - URL that incoming events will be pushed to in internal kyma format (default "http://event-bus-publish.kyma-system.svc.cluster.local:8080/v1/events")
  - **log-level** (string) - log level that should be used (can be ERROR, WARN, INFO, DEBUG, TRACE). Trace logs full events and requests  (default "ERROR")
  - **timeout** (int) - timeout for forwarding requests to the event bus (default 2000)
  - **topic-conf** (string) - location of the topic mapper configuration file (default "conf/topic_config.json")

## Build

```
docker build -t <username>/qualtrics-event-gw:<version> .
docker push <username>/qualtrics-event-gw:<version>
```

## Local Test

```
docker run -d -p 8080:8080 -p 8081:8081 -v "$(pwd)"/conf/:/conf --rm <username>/qualtrics-event-gw:<version> -kyma-eventurl http://httpbin.org/anything
```

Then you can test:

```
curl -X POST \
  http://localhost:8080 \
  -H 'Content-Type: application/x-www-form-urlencoded' \
  -d 'HMAC=5db5cd1e0200e9ce6831ea0f1924beb97f136939e367c53f5af8f7d79c6727af495cb90d6c744fe72ab45efd9f9f2ba7c15b13fdf6565af3d5a4f50ed2c0c7d4&MSG=%7B%22Status%22%3A%22Complete%22%2C%22SurveyID%22%3A%22SV_22VlHYNeldrrrkp%22%2C%22ResponseID%22%3A%22R_1l9mcVuXubb4aGm%22%2C%22CompletedDate%22%3A%222019-06-26%2013%3A28%3A11%22%2C%22BrandID%22%3A%22sapdevelopment%22%7D&Topic=sapdevelopment.surveyengine.completedResponse.SV_22VlHYNeldrrrkp&undefined='
```

## Kyma



After deployment you can import a Grafana Dashboard: `dashboard/Qualtrics Event GW Dashboard.json`.


## Event Registration on Qualtrics

```
curl -X POST \
  https://<qualtrics host>/API/v3/eventsubscriptions/ \
  -H 'X-API-TOKEN: ...' \
  -H 'Content-Type: application/json' \
  -d '{
    "topics": "*",
    "publicationUrl": "https://qualtrics-event.<kymahost>",
    "encrypt": false,
    "sharedKey": "sharedkey"
}'
```

For details see: https://api.qualtrics.com/reference#create-subscription

## Load Generation:

Install loadtest: https://www.npmjs.com/package/loadtest

To further simplify `testing/loadtest.sh` can be executed which supports "warmup".
```
loadtest -c 1 --rps 10 \
    https://<event-gw-url> \
       -T 'application/x-www-form-urlencoded' \
       -P 'HMAC=5db5cd1e0200e9ce6831ea0f1924beb97f136939e367c53f5af8f7d79c6727af495cb90d6c744fe72ab45efd9f9f2ba7c15b13fdf6565af3d5a4f50ed2c0c7d4&MSG=%7B%22Status%22%3A%22Complete%22%2C%22SurveyID%22%3A%22SV_22VlHYNeldrrrkp%22%2C%22ResponseID%22%3A%22R_1l9mcVuXubb4aGm%22%2C%22CompletedDate%22%3A%222019-06-26%2013%3A28%3A11%22%2C%22BrandID%22%3A%22sapdevelopment%22%7D&Topic=sapdevelopment.surveyengine.completedResponse.SV_22VlHYNeldrrrkp&undefined='
```



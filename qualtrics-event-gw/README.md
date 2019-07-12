# qualtrics-event-gw

## Build

```
docker build -t <username>/qualtrics-event-gw:<version> .
docker push <username>/qualtrics-event-gw:<version>
```

## Local Test

```
docker run -d -p 8080:8080 -p 8081:8081 --rm <username>//qualtrics-event-gw:<version> -kyma-eventurl http://httpbin.org/anything
```

Then you can test:

```
curl -X POST \
  http://localhost:8080 \
  -H 'Content-Type: application/x-www-form-urlencoded' \
  -d 'HMAC=5db5cd1e0200e9ce6831ea0f1924beb97f136939e367c53f5af8f7d79c6727af495cb90d6c744fe72ab45efd9f9f2ba7c15b13fdf6565af3d5a4f50ed2c0c7d4&MSG=%7B%22Status%22%3A%22Complete%22%2C%22SurveyID%22%3A%22SV_22VlHYNeldrrrkp%22%2C%22ResponseID%22%3A%22R_1l9mcVuXubb4aGm%22%2C%22CompletedDate%22%3A%222019-06-26%2013%3A28%3A11%22%2C%22BrandID%22%3A%22sapdevelopment%22%7D&Topic=sapdevelopment.surveyengine.completedResponse.SV_22VlHYNeldrrrkp&undefined='
```

## Kyma

```
kubectl create cm -n qualtrics qualtrics-event-gw-config --from-file conf/
kubectl apply -f kubernetes/package.yaml -n <yournamespace>
kubectl apply -f kubernetes/service_monitor.yaml
```

After that import the `kubernetes/Qualtrics Event GW Dashboard.json` file into Grafana.

## Adaptations for cluster required

In file kubernetes/package.yaml` see things in curly braces:

```
apiVersion: v1
kind: Secret`
metadata:
  name: qualtrics-event-gw-hmac
type: Opaque
data:
  sharedkey: {base64key}
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: qualtrics-event-gw-params
  labels:
    app: qualtrics-event-gw
data:
  application-name: "{app-name}"
  timeout-mil: "2000"
  kyma-eventurl: "http://event-bus-publish.kyma-system.svc.cluster.local:8080/v1/events"
```

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


Analyze load w/o Grafana: 
```
kubectl top pods -n qualtrics -l app=qualtrics-event-gw --containers
```

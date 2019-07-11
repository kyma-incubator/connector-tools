# qualtrics-webhook-register

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

## Kyma

```
kubectl create cm -n qualtrics qualtrics-webhook-registration-config --from-file conf/
kubectl apply -f kubernetes/package.yaml -n <yournamespace>

```

## Adaptations for cluster required

In file `kubernetes/package.yaml` see things in curly braces:

```

---
apiVersion: v1
kind: ConfigMap
metadata:
  name: qualtrics-webhook-registration-params
  labels:
    app: qualtrics-webhook-registration
data:
  event-gateway-base-url: "http://{app-name}-event-service-external-api.kyma-integration.svc.cluster.local:8081"
  application-name: "{app-name}"
  timeout-mil: "2000"
  qualtrics-base-url: "https://env.qualtrics.com"
  subscription-url: "{event service name}"
  refresh-interval: "60"
  refresh-cycle: "10"
---
apiVersion: v1
kind: Secret
metadata:
  name: qualtrics-webhook-registration-secret
type: Opaque
data:
  sharedkey: {base64key}
  qualtrics-apikey: {qualtricsapikey, see https://api.qualtrics.com/docs/api-key-authentication}

```


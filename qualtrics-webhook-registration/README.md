# qualtrics-webhook-register

## Build

```
docker build -t <username>/qualtrics-webhook-registration:<version> .
docker push <username>/qualtrics-webhook-registration:<version>
```

## Local Test

```
docker run -d -p 8080:8080  --rm <username>/qualtrics-webhook-registration:<version> \
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


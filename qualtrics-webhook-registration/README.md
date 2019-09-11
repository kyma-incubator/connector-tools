# qualtrics-webhook-registration

## About

This "Job" automates the registration of subscriptions via the Qualtrics API (https://api.qualtrics.com/reference#create-subscription). The determination which subscriptions need to be registered/de-registered leverages the Kyma Event API (https://kyma-project.io/assets/docs/1.3/application-connector/docs/assets/eventsapi.yaml) through the GET `/{application}/v1/events/subscribed` operation. The file `conf/topic-config.json` specifies the mapping from Kyma Event Types to Qualtrics topic subscriptions.

## Command Line Parameters

The application uses the following command line arguments to start: 


  - **event-gateway-label-selector** (string) - kubernetes label selector used to identify standard event gateway service inside the kyma cluster (optional, as otherwise default will be used)
  - **event-gateway-namespace** (string) - namespace for discovery of standard event gateway service inside the kyma cluster (default "kyma-integration")
  - **kubeconfig** (string) - path pointing towards kubeconfig file to be used for local testing
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

The hostname of the standard Event Gateway is determined using Kubernetes service discovery based on labels. Hence a cluster internal url is going to be resolved. To enable local testing, the service needs to be made available prior to testing using a port forward:

```
export QUALTRICS_APPNAME=<application-name>
export QUALTRICS_SVC=$(kubectl get svc -n kyma-integration -l application=$QUALTRICS_APPNAME,heritage=Tiller-event-service -o jsonpath="{range .items[*]}{@.metadata.name}{end}")
kubectl port-forward -n kyma-integration svc/$QUALTRICS_SVC 8081
```

Now this needs to be mapped to the "right hostname". To do that, put the output of `echo "127.0.0.1     $QUALTRICS_SVC.kyma-integration.svc.cluster.local"` into your `/etc/hosts` file.

Now you can run

```
go run main.go healthz.go --kubeconfig <your kubeconfig> \
-application-name qualtrics -timeout-mil 60000 -qualtrics-apikey <your apikey> \
-qualtrics-base-url https://env.qualtrics.com \
-subscription-url <https://<gw>.kymahost> \
-shared-key <something secret> -config-file conf/topic-config.json -log-level TRACE \
-refresh-interval 60 -refresh-cycle 10
```

Kubeconfig cannot rely on auth provider. It is recommended to download the file from your kyma cluster.

## Kubernetes

If you deploy this gateway inside a kyma cluster (as it is intendend), you must ensure that the right cluster roles and role bindings are in place. The service requires `list` access to the `services` resource. The below example specifies such a Cluster Role and maps it to the default service account of a namespace.

```
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: qualtrics-role
  labels:
    app: qualtrics
rules:
  - apiGroups: ["*"]
    resources: ["services"]
    verbs: ["list"]
---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: qualtrics-rolebinding
  labels:
    app: qualtrics
subjects:
  - kind: User
    name: system:serviceaccount:<namespace>:default
    apiGroup: rbac.authorization.k8s.io
roleRef:
  kind: ClusterRole
  name: qualtrics-role
  apiGroup: rbac.authorization.k8s.io

```


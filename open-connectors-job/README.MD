# open-connectors-job

## About

This is a batch job that synchronizes the state between one service instance of [SAP Cloud Platform Open Connectors](https://help.sap.com/viewer/product/OPEN_CONNECTORS/Cloud/en-US) and a [Compass](https://github.com/kyma-incubator/compass) runtime. It is meant to be scheduled by a higher level automation. In compass terminology the job acts as an Integration System.

## Retrieve Compass Integration System Credentials

Currently Compass does not offer a UI for managing [Integration Systems](https://github.com/kyma-incubator/compass/blob/master/docs/integration-system/integration-systems.md). Hence the API calls should be made using a GraphQl client.

Create Integration System:

```
mutation {
  createIntegrationSystem (in: {name: "your name", description:"your desc"}) {
    id
  }
}
```

With the ID of the previously created integration system, one can generate the client credentials:


```
mutation {
  generateClientCredentialsForIntegrationSystem (id: "integration system id") {
    auth {
      credential {
        ... on OAuthCredentialData {
          clientId
          clientSecret
          url
        }
      }
    }
  }
}
```

## Command Line Parameters

The application uses the following command line arguments to start: 
  - **config** (string) - location of the config file (default "/config/config.json")
  
The config file should follow this structure as a json file:

```
type JobConfig struct {
	LogLevel       string               `json:"logLevel"` //Log level
	Compass        CompassConfig        `json:"compass"`
	OpenConnectors OpenConnectorsConfig `json:"openConnectors"`
}

type CompassConfig struct {
	DirectorURL       string `json:"directorURL"`       //Url pointing to OAuth2 secured Director URL of Compass
	TenantId          string `json:"tenantId"`          //Url pointing to OAuth2 secured Director URL of Compass
	TokenUrl          string `json:"tokenUrl"`          //Url pointing towards client credentials token provider
	ClientID          string `json:"clientID"`          //Client id of the Compass Integration System
	ClientSecret      string `json:"clientSecret"`      //Client Secret of the Compass Integration System
	ApplicationPrefix string `json:"applicationPrefix"` //Prefix assigned to applications created
	TimeoutMills      int    `json:"timeoutMills"`      //Timeout in milliseconds for Compass API Calls
}

type OpenConnectorsConfig struct {
	Hostname           string   `json:"hostname"`           //Hostname of the Open Connectors instance
	OrganizationSecret string   `json:"organizationSecret"` //Organization Secret of the Open Connectors instance
	UserSecret         string   `json:"userSecret"`         //User Secret of the Open Connectors instance
	TimeoutMills       int      `json:"timeoutMills"`       //Timeout in milliseconds for Open Connectors API Calls
	Tags               []string `json:"tags"`               //Tags to select instances in open connectors, empty means
	//all
}
```

An example is:

```
{
  "logLevel": "Trace",
  "compass":{
    "directorURL":"https://compass-gateway-auth-oauth.mycompass.org/director/graphql",
    "tenantId":"3e64ebae-38b5-46a0-b1ed-9ccee153a0ae",
    "tokenUrl":"https://oauth2.mycompass.org/oauth2/token",
    "clientID":"clientid",
    "clientSecret":"clientsecret",
    "applicationPrefix":"app-prefix",
    "timeoutMills":6000
  },
  "openConnectors":{
    "hostname":"api.openconnectors.ext.hanatrial.ondemand.com",
    "organizationSecret":"orgsecret",
    "userSecret":"usersecret",
    "timeoutMills":6000
  }
}

```


## Build (Docker)

```
docker build -t <username>/open-connectors-job:<version> .
docker push <username>/qualtrics-event-gw:<version>
```

## Local Test

To test locally 
```
go run main.go -config <path to config>
```


## Kubernetes

If you deploy this job inside a Kubernetes cluster, you should schedule this as a [Cronjob](https://kubernetes.io/docs/concepts/workloads/controllers/cron-jobs/).

A sample manifest would be:

```
apiVersion: batch/v1beta1
kind: CronJob
metadata:
  name: open-connectors-job
  labels:
    app: open-connectors-job
spec:
  schedule: "*/1 * * * *"
  concurrencyPolicy: Forbid
  jobTemplate:
    metadata:
      labels:
        app: open-connectors-job
    spec:
      template:
        metadata:
          labels:
            app: open-connectors-job
          annotations:
            sidecar.istio.io/inject: "false"
        spec:
          restartPolicy: Never
          containers:
            - name: open-connectors-job
              image: <your docker image>
              imagePullPolicy: Always
              command: ['./open-connectors-job', '-config', '/config/config.json']
              volumeMounts:
                - name: config
                  mountPath: "/config"
                  readOnly: true
          volumes:
            - name: config
              secret:
                secretName: configfiles

---
apiVersion: v1
kind: Secret
metadata:
  name: configfiles
  labels:
    app: open-connectors-job
type: Opaque
data:
  config.json: <base64 version of config file>
```



#  Connector Tools

`Connector-tools` contain projects required to build a `Kyma`connector bundle. It is meant to be used for connecting an application to kyma having no built-in Kyma pairing logic provided.

The repository does not contain the actual bundle definition, it is providing the source code for the different tools only.
Samples for actual bundles will be provided soon.

## Components of the Bundle

![Architecture Diagram](assets/architecture.svg)

|Component|Description|
|---|---|
|MQTT Event Bridge|NodeJS application which consumes MQTT messages, transforms them into a JSON payload and forwards them to the Kyma event bus|
|OAuth2 Server|NodeJS based OAuth2 server with configurable client secrets for the client_credentials grant used for MQTT message authentication|
|API Registration Job|Golang application to register configured ODATA services and events types to a configured Kyma Application. It will check if the API will respond before it gets registered.|

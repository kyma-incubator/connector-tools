# Kyma Generic Event Gateway
This app consumes webhooks from external applications and transforms their data into a kyma event.


## Parameters
| Name              | Description                                                              | Type   | Info                                                                                   |
| ----------------- | :----------------------------------------------------------------------- | :----- | :------------------------------------------------------------------------------------- |
| app-name          | The name of the registered application in kyma                           | string |                                                                                        |
| username          | Username to authenticate                                                 | string |                                                                                        |
| password          | Password to authenticate                                                 | string |                                                                                        |
| event-type-query  | The query string used to determine the event type from the received data | string | For examples refer to the get function of https://github.com/tidwall/gjson             |
| event-publish-url | The kyma event bus url                                                   | string | Defaults to: http://event-publish-service.kyma-system.svc.cluster.local:8080/v1/events |

## Example Usage
go run cmd/gen-event-gw/main.go --app-name=myapp --username=testuser --password=testpw --event-type-query="param2.eventtype" --event-publish-url=http://httpbin.org/anything

curl -X POST -H "Content-Type: application/json" --user testuser:testpw http://localhost:8080/events --data '{"param1":"xyz","param2": {"eventtype": "myTestEvent"}}'

## Docker
docker run -p 8080:8080 -d jcawley5/gen-event-gw --app-name=myapp --username=testuser --password=testpw --event-type-query="param2.eventtype" --event-publish-url=http://httpbin.org/anything
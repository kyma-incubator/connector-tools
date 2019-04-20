# API Registration Job

## Configuration
 Will read API configuration from a file `files/apis.json`.
 The file needs to be in a format like this:
 ```
 [
      {
        "path": "base path to the API",
        "name": "Name of the API as shown in the Service Catalog",
        "description": "Description of the API as shown in the Service Catalog"
      },
      {
        ...
      }
    ]
```

Will read event types from a file `files/events.json`
 The file needs to be in a format like this:
 ´´´
 {
    "name": "Name as shown in the Service Catalog",
    "provider": "Provider as shown in the Service Catalog",
    "description": "Description as shown in the Service Catalog",
    "events": {
        //Event types in AsyncAPI spec//
    }
 }
´´´

Furthermore the following ENV variables needs to be supplied:
- APPLICATION_NAME: Name of the Application at which the APIs should get registered to
- SYSTEM_URL: Base URL to the system whose APIs will be registered
- BASIC_USER: Basic Auth Username used for protecting the APIs
- BASIC_PASSWORD: Basic Auth Password used for protecting the APIs
- PROVIDER_NAME: Provider name as shown in the Service Catalog 
- PRODUCT_NAME: Product name of the connected system as shown in the Service Catalog

## Development

To run it local, execute:
```
export APPLICATION_NAME=
export SYSTEM_URL=
export BASIC_USER=
export BASIC_PASSWORD=
export PROVIDER_NAME=
export PRODUCT_NAME=

go build -o registration_app
./registration_app
```

## Testing

```
export APPLICATION_NAME=
export SYSTEM_URL=
export BASIC_USER=
export BASIC_PASSWORD=
export PROVIDER_NAME=
export PRODUCT_NAME=

go test
```

## Running local with docker

Either use latest master branch revision:
```
docker run -p8080:8080 -e APPLICATION_NAME=myApp -e OAUTH_URL=http://oauth2server:8080 --link oauth2server eu.gcr.io/kyma-project/incubator/mqtt-event-bridge:master
```

Or local revision:
```
docker build -t mqtt-event-bridge:latest .
docker run -p8080:8080 -e APPLICATION_NAME=myApp -e OAUTH_URL=http://oauth2server:8080 --link oauth2server mqtt-event-bridge:latest
```
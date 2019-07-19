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
- APP_KIND: Determines the API type and authentication mechanism used. 
  - `odata-with-basicauth` (default)
  - `rest-with-apikey`

## Development

To run it local, execute:
```
export APPLICATION_NAME=
export SYSTEM_URL=
export BASIC_USER=
export BASIC_PASSWORD=
export PROVIDER_NAME=
export PRODUCT_NAME=
export APP_KIND=

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
export APP_KIND=

go test
```

## Build locally

```bash
docker build -t api-registration-job:latest .
```

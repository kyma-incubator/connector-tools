# MQTT Event Bridge

## Configuration
The following ENV variables needs to be supplied:
- APPLICATION_NAME: value for the source-id attribte of the events, usually the identifier of the connected system
- OAUTH_URL: URL to the oauth server used for validating the OAuth2 token attached to the message

## Development

To run it local, execute:
```
export APPLICATION_NAME=
export OAUTH_URL=

npm install
npm start
```

## Testing

```
export MQTT_URL=ws://localhost:8080
export OAUTH_URL=http://localhost:9000
export OAUTH_CLIENT_ID=XXX
export OAUTH_CLIENT_SECRET=YYY

npm install
npm test
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

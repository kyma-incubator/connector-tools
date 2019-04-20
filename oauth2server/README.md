# OAuth2 Server

## Configuration
The following ENV variables needs to be supplied:
- OAUTH_CLIENT_ID = Pre-configured clientId for client_credentials grant
- OAUTH_CLIENT_SECRET = Pre-configured clientSecret for client_credentials grant

## Development

To run it local, execute:
```
export OAUTH_CLIENT_ID=
export OAUTH_CLIENT_SECRET=

npm install
npm start
```

## Running local with docker

Either use latest master branch revision:
```
docker run --name oauth2server -p9000:8080 -e OAUTH_CLIENT_ID=XXX -e OAUTH_CLIENT_SECRET=YYY eu.gcr.io/kyma-project/incubator/oauth2server:master
```

Or local revision:
```
docker build -t oauth2server:latest .
docker run --name oauth2server -p9000:8080 -e OAUTH_CLIENT_ID=XXX -e OAUTH_CLIENT_SECRET=YYY oauth2server:latest
```
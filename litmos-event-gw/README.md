# Overview

## Test Locally

```shell script

#build and start docker container
make run-docker-local

#send an event
curl -X POST -d @./assets/test-event.json http://localhost:8080/events -v

```
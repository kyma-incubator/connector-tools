[![codecov](https://codecov.io/gh/abbi-gaurav/connector-tools/branch/master/graph/badge.svg)](https://codecov.io/gh/abbi-gaurav/connector-tools)

# Overview

This is a HTTP gateway that takes webhook payloads from [Litmos](https://support.litmos.com/hc/en-us/articles/360022948994-Webhooks) and converts them to Kyma events.

Webhook configurations in Litmos are managed manually.
 
## Command line parameters

* **--verbose** (boolean) - LOG each incoming request headers. Useful for debugging.
* **--app-name** (string) - Name of the Kyma application to which Litmos tenant is bound. **REQUIRED**
* **--event-publish-url** (string) - Kyma internal service URL to which Kyma events will be published.
* **--base-topic** (string) - Base topic name as used in the Async API
* **--skip-tls-verify** (boolean) - Skip TLS verify. Used for local testing. **Not recommended for production**.

## Make Commands

* Build Locally

    ```shell script
    make build
    ```

* Build Docker image

    ```shell script
    make build-image
    ```

* Run Local docker

    ```shell script
    make run-docker-local
    ```

* Test against a local running instance on localhost:8080 (docker or direct)

    ```shell script
    make test-local
    ```

## Code coverage

* Managed manually at present (TODO: integrate with CI)
* To upload code coverage. Contact `gaurav.abbi@sap.com` for token.

```shell script
source codecov-token.txt && bash <(curl -s https://codecov.io/bash)
```

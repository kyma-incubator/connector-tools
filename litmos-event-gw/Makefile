RELEASE?=master
APP?=litmos-event-gw
GOOS?=linux
GOARCH?=amd64
DOCKER_REG_PATH?=eu.gcr.io/kyma-project/incubator
CONTAINER_IMAGE?=${DOCKER_REG_PATH}/${APP}
EP_URL?="https://passport.free.beeceptor.com/my/api/path"

clean:
	rm -rf ${APP}

vet:
	GO111MODULE=on go vet .

test:
	GO111MODULE=on go test -race -v ./... -coverprofile=coverage.txt -covermode=atomic

build: clean vet test
	GO111MODULE=on CGO_ENABLED=1 go build -race \
		-ldflags "-s -w" \
		-o ${APP}

build-image: build
	docker build -t $(CONTAINER_IMAGE):$(RELEASE) .

push-image: build-image
	docker push $(CONTAINER_IMAGE):$(RELEASE)

tidy:
	GO111MODULE=on go mod tidy

run-docker-local: build-image
	docker run -d -p 8080:8080 $(CONTAINER_IMAGE):$(RELEASE) --app-name=litmos --event-publish-url=$(EP_URL) --skip-tls-verify=true

test-local:
	curl -X POST -d @./assets/test-event.json http://localhost:8080/events -v

RELEASE?=0.0.2
APP?=litmos-event-gw
GOOS?=linux
GOARCH?=amd64
DOCKER_ACCOUNT?=gabbi
CONTAINER_IMAGE?=${DOCKER_ACCOUNT}/${APP}
EP_URL?="https://passport.free.beeceptor.com/my/api/path"

clean:
	rm -rf ${APP}

vet:
	GO111MODULE=on go vet .

test:
	GO111MODULE=on go test -v ./...

build: clean vet test
	GO111MODULE=on CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build \
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

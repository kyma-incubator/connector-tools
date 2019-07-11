# Build the manager binary
FROM golang:1.12.5 as builder

WORKDIR /workspace
# Copy the go source
COPY main.go main.go
COPY healthz.go healthz.go
COPY pkg/ pkg/

# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum



# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o qualtrics-webhook-registration

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:latest
WORKDIR /
COPY --from=builder /workspace/qualtrics-webhook-registration .
ENTRYPOINT ["/qualtrics-webhook-registrationREADME.md"]

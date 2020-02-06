FROM golang:1.13 as builder

ENV GO111MODULE=on

WORKDIR /app
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY cmd               ./cmd

RUN ls /app/
RUN CGO_ENABLED=0 GOOS=linux go build -v -a -installsuffix cgo -o kyma-oauth-proxy ./cmd/kyma-oauth-proxy

FROM alpine:3.8 as certs

RUN apk add -U --no-cache ca-certificates

FROM scratch
WORKDIR /app

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/kyma-oauth-proxy /app/

EXPOSE 8080
ENTRYPOINT ["/app/kyma-oauth-proxy"]
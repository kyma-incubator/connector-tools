FROM golang:1.13 as builder

ENV GO111MODULE=on

WORKDIR /app
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY pkg               ./pkg
COPY cmd               ./cmd

RUN ls /app/
RUN CGO_ENABLED=0 GOOS=linux go build -v -a -installsuffix cgo -o gen-event-gw ./cmd/gen-event-gw

FROM scratch
WORKDIR /app
COPY --from=builder /app/gen-event-gw /app/

EXPOSE 8080
ENTRYPOINT ["/app/gen-event-gw"]
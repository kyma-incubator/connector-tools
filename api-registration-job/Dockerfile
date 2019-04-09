FROM golang:alpine AS build-env
WORKDIR /src
COPY main.go .

RUN go build -o registration_app

# final stage
FROM alpine
RUN apk --no-cache upgrade && apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=build-env /src/registration_app /app/

CMD ["./registration_app"]
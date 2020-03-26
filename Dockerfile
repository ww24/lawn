FROM golang:1.14-alpine3.11 AS build

WORKDIR $GOPATH/src/github.com/ww24/lawn
COPY . .
RUN CGO_ENABLED=0 go build -o /usr/local/bin/lawn ./cmd/lawn


FROM alpine:3.11

RUN apk add --no-cache tzdata ca-certificates
COPY --from=build /usr/local/bin/lawn /usr/local/bin/lawn
ENTRYPOINT [ "lawn" ]

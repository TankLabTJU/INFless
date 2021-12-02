FROM golang:1.13-alpine AS builder

WORKDIR /go/src/github.com/openfaas-incubator/faas-idler

ENV GO111MODULE=off
ENV CGO_ENABLED=0

COPY types      types
COPY main.go    main.go
COPY vendor     vendor

RUN go build -o /usr/bin/faas-idler .

FROM alpine:3.11

RUN addgroup -S app && adduser -S -g app app
RUN mkdir -p /home/app

WORKDIR /home/app

COPY --from=builder /usr/bin/faas-idler /home/app/

RUN chown -R app /home/app
USER app

EXPOSE 8080
VOLUME /tmp

ENTRYPOINT ["/home/app/faas-idler"]

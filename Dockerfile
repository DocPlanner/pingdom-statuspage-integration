FROM golang:1.13-alpine AS builder
RUN apk update && apk add alpine-sdk git && rm -rf /var/cache/apk/*
WORKDIR /go/src/app

ADD ./ /go/src/app/
RUN go mod download
RUN go test
RUN go build -o ./pingdom-statuspage-integration .

FROM alpine:3.11
EXPOSE 80
ENTRYPOINT ["./pingdom-statuspage-integration"]
ENV GIN_MODE=release
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

COPY --from=builder /go/src/app/pingdom-statuspage-integration .
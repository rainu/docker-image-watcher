FROM golang:1.11 as buildContainer

ENV CGO_ENABLED=0
ENV GO111MODULE=on
ENV GOOS=linux
ENV GOPATH=/

COPY . /src/docker-image-watcher
WORKDIR /src/docker-image-watcher

RUN go get ./... &&\
    go build -ldflags -s -a -installsuffix cgo -o docker-image-watcher ./cmd/watcher/


FROM alpine

COPY --from=buildContainer /src/docker-image-watcher/docker-image-watcher /docker-image-watcher

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

USER 10000:10000
EXPOSE 8080

ENTRYPOINT ["/docker-image-watcher"]
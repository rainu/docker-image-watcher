# docker image watcher

Docker-Image-Watcher is a server component which can observe docker images and notify some listeners. So you can
use it to trigger re-build you docker images if a base docker images was updated.

## Installation

First of all you need a running postgres database:
```bash
docker run --rm --name postgres \
    -p 5432:5432 \
    -e POSTGRES_PASSWORD=postgres \
    -e POSTGRES_USER=postgres \
    -e POSTGRES_DB=postgres \
    postgres:alpine
```

Then you can start the application:
```bash
docker run -p 8080:8080 \
    --link postgres \
    -e DATABASE_HOST=postgres \ 
    rainu/docker-image-watcher
```

Or if you have [docker-compose](https://docs.docker.com/compose/gettingstarted/) installed. You can just use this command:
```bash
docker run -p 8080:8080 \
    --link postgres \
    -e DATABASE_HOST=postgres \ 
    rainu/docker-image-watcher
```

## Usage example

The following curl let observe the **alpine:latest** image and trigger the **example.org/trigger** endpoint:
```bash
curl -v -XPOST localhost:8080/api/v1/registry -d '{"image": "library/alpine", "trigger": {"name": "test", "url": "http://example.org/trigger"}}'
```

## Documentation

### Configuration

| ENV-Variable                   | Default-Value | required | Description  |
| ------------------------------ |:-------------:|:--------:| -------------|
| BIND_PORT                      | 8080          | false    | The port where the service listen on |
| OBSERVATION_INTERVAL           | 1m            | false    | The interval for checking the observations (registry lookup) |
| OBSERVATION_DISPATCH_INTERVAL  | 30s           | false    | The interval for looking processable observations |
| OBSERVATION_LIMIT              | 10            | false    | How many observations should be check simultaneously |
| NOTIFICATION_DISPATCH_INTERVAL | 1m            | false    | The interval looking for overdue listeners |
| NOTIFICATION_LIMIT             | 60            | false    | How many listeners should be notify simultaneously |
| NOTIFICATION_TIMEOUT           | 10s           | false    | The connection timeout for each listener |
| DATABASE_SSL                   | disable       | false    | The database ssl mode: enable/disable |
| DATABASE_HOST                  | localhost     | false    | The database host |
| DATABASE_PORT                  | 5432          | false    | The database port |
| DATABASE_SCHEMA                | postgres      | false    | The database schema |
| DATABASE_USER                  | postgres      | false    | The database user |
| DATABASE_PASSWORD              | postgres      | false    | The database password |

### API

See the [API-Documentation](./api_doc.yml)

## Development setup

The following scriptlet shows how to setup the project and build from source code.

```sh
mkdir -p ./workspace/src
export GOPATH=./workspace

cd ./workspace/src
git clone git@github.com:rainu/docker-image-watcher.git

cd docker-image-watcher
go get ./...
go build -ldflags -s -a -installsuffix cgo ./cmd/watcher/
```

## Release History
* 0.0.1 The first implementation

## Meta

Distributed under the MIT license. See ``LICENSE`` for more information.

### Intention

I searched a solution to automatically re-build my own docker images after the base image (such like alpine) was updated.
The official docker-hub dont offer such a solution. You can only expose a trigger api for your own docker image. But the
triggering is on your own. So this application is the link between observation and triggering.

## Contributing

1. Fork it (<https://github.com/rainu/docker-image-watcher/fork>)
2. Create your feature branch (`git checkout -b feature/fooBar`)
3. Commit your changes (`git commit -am 'Add some fooBar'`)
4. Push to the branch (`git push origin feature/fooBar`)
5. Create a new Pull Request

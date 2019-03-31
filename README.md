# docker image watcher

Docker-Image-Watcher is a server component which can observe docker images and notify some listeners. So you can
use it to trigger re-build you docker images if a base docker images was updated.

## Installation

Docker:

```sh
docker run --rm --name postgres \
           -p 5432:5432 \
           -e POSTGRES_PASSWORD=postgres \
           -e POSTGRES_USER=postgres \
           -e POSTGRES_DB=postgres \
           postgres
docker run -p 8080:8080 rainu/docker-image-watcher
```

## Usage example

## Documentation

### Configuration

| ENV-Variable                   | CLI-Option-Name                  | Default-Value | required | Description  |
| ------------------------------ |----------------------------------|:-------------:|:--------:| -------------|
| BIND_PORT                      | --bind-port                      | 8080          | false    | The port where the service listen on |
| OBSERVATION_INTERVAL           | --observation-interval           | 1m            | false    | The interval for checking the observations (registry lookup) |
| OBSERVATION_DISPATCH_INTERVAL  | --observation-dispatch-interval  | 30s           | false    | The interval for looking processable observations |
| OBSERVATION_LIMIT              | --observation-limit              | 10            | false    | How many observations should be check simultaneously |
| NOTIFICATION_DISPATCH_INTERVAL | --notification-dispatch-interval | 1m            | false    | The interval looking for overdue listeners |
| NOTIFICATION_LIMIT             | --notification-limit             | 60            | false    | How many listeners should be notify simultaneously |
| NOTIFICATION_TIMEOUT           | --notification-timeout           | 10s           | false    | The connection timeout for each listener |
| DATABASE_SSL                   | --database-ssl                   | disable       | false    | The database ssl mode: enable/disable |
| DATABASE_HOST                  | --database-host                  | localhost     | false    | The database host |
| DATABASE_PORT                  | --database-port                  | 5432          | false    | The database port |
| DATABASE_SCHEMA                | --database-schema                | postgres      | false    | The database schema |
| DATABASE_USER                  | --database-user                  | postgres      | false    | The database user |
| DATABASE_PASSWORD              | --database-password              | postgres      | false    | The database password |

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
* 0.0.1

## Meta

Distributed under the MIT license. See ``LICENSE`` for more information.

### Intention

tbd

## Contributing

1. Fork it (<https://github.com/rainu/docker-image-watcher/fork>)
2. Create your feature branch (`git checkout -b feature/fooBar`)
3. Commit your changes (`git commit -am 'Add some fooBar'`)
4. Push to the branch (`git push origin feature/fooBar`)
5. Create a new Pull Request

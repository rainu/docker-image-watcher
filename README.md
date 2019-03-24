# docker image watcher

tbd

## Installation

Docker:

```sh
docker run -p 8080:8080 rainu/docker-image-watcher
```

## Usage example

## Documentation

### Configuration

| ENV-Variable        | CLI-Option-Name      | Default-Value | required | Description  |
| ------------------- |----------------------|:-------------:|:--------:| -------------|
| BIND_PORT           | --bind-port          | 8080          | false    | The port where the service listen on |

### API

| Method  | Path      | Variables     | Body |  Description  |
| ------- | --------- | ------------- | ---- | ------------- |

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

# staticstorage - Golang microservice to cache and serve external assets

## Why

~~Just for fun~~ Serve external assets without keeping them locally or in repository

## Features

* Cache
* Separate cache warm up and start
* Healthcheck
* Graceful shutdown


## Usage

### Basic start

```shell
staticstorage -c <CONFIG_PATH> # Warm up cache and start server
```

### Separate cache warm up and start

Separating cache warm up and start especially handful with initContainers in Kubernetes

```shell
staticstorage -c <CONFIG_PATH> warm-cache # Warm up cache without start server
```

Then start server

```shell
staticstorage -c <CONFIG_PATH> server --no-warm-up-cache
```

## Configuration

Make your error pages full of `http.cat`

```json
{
    "cache": "/var/staticstorage",
    "port": "8080",
    "insecure": true,
    "handlers": [
        {
            "path": "/status/401",
            "asset": "401cat"
        },
        {
            "path": "/status/403",
            "asset": "403cat"
        },
        {
            "path": "/status/404",
            "asset": "404cat"
        }
    ],
    "externalAssets": [
        {
            "name": "401cat",
            "url": "https://http.cat/401"
        },
        {
            "name": "403cat",
            "url": "https://http.cat/403"
        },
        {
            "name": "404cat",
            "url": "https://http.cat/404"
        }
    ]
}
```

## Healthcheck

At `/health` available route for service healthcheck useful to use as K8S ready and liveness probe since starts only after cache warm up

```shell
curl http://localhost/health
```

## Build and CI

### Build

Build service docker image

```shell
make 
```

### CI

Lint service

```shell
make check
```
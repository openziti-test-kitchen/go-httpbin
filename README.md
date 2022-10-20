# go-httpbin with embedded OpenZiti

A reasonably complete and well-tested golang port of [Kenneth Reitz][kr]'s
[httpbin][httpbin-org] service that uses the OpenZiti SDK to listen for HTTP requests on an OpenZiti zero trust overlay instead of IP.

## Tutorial

You may follow the steps in [the end-to-end tutorial](./examples/ziti-end-to-end#readme) to run a local Ziti stack built from the main Ziti source repo [![OpenZiti github stars badge](https://img.shields.io/github/stars/openziti/ziti?style=flat)](https://github.com/openziti/ziti/stargazers).

## Usage

### Configuration

go-httpbin can be configured via either command line arguments or environment
variables (or a combination of the two):

| Argument| Env var | Documentation | Default |
| - | - | - | - |
| `-host` | `HOST` | Host to listen on | "0.0.0.0" |
| `-https-cert-file` | `HTTPS_CERT_FILE` | HTTPS Server certificate file | |
| `-https-key-file` | `HTTPS_KEY_FILE` | HTTPS Server private key file | |
| `-max-body-size` | `MAX_BODY_SIZE` | Maximum size of request or response, in bytes | 1048576 |
| `-max-duration` | `MAX_DURATION` | Maximum duration a response may take | 10s |
| `-port` | `PORT` | Port to listen on | 8080 |
| `-use-real-hostname` | `USE_REAL_HOSTNAME` | Expose real hostname as reported by os.Hostname() in the /hostname endpoint | false |
| `-ziti` | `ENABLE_ZITI` | Enable using a ziti network | false|
| `-ziti-identity` | `ZITI_IDENTITY` | Ziti identity json file location | - |
| `-ziti-name` | `ZITI_SERVICE_NAME` | Name of Ziti Service to bind against | - |

**Notes:**

* Command line arguments take precedence over environment variables.
* As an alternative to supplying `ziti-identity` as a file you may define environment variable `ZITI_IDENTITY_JSON`.

### Standalone binary

Follow the [Installation](#installation) instructions to install go-httpbin as
a standalone binary. (This currently requires a working Go runtime.)

Examples:

[![Ziti Reference](https://github.com/openziti/ziti#readme)]
[![Ziti Http Reference](https://github.com/openziti-test-kitchen/go-http#readme)]
This example assumes you are familiar with spinning up a ziti network and have a network with a service named "httpbin".

```bash
# Ensure your ziti network is spun up.
# Run http server
$ go-httpbin -ziti -ziti-identity ./my-ziti-identity.json -ziti-name "my httpbin service"

#Run https server
$ openssl genrsa -out server.key 2048
$ openssl ecparam -genkey -name secp384r1 -out server.key
$ openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
$ go-httpbin -https-cert-file ./server.crt -https-key-file ./server.key -ziti -ziti-identity ${ZITI_IDENTITY} -ziti-name httpbin
```

### Docker

A multi-platform Docker image is published to [Docker Hub](https://hub.docker.com/r/openziti/go-httpbin):

Run the included Compose project:

```bash
# httpbinz-server1.json exists in the same dir as docker-compose.yml
ZITI_IDENTITY_JSON="$(< ./my-ziti-identity.json)" \
ZITI_SERVICE_NAME="my httpbin service" \
    docker compose run httpbin                                      
```

Run without Compose:

```bash
# Run http server with my-ziti-identity.json in current working dir
docker run \
    -e ENABLE_ZITI=true \
    -e ZITI_IDENTITY_JSON="$(< ./my-ziti-identity.json)" \
    -e ZITI_SERVICE_NAME="my httpbin service" \
    openziti/go-httpbin

# Run https server with my-ziti-identity.json in current working dir
docker run \
    -e HTTPS_CERT_FILE='/tmp/server.crt' \
    -e HTTPS_KEY_FILE='/tmp/server.key' \
    -v /tmp:/tmp \
    -e ENABLE_ZITI=true \
    -e ZITI_IDENTITY_JSON="$(< ./my-ziti-identity.json)" \
    -e ZITI_SERVICE_NAME="my httpbin service" \
    openziti/go-httpbin
```

Build the Container Image for your Platform

```bash
docker compose build httpbin
```

## Motivation & prior art

I've been a longtime user of [Kenneith Reitz][kr]'s original
[httpbin.org][httpbin-org], and wanted to write a golang port for fun and to
see how far I could get using only the stdlib.

When I started this project, there were a handful of existing and incomplete
golang ports, with the most promising being [ahmetb/go-httpbin][ahmet]. This
project showed me how useful it might be to have an `httpbin` _library_
available for testing golang applications.

### Known differences from other httpbin versions

Compared to [the original][httpbin-org]:

- No `/brotli` endpoint (due to lack of support in Go's stdlib)
- The `?show_env=1` query param is ignored (i.e. no special handling of
runtime environment headers)
- Response values which may be encoded as either a string or a list of strings
will always be encoded as a list of strings (e.g. request headers, query
params, form values)

Compared to [ahmetb/go-httpbin][ahmet]:

- No dependencies on 3rd party packages
- More complete implementation of endpoints

## Development

```bash
# local development
make
make test
make testcover
make run

# building & pushing docker images
make image
make imagepush
```

[kr]: https://github.com/kennethreitz
[httpbin-org]: https://httpbin.org/
[httpbin-repo]: https://github.com/kennethreitz/httpbin
[ahmet]: https://github.com/ahmetb/go-httpbin
[docker-hub]: https://hub.docker.com/r/mccutchen/go-httpbin/
[observer]: https://pkg.go.dev/github.com/mccutchen/go-httpbin/v2/httpbin#Observer
[custom-instrumentation]: ./examples/custom-instrumentation/

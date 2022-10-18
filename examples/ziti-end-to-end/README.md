# Setting up an End-to-End Ziti network
## Prerequisites:

Ziti Binaries and relevant configs.

[OpenZiti](https://github.com/openziti/ziti/tree/main) contains downloads for the latest builds or you can build it yourself. There are many example configs [in the etc dir](https://github.com/openziti/ziti/tree/main/etc).

I have set up helpful environment variable for how I chose to get setup. I cloned the Ziti repository for local builds and set `$ZITI_SOURCE` to that directory.

## Starting up the Ziti Network:

### Controller Setup

First we run the ziti controller with the arguent being which ever config file we chose. For this example I am using the `ctrl.with.edge.yml` file bundled for examples in the ziti repository.
```
ziti-controller run --log-formatter pfxlog $ZITI_SOURCE/ziti/etc/ctrl.with.edge.yml
```

We then need to initialize the controller with an admin user. For ease of testing I just used `admin admin admin` for my short local test.

```
ziti agent controller init <username> <password> <name-of-user>
```

You can then login to the edge via `ziti edge login` and putting in your credentials.

### Policies
Next we create all of the relevant service policies:

```
ziti edge create service-policy dial-simple Dial \
    --service-roles '#simple' \
    --identity-roles '#simple-client'

ziti edge create service-policy bind-simple Bind \
    --service-roles '#simple' \
    --identity-roles '#simple-server'
```

And the Edge router policies:

```
ziti edge create edge-router-policy simple-client \
    --identity-roles '#simple-client' \
    --edge-router-roles '#all'

ziti edge create edge-router-policy simple-server \
    --identity-roles '#simple-server' \
    --edge-router-roles '#all'

ziti edge create service-edge-router-policy simple \
    --service-roles '#simple' \
    --edge-router-roles '#all'
```

Now we create and enroll our client user

```
ziti edge create identity service simple-client \
    --jwt-output-file simple-client.jwt \
    --role-attributes simple-client,simple

ziti edge enroll \
    --jwt simple-client.jwt \
    --out simple-client.json
```

And then create and enroll the server

```
ziti edge create identity service simple-server \
    --jwt-output-file simple-server.jwt \
    --role-attributes simple-server,simple

ziti edge enroll \
    --jwt simple-server.jwt \
    --out simple-server.json
```

All that's left to do for this portion is to create the service
```
ziti edge create service echo --role-attributes simple
```

Now you can optionally delete the two jwt's we generated. We will not be using them in the rest of the example.
### Edge Router
Next we create and run the edge router:

```
ziti edge create edge-router edge-router \
    --jwt-output-file edge-router.jwt \
    --tunneler-enabled
```

Just like with the controller I'm using an example `edge.router.yml` found in the source repository.
```
ziti-router enroll --jwt edge-router.jwt ${ZITI_SOURCE}/ziti/etc/edge.router.yml
```

Now we run the edge router like so! The controller host and port info is found in the controller config file you used.
```
CONTROLLER_HOST=localhost \
CONTROLLER_PORT=6262 \
ZITI_EDGE_PORT=3022 \
LINK_LISTENER_PORT=4022 \
exec $GOPATH/bin/ziti-router run \
    --debug-ops \
    --verbose \
    --log-formatter pfxlog \
    ${ZITI_SOURCE}/ziti/etc/edge.router.yml
```

## go-httpbin Server
All we need to do to run the server is pass in the relevant flags. This assumes you ran all above commands in the same directory, which you should see `simple-server.json`. If not then point the `-ziti-identity` flag to that file wherever it is.
```
go-httpbin -ziti -ziti-identity ${PWD}/simple-server.json -ziti-name echo
```

## Client
Running the client is all the same, just pass in the relevant files like above.
```
go-httpbin-client \
    -header k=v \
    -header k=v2 \
    -query y=m \
    -ziti \
    -ziti-identity ${PWD}/simple-client.json \
    -ziti-name echo \
    post test
```

That will give us the output

```json
{
    "args": {
        "y": [
            "m"
        ]
    },
    "headers": {
        "Accept-Encoding": [
            "gzip"
        ],
        "Host": [
            "echo"
        ],
        "K": [
            "v",
            "v2"
        ],
        "User-Agent": [
            "Go-http-client/1.1"
        ]
    },
    "origin": "ziti-edge-router connId=2147483648, logical=ziti-sdk[router=tls://127.0.0.1:3022]",
    "url": "http://echo/post?y=m",
    "data": "test",
    "files": null,
    "form": null,
    "json": null
}
```

Now you're able to use this basic client againt the go-httpbin server over OpenZiti! 
package zitiserver

import (
	"context"
	"net"
	"net/http"

	"github.com/openziti/sdk-golang/ziti"
	"github.com/openziti/sdk-golang/ziti/config"
)

type ZitiServer struct {
	listener     net.Listener
	serviceName  string
	identityJson string

	handler http.Handler
}

func New(serviceName, identityJson string, handler http.Handler) *ZitiServer {
	return &ZitiServer{
		serviceName:  serviceName,
		identityJson: identityJson,
		handler:      handler,
	}
}

func (zs *ZitiServer) SetKeepAlivesEnabled(keepAlive bool) {
	//TODO
}

func (zs *ZitiServer) Shutdown(context.Context) error {
	//TODO more graceful shutdown
	return zs.listener.Close()
}

func (zs *ZitiServer) ListenAndServe() error {
	var err error
	if zs.listener, err = zs.configAndGetListener(); err != nil {
		return err
	}

	return http.Serve(zs.listener, zs.handler)
}

func (zs *ZitiServer) ListenAndServeTLS(tlsCertFile, tlsKeyFile string) error {
	var err error
	if zs.listener, err = zs.configAndGetListener(); err != nil {
		return err
	}
	return http.ServeTLS(zs.listener, zs.handler, tlsCertFile, tlsKeyFile)
}

func (zs *ZitiServer) configAndGetListener() (net.Listener, error) {
	config, err := config.NewFromFile(zs.identityJson)
	if err != nil {
		return nil, err
	}

	zitiContext := ziti.NewContextWithConfig(config)
	return zitiContext.Listen(zs.serviceName)
}

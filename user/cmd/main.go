package main

import (
	"fmt"
	"net"

	"github.com/kkyr/fig"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/salman-pathan/go-micro-arch/common/config"
	userservice "github.com/salman-pathan/go-micro-arch/user"
	user "github.com/salman-pathan/go-micro-arch/user/pb"
)

var (
	logger log.Logger
	cfg    config.Config
)

const (
	serverCert = "ssl/server.crt"
	serverKey  = "ssl/server.pem"
)

func main() {

	// configuration
	if err := fig.Load(&cfg); err != nil {
		panic("failed to load config")
	}

	//	logger
	logger.SetFormatter(&log.JSONFormatter{})
	logger.SetReportCaller(true)
	logger.SetLevel(log.DebugLevel)

	tlsCredentials, err := credentials.NewServerTLSFromFile(serverCert, serverKey)
	if err != nil {
		panic(err)
	}

	//	services
	userService := userservice.NewService(logger)

	//	tcp listener
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Server.Port))
	if err != nil {
		panic(err)
	}

	//	server options
	opts := []grpc.ServerOption{
		grpc.Creds(tlsCredentials),
	}

	grpcServer := grpc.NewServer(opts...)
	user.RegisterUserServer(grpcServer, userService)

	//	run gRPC server
	if err := grpcServer.Serve(listener); err != nil {
		panic(err)
	}

}

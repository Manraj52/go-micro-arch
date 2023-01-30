package main

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/kkyr/fig"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	userservice "github.com/salman-pathan/go-micro-arch/user"
	"github.com/salman-pathan/go-micro-arch/user/cmd/config"
	user "github.com/salman-pathan/go-micro-arch/user/pb"
	mongorepository "github.com/salman-pathan/go-micro-arch/user/repositories/mongo"
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

	//	mongo
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.Timeout)
	defer cancel()

	mongoUri := fmt.Sprintf("mongodb://%s:%s@%s:%d", cfg.Mongo.Username, cfg.Mongo.Password, cfg.Mongo.Host, cfg.Mongo.Port)
	fmt.Println(mongoUri)
	mongoOptions := options.Client().ApplyURI(mongoUri)
	mongoClient, err := mongo.Connect(ctx, mongoOptions)
	if err != nil {
		log.Fatal(err)
	}

	if err = mongoClient.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err = mongoClient.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	//	logger
	logger = *log.New()
	logger.SetFormatter(&log.JSONFormatter{})
	logger.SetReportCaller(true)
	logger.SetLevel(log.DebugLevel)

	tlsCredentials, err := credentials.NewServerTLSFromFile(serverCert, serverKey)
	if err != nil {
		panic(err)
	}

	//	repositories
	userRepository, err := mongorepository.NewUserRepository(cfg.Mongo.Database, *mongoClient)
	if err != nil {
		log.Fatal(err)
	}

	//	location
	location, err := time.LoadLocation(cfg.Location)
	if err != nil {
		log.Fatal(err)
	}

	//	services
	userService := userservice.NewService(logger, location, userRepository)

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

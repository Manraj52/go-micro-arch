//  Copyright 2023 Salman Khan. All rights reserved.
//  Use of this source code is governed by GPL style
//  license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kkyr/fig"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/salman-pathan/go-micro-arch/common/auth"
	"github.com/salman-pathan/go-micro-arch/common/config"
	"github.com/salman-pathan/go-micro-arch/common/middlewares"
	"github.com/salman-pathan/go-micro-arch/common/response"
	"github.com/salman-pathan/go-micro-arch/user"
	pb "github.com/salman-pathan/go-micro-arch/user/pb"
)

var cfg config.Config

const (
	certFile = "ssl/ca.crt"
)

func main() {

	//	configuration
	if err := fig.Load(&cfg); err != nil {
		log.Fatal(err)
	}

	//	logger
	log.SetFormatter(&log.JSONFormatter{})
	log.SetReportCaller(true)
	log.SetLevel(log.DebugLevel)

	//	credentials
	creds, err := credentials.NewClientTLSFromFile(certFile, "")
	if err != nil {
		panic(err)
	}

	//	auth helper
	authHelper, err := auth.NewPasetoAuth(cfg.Server.SecretKey)
	if err != nil {
		panic(err)
	}

	serverOptions := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
	}

	//	middleware
	authMiddleware := middlewares.NewAuthMiddleware(authHelper)

	//	gRPC client
	userConn, err := grpc.Dial("user-svc:8001", serverOptions...)
	if err != nil {
		panic(err)
	}
	defer userConn.Close()

	userClient := pb.NewUserClient(userConn)

	crLogger := log.WithField("service.name", "ResponseHelper")
	respHelper := response.NewCustomResponse(crLogger)

	//  router
	router := gin.Default()
	router.Use(cors.Default())

	v1 := router.Group("/api/v1")
	userRouteGroup := v1.Group("/user")
	{
		var userRoutes user.Routes
		{
			userLogger := log.WithFields(log.Fields{
				"service": "user",
			})
			userRoutes = user.NewUserRoutes(userLogger, userRouteGroup, userClient, authMiddleware, respHelper)
			userRoutes.RegisterRoutes()
		}
	}

	router.Run(fmt.Sprintf(":%d", cfg.Server.Port))
}

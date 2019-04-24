package main

import (
	"2019_1_OPG_plus_2/internal/app/server"
	"2019_1_OPG_plus_2/internal/pkg/authservice"
	"2019_1_OPG_plus_2/internal/pkg/tsLogger"
	authService "2019_1_OPG_plus_2/internal/proto"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

// @title Colors service API by OPG+2
// @version 1.0
// @description Game based on filling field with color cells
// @license.name Apache 2.0

// @contact.name @DanikNik
// @contact.email nikolsky.dan@gmail.com

// @host localhost:8002
// @BasePath /api

func main() {

	go func() {
		lis, err := net.Listen("tcp", ":50242")
		if err != nil {
			log.Fatalln("cant listen port", err)
		}

		authserver := grpc.NewServer()

		authService.RegisterAuthServiceServer(authserver, authservice.NewServer())

		fmt.Println("starting server at :50242")
		authserver.Serve(lis)
	}()

	params := server.Params{Port: os.Getenv("PORT")}
	if params.Port == "" {
		params.Port = "8002"
	}

	err := server.StartApp(params)
	if err != nil {
		server.StopApp()
		tsLogger.LogFatal("%s", err)
	}

}

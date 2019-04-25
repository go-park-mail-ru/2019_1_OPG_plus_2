package authserver

import (
	"2019_1_OPG_plus_2/internal/pkg/authservice"
	"2019_1_OPG_plus_2/internal/pkg/config"
	"2019_1_OPG_plus_2/internal/pkg/db"
	authService "2019_1_OPG_plus_2/internal/proto"
	"google.golang.org/grpc"
	"net"
	"os"
)

func Start() error {
	if os.Getenv("COLORS_AUTH_MODE") == "TEST" {
		db.AuthDbName = config.Db.AuthTestDb
		db.CoreDbName = config.Db.CoreTestDb
	}

	serv := authservice.NewServer()

	serv.Log.Run()

	if err := db.Open(); err != nil {
		serv.Log.LogErr("%v", err)
	}

	lis, err := net.Listen("tcp", ":50242")
	if err != nil {
		serv.Log.LogFatal("cant listen port: %s", err)
	}

	server := grpc.NewServer()

	authService.RegisterAuthServiceServer(server, serv)

	serv.Log.LogTrace("starting server at :50242")
	return server.Serve(lis)
}

//func Stop() {
//
//}
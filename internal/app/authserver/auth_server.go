package authserver

import (
	"2019_1_OPG_plus_2/internal/pkg/authproto"
	"2019_1_OPG_plus_2/internal/pkg/authservice"
	"2019_1_OPG_plus_2/internal/pkg/config"
	"2019_1_OPG_plus_2/internal/pkg/db"
	"google.golang.org/grpc"
	"net"
	"os"
)

//func init(){
//	config.Init()
//	auth.Init()
//}

func Start() error {
	//config.Init()
	serv := authservice.NewService()

	serv.Log.Run()
	//auth.Init()

	if os.Getenv("COLORS_AUTH_MODE") == "TEST" {
		db.AuthDbName = config.Db.AuthTestDb
		db.CoreDbName = config.Db.CoreTestDb
		serv.Log.LogTrace("AUTH: TESTING MODE")
	}
	if err := db.Open(); err != nil {
		serv.Log.LogErr("%v", err)
	}

	lis, err := net.Listen("tcp", ":"+config.Auth.AuthPort)
	if err != nil {
		serv.Log.LogFatal("AUTH: cant listen port: %s", err)
	}

	server := grpc.NewServer()

	authproto.RegisterAuthServiceServer(server, serv)

	serv.Log.LogTrace("AUTH: starting core at %v:%v", config.Auth.AuthServiceLocation, config.Auth.AuthPort)
	return server.Serve(lis)
}

//func Stop() {
//
//}

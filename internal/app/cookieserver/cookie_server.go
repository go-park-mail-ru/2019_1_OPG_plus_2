package cookieserver

import (
	"2019_1_OPG_plus_2/internal/pkg/config"
	"2019_1_OPG_plus_2/internal/pkg/cookiecheckerproto"
	"2019_1_OPG_plus_2/internal/pkg/cookieservice"
	"2019_1_OPG_plus_2/internal/pkg/tsLogger"
	"google.golang.org/grpc"
	"net"
)

func Start() error {
	serv := cookieservice.NewServer(tsLogger.NewLogger())
	serv.Log.Run()

	lis, err := net.Listen("tcp", ":"+config.Auth.CookieServicePort)
	if err != nil {
		serv.Log.LogFatal("COOKIE: cant listen port: %s", err)
	}

	server := grpc.NewServer()
	cookiecheckerproto.RegisterCookieCheckerServer(server, serv)
	serv.Log.LogTrace("Serving cookie at %v:%v", config.Auth.CookieServiceLocation, config.Auth.CookieServicePort)
	return server.Serve(lis)
}

//func Stop() {}

package cookieserver

import (
	"2019_1_OPG_plus_2/internal/pkg/config"
	"2019_1_OPG_plus_2/internal/pkg/cookiecheckerproto"
	"2019_1_OPG_plus_2/internal/pkg/cookieservice"
	"google.golang.org/grpc"
	"log"
	"net"
)

func Start() error {
	serv := cookieservice.Server{}

	lis, err := net.Listen("tcp", ":"+config.Auth.CookieServicePort)
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()
	cookiecheckerproto.RegisterCookieCheckerServer(server, serv)
	log.Println("Serving cookie at ", config.Auth.CookieServiceLocation, ":", config.Auth.CookieServicePort)
	return server.Serve(lis)
}

func Stop() {}

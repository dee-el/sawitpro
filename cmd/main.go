package main

import (
	"os"

	"github.com/labstack/echo/v4"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/handler"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/SawitProRecruitment/UserService/service"
)

func main() {

	e := echo.New()
	var server generated.ServerInterface = newServer()
	generated.RegisterHandlers(e, server)
	e.Logger.Fatal(e.Start(":1323"))
}

func newServer() *handler.Server {

	repo := repository.NewRepository(os.Getenv("DATABASE_URL"))
	service := service.NewUserService(service.UserServiceOptions{
		Repository: repo,
		AuthService: service.NewJWTAuthService(service.JWTAuthServiceOptions{
			PrivateKey: os.Getenv("PRIVATE_CERT"),
			PublicKey:  os.Getenv("PUBLIC_CERT"),
		}),
	})

	return handler.NewServer(service)
}

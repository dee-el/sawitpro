package handler

import (
	"github.com/SawitProRecruitment/UserService/service"
)

type Server struct {
	userService service.Service
}

func NewServer(userService service.Service) *Server {
	return &Server{
		userService: userService,
	}
}

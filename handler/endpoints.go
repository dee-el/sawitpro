package handler

import (
	"net/http"
	"strings"

	"github.com/SawitProRecruitment/UserService/common/errors"
	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/service"
	"github.com/labstack/echo/v4"
)

func (s *Server) AuthRegister(ctx echo.Context) error {
	var req generated.RegisterRequest
	err := ctx.Bind(&req)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: "Invalid Format",
		})
	}

	id, err := s.userService.Register(ctx.Request().Context(), service.RegisterRequest{
		FullName: req.Fullname,
		Password: req.Password,
		Phone:    req.Phone,
	})
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	var resp generated.RegisterResponse
	resp.Data.Id = id
	return ctx.JSON(http.StatusCreated, resp)
}

func (s *Server) AuthLogin(ctx echo.Context) error {
	var req generated.LoginRequest
	err := ctx.Bind(&req)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: "Invalid Format",
		})
	}

	lr, err := s.userService.Login(ctx.Request().Context(), service.LoginRequest{
		Password: req.Password,
		Phone:    req.Phone,
	})
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	var resp generated.LoginResponse
	resp.Data.Id = lr.ID
	resp.Data.Token = lr.Token
	return ctx.JSON(http.StatusOK, resp)
}

func (s *Server) AuthMe(ctx echo.Context) error {
	token := s.getToken(ctx)
	if token == "" {
		return ctx.JSON(http.StatusForbidden, "")
	}

	usr, err := s.userService.GetMe(ctx.Request().Context(), token)
	if err != nil {
		return ctx.JSON(http.StatusForbidden, err)
	}

	var resp generated.MeResponse
	resp.Data.Fullname = usr.FullName
	resp.Data.Phone = usr.Phone
	return ctx.JSON(http.StatusOK, resp)
}

func (s *Server) AuthMeUpdate(ctx echo.Context) error {
	token := s.getToken(ctx)
	if token == "" {
		return ctx.JSON(http.StatusForbidden, "")
	}

	var req generated.MeUpdateRequest
	err := ctx.Bind(&req)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: "Invalid Format",
		})
	}

	err = s.userService.UpdateMe(ctx.Request().Context(), token, service.UpdateRequest{
		Phone:    req.Phone,
		FullName: req.Fullname,
	})
	if err != nil {
		switch ve := err.(type) {
		case *errors.Error:
			if ve.Type == errors.TypeBadRequestError {
				return ctx.JSON(http.StatusBadRequest, err)
			} else if ve.Type == errors.TypeConflictedError {
				return ctx.JSON(http.StatusConflict, err)
			}
		}

		return ctx.JSON(http.StatusInternalServerError, err)
	}

	var resp generated.MeUpdateResponse
	return ctx.JSON(http.StatusOK, resp)
}

func (s *Server) getToken(ctx echo.Context) string {
	key := "Bearer "
	reqToken := ctx.Request().Header.Get("Authorization")
	if !strings.HasPrefix(reqToken, key) {
		return ""
	}

	token := strings.TrimPrefix(reqToken, key)
	return token
}

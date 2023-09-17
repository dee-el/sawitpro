package handler_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/labstack/echo/v4"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/handler"
	"github.com/SawitProRecruitment/UserService/service"
)

func TestEndpoint_Register(t *testing.T) {
	e := echo.New()

	reqBody := []byte(`{"fullname": "David","password": "aaaaaaA12!","phone": "+6212345698"}`)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	var id int64 = 1
	svc := service.NewMockService(gomock.NewController(t))
	svc.EXPECT().Register(gomock.Any(), gomock.Any()).Return(id, nil)

	var server generated.ServerInterface = handler.NewServer(svc)
	generated.RegisterHandlers(e, server)

	c := e.NewContext(req, rec)

	err := server.AuthRegister(c)
	if err != nil {
		t.Fatal(err)
	}

	diff := cmp.Diff(http.StatusCreated, rec.Code)
	if diff != "" {
		t.Fatalf("[%s] mismatch (-want +got):\n%s", t.Name(), diff)
	}
}

func TestEndpoint_Login(t *testing.T) {
	e := echo.New()

	reqBody := []byte(`{"password": "aaaaaaA12!","phone": "+6212345698"}`)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	svc := service.NewMockService(gomock.NewController(t))
	svc.EXPECT().Login(gomock.Any(), gomock.Any()).Return(service.LoginResponse{
		ID:    1,
		Token: "hello",
	}, nil)

	var server generated.ServerInterface = handler.NewServer(svc)
	generated.RegisterHandlers(e, server)

	c := e.NewContext(req, rec)

	err := server.AuthLogin(c)
	if err != nil {
		t.Fatal(err)
	}

	diff := cmp.Diff(http.StatusOK, rec.Code)
	if diff != "" {
		t.Fatalf("[%s] mismatch (-want +got):\n%s", t.Name(), diff)
	}
}

func TestEndpoint_GetMe(t *testing.T) {
	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer hello")

	rec := httptest.NewRecorder()

	svc := service.NewMockService(gomock.NewController(t))
	svc.EXPECT().GetMe(gomock.Any(), gomock.Any()).Return(service.GetMeResponse{
		FullName: "Nivea",
		Phone:    "+6212345698",
	}, nil)

	var server generated.ServerInterface = handler.NewServer(svc)
	generated.RegisterHandlers(e, server)

	c := e.NewContext(req, rec)

	err := server.AuthMe(c)
	if err != nil {
		t.Fatal(err)
	}

	diff := cmp.Diff(http.StatusOK, rec.Code)
	if diff != "" {
		t.Fatalf("[%s] mismatch (-want +got):\n%s", t.Name(), diff)
	}
}

func TestEndpoint_UpdateMe(t *testing.T) {
	e := echo.New()

	reqBody := []byte(`{"fullname": "Nivea2","phone": "+6212345698"}`)
	req := httptest.NewRequest(http.MethodPatch, "/auth/me", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer hello")

	rec := httptest.NewRecorder()

	svc := service.NewMockService(gomock.NewController(t))
	svc.EXPECT().UpdateMe(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	var server generated.ServerInterface = handler.NewServer(svc)
	generated.RegisterHandlers(e, server)

	c := e.NewContext(req, rec)

	err := server.AuthMeUpdate(c)
	if err != nil {
		t.Fatal(err)
	}

	diff := cmp.Diff(http.StatusOK, rec.Code)
	if diff != "" {
		t.Fatalf("[%s] mismatch (-want +got):\n%s", t.Name(), diff)
	}
}

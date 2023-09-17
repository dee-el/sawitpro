package service

import "context"

type Service interface {
	Register(ctx context.Context, params RegisterRequest) (int64, error)
	Login(ctx context.Context, params LoginRequest) (LoginResponse, error)
	GetMe(ctx context.Context, token string) (GetMeResponse, error)
	UpdateMe(ctx context.Context, token string, params UpdateRequest) error
}

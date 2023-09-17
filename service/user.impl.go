package service

import (
	"context"
	"fmt"
	"time"

	"github.com/SawitProRecruitment/UserService/common/errors"
	"github.com/SawitProRecruitment/UserService/common/validator"
	"golang.org/x/crypto/bcrypt"

	"github.com/SawitProRecruitment/UserService/repository"
)

type UserService struct {
	repository  repository.RepositoryInterface
	authService AuthService
}

type UserServiceOptions struct {
	Repository  repository.RepositoryInterface
	AuthService AuthService
}

func NewUserService(opts UserServiceOptions) *UserService {
	return &UserService{
		repository:  opts.Repository,
		authService: opts.AuthService,
	}
}

type RegisterRequest struct {
	FullName string `json:"fullname" validate:"required,min=3,max=60"`
	Password string `json:"password" validate:"required,password"`
	Phone    string `json:"phone" validate:"required,phone"`
}

func (s *UserService) Register(ctx context.Context, params RegisterRequest) (int64, error) {
	err := validator.Validate(validator.ContentTypeJSON, params)
	if err != nil {
		return 0, err
	}

	usr, err := s.repository.FindUserByPhone(ctx, params.Phone)
	if err != nil {
		return 0, err
	}

	if usr != nil {
		e := errors.Conflicted("")
		e.AddField("phone", fmt.Sprintf("[%s] already registered", params.Phone))
		return 0, e
	}

	hashed, err := s.hashed(params.Password)
	if err != nil {
		return 0, err
	}

	now := time.Now()
	return s.repository.CreateUser(ctx, &repository.User{
		FullName:  params.FullName,
		Password:  hashed,
		Phone:     params.Phone,
		CreatedAt: now,
		UpdatedAt: now,
	})
}

type LoginRequest struct {
	Password string `json:"password" validate:"required,password"`
	Phone    string `json:"phone" validate:"required,phone"`
}

type LoginResponse struct {
	ID    int64
	Token string
}

func (s *UserService) Login(ctx context.Context, params LoginRequest) (LoginResponse, error) {
	lr := LoginResponse{}

	err := validator.Validate(validator.ContentTypeJSON, params)
	if err != nil {
		return lr, err
	}

	usr, err := s.repository.FindUserByPhone(ctx, params.Phone)
	if err != nil {
		return lr, err
	}

	if usr == nil {
		e := errors.BadRequest("phone or password is wrong")
		return lr, e
	}

	// compare
	err = bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(params.Password))
	if err != nil {
		e := errors.BadRequest("phone or password is wrong")
		return lr, e
	}

	now := time.Now()
	err = s.repository.CreateUserAttendance(ctx, &repository.UserAttendance{
		UserID:  usr.ID,
		LoginAt: now,
	})
	if err != nil {
		return lr, err
	}

	err = s.repository.SaveUserAttendanceSummary(ctx, usr.ID)
	if err != nil {
		return lr, err
	}

	token, err := s.authService.CreateToken(usr, now)
	if err != nil {
		return lr, err
	}

	lr.ID = usr.ID
	lr.Token = token
	return lr, nil
}

type GetMeResponse struct {
	FullName string
	Phone    string
}

func (s *UserService) GetMe(ctx context.Context, token string) (GetMeResponse, error) {
	resp := GetMeResponse{}

	id, err := s.authService.ValidateToken(token)
	if err != nil {
		return resp, err
	}

	usr, err := s.repository.FindUserByID(ctx, id)
	if err != nil {
		return resp, err
	}

	if usr == nil {
		e := errors.BadRequest("User Not Found")
		return resp, e
	}

	resp.FullName = usr.FullName
	resp.Phone = usr.Phone
	return resp, nil
}

type UpdateRequest struct {
	Phone    *string `json:"phone" validate:"omitempty,phone"`
	FullName *string `json:"fullname" validate:"omitempty,min=3,max=60"`
}

func (s *UserService) UpdateMe(ctx context.Context, token string, params UpdateRequest) error {
	id, err := s.authService.ValidateToken(token)
	if err != nil {
		return err
	}

	err = validator.Validate(validator.ContentTypeJSON, params)
	if err != nil {
		return err
	}

	usr, err := s.repository.FindUserByID(ctx, id)
	if err != nil {
		return err
	}

	if usr == nil {
		e := errors.BadRequest("User Not Found")
		return e
	}

	if usr.Phone != *params.Phone {
		pUsr, err := s.repository.FindUserByPhone(ctx, *params.Phone)
		if err != nil {
			return err
		}

		if pUsr != nil {
			e := errors.Conflicted("")
			e.AddField("phone", fmt.Sprintf("[%s] already registered", *params.Phone))
			return e
		}
	}

	if params.FullName != nil && *params.FullName != "" {
		usr.FullName = *params.FullName
	}

	usr.Phone = *params.Phone
	return s.repository.UpdateUser(ctx, usr)
}

func (s *UserService) hashed(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 6)
	if err != nil {
		return "", err
	}

	return string(hashed), nil
}

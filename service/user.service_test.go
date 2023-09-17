package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/SawitProRecruitment/UserService/common/errors"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/SawitProRecruitment/UserService/service"
	gomock "github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
)

func TestUserService_Register(t *testing.T) {
	scenarios := []struct {
		name     string
		init     func() *service.UserService
		input    func() (context.Context, service.RegisterRequest)
		expected func() (int64, error)
	}{
		{
			name: "OK_Create",
			init: func() (srv *service.UserService) {
				repo := repository.NewMockRepositoryInterface(gomock.NewController(t))

				var id int64 = 1
				repo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(id, nil)
				repo.EXPECT().FindUserByPhone(gomock.Any(), gomock.Any()).Return(nil, nil)
				srv = service.NewUserService(service.UserServiceOptions{
					Repository:  repo,
					AuthService: service.NewMockAuthService(gomock.NewController(t)),
				})

				return srv
			},
			input: func() (context.Context, service.RegisterRequest) {
				ctx := context.Background()

				args := service.RegisterRequest{}
				args.FullName = "Nivea"
				args.Password = "wh4t_yoU_could_do"
				args.Phone = "+6212345678"

				return ctx, args
			},
			expected: func() (int64, error) {
				return 1, nil
			},
		},
		{
			name: "Failed_ConflictedPhone",
			init: func() (srv *service.UserService) {
				repo := repository.NewMockRepositoryInterface(gomock.NewController(t))

				repo.EXPECT().FindUserByPhone(gomock.Any(), gomock.Any()).Return(&repository.User{
					ID: 10,
				}, nil)
				srv = service.NewUserService(service.UserServiceOptions{
					Repository:  repo,
					AuthService: service.NewMockAuthService(gomock.NewController(t)),
				})

				return srv
			},
			input: func() (context.Context, service.RegisterRequest) {
				ctx := context.Background()

				args := service.RegisterRequest{}
				args.FullName = "Nivea"
				args.Password = "wh4t_yoU_could_do"
				args.Phone = "+62123456711"

				return ctx, args
			},
			expected: func() (int64, error) {
				e := errors.Conflicted("")
				e.AddField("phone", fmt.Sprintf("[%s] already registered", "+62123456711"))
				return 0, e
			},
		},
		{
			name: "Failed_ValidateName",
			init: func() (srv *service.UserService) {
				repo := repository.NewMockRepositoryInterface(gomock.NewController(t))

				srv = service.NewUserService(service.UserServiceOptions{
					Repository:  repo,
					AuthService: service.NewMockAuthService(gomock.NewController(t)),
				})

				return srv
			},
			input: func() (context.Context, service.RegisterRequest) {
				ctx := context.Background()

				args := service.RegisterRequest{}
				args.FullName = ""
				args.Password = "wh4t_yoU_could_do"
				args.Phone = "+62123456711"

				return ctx, args
			},
			expected: func() (int64, error) {
				e := errors.BadRequest("Bad request")
				e.AddField("fullname", "is a required field")
				return 0, e
			},
		},
	}

	for _, scn := range scenarios {
		t.Run(scn.name, func(t *testing.T) {
			s := scn.init()

			id, err := s.Register(scn.input())
			expectedId, expectedErr := scn.expected()

			diff := cmp.Diff(expectedId, id)
			if diff != "" {
				t.Fatalf("[%s] mismatch (-want +got):\n%s", t.Name(), diff)
			}

			diff = cmp.Diff(expectedErr, err)
			if diff != "" {
				t.Fatalf("[%s] mismatch (-want +got):\n%s", t.Name(), diff)
			}
		})
	}
}

func TestUserService_Login(t *testing.T) {
	scenarios := []struct {
		name     string
		init     func() *service.UserService
		input    func() (context.Context, service.LoginRequest)
		expected func() (service.LoginResponse, error)
	}{
		{
			name: "OK_Login",
			init: func() (srv *service.UserService) {
				repo := repository.NewMockRepositoryInterface(gomock.NewController(t))

				repo.EXPECT().FindUserByPhone(gomock.Any(), gomock.Any()).Return(&repository.User{
					ID:       10,
					FullName: "Hello",
					Password: "$2a$06$.5OUa9UMKS5qoICzUWi4XOmWDhDGR9wlNu2pNjhdpMTiCTgdhOpBW",
					Phone:    "+6212345678",
				}, nil)

				repo.EXPECT().CreateUserAttendance(gomock.Any(), gomock.Any()).Return(nil)
				repo.EXPECT().SaveUserAttendanceSummary(gomock.Any(), gomock.Any()).Return(nil)

				auth := service.NewMockAuthService(gomock.NewController(t))
				auth.EXPECT().CreateToken(gomock.Any(), gomock.Any()).Return("hello", nil)
				srv = service.NewUserService(service.UserServiceOptions{
					Repository:  repo,
					AuthService: auth,
				})

				return srv
			},
			input: func() (context.Context, service.LoginRequest) {
				ctx := context.Background()

				args := service.LoginRequest{}
				args.Password = "aaaaaaA12!"
				args.Phone = "+6212345678"

				return ctx, args
			},
			expected: func() (service.LoginResponse, error) {
				return service.LoginResponse{
					ID:    10,
					Token: "hello",
				}, nil
			},
		},
		{
			name: "Failed_UserNotFound",
			init: func() (srv *service.UserService) {
				repo := repository.NewMockRepositoryInterface(gomock.NewController(t))

				repo.EXPECT().FindUserByPhone(gomock.Any(), gomock.Any()).Return(nil, nil)

				auth := service.NewMockAuthService(gomock.NewController(t))

				srv = service.NewUserService(service.UserServiceOptions{
					Repository:  repo,
					AuthService: auth,
				})

				return srv
			},
			input: func() (context.Context, service.LoginRequest) {
				ctx := context.Background()

				args := service.LoginRequest{}
				args.Password = "aaaaaaA12!"
				args.Phone = "+6212345678"

				return ctx, args
			},
			expected: func() (service.LoginResponse, error) {
				return service.LoginResponse{
					ID:    0,
					Token: "",
				}, errors.BadRequest("phone or password is wrong")
			},
		},
		{
			name: "Failed_UserIncorrectPassword",
			init: func() (srv *service.UserService) {
				repo := repository.NewMockRepositoryInterface(gomock.NewController(t))

				repo.EXPECT().FindUserByPhone(gomock.Any(), gomock.Any()).Return(&repository.User{
					ID:       10,
					FullName: "Hello",
					Password: "$2a$06$.5OUa9UMKS5qoICzUWi4XOmWDhDGR9wlNu2pNjhdpMTiCTgdhOpBW",
					Phone:    "+6212345678",
				}, nil)

				auth := service.NewMockAuthService(gomock.NewController(t))

				srv = service.NewUserService(service.UserServiceOptions{
					Repository:  repo,
					AuthService: auth,
				})

				return srv
			},
			input: func() (context.Context, service.LoginRequest) {
				ctx := context.Background()

				args := service.LoginRequest{}
				args.Password = "aaabaaA12!"
				args.Phone = "+6212345678"

				return ctx, args
			},
			expected: func() (service.LoginResponse, error) {
				return service.LoginResponse{
					ID:    0,
					Token: "",
				}, errors.BadRequest("phone or password is wrong")
			},
		},
	}

	for _, scn := range scenarios {
		t.Run(scn.name, func(t *testing.T) {
			s := scn.init()

			res, err := s.Login(scn.input())
			expectedRes, expectedErr := scn.expected()

			diff := cmp.Diff(expectedRes, res)
			if diff != "" {
				t.Fatalf("[%s] mismatch (-want +got):\n%s", t.Name(), diff)
			}

			diff = cmp.Diff(expectedErr, err)
			if diff != "" {
				t.Fatalf("[%s] mismatch (-want +got):\n%s", t.Name(), diff)
			}
		})
	}
}

func TestUserService_GetMe(t *testing.T) {
	scenarios := []struct {
		name     string
		init     func() *service.UserService
		input    func() (context.Context, string)
		expected func() (service.GetMeResponse, error)
	}{
		{
			name: "OK_GetMe",
			init: func() (srv *service.UserService) {
				repo := repository.NewMockRepositoryInterface(gomock.NewController(t))

				repo.EXPECT().FindUserByID(gomock.Any(), gomock.Any()).Return(&repository.User{
					ID:       10,
					FullName: "Nivea",
					Password: "$2a$06$.5OUa9UMKS5qoICzUWi4XOmWDhDGR9wlNu2pNjhdpMTiCTgdhOpBW",
					Phone:    "+6212345678",
				}, nil)

				var id int64 = 10
				auth := service.NewMockAuthService(gomock.NewController(t))
				auth.EXPECT().ValidateToken(gomock.Any()).Return(id, nil)
				srv = service.NewUserService(service.UserServiceOptions{
					Repository:  repo,
					AuthService: auth,
				})

				return srv
			},
			input: func() (context.Context, string) {
				ctx := context.Background()

				return ctx, "hello"
			},
			expected: func() (service.GetMeResponse, error) {
				return service.GetMeResponse{
					FullName: "Nivea",
					Phone:    "+6212345678",
				}, nil
			},
		},
		{
			name: "Failed_UserNotFound",
			init: func() (srv *service.UserService) {
				repo := repository.NewMockRepositoryInterface(gomock.NewController(t))

				repo.EXPECT().FindUserByID(gomock.Any(), gomock.Any()).Return(nil, nil)

				var id int64 = 10
				auth := service.NewMockAuthService(gomock.NewController(t))
				auth.EXPECT().ValidateToken(gomock.Any()).Return(id, nil)
				srv = service.NewUserService(service.UserServiceOptions{
					Repository:  repo,
					AuthService: auth,
				})

				return srv
			},
			input: func() (context.Context, string) {
				ctx := context.Background()

				return ctx, "hello"
			},
			expected: func() (service.GetMeResponse, error) {
				return service.GetMeResponse{
					FullName: "",
					Phone:    "",
				}, errors.BadRequest("User Not Found")
			},
		},
	}

	for _, scn := range scenarios {
		t.Run(scn.name, func(t *testing.T) {
			s := scn.init()

			res, err := s.GetMe(scn.input())
			expectedRes, expectedErr := scn.expected()

			diff := cmp.Diff(expectedRes, res)
			if diff != "" {
				t.Fatalf("[%s] mismatch (-want +got):\n%s", t.Name(), diff)
			}

			diff = cmp.Diff(expectedErr, err)
			if diff != "" {
				t.Fatalf("[%s] mismatch (-want +got):\n%s", t.Name(), diff)
			}
		})
	}
}

func TestUserService_UpdateMe(t *testing.T) {
	scenarios := []struct {
		name     string
		init     func() *service.UserService
		input    func() (context.Context, string, service.UpdateRequest)
		expected func() error
	}{
		{
			name: "OK_UpdateMe",
			init: func() (srv *service.UserService) {
				repo := repository.NewMockRepositoryInterface(gomock.NewController(t))

				repo.EXPECT().FindUserByID(gomock.Any(), gomock.Any()).Return(&repository.User{
					ID:       10,
					FullName: "Nivea",
					Password: "$2a$06$.5OUa9UMKS5qoICzUWi4XOmWDhDGR9wlNu2pNjhdpMTiCTgdhOpBW",
					Phone:    "+6212345678",
				}, nil)

				repo.EXPECT().FindUserByPhone(gomock.Any(), gomock.Any()).Return(nil, nil)
				repo.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Return(nil)

				var id int64 = 10
				auth := service.NewMockAuthService(gomock.NewController(t))
				auth.EXPECT().ValidateToken(gomock.Any()).Return(id, nil)
				srv = service.NewUserService(service.UserServiceOptions{
					Repository:  repo,
					AuthService: auth,
				})

				return srv
			},
			input: func() (context.Context, string, service.UpdateRequest) {
				ctx := context.Background()

				ps := "+6212345679"
				return ctx, "hello", service.UpdateRequest{
					Phone: &ps,
				}
			},
			expected: func() error {
				return nil
			},
		},
		{
			name: "Failed_UserNotFound",
			init: func() (srv *service.UserService) {
				repo := repository.NewMockRepositoryInterface(gomock.NewController(t))

				repo.EXPECT().FindUserByID(gomock.Any(), gomock.Any()).Return(nil, nil)

				var id int64 = 10
				auth := service.NewMockAuthService(gomock.NewController(t))
				auth.EXPECT().ValidateToken(gomock.Any()).Return(id, nil)
				srv = service.NewUserService(service.UserServiceOptions{
					Repository:  repo,
					AuthService: auth,
				})

				return srv
			},
			input: func() (context.Context, string, service.UpdateRequest) {
				ctx := context.Background()

				ps := "+6212345679"
				return ctx, "hello", service.UpdateRequest{
					Phone: &ps,
				}
			},
			expected: func() error {
				return errors.BadRequest("User Not Found")
			},
		},
		{
			name: "Failed_ConflictedPhone",
			init: func() (srv *service.UserService) {
				repo := repository.NewMockRepositoryInterface(gomock.NewController(t))

				repo.EXPECT().FindUserByID(gomock.Any(), gomock.Any()).Return(&repository.User{
					ID:       10,
					FullName: "Nivea",
					Password: "$2a$06$.5OUa9UMKS5qoICzUWi4XOmWDhDGR9wlNu2pNjhdpMTiCTgdhOpBW",
					Phone:    "+6212345678",
				}, nil)

				repo.EXPECT().FindUserByPhone(gomock.Any(), gomock.Any()).Return(&repository.User{
					ID:       11,
					FullName: "Nivea",
					Password: "$2a$06$.5OUa9UMKS5qoICzUWi4XOmWDhDGR9wlNu2pNjhdpMTiCTgdhOpBW",
					Phone:    "+6212345679",
				}, nil)

				var id int64 = 10
				auth := service.NewMockAuthService(gomock.NewController(t))
				auth.EXPECT().ValidateToken(gomock.Any()).Return(id, nil)
				srv = service.NewUserService(service.UserServiceOptions{
					Repository:  repo,
					AuthService: auth,
				})

				return srv
			},
			input: func() (context.Context, string, service.UpdateRequest) {
				ctx := context.Background()

				ps := "+6212345679"
				return ctx, "hello", service.UpdateRequest{
					Phone: &ps,
				}
			},
			expected: func() error {
				e := errors.Conflicted("")
				e.AddField("phone", fmt.Sprintf("[%s] already registered", "+6212345679"))
				return e
			},
		},
	}

	for _, scn := range scenarios {
		t.Run(scn.name, func(t *testing.T) {
			s := scn.init()

			err := s.UpdateMe(scn.input())
			expectedErr := scn.expected()

			diff := cmp.Diff(expectedErr, err)
			if diff != "" {
				t.Fatalf("[%s] mismatch (-want +got):\n%s", t.Name(), diff)
			}
		})
	}
}

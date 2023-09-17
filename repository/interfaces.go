// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import "context"

type RepositoryInterface interface {
	CreateUser(ctx context.Context, u *User) (int64, error)
	FindUserByPhone(ctx context.Context, phone string) (*User, error)
	FindUserByID(ctx context.Context, id int64) (*User, error)
	CreateUserAttendance(ctx context.Context, ua *UserAttendance) error
	SaveUserAttendanceSummary(ctx context.Context, userID int64) error

	UpdateUser(ctx context.Context, u *User) error
}

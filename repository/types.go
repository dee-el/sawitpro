// This file contains types that are used in the repository layer.
package repository

import "time"

type GetTestByIdInput struct {
	Id string
}

type GetTestByIdOutput struct {
	Name string
}

type User struct {
	ID        int64
	FullName  string
	Password  string
	Phone     string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
type UserAttendance struct {
	UserID  int64
	LoginAt time.Time
}

type UserAttendanceSummary struct {
	UserID        int64
	TotalLoggedIn int
}

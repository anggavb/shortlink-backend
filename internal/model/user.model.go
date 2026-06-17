package model

import "time"

type User struct {
	Id             int        `db:"id"`
	Email          string     `db:"email"`
	Password       string     `db:"password"`
	FirstName      *string    `db:"first_name"`
	LastName       *string    `db:"last_name"`
	Role           *string    `db:"role"`
	Workplace      *string    `db:"workplace"`
	IsMember       bool       `db:"is_member"`
	IsReceiveEmail bool       `db:"is_receive_email"`
	Photo          *string    `db:"photo"`
	VerifiedAt     *string    `db:"verified_at"`
	CreatedAt      *time.Time `db:"created_at"`
	UpdatedAt      *time.Time `db:"updated_at"`
}

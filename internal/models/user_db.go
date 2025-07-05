package models

import "time"

type UserDB struct {
	ID           int
	Username     string
	PasswordHash string
	CreatedAt    time.Time
}

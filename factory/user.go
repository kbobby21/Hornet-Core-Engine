package factory

import "time"

type User struct {
	Email          string `json:"email"`
	Password       string `json:"password"`
	Organization   string `json:"organization"`
	HashedPassword string `json:"hashed_password"`
	Token          string `json:"token"`
	ValidTill  time.Time `json:"valid_till"`
}

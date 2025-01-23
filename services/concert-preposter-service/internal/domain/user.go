package domain

import "time"

type User struct {
	Id				int
	Name 			string
	Email 			string
	Hashed_password string
	Phone_number 	string
	Birthday 		time.Time
}


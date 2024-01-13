package user

import "time"

type User struct {
	ID             string
	FirstName      string
	SecondName     string
	Birthdate      time.Time
	Biography      string
	City           string
	HashedPassword string
}

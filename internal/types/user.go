package types

import "time"

type User struct {
	ID        int
	Name      string
	Lastname  string
	Email     string
	Password  string
	CreatedAt time.Time
}

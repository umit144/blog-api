package types

import (
	"encoding/json"
	"time"
)

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Lastname  string    `json:"lastname"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
}

func (u User) MarshalJSON() ([]byte, error) {
	type Alias User
	return json.Marshal(&struct {
		*Alias
		CreatedAt string `json:"createdAt"`
	}{
		Alias:     (*Alias)(&u),
		CreatedAt: u.CreatedAt.Format(time.RFC1123),
	})
}

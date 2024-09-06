package types

import (
	"encoding/json"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	Id        string    `json:"id"`
	Name      string    `json:"name" validate:"required,min=3,max=50"`
	Lastname  string    `json:"lastname" validate:"required,min=3,max=50"`
	Email     string    `json:"email" validate:"required,email"`
	Password  string    `json:"password,omitempty" validate:"required,min=8,max=150"`
	CreatedAt time.Time `json:"createdAt"`
}

func (u User) Validate() map[string]string {
	v := validator.New()
	err := v.Struct(u)
	if err == nil {
		return nil
	}

	errorsMap := make(map[string]string)

	for _, err := range err.(validator.ValidationErrors) {
		errorsMap[err.Field()] = err.Tag()
	}

	return errorsMap
}

func (u User) MarshalJSON() ([]byte, error) {
	type Alias User
	return json.Marshal(&struct {
		*Alias
		Password  interface{} `json:"password,omitempty"` // Exclude Password field from JSON
		CreatedAt string      `json:"createdAt"`
	}{
		Alias:     (*Alias)(&u),
		Password:  nil,
		CreatedAt: u.CreatedAt.Format(time.RFC1123),
	})
}

func (u User) CreateToken() (*string, error) {
	var secretKey []byte = []byte(os.Getenv("JWT_SECRET"))
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": u.Id,
		"iss": "go-blog",
		"aud": "user-role",
		"exp": time.Now().Add(time.Hour).Unix(),
		"iat": time.Now().Unix(),
	})

	tokenString, err := claims.SignedString(secretKey)
	if err != nil {
		return nil, err
	}

	return &tokenString, nil
}

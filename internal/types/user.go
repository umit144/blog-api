package types

import (
	"encoding/json"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	Id             string    `json:"id,omitempty" db:"id"`
	Name           string    `json:"name,omitempty" validate:"required,min=3,max=50" db:"name"`
	Lastname       string    `json:"lastname,omitempty" validate:"omitempty,min=3,max=50" db:"lastname"`
	Email          string    `json:"email,omitempty" validate:"required,email" db:"email"`
	Password       string    `json:"password,omitempty" validate:"omitempty,min=8,max=150" db:"password"`
	GoogleID       string    `json:"google_id,omitempty" db:"google_id"`
	ProfilePicture string    `json:"profile_picture,omitempty" db:"profile_picture"`
	AuthProvider   string    `json:"auth_provider,omitempty" db:"auth_provider"`
	CreatedAt      time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at,omitempty" db:"updated_at"`
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
	aux := struct {
		*Alias
		Password       interface{} `json:"password,omitempty"`
		CreatedAt      *string     `json:"created_at,omitempty"`
		UpdatedAt      *string     `json:"updated_at,omitempty"`
		GoogleID       interface{} `json:"google_id,omitempty"`
		ProfilePicture interface{} `json:"profile_picture,omitempty"`
	}{
		Alias:    (*Alias)(&u),
		Password: nil,
	}

	if !u.CreatedAt.IsZero() {
		createdAt := u.CreatedAt.Format(time.RFC3339)
		aux.CreatedAt = &createdAt
	}

	if !u.UpdatedAt.IsZero() {
		updatedAt := u.UpdatedAt.Format(time.RFC3339)
		aux.UpdatedAt = &updatedAt
	}

	if u.GoogleID == "" {
		aux.GoogleID = nil
	}

	if u.ProfilePicture == "" {
		aux.ProfilePicture = nil
	}

	return json.Marshal(aux)
}

func (u User) CreateToken() (*string, error) {
	var secretKey []byte = []byte(os.Getenv("JWT_SECRET"))
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":           u.Id,
		"iss":           "go-blog",
		"aud":           "user-role",
		"exp":           time.Now().Add(time.Hour * 24).Unix(),
		"iat":           time.Now().Unix(),
		"auth_provider": u.AuthProvider,
	})

	tokenString, err := claims.SignedString(secretKey)
	if err != nil {
		return nil, err
	}

	return &tokenString, nil
}

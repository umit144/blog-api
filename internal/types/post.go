package types

import (
	"github.com/go-playground/validator/v10"
	"time"
)

type Post struct {
	Id        string    `json:"id,omitempty"`
	Title     string    `json:"title,omitempty" validate:"required,min=3,max=50"`
	Slug      string    `json:"slug,omitempty"`
	Content   string    `json:"content,omitempty"  validate:"required,min=3"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	Author    User      `json:"author,omitempty" validate:"-"`
}

func (p Post) Validate() map[string]string {
	v := validator.New()
	err := v.Struct(p)
	if err == nil {
		return nil
	}

	errorsMap := make(map[string]string)

	for _, err := range err.(validator.ValidationErrors) {
		errorsMap[err.Field()] = err.Tag()
	}

	return errorsMap
}

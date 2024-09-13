package types

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type Category struct {
	Id        string    `json:"id,omitempty"`
	Title     string    `json:"title,omitempty" validate:"required,min=3,max=50"`
	Slug      string    `json:"slug,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
}

func (c Category) Validate() map[string]string {
	v := validator.New()
	err := v.Struct(c)
	if err == nil {
		return nil
	}

	errorsMap := make(map[string]string)

	for _, err := range err.(validator.ValidationErrors) {
		errorsMap[err.Field()] = err.Tag()
	}

	return errorsMap
}

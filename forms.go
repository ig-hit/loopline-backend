package main

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type RegisterForm struct {
	Name string `json:"name"`
}

func (f RegisterForm) Validate() error {
	return validation.ValidateStruct(&f,
		validation.Field(&f.Name, validation.NotNil, validation.Length(3, 20)))
}

type NotebookForm struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

func (f NotebookForm) Validate() error {
	return validation.ValidateStruct(&f,
		validation.Field(&f.Title, validation.NotNil, validation.Length(3, 20)))
}

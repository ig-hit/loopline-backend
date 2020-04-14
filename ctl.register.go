package main

import (
	"github.com/labstack/echo/v4"
	"math/rand"
	"net/http"
	"time"
)

type RegisterStore map[string]string

type RegisterController struct {
	*BaseController
	store RegisterStore
}

func (ctl *RegisterController) Register(c echo.Context) error {
	form := &RegisterForm{}
	if err := c.Bind(&form); err != nil {
		return ErrInvalidForm
	}

	if err := form.Validate(); err != nil {
		return err
	}

	name := form.Name
	rand.Seed(time.Now().UnixNano())

	const bytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	token := make([]byte, 36)
	for i := range token {
		token[i] = bytes[rand.Int63() % int64(len(bytes))]
	}

	res := string(token)

	ctl.store[res] = name

	_ = c.JSON(http.StatusCreated, res)

	return nil
}

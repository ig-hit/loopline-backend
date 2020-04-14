package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func MustHaveToken(tokenStorage RegisterStore) echo.MiddlewareFunc {
	return middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		Validator: func(token string, context echo.Context) (b bool, e error) {
			name, found := tokenStorage[token]

			if !found {
				return false, ErrInvalidToken
			}

			context.Set("username", name)
			context.Set("token", token)

			return true, nil
		},
	})
}

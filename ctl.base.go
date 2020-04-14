package main

import "github.com/labstack/echo/v4"

type Action func(echo.Context) error

type BaseController struct {
}

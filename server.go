package main

import (
	"flag"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	port := flag.Int("port", 64567, "listen at port")
	flag.Parse()

	host := NewHost()

	e := echo.New()
	e.Use(middleware.CORS())
	e.Any("/*", func(c echo.Context) (err error) {
		req := c.Request()
		res := c.Response()
		host.Echo.ServeHTTP(res, req)
		return
	})
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", *port)))
}

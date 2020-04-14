package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	RegisterRoute       = "register"
	CreateNotebookRoute = "create"
	ReadNotebookRoute   = "read"
	UpdateNotebookRoute = "update"
	DeleteNotebookRoute = "delete"
	ListNotebookRoute   = "list"
)

func InitRoutes() Routes {
	var (
		registerCtl = MustRegisterController()
		notebookCtl = MustNotebookController()
	)

	routes := Routes{
		POST("/register", registerCtl.Register, RegisterRoute),

		POST("/notebooks", notebookCtl.Create, CreateNotebookRoute),
		GET("/notebooks/:id", notebookCtl.Read, ReadNotebookRoute),
		PATCH("/notebooks/:id", notebookCtl.Update, UpdateNotebookRoute),
		DELETE("/notebooks/:id", notebookCtl.Delete, DeleteNotebookRoute),
		GET("/notebooks", notebookCtl.List, ListNotebookRoute),
	}

	return routes
}

func public() []string {
	return []string{
		RegisterRoute,
	}
}

func protected() []string {
	return []string{
		CreateNotebookRoute,
		ReadNotebookRoute,
		UpdateNotebookRoute,
		DeleteNotebookRoute,
		ListNotebookRoute,
	}
}

func NewHost() *Host {
	host := echo.New()
	host.HTTPErrorHandler = ErrorHandler

	commonMw := []echo.MiddlewareFunc{
		middleware.Logger(),
		middleware.Recover(),
		middleware.CORS(),
	}

	routes := InitRoutes()

	publicGroup := &Group{Group: host.Group("")}
	publicGroup.Use(commonMw...)

	protectedGroup := &Group{Group: host.Group("")}
	protectedGroup.Use(commonMw...)
	protectedGroup.Use(CreateTokenMiddleware())

	var rs []*Route
	rs = routes.Register(public(), publicGroup)
	rs = append(rs, routes.Register(protected(), protectedGroup)...)

	return &Host{
		Echo:   host,
		Routes: Map(rs),
	}
}

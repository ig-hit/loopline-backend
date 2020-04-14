package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"strings"
)

type Route struct {
	Name       string
	Path       string
	Method     string
	Action     Action
	Middleware []echo.MiddlewareFunc
}

type Routes []*Route

type Group struct {
	*echo.Group
}

type Host struct {
	Echo   *echo.Echo
	Routes map[string]*Route
}

func (g *Group) Register(route *Route) {
	if route.Method == "" {
		route.Method = "GET"
	}

	path := fmt.Sprintf("/%s", strings.Trim(route.Path, "/"))
	registered := g.Add(route.Method, path, echo.HandlerFunc(route.Action), route.Middleware...)
	registered.Name = route.Name
}

func (rs *Routes) Register(names []string, g *Group) []*Route {
	all := Map(*rs)
	var res []*Route

	for _, name := range names {
		route, exists := all[name]
		if !exists {
			fmt.Printf("!Route %s is referenced by name but not found in initiated route list!\n", name)
			continue
		}
		g.Register(route)
		res = append(res, route)
	}

	return res
}

func Map(rs Routes) map[string]*Route {
	res := make(map[string]*Route)
	for _, r := range rs {
		res[r.Name] = r
	}

	return res
}

func GET(path string, a Action, name string) *Route {
	return METHOD("GET", path, a, name)
}

func POST(path string, a Action, name string) *Route {
	return METHOD("POST", path, a, name)
}

func PUT(path string, a Action, name string) *Route {
	return METHOD("PUT", path, a, name)
}

func PATCH(path string, a Action, name string) *Route {
	return METHOD("PATCH", path, a, name)
}

func DELETE(path string, a Action, name string) *Route {
	return METHOD("DELETE", path, a, name)
}

func METHOD(method, path string, a Action, name string) *Route {
	return &Route{
		Path:   path,
		Method: method,
		Action: a,
		Name:   name,
	}
}

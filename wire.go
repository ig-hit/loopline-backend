//+build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/labstack/echo/v4"
	"sync"
)

var onceRegisterStore sync.Once
var registerStore RegisterStore

var onceNotebookStore sync.Once
var notebookStore NotebookStore

var onceBase sync.Once
var baseCtl *BaseController

func CreateRegisterStore() RegisterStore {
	onceRegisterStore.Do(func() {
		registerStore = make(RegisterStore)
	});
	return registerStore
}

func CreateNotebookStore() NotebookStore {
	onceNotebookStore.Do(func() {
		notebookStore = make(NotebookStore)
	});
	return notebookStore
}

func CreateTokenMiddleware() echo.MiddlewareFunc {
	panic(wire.Build(
		CreateRegisterStore,
		MustHaveToken))
}

func createBase(c ...string) (*BaseController, error) {
	panic(wire.Build(
		wire.Struct(new(BaseController), "*")))
}

func MustBase(c ...string) *BaseController {
	var err error
	onceBase.Do(func() {
		baseCtl, err = createBase(c...)
		if err != nil {
			panic(err)
		}
	})

	return baseCtl
}

func createNotebookController(c ...string) (*NotebookController, error) {
	panic(wire.Build(
		MustBase,
		CreateNotebookStore,
		wire.Struct(new(NotebookController), "*")))
}

func MustNotebookController(c ...string) *NotebookController {
	ctl, err := createNotebookController(c...)
	if err != nil {
		panic(err)
	}
	return ctl
}

func createRegisterController(c ...string) (*RegisterController, error) {
	panic(wire.Build(
		MustBase,
		CreateRegisterStore,
		wire.Struct(new(RegisterController), "*")))
}

func MustRegisterController(c ...string) *RegisterController {
	ctl, err := createRegisterController(c...)
	if err != nil {
		panic(err)
	}
	return ctl
}

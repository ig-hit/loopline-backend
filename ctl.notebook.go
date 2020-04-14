package main

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"time"
)

type NotebookStore map[string]map[int]*Notebook

type NotebookController struct {
	*BaseController
	store NotebookStore
}

type Notebook struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Text  string `json:"text"`
}

func (ctl *NotebookController) List(c echo.Context) error {
	token := fmt.Sprintf("%v", c.Get("token"))
	books := ctl.store[token]

	if books == nil {
		books = make(map[int]*Notebook, 0)
	}

	res := make([]*Notebook, 0)
	for _, book := range books {
		res = append(res, book)
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].ID < res[j].ID
	})

	err := c.JSON(http.StatusOK, res)
	if err != nil {
		c.Logger().Error(err)
	}

	return nil
}

func (ctl *NotebookController) Create(c echo.Context) error {
	token := fmt.Sprintf("%v", c.Get("token"))
	books := ctl.store[token]
	total := len(books)

	form := &NotebookForm{}
	if err := c.Bind(&form); err != nil {
		return ErrInvalidForm
	}

	if err := form.Validate(); err != nil {
		return err
	}

	rand.Seed(int64(time.Now().Nanosecond()))
	view := &Notebook{
		ID:    total + 1,
		Title: form.Title,
		Text:  form.Text,
	}

	if books == nil {
		books = make(map[int]*Notebook)
	}

	books[view.ID] = view
	ctl.store[token] = books

	err := c.JSON(http.StatusCreated, view)
	if err != nil {
		c.Logger().Error(err)
	}

	return nil
}

func (ctl *NotebookController) Update(c echo.Context) error {
	token := fmt.Sprintf("%v", c.Get("token"))
	books := ctl.store[token]

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	book, found := books[id]
	if !found {
		return ErrNotFound
	}

	form := &NotebookForm{}

	bookJs, _ := json.Marshal(book)
	err = json.Unmarshal(bookJs, form)
	if err != nil {
		return ErrInvalidForm
	}

	if err := c.Bind(form); err != nil {
		return ErrInvalidForm
	}

	if err := form.Validate(); err != nil {
		return err
	}

	book.Title = form.Title
	book.Text = form.Text

	books[id] = book

	err = c.JSON(http.StatusOK, book)
	if err != nil {
		c.Logger().Error(err)
	}

	return nil
}

func (ctl *NotebookController) Read(c echo.Context) error {
	token := fmt.Sprintf("%v", c.Get("token"))
	books := ctl.store[token]

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	book, found := books[id]
	if !found {
		return ErrNotFound
	}

	err = c.JSON(http.StatusOK, book)
	if err != nil {
		c.Logger().Error(err)
	}

	return nil
}

func (ctl *NotebookController) Delete(c echo.Context) error {
	token := fmt.Sprintf("%v", c.Get("token"))
	books := ctl.store[token]

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	_, found := books[id]
	if !found {
		return ErrNotFound
	}

	delete(books, id)

	err = c.NoContent(http.StatusOK)
	if err != nil {
		c.Logger().Error(err)
	}

	return nil
}

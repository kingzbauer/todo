package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func init() {
	_, err := GetDatabase()
	if err != nil {
		log.Panic(err)
	}
}

func list(c echo.Context) error {
	todos := List()
	return c.JSON(http.StatusOK, todos)
}

func create(c echo.Context) (err error) {
	todo := &Todo{}

	if err = c.Bind(todo); err != nil {
		return
	}

	if err = c.Validate(todo); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	_, err = InsertOne(todo)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, todo)
}

func retrieve(c echo.Context) error {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.String(http.StatusNotFound, "Not found")
	}

	filter := bson.D{{"_id", objectID}}
	v, err := Retrieve(filter, Todo{})
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, v)
}

func deleteTodo(c echo.Context) error {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.String(http.StatusNotFound, "Not found")
	}

	filter := bson.D{{"_id", objectID}}
	v, err := Retrieve(filter, Todo{})
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}

	err = Delete(v.(FilterID))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}

func start(port int) {
	e := echo.New()

	e.Use(middleware.CORS())

	e.Validator = &CustomValidator{validator: validator.New()}
	e.GET("/api/v1/todos/", list)
	e.POST("/api/v1/todos/", create)
	e.GET("/api/v1/todos/:id", retrieve)
	e.DELETE("/api/v1/todos/:id", deleteTodo)

	e.Logger.Fatal(e.Start(fmt.Sprintf("0.0.0.0:%d", port)))
}

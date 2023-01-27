package main

import (
	"context"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/toshiykst/go-layerd-architecture/app/handler"
	"github.com/toshiykst/go-layerd-architecture/app/infrastructure/database"
	"github.com/toshiykst/go-layerd-architecture/app/usecase"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	ctx := context.Background()

	// TODO: Use env config
	db := database.NewDBRepository(ctx, database.Config{
		User:     os.Getenv("MYSQL_USER"),
		Password: os.Getenv("MYSQL_PASSWORD"),
		Host:     os.Getenv("MYSQL_HOST"),
		DBName:   os.Getenv("MYSQL_DATABASE"),
		Debug:    os.Getenv("MYSQL_DEBUG") == "true",
	})

	uuc := usecase.NewUserUsecase(db)
	uh := handler.NewUserHandler(uuc)

	e.GET("/users/:id", uh.GetUser)

	e.Logger.Fatal(e.Start(":8080"))
}

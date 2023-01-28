package main

import (
	"context"
	"log"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/toshiykst/go-layerd-architecture/app/env"
	"github.com/toshiykst/go-layerd-architecture/app/handler"
	"github.com/toshiykst/go-layerd-architecture/app/infrastructure/database"
	"github.com/toshiykst/go-layerd-architecture/app/usecase"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	ctx := context.Background()

	c, err := env.NewConfig()
	if err != nil {
		log.Fatal(err.Error())
	}

	db := database.NewDBRepository(ctx, database.Config{
		User:     c.DBUser,
		Password: c.DBPassword,
		Host:     c.DBHost,
		DBName:   c.DBName,
		Debug:    c.DBDebug,
	})

	uuc := usecase.NewUserUsecase(db)
	uh := handler.NewUserHandler(uuc)

	e.GET("/users/:id", uh.GetUser)

	e.Logger.Fatal(e.Start(":8080"))
}

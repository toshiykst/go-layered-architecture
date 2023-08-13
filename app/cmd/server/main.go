package main

import (
	"context"
	"log"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/toshiykst/go-layerd-architecture/app/domain/domainservice"
	"github.com/toshiykst/go-layerd-architecture/app/domain/factory"
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

	uf := factory.NewUserFactory()
	us := domainservice.NewUserService(db)
	uuc := usecase.NewUserUsecase(db, uf, us)
	uh := handler.NewUserHandler(uuc)

	e.POST("/users", uh.CreateUser)
	e.GET("/users/:id", uh.GetUser)
	e.GET("/users", uh.GetUsers)
	e.PUT("/users/:id", uh.UpdateUser)
	e.DELETE("/users/:id", uh.DeleteUser)

	gf := factory.NewGroupFactory()
	gs := domainservice.NewGroupService(db)
	guc := usecase.NewGroupUsecase(db, gf, gs, us)
	gh := handler.NewGroupHandler(guc)

	e.POST("/groups", gh.CreateGroup)
	e.GET("/groups/:id", gh.GetGroup)
	e.GET("/groups", gh.GetGroups)
	e.PUT("/groups/:id", gh.UpdateGroup)

	e.Logger.Fatal(e.Start(":8080"))
}

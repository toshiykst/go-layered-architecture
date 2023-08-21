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
	gf := factory.NewGroupFactory()

	us := domainservice.NewUserService(db)
	gs := domainservice.NewGroupService(db)

	uuc := usecase.NewUserUsecase(db, uf, us, gs)
	uh := handler.NewUserHandler(uuc)

	guc := usecase.NewGroupUsecase(db, gf, gs, us)
	gh := handler.NewGroupHandler(guc)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/users", uh.CreateUser)
	e.GET("/users/:id", uh.GetUser)
	e.GET("/users", uh.GetUsers)
	e.PUT("/users/:id", uh.UpdateUser)
	e.DELETE("/users/:id", uh.DeleteUser)

	e.POST("/groups", gh.CreateGroup)
	e.GET("/groups/:id", gh.GetGroup)
	e.GET("/groups", gh.GetGroups)
	e.PUT("/groups/:id", gh.UpdateGroup)
	e.DELETE("/groups/:id", gh.DeleteGroup)

	e.Logger.Fatal(e.Start(":8080"))
}

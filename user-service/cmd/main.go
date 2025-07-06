package main

import (
	"github.com/Gergenus/commerce/user-service/internal/config"
	"github.com/Gergenus/commerce/user-service/internal/handlers"
	"github.com/Gergenus/commerce/user-service/internal/repository"
	"github.com/Gergenus/commerce/user-service/internal/service"
	dbpkg "github.com/Gergenus/commerce/user-service/pkg/db"
	"github.com/Gergenus/commerce/user-service/pkg/jwtpkg"
	"github.com/Gergenus/commerce/user-service/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg := config.InitConfig()
	db := dbpkg.InitDB(cfg.PostgresURL)
	defer db.DB.Close()
	log := logger.SetUp(cfg.LogLevel)

	repo := repository.NewPostgresRepository(db)
	jwtstr := jwtpkg.NewUserJWTpkg(cfg.JWTSecret)
	srv := service.NewUserService(log, &repo, jwtstr)
	handler := handlers.NewUserHandler(&srv)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	group := e.Group("/api/v1/users")
	{
		group.POST("/auth/register", handler.Register)
		group.POST("/auth/login", handler.Login)
		group.POST("/auth/refresh", handler.Refresh)
	}
	e.Start(":" + cfg.HTTPPort)
}

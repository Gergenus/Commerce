package main

import (
	"context"

	"github.com/Gergenus/commerce/product-service/internal/config"
	"github.com/Gergenus/commerce/product-service/internal/handlers"
	"github.com/Gergenus/commerce/product-service/internal/repository"
	"github.com/Gergenus/commerce/product-service/internal/service"
	dbpkg "github.com/Gergenus/commerce/product-service/pkg/db"
	"github.com/Gergenus/commerce/product-service/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg := config.InitConfig()
	db := dbpkg.InitDB(cfg.PostgresURL)
	defer db.DB.Close(context.Background())
	log := logger.SetupLogger(cfg.LogLevel)

	repo := repository.NewPostgresRepository(db)
	serv := service.NewProductService(log, &repo)
	hand := handlers.NewProductHandler(&serv)
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	group := e.Group("/api/v1/products")
	{
		group.POST("/", hand.AddCategory)
		group.POST("/create", hand.CreateProduct)   // create product
		group.GET("/", hand.GetProductByID)         // get product by id
		group.POST("/stock/add", hand.AddStockByID) // add stock by id
		group.GET("/stock", hand.GetStockByID)      // get stock by id

	}

	e.Start(":" + cfg.HTTPPort)

}

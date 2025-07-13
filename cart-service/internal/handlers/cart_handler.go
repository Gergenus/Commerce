package handlers

import (
	"net/http"

	"github.com/Gergenus/commerce/cart-service/internal/models"
	"github.com/Gergenus/commerce/cart-service/internal/service"
	"github.com/labstack/echo/v4"
)

type CartHandler struct {
	srv service.CartServiceInterface
}

func NewCartHandler(srv service.CartServiceInterface) *CartHandler {
	return &CartHandler{srv: srv}
}

func (ch *CartHandler) AddToCart(c echo.Context) error {
	UUID := c.Get("uuid")
	UUIDString, ok := UUID.(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "internal error",
		})
	}

	var data models.AddCartRequest
	err := c.Bind(&data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "invalid payload",
		})
	}

	err = ch.srv.AddToCart(c.Request().Context(), UUIDString, data.ProductId, data.Stock)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "internal error",
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
	})
}

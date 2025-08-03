package handlers

import (
	"errors"
	"net/http"

	"github.com/Gergenus/commerce/order-service/internal/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type OrderHandler struct {
	srv service.OrderServiceInterface
}

func NewOrderHandler(srv service.OrderServiceInterface) OrderHandler {
	return OrderHandler{srv: srv}
}

func (o OrderHandler) CreateOrder(c echo.Context) error {
	uid := c.Get("uuid").(string)

	orderID, err := o.srv.CreateOrder(c.Request().Context(), uuid.MustParse(uid))
	if err != nil {
		if errors.Is(err, service.ErrOrderNotReserved) {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "insufficient stock",
			})
		}
		if errors.Is(err, service.ErrNoCartFound) {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "no cart found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal server error",
		})
	}
	return c.JSON(http.StatusOK, map[string]any{
		"order_id": orderID,
	})
}

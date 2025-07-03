package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/Gergenus/commerce/user-service/internal/models"
	"github.com/Gergenus/commerce/user-service/internal/service"
	"github.com/labstack/echo/v4"
)

const AccessTokenDuration = 60 * 24 * 60 * 60
const RefreshTokenDuration = 60 * 24 * 60 * 60

type UserHandler struct {
	srv service.UserServiceInterface
}

func NewUserHandler(srv service.UserServiceInterface) UserHandler {
	return UserHandler{srv: srv}
}

func (u *UserHandler) Register(c echo.Context) error {
	var user models.User

	err := c.Bind(&user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "invalid payload",
		})
	}
	uid, err := u.srv.AddUser(c.Request().Context(), user)
	if err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": "user already exists",
			})
		}
		if errors.Is(err, service.ErrIncorrectEmail) {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": "incorrect email",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "internal error",
		})
	}
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"uid": uid.String(),
	})
}

func (u *UserHandler) Login(c echo.Context) error {
	var loginReq models.LoginRequest

	err := c.Bind(&loginReq)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "invalid payload",
		})
	}
	ip := c.RealIP()
	AccessToken, RefreshToken, err := u.srv.Login(c.Request().Context(), loginReq.Email, loginReq.Password, c.Request().UserAgent(), ip)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "internal error",
		})
	}
	err = setCookie(c, "AccessToken", AccessToken, AccessTokenDuration)
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}

	err = setCookie(c, "RefreshToken", RefreshToken, RefreshTokenDuration)
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}
	return c.JSON(http.StatusAccepted, map[string]interface{}{
		"AccessToken":  AccessToken,
		"RefreshToken": RefreshToken,
	})
}

func setCookie(c echo.Context, key, value string, duration int) error {
	cookie := http.Cookie{
		Name:     key,
		HttpOnly: true,
		Secure:   true,
		MaxAge:   duration,
		Value:    value,
		Path:     "/api/v1",
	}
	c.SetCookie(&cookie)
	return nil
}

func (u *UserHandler) Test(c echo.Context) error {
	AccessToken, err := c.Cookie("AccessToken")
	if err != nil {
		return err
	}
	RefreshToken, err := c.Cookie("RefreshToken")
	if err != nil {
		return err
	}
	return c.JSON(200, map[string]interface{}{
		"AccessToken":  AccessToken.Value,
		"RefreshToken": RefreshToken.Value,
	})
}

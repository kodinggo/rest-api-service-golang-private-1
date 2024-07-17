package httpsvc

import (
	"net/http"

	"github.com/kodinggo/rest-api-service-golang-private-1/internal/model"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authUsecase model.AuthUsecase
}

func NewAuthHandler(authUsecase model.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUsecase: authUsecase}
}

func (h *AuthHandler) RegisterRoutes(e *echo.Echo) {
	g := e.Group("/auth")

	g.Use(checkErrorMiddleware)

	g.POST("/login", h.login, basicAuthMiddleware)
}

func (h *AuthHandler) login(c echo.Context) error {
	authMap, ok := c.Request().Context().Value(model.LoginKey).(map[string]string)
	if !ok {
		return model.NewErrorUnAuthorized("invalid token")
	}

	username := authMap["username"]
	password := authMap["password"]
	if username == "" || password == "" {
		return model.NewErrorUnAuthorized("invalid token")
	}

	token, err := h.authUsecase.Login(c.Request().Context(), username, password)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{
		"access_token": token,
	})
}

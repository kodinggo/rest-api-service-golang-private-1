package httpsvc

import (
	"context"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/kodinggo/rest-api-service-golang-private-1/internal/model"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func authMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// ambil token dari header
		authHeader := c.Request().Header.Get(echo.HeaderAuthorization)
		splitAuth := strings.Split(authHeader, " ") // Basic token
		if len(splitAuth) != 2 {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
		}
		if splitAuth[0] != "Basic" {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
		}

		// decode base64 ke string
		decodedToken, err := base64.StdEncoding.DecodeString(splitAuth[1])
		if err != nil {
			log.Errorf("failed when decode token from base64")
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
		}

		// split untuk mendapatkan username & password
		splitDecodedToken := strings.Split(string(decodedToken), ":") // base64(username:password)
		if len(splitDecodedToken) != 2 {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
		}

		username := splitDecodedToken[0]
		password := splitDecodedToken[1]

		ctx := context.WithValue(c.Request().Context(), model.LoginKey, map[string]string{
			"username": username,
			"password": password,
		})

		req := c.Request().WithContext(ctx)
		c.SetRequest(req)

		return next(c)
	}
}

package httpsvc

import (
	"context"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kodinggo/rest-api-service-golang-private-1/internal/config"
	"github.com/kodinggo/rest-api-service-golang-private-1/internal/model"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func basicAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
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

func bearerAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// ambil token dari header
		authHeader := c.Request().Header.Get(echo.HeaderAuthorization)
		splitAuth := strings.Split(authHeader, " ") // Basic token
		if len(splitAuth) != 2 {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
		}
		if splitAuth[0] != "Bearer" {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
		}

		// decode jwt ke custom claims
		strToken := splitAuth[1]
		token, err := jwt.ParseWithClaims(strToken, &model.CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(config.JWTSigningKey()), nil
		})
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
		}

		if !token.Valid {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
		}

		customClaims, ok := token.Claims.(*model.CustomClaims)
		if !ok {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
		}

		ctx := context.WithValue(c.Request().Context(), model.JWTKey, customClaims)

		req := c.Request().WithContext(ctx)
		c.SetRequest(req)

		return next(c)
	}
}

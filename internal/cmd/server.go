package cmd

import (
	"net/http"

	"github.com/kodinggo/rest-api-service-golang-private-1/internal/config"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
)

var serverCMD = &cobra.Command{
	Use:   "serve",
	Short: "serve is a command to run the service",
	Run: func(cmd *cobra.Command, args []string) {
		// Run server
		e := echo.New()

		// It's used to verify
		e.GET("/ping", func(c echo.Context) error {
			return c.String(http.StatusOK, "pong!")
		})

		e.Logger.Fatal(e.Start(config.Port()))
	},
}

func init() {
	rootCMD.AddCommand(serverCMD)
}

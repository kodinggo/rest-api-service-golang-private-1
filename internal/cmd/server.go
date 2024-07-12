package cmd

import (
	"net/http"

	"github.com/kodinggo/rest-api-service-golang-private-1/internal/config"
	"github.com/kodinggo/rest-api-service-golang-private-1/internal/db"
	"github.com/kodinggo/rest-api-service-golang-private-1/internal/delivery/httpsvc"
	"github.com/kodinggo/rest-api-service-golang-private-1/internal/repository"
	"github.com/kodinggo/rest-api-service-golang-private-1/internal/usecase"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
)

var serverCMD = &cobra.Command{
	Use:   "serve",
	Short: "serve is a command to run the service",
	Run: func(cmd *cobra.Command, args []string) {
		// Init DB connection
		dbConn := db.InitMySQLConn()
		// Run server
		e := echo.New()

		// It's used to verify
		e.GET("/ping", func(c echo.Context) error {
			return c.String(http.StatusOK, "pong!")
		})

		// Init Repository
		storyRepo := repository.NewStoryRepository(dbConn)
		userRepo := repository.NewUserRepository(dbConn)

		// Init Usecase
		storyUsecase := usecase.NewStoryUsecase(storyRepo, userRepo)

		// Init HTTP Handler
		h := httpsvc.NewStoryHandler(storyUsecase)
		h.RegisterRoutes(e)

		e.Logger.Fatal(e.Start(config.Port()))
	},
}

func init() {
	rootCMD.AddCommand(serverCMD)
}

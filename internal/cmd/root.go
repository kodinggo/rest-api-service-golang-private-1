package cmd

import (
	"log"

	"github.com/kodinggo/rest-api-service-golang-private-1/internal/config"
	"github.com/spf13/cobra"
)

var rootCMD = &cobra.Command{
	Use:   "rest-api",
	Short: "rest-api is a rest api service",
}

func init() {
	config.InitConfig()
}

func Execute() {
	if err := rootCMD.Execute(); err != nil {
		log.Panicf("failed running cmd, error: %v", err)
	}
}

package cmd

import (
	"context"

	"github.com/kodinggo/rest-api-service-golang-private-1/internal/db"
	"github.com/labstack/gommon/log"
	"github.com/spf13/cobra"
)

var esIndexCMD = &cobra.Command{
	Use:   "esindex",
	Short: "esindex is a command to create index elasticsearch",
	Run: func(cmd *cobra.Command, args []string) {
		esClient := db.NewESClient()

		esIndex, err := esClient.CreateIndex(db.IndexName).BodyString(db.Mapping).Do(context.Background())
		if err != nil {
			log.Panicf("failed when create elasticsearch index, error: %v", err)
		}
		if !esIndex.Acknowledged {
			log.Warn("unsuccesfully creating index")
		} else {
			log.Info("succesfully creating index")
		}
	},
}

func init() {
	rootCMD.AddCommand(esIndexCMD)
}

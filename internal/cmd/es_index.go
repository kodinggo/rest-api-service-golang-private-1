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

		const mapping = `{}`
		const indexName = "my_index"

		esIndex, err := esClient.CreateIndex(indexName).BodyString(mapping).Do(context.Background())
		if err != nil {
			log.Panicf("failed when create elasticsearch index, error: %v", err)
		}
		if !esIndex.Acknowledged {
			log.Warn("unsuccesfully creating index")
		} else {
			log.Warn("succesfully creating index")
		}
	},
}

func init() {
	rootCMD.AddCommand(esIndexCMD)
}

package db

import (
	"github.com/olivere/elastic/v7"
	log "github.com/sirupsen/logrus"
)

var elasticHost = "http://localhost:9200" // TODO: Move to config

const (
	Mapping   = `{}`
	IndexName = "my_index"
)

func NewESClient() *elastic.Client {
	esClient, err := elastic.NewClient(elastic.SetURL(elasticHost), elastic.SetSniff(false))
	if err != nil {
		log.Panicf("failed when connect to elasticsearch, error: %v", err)
	}
	return esClient
}

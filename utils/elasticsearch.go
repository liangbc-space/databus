package utils

import (
	"databus/system"
	"github.com/olivere/elastic/v7"
	"log"
	"os"
)

var ElasticsearchClient *elastic.Client

func InitElasticsearch() {
	config := system.ApplicationCfg.ElasticsearchConfig

	client, err := elastic.NewClient(
		elastic.SetURL(config.Urls...),
		elastic.SetBasicAuth(config.Username, config.Password),

		elastic.SetSniff(false),
		elastic.SetGzip(true),
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)),
	)

	if err != nil {
		panic(err)
	}

	ElasticsearchClient = client
}

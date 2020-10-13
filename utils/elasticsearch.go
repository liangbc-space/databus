package utils

import (
	"github.com/liangbc-space/databus/system"
	"github.com/liangbc-space/databus/utils/exception"
	"github.com/olivere/elastic/v7"
	"go.uber.org/zap"
	"log"
	"os"
)

var ElasticsearchClient *elastic.Client

func InitElasticsearch() {
	config := system.ApplicationCfg.ElasticsearchConfig

	connOptions := []elastic.ClientOptionFunc{
		elastic.SetURL(config.Urls...),

		elastic.SetSniff(false),
		elastic.SetGzip(true),
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
	}

	if config.Username != "" {
		connOptions = append(connOptions, elastic.SetBasicAuth(config.Username, config.Password))
	}

	if system.ApplicationCfg.Debug {
		connOptions = append(connOptions, elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)))
	}

	client, err := elastic.NewClient(connOptions...)

	if err != nil {
		logger := NewDefaultLogger()
		defer logger.Sync()

		logger.Panic("ES连接失败："+err.Error(), zap.Reflect("connOptions", connOptions))
		exception.Throw("ES连接失败："+err.Error(), 1)
	}

	ElasticsearchClient = client
}

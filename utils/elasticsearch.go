package utils

import (
	"github.com/liangbc-space/databus/system"
	"github.com/olivere/elastic/v7"
	"go.uber.org/zap"
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
		logger := NewDefaultLogger()
		defer logger.Sync()

		logger.Panic("初始化ES连接失败："+err.Error(),
			zap.Strings("urls", config.Urls),
			zap.String("username", config.Username),
			zap.String("password", config.Password),
		)
	}

	ElasticsearchClient = client
}

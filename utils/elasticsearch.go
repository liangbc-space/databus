package utils

import (
	"databus/system"
	"github.com/elastic/go-elasticsearch/v7"
)

func ElasticsearchInstance() {
	config := system.SystemConfig.ElasticsearchConfig
	cfg := elasticsearch.Config{
		Addresses: config.Addresses,
	}
	elasticsearch.NewClient(cfg)
}

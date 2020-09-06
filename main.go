package main

import (
	"databus/models"
	mysql_elasticsearch "databus/mysql-elasticsearch"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	//"databus/routers"
	"databus/system"
	"flag"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"os"
)

func init() {
	//	初始化配置
	configPath := flag.String("systemConfig", "conf/conf.yaml", "system config file path")
	flag.Parse()

	err := system.LoadConfiguration(*configPath)
	if err != nil {
		panic(err)
		return
	}

	//	初始化数据库
	_, err = models.InitDB()
	if err != nil {
		panic(err)
		return
	}

	//	初始化redis连接池
	//utils.InitRedis()
}

func main() {
	//	释放数据库连接
	defer models.DB.Close()

	mysql_elasticsearch.Run()

	/*signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan,
		os.Kill,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)*/

	/*done := make(chan bool)

	go func() {
		run()
		done <- true
	}()

	select {
	case <-done:
		fmt.Printf("退出")
	}*/

	//request := routers.InitRouter()

	//if err := request.Run(fmt.Sprintf(":%d", system.SystemConfig.Port)); err != nil {
	//	panic(err)
	//}

}


func canal2() {

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":     strings.Join(system.SystemConfig.KafkaConfig.Brokers, ","),
		"broker.address.family": system.SystemConfig.KafkaConfig.BrokerAddressFamily,
		"group.id":              "test1",
		"session.timeout.ms":    6000,
		"auto.offset.reset":     "earliest",
		"enable.auto.commit":    false,
	})

	metas, err := c.GetMetadata(nil, false, 3000)
	var topics []string
	reg := regexp.MustCompile("^cn01_db.z_goods_(\\d{2})$")
	for _, topicMetadata := range metas.Topics {
		if reg.MatchString(topicMetadata.Topic) {
			topics = append(topics, topicMetadata.Topic)
		}
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create consumer: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created Consumer %v\n", c)

	err = c.SubscribeTopics(topics, nil)

	run := true

	for run == true {
		select {
		case sig := <-signalChan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			ev := c.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				fmt.Printf("%% Message on %s:\n%s\n",
					e.TopicPartition, string(e.Value))
				if e.Headers != nil {
					fmt.Printf("%% Headers: %v\n", e.Headers)
				}
			case kafka.Error:
				// Errors should generally be considered
				// informational, the client will try to
				// automatically recover.
				// But in this example we choose to terminate
				// the application if all brokers are down.
				fmt.Fprintf(os.Stderr, "%% Error: %v: %v\n", e.Code(), e)
				if e.Code() == kafka.ErrAllBrokersDown {
					run = false
				}
			default:
				fmt.Printf("Ignored %v\n", e)
			}
		}
	}

	fmt.Printf("Closing consumer\n")
	c.Close()
}

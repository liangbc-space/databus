package main

import (
	"databus/models"
	mysql_elasticsearch "databus/mysql-elasticsearch"
	"os/signal"
	"regexp"
	"strings"
	"sync"
	"syscall"
	"time"

	//"databus/routers"
	"databus/system"
	"flag"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/panjf2000/ants/v2"
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

func test1() {

	p, err := ants.NewPool(10)
	if err != nil {
		panic(err)
	}
	defer p.Release()

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan,
		os.Interrupt,
		os.Kill,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	counter := make(chan int,100)

	wg := new(sync.WaitGroup)

	p.Submit(func() {
		defer wg.Done()
		for {

			select {
			case sig, ok := <-signalChan:
				fmt.Printf("收到信号：%v\n", sig)
				if ok {
					close(signalChan)
				}
				return
			default:
				i1 := <-counter
				time.Sleep(time.Second * 1)
				fmt.Printf("当前值：%d\n", i1)
			}
		}

	})

	for i := 1; i <= 5; i++ {
		wg.Add(1)
		counter <- i

	}

	wg.Wait()

}

type done struct {
	c  chan os.Signal
	wg *sync.WaitGroup
}

func (done done) doneEvent() {
	signal.Notify(done.c, os.Kill,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	go func() {
		<-done.c
		close(done.c)
	}()
}

func newDoneControl() *done {
	done := done{
		c:  make(chan os.Signal),
		wg: new(sync.WaitGroup),
	}

	done.doneEvent()

	return &done
}

func test2() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	wg := new(sync.WaitGroup)

	//done := newDoneControl()
	for i := 1; i <= 2; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for {
				select {
				case _, ok := <-signalChan:
					fmt.Println("exit    ", i)
					if ok {
						close(signalChan)
					}
					fmt.Println("123")
					return
				default:
					if i == 2 {
						time.Sleep(5 * time.Second)
					} else {
						time.Sleep(2 * time.Second)
					}

					//time.Sleep(1 * time.Second)

					fmt.Printf("当前值%d\n", i)
				}
			}
		}(i)
	}

	wg.Wait()

}

func main() {

	test1()

	return

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

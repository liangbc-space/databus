package mysql_elasticsearch

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/liangbc-space/databus/models"
	"github.com/liangbc-space/databus/system"
	"github.com/liangbc-space/databus/utils"
	"github.com/liangbc-space/databus/utils/exception"
	"github.com/natefinch/lumberjack"
	"github.com/panjf2000/ants/v2"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func Run() {
	//	获取topics
	topics := createConsumerInstance().getTopics()

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	listener := utils.NewSignalListener(signalChan)
	listener.SignalEvent(func() {
		sig, _ := <-listener.SignalChan
		fmt.Printf("收到信号【%v】程序退出\n", sig)
	})

	p, err := ants.NewPoolWithFunc(len(topics), func(args interface{}) {
		exception.Try(func() {
			execute(args)
		}).Catch(func(ex exception.Exception) {
			fmt.Printf("消费kafka数据同步到ES异常：%s	文件：%s:%d\n",
				ex.Message(), ex.File(), ex.Line(),
			)
			os.Exit(ex.Code())
		})
	})
	if err != nil {
		exception.Throw("创建协程池失败："+err.Error(), 1)
	}

	for _, topic := range topics {
		listener.Add(1)

		args := make(map[string]interface{}, 2)
		args["topic"] = topic
		args["signal_listener"] = listener

		if err := p.Invoke(args); err != nil {
			exception.Throw("协程执行任务失败："+err.Error(), 1)
		}
	}

	listener.Wait()

}

func execute(args interface{}) {
	argsMap := args.(map[string]interface{})

	topic := argsMap["topic"].(string)
	control := argsMap["signal_listener"].(*utils.SignalControl)
	defer control.Done()

	consumer := createConsumerInstance()
	if err := consumer.Subscribe(topic, nil); err != nil {
		exception.Throw("订阅kafka消息失败："+err.Error(), 1)
	}
	defer consumer.Close()

	reg := regexp.MustCompile(`^cn01_db.z_goods_(\d{2,3})$`)
	matches := reg.FindStringSubmatch(topic)
	if len(matches) < 1 {
		return
	}
	tableHash := matches[1]
	logger := logger()
	defer logger.Sync()

	allOptionData := make([]map[string]interface{}, 0)
	saveOptionData := make([]map[string]interface{}, 0)
	for {

		select {
		case <-control.SignalChan:
			if len(allOptionData) > 0 {
				sync(&allOptionData, &saveOptionData)
				logger.Info("收到退出信号，所有消费数据处理完成")
			}
			return
		default:
			//	获取kafka消息	超时时间100
			message := consumer.pullMessages(100)
			if message == nil {
				if len(allOptionData) > 0 {
					sync(&allOptionData, &saveOptionData)
				}
				continue
			}
			//	提交消费偏移量
			consumer.Commit()

			if system.ApplicationCfg.KafkaConfig.ConsumerLogs {
				gz, err := utils.GzipEncode(message.Value)
				if err != nil {
					utils.NewDefaultLogger().Panic("gzip数据压缩错误"+err.Error(),
						zap.ByteString("gzData", message.Value),
					)
					exception.Throw("gzip数据压缩错误："+err.Error(), 2)
				}
				logger.Info(consumer.String()+"成功获取到数据",
					zap.String("message", base64.StdEncoding.EncodeToString(gz)),
				)
			}

			//optionData := make(map[string]interface{})
			optionData := new(struct {
				Data       interface{} `json:"data"`
				OptionType string      `json:"type"`
			})
			if err := json.Unmarshal(message.Value, &optionData); err != nil {
				exception.Throw("kafka消息解析失败："+err.Error(), 1)
			}

			data := optionData.Data.([]interface{})

			for _, item := range data {
				item := item.(map[string]interface{})
				goodsData := make(map[string]interface{})

				goodsData["goods_id"] = item["id"]
				goodsData["store_id"] = item["store_id"]
				goodsData["_unq"] = item["store_id"]
				goodsData["operation_type"] = strings.ToUpper(optionData.OptionType)
				goodsData["table_hash"] = tableHash

				allOptionData = append(allOptionData, goodsData)
				if goodsData["operation_type"] != "DELETE" {
					saveOptionData = append(saveOptionData, goodsData)
				}
			}

			if len(allOptionData) >= 100 {
				sync(&allOptionData, &saveOptionData)
			}

		}

	}

}

func sync(allOptionData *[]map[string]interface{}, saveOptionData *[]map[string]interface{}) {
	millisecond := time.Now().UnixNano() / 1e6

	failedIds := make([]uint32, 0)
	for _, item := range *allOptionData {
		goodsId, err := strconv.Atoi(item["goods_id"].(string))
		if err != nil {
			exception.Throw("数据同步到ES异常："+err.Error(), 1)
		}
		failedIds = append(failedIds, uint32(goodsId))
	}

	tableName := fmt.Sprintf("z_goods_%s", (*allOptionData)[0]["table_hash"].(string))

	exception.Try(func() {

		//	构建es商品数据
		goodsLists := make(map[string]esGoods)
		if len(*saveOptionData) > 0 {
			goodsLists = buildEsGoods(*saveOptionData)
		}

		//	数据更新到es
		failedIds = pushToElasticsearch(*allOptionData, goodsLists)
		if len(failedIds) > 0 {
			exception.Throw(fmt.Sprintf("ES数据更新%d失败", len(failedIds)), 2)
		}

	}).Catch(func(ex exception.Exception) {
		if len(failedIds) > 0 {
			models.DB.Table(tableName).
				Where("id in(?)", failedIds).
				Update("modify_time", time.Now().Unix()+1)
		}
		//TODO：删除数据失败

	}).Finally(func() {
		fmt.Printf("%s： %s		成功%d条数据	失败%d条数据		耗时%dms\n",
			time.Now().Format("2006/01/02 15:04:05.000"),
			tableName,
			len(*allOptionData)-len(failedIds),
			len(failedIds),
			time.Now().UnixNano()/1e6-millisecond,
		)

		millisecond = time.Now().UnixNano() / 1e6
		*allOptionData = (*allOptionData)[0:0]
		*saveOptionData = (*saveOptionData)[0:0]
	})

}

func buildEsGoods(optionDatas []map[string]interface{}) map[string]esGoods {
	tableHash := optionDatas[0]["table_hash"].(string)
	//	获取商品的基本信息
	list := make(goodsLists, 0)
	list = models.GetGoods(tableHash, optionDatas)

	goodsIds := make([]string, 0)
	storeIds := make([]string, 0)
	categoryIds := make([]string, 0)
	for _, goods := range list {
		goodsIds = append(goodsIds, strconv.Itoa(int(goods.Id)))
		storeIds = append(storeIds, strconv.Itoa(int(goods.StoreId)))
		categoryIds = append(categoryIds, strings.Split(goods.CategoryPath, ",")...)
	}

	//	去重
	goodsIds = utils.RemoveRepeat(goodsIds)
	storeIds = utils.RemoveRepeat(storeIds)
	categoryIds = utils.RemoveRepeat(categoryIds)

	//  获取商品tag信息
	GoodsTags = models.GetGoodsTags(goodsIds, storeIds)

	//  获取商品推荐信息
	GoodsRecommends = models.GetGoodsRecommends(goodsIds, storeIds)

	//  获取商品分类信息
	GoodsCategories = models.GetGoodsCategories(tableHash, categoryIds)

	//  获取商品附属分类信息
	GoodsSubCategories = models.GetGoodsSubCategories(tableHash, goodsIds, storeIds)

	//  获取商品图片信息
	GoodsOtherImages = models.GetGoodsOtherImages(tableHash, goodsIds, storeIds)

	//  获取商品销量属性信息
	GoodsSaleProperties = models.GetGoodsSaleProperties(tableHash, goodsIds, storeIds)

	//  获取商品属性信息
	GoodsProperties = models.GetGoodsProperties(tableHash, goodsIds, storeIds)

	//	组装es商品数据
	return list.buildElasticsearchGoods(tableHash)

}

func logger() *zap.Logger {

	loggerCfg := utils.LoggerCfg{
		Level: zap.InfoLevel,
		Hook: lumberjack.Logger{
			Filename:   "logs/kafka-messages.log",
			MaxAge:     5,
			MaxBackups: 10,
			MaxSize:    512,
			Compress:   true,
			LocalTime:  true,
		},
		WithCaller: false,
	}

	return loggerCfg.NewLogger()
}

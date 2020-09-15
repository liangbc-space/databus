package mysql_elasticsearch

import (
	"databus/models"
	"databus/utils"
	"encoding/json"
	"fmt"
	"github.com/panjf2000/ants/v2"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"
)

func Run() {
	//	获取topics
	topics := getTopics()

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	control := utils.SignalEvent(signalChan, signalHandleFunc)

	p, err := ants.NewPoolWithFunc(len(topics), execute)
	if err != nil {
		panic(err)
	}

	for _, topic := range topics {
		control.Add(1)

		args := make(map[string]interface{}, 2)
		args["topic"] = topic
		args["signal_control"] = control

		if err := p.Invoke(args); err != nil {
			panic(err)
		}
	}

	control.Wait()

}

func signalHandleFunc(control *utils.SignalControl) {
	<-control.SignalChan
}

func execute(args interface{}) {
	argsMap := args.(map[string]interface{})

	topic := argsMap["topic"].(string)
	control := argsMap["signal_control"].(*utils.SignalControl)

	consumer := createConsumerInstance()
	if err := consumer.Subscribe(topic, nil); err != nil {
		panic(err)
	}

	reg := regexp.MustCompile(`^cn01_db.z_goods_(\d{2})$`)
	tableHash := reg.FindStringSubmatch(topic)[1]

	defer consumer.Close()
	defer control.Done()

	allOptionData := make([]map[string]interface{}, 0)
	saveOptionData := make([]map[string]interface{}, 0)
	for {

		select {
		case <-control.SignalChan:
			fmt.Println("收到退出信号")
			return
		default:
			//	获取kafka消息
			message := pullMessages(consumer)
			if message == nil {
				continue
			}

			//optionData := make(map[string]interface{})
			type Message struct {
				Data       interface{} `json:"data"`
				OptionType string      `json:"type"`
			}
			optionData := new(Message)
			if err := json.Unmarshal(message.Value, &optionData); err != nil {
				panic(err)
			}

			data := optionData.Data.([]interface{})

			for _, item := range data {
				item := item.(map[string]interface{})
				goodsData := make(map[string]interface{})

				goodsData["goods_id"] = item["id"]
				goodsData["store_id"] = item["store_id"]
				goodsData["operation_type"] = strings.ToUpper(optionData.OptionType)

				allOptionData = append(allOptionData, goodsData)
				if goodsData["operation_type"] != "DELETE" {
					saveOptionData = append(saveOptionData, goodsData)
				}

			}

			if len(allOptionData) >= 100 {
				goodsLists := elasticsearchGoodsData(tableHash, saveOptionData)
				pushToElasticsearch(allOptionData, goodsLists)

				allOptionData = make([]map[string]interface{}, 0)
				saveOptionData = make([]map[string]interface{}, 0)

				/*if *message.TopicPartition.Topic == "cn01_db.z_goods_00" {
					time.Sleep(time.Second * 1)
				} else {
					time.Sleep(time.Second * 5)
				}*/
				/*bytes, _ := json.Marshal(list)
				fmt.Println(string(bytes))*/


			}


		}

	}

}

func elasticsearchGoodsData(tableHash string, optionDatas []map[string]interface{}) map[string]esGoods {
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

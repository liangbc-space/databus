package mysql_elasticsearch

import (
	"context"
	"fmt"
	"github.com/liangbc-space/databus/models"
	"github.com/liangbc-space/databus/utils"
	"github.com/liangbc-space/databus/utils/exception"
	"github.com/olivere/elastic/v7"
	"reflect"
	"strconv"
	"strings"
)

const (
	INDEX = "goods_base_test"
	ALIAS = "alias_goods_base_test"
)

type (
	//elasticsearch商品
	esGoods map[string]interface{}

	//商品基础信息
	goodsLists []models.Godos
)

var (
	//商品tag
	GoodsTags []models.GoodsTag

	//商品推荐信息
	GoodsRecommends []models.GoodsRecommend

	//商品分类
	GoodsCategories map[uint32]models.GoodsCategory

	//商品附属分类
	GoodsSubCategories map[string][]models.GoodsSubCategory

	//商品图片
	GoodsOtherImages []models.GoodsOtherImage

	//商品销量属性
	GoodsSaleProperties []models.GoodsSaleProperty

	//商品属性
	GoodsProperties []models.GoodsProperty
)

func uniqueId(goodsId interface{}, storeId interface{}) string {
	return fmt.Sprintf("%v-%v", storeId, goodsId)
}

func (list goodsLists) buildElasticsearchGoods(tableHash string) map[string]esGoods {
	goodsLists := make(map[string]esGoods)
	for _, goods := range list {
		goodsItem := make(esGoods)
		goodsItem["mysql_table_name"] = fmt.Sprintf("z_goods-%s", tableHash)

		uniqueId := uniqueId(goods.Id, goods.StoreId)

		(&goodsItem).
			//	商品基础信息
			buildBaseGoods(goods).
			//	商品标签
			buildGoodsTags(uniqueId).
			//	商品推荐信息
			buildGoodsRecommends(uniqueId).
			//	商品分类
			buildGoodsCategories(strings.Split(goods.CategoryPath, ",")).
			//	商品子分类
			buildGoodsSubCategories(uniqueId).
			//	商品图片
			buildGoodsOtherImages(uniqueId).
			//	商品sale属性
			buildGoodsSaleProps(uniqueId).
			//	商品属性
			buildGoodsProps(uniqueId).
			//	搜索关键词
			initSearchKeywords()

		goodsLists[uniqueId] = goodsItem
	}

	return goodsLists
}

func (p *esGoods) buildBaseGoods(goods models.Godos) *esGoods {
	goodsItem := *p

	rValue := reflect.ValueOf(goods)
	for i := 0; i < rValue.NumField(); i++ {
		field := rValue.Type().Field(i).Tag.Get("gorm")

		goodsItem[field] = rValue.Field(i).Interface()
	}
	delete(goodsItem, "user_group_id_values")

	userGroupIds := make([]uint32, 0)
	if goods.UserGroupIdValues != "" {
		ids := strings.Split(goods.UserGroupIdValues, ",")

		for _, id := range utils.RemoveRepeat(ids) {
			id, _ := strconv.Atoi(id)

			userGroupIds = append(userGroupIds, uint32(id))
		}
	}
	goodsItem["user_group_ids"] = userGroupIds

	return p
}

func (p *esGoods) buildGoodsTags(unqId string) *esGoods {
	goodsItem := *p

	tagIds := make([]uint32, 0)
	tagNames := make([]string, 0)

	for _, tag := range GoodsTags {
		if unqId == tag.UniqueId {
			tagIds = append(tagIds, tag.TagId)
			tagNames = append(tagNames, tag.TagName)
		}
	}
	goodsItem["tag_ids"] = tagIds
	goodsItem["tag_names"] = tagNames

	return p
}

func (p *esGoods) buildGoodsRecommends(unqId string) *esGoods {
	goodsItem := *p

	recIds := make([]uint32, 0)
	recNames := make([]string, 0)

	for _, recommend := range GoodsRecommends {
		if unqId == recommend.UniqueId {

			recIds = append(recIds, recommend.RecId)
			recNames = append(recNames, recommend.RecName)
			//goodsItem[fmt.Sprintf("up_time_%d", recommend.RecIndex)] = recommend.RecUpTime

		}
	}

	goodsItem["rec_ids"] = recIds
	goodsItem["rec_names"] = recNames

	return p
}

func (p *esGoods) buildGoodsCategories(categoryIds []string) *esGoods {

	categories := make([]models.GoodsCategory, 0)
	for _, id := range categoryIds {
		id, _ := strconv.Atoi(id)

		if _, ok := GoodsCategories[uint32(id)]; ok {
			categories = append(categories, GoodsCategories[uint32(id)])
		}
	}

	p.initCategory(categories)

	return p

}

func (p *esGoods) buildGoodsSubCategories(uniqueId string) *esGoods {

	subCategories, ok := GoodsSubCategories[uniqueId]
	if ok {
		p.initCategory(subCategories)
	}

	return p
}

func (p *esGoods) buildGoodsOtherImages(unqId string) *esGoods {
	goodsItem := *p

	images := make([]string, 0)
	for _, image := range GoodsOtherImages {
		if unqId == image.UniqueId {
			images = append(images, image.Image)
		}
	}
	goodsItem["images_other"] = images

	return p
}

func (p *esGoods) buildGoodsSaleProps(unqId string) *esGoods {
	goodsItem := *p

	propertyNames := make([]string, 0)
	propertyImages := make([]string, 0)
	for _, property := range GoodsSaleProperties {
		if unqId == property.UniqueId {
			propertyNames = append(propertyNames, property.BaseName)
			propertyImages = append(propertyImages, property.Image)
		}
	}
	goodsItem["main_prop_name"] = propertyNames
	goodsItem["main_prop_image"] = propertyNames

	return p
}

func (p *esGoods) buildGoodsProps(unqId string) *esGoods {
	goodsItem := *p

	propertyIds := make([]uint32, 0)
	for _, property := range GoodsProperties {
		if unqId == property.UniqueId {
			propertyIds = append(propertyIds, uint32(property.ValueId))
		}
	}
	goodsItem["property_ids"] = propertyIds

	return p
}

func (p *esGoods) initCategory(categories interface{}) {
	goodsItem := *p

	categoryIds := make([]uint32, 0)
	if _, ok := goodsItem["category_ids"]; !ok {
		goodsItem["category_ids"] = make([]uint32, 0)
	}

	rValue := reflect.ValueOf(goodsItem["category_ids"])
	if rValue.Kind() != reflect.Slice && rValue.Kind() != reflect.Array {
		panic("商品解析出错，category_ids类型不正确")
	}
	for i := 0; i < rValue.Len(); i++ {
		categoryIds = append(categoryIds, rValue.Index(i).Interface().(uint32))
	}

	categoryNames := make([]string, 0)
	if _, ok := goodsItem["category_names"]; !ok {
		goodsItem["category_names"] = make([]string, 0)
	}
	rValue = reflect.ValueOf(goodsItem["category_names"])
	if rValue.Kind() != reflect.Slice && rValue.Kind() != reflect.Array {
		panic("商品解析出错，category_names类型不正确")
	}
	for i := 0; i < rValue.Len(); i++ {
		categoryNames = append(categoryNames, rValue.Index(i).Interface().(string))
	}

	rValue = reflect.ValueOf(categories)
	if rValue.Kind() != reflect.Slice && rValue.Kind() != reflect.Array {
		exception.Throw("商品解析出错，categories类型不正确", 1)
	}

	for i := 0; i < rValue.Len(); i++ {
		switch rValue.Index(i).Type() {
		case reflect.TypeOf(models.GoodsCategory{}):
			goodsCategory := rValue.Index(i).Interface().(models.GoodsCategory)
			categoryIds = append(categoryIds, goodsCategory.GoodsCategoryId)
			categoryNames = append(categoryNames, goodsCategory.GoodsCategoryName)
		case reflect.TypeOf(models.GoodsSubCategory{}):
			goodsSubCategory := rValue.Index(i).Interface().(models.GoodsSubCategory)
			categoryIds = append(categoryIds, goodsSubCategory.GoodsCategoryId)
			categoryNames = append(categoryNames, goodsSubCategory.GoodsCategoryName)
		}
	}

	goodsItem["category_ids"] = categoryIds
	goodsItem["category_names"] = categoryNames

}

func (p *esGoods) initSearchKeywords() {
	goodsItem := *p

	searchKeywords := make([]string, 0)

	searchKeywords = append(searchKeywords, goodsItem["base_name"].(string))
	searchKeywords = append(searchKeywords, goodsItem["codeno"].(string))

	rType := reflect.TypeOf(goodsItem["tag_names"])
	if rType.Kind() == reflect.Slice || rType.Kind() == reflect.Array {
		//	标签名称
		tagNames := goodsItem["tag_names"].([]string)
		for _, tagName := range tagNames {
			searchKeywords = append(searchKeywords, tagName)
		}
	}

	goodsItem["search_keywords"] = searchKeywords
}

func pushToElasticsearch(allOptionData []map[string]interface{}, goodsLists map[string]esGoods) (failedGoodsIds []uint32) {

	bulk := utils.ElasticsearchClient.Bulk()
	for _, optionData := range allOptionData {
		unqId := uniqueId(optionData["goods_id"], optionData["store_id"])
		goods, ok := goodsLists[unqId]

		if optionData["operation_type"] == "DELETE" {
			request := elastic.NewBulkDeleteRequest().Index(INDEX).Id(unqId)
			bulk.Add(request)
		} else if ok {
			request := elastic.NewBulkIndexRequest().Index(INDEX).Id(unqId).Doc(goods)
			bulk.Add(request)
		}
	}

	if bulk.NumberOfActions() > 0 {
		response, err := bulk.Do(context.TODO())
		if err != nil {
			exception.Throw("写入ES失败："+err.Error(), 1)
		}

		for _, fail := range response.Failed() {
			if fail.Error != nil {
				goodsId := strings.Split(fail.Id, "-")
				pkId, _ := strconv.Atoi(goodsId[1])
				failedGoodsIds = append(failedGoodsIds, uint32(pkId))
			}
		}
	}

	return failedGoodsIds
}

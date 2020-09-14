package mysql_elasticsearch

import (
	"databus/models"
	"databus/utils"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

//elasticsearch商品
type esGoods map[string]interface{}

//elasticsearch商品列表
type esGoodsLists []esGoods

//商品基础信息
type goodsLists []models.Godos

//商品tag
var GoodsTags []models.GoodsTag

//商品推荐信息
var GoodsRecommends []models.GoodsRecommend

//商品分类
var GoodsCategories map[uint32]models.GoodsCategory

//商品附属分类
var GoodsSubCategories map[string][]models.GoodsSubCategory

//商品图片
var GoodsOtherImages []models.GoodsOtherImage

//商品销量属性
var GoodsSaleProperties []models.GoodsSaleProperty

//商品属性
var GoodsProperties []models.GoodsProperty

func uniqueId(goodsId interface{}, storeId interface{}) string {
	return fmt.Sprintf("%d-%d", storeId, goodsId)
}

func (list goodsLists) buildElasticsearchGoods(tableHash string) (esGoodsLists esGoodsLists) {
	for _, goods := range list {
		goodsItem := make(esGoods)

		goodsItem = goodsItem.buildBaseGoods(goods)
		goodsItem["mysql_table_name"] = fmt.Sprintf("z_goods-%s", tableHash)

		goodsItem = goodsItem.buildGoodsTags(goods.Id, goods.StoreId)

		goodsItem = goodsItem.buildGoodsRecommends(goods.Id, goods.StoreId)

		goodsItem = goodsItem.buildGoodsCategories(strings.Split(goods.CategoryPath, ","))

		goodsItem = goodsItem.buildGoodsSubCategories(uniqueId(goods.Id, goods.StoreId))

		esGoodsLists = append(esGoodsLists, goodsItem)
	}

	return esGoodsLists
}

func (goodsItem esGoods) buildBaseGoods(goods models.Godos) esGoods {
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

	return goodsItem
}

func (goodsItem esGoods) buildGoodsTags(goodsId uint32, storeId uint32) esGoods {
	tagIds := make([]uint32, 0)
	tagNames := make([]string, 0)

	for _, tag := range GoodsTags {
		if tag.StoreId == storeId && tag.GoodsId == goodsId {
			tagIds = append(tagIds, tag.TagId)
			tagNames = append(tagNames, tag.TagName)
		}
	}
	goodsItem["tag_ids"] = tagIds
	goodsItem["tag_names"] = tagNames

	return goodsItem
}

func (goodsItem esGoods) buildGoodsRecommends(goodsId uint32, storeId uint32) esGoods {
	recIds := make([]uint32, 0)
	recNames := make([]string, 0)

	for _, recommend := range GoodsRecommends {
		if recommend.GoodsId == goodsId && recommend.StoreId == storeId {

			recIds = append(recIds, recommend.RecId)
			recNames = append(recNames, recommend.RecName)
			goodsItem[fmt.Sprintf("up_time_%d", recommend.RecIndex)] = recommend.RecUpTime

		}
	}

	goodsItem["rec_ids"] = recIds
	goodsItem["rec_names"] = recNames

	return goodsItem
}

func (goodsItem esGoods) buildGoodsCategories(categoryIds []string) esGoods {
	/*ids := make([]uint32, 0)
	names := make([]string, 0)
	for _, id := range categoryIds {
		id, _ := strconv.Atoi(id)

		if _, ok := GoodsCategories[uint32(id)]; ok {
			ids = append(ids, GoodsCategories[uint32(id)].GoodsCategoryId)
			names = append(names, GoodsCategories[uint32(id)].GoodsCategoryName)
		}
	}
	goodsItem["category_ids"] = ids
	goodsItem["category_names"] = names

	return goodsItem*/

	if _,ok:=goodsItem["category_ids"];!ok {
		goodsItem["category_ids"] = make([]uint32,0)
	}

	rType := reflect.ValueOf(goodsItem["category_ids"])
	for _, id := range categoryIds {
		id, _ := strconv.Atoi(id)

		goodsCategory, ok := GoodsCategories[uint32(id)]
		if ok {
			switch rType.Kind() {
			case reflect.Slice,reflect.Array:
				fmt.Println(123)
				//fmt.Println(reflect.ValueOf(goodsCategory).FieldByName("GoodsCategoryId"))
				reflect.Append(reflect.ValueOf(goodsItem["category_ids"]).Elem(), reflect.ValueOf(goodsCategory).FieldByName("GoodsCategoryId"))
				/*for i:=0;i<rValue.Len();i++{

					fmt.Println(rValue.Index(i).Type())
				}*/
			}
			/*switch rValue.Interface().(type) {
			case []int8,[]int16,[]int32,[]int64,[]uint8,[]uint16,[]uint32,[]uint64:
				reflect.Append(rValue,reflect.ValueOf(GoodsCategories[uint32(id)].GoodsCategoryId))

			}*/
			/*ids = append(ids, GoodsCategories[uint32(id)].GoodsCategoryId)
			names = append(names, GoodsCategories[uint32(id)].GoodsCategoryName)*/
		}
	}
	fmt.Println(goodsItem["category_ids"])
	/*goodsItem["category_ids"] = ids
	goodsItem["category_names"] = names*/
	os.Exit(1)

	return goodsItem
}

func (goodsItem esGoods) buildGoodsSubCategories(uniqueId string) esGoods {
	subCategories, ok := GoodsSubCategories[uniqueId]
	if ok {
		for _, category := range subCategories {
			rValue := reflect.ValueOf(goodsItem["category_names"])
			switch rValue.Kind() {
			case reflect.Slice, reflect.Array:
				fmt.Println(rValue.Index(0))
			}
			fmt.Println(category)
		}
	}

	if len(subCategories) > 0{
		os.Exit(1)
	}
	return goodsItem
}

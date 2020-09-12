package models

import (
	"databus/utils"
	"fmt"
	"strings"
)

//	商品基础信息
type Godos struct {
	UniqueId          string  `gorm:"uniqueeid" json:"unique_id"`
	Id                uint32  `gorm:"id"`
	StoreId           uint32  `gorm:"store_id"`
	BaseName          string  `gorm:"base_name"`
	BrandId           uint32  `gorm:"brand_id"`
	GoodsTypeId       int32   `gorm:"goods_type_id"`
	CategoryPath      string  `gorm:"category_path"`
	CategoryId        uint32  `gorm:"category_id"`
	Codeno            string  `gorm:"codeno"`
	Image             string  `gorm:"image"`
	Price             float64 `gorm:"price"`
	CostPrice         float64 `gorm:"cost_price"`
	MarketPrice       float64 `gorm:"market_price"`
	ListOrder         int32   `gorm:"listorder"`
	Status            int8    `gorm:"status"`
	UpTime            int32   `gorm:"up_time"`
	DownTime          int32   `gorm:"down_time"`
	CreateTime        int32   `gorm:"create_time"`
	ModifyTime        int32   `gorm:"modify_time"`
	TemplatePage      string  `gorm:"template_page"`
	VisitCounts       int32   `gorm:"visit_counts"`
	BuyCounts         int32   `gorm:"buy_counts"`
	WishlistCounts    int32   `gorm:"wishlist_counts"`
	CommentCounts     int32   `gorm:"comment_counts"`
	CommentValue      int32   `gorm:"comment_value"`
	StockNums         int32   `gorm:"stock_nums"`
	SaleMode          int8    `gorm:"sale_mode"`
	SpecMode          int8    `gorm:"spec_mode"`
	IsDiyRemark       int8    `gorm:"is_diy_remark"`
	Weight            float32 `gorm:"weight"`
	StartTime         uint32  `gorm:"start_time"`
	EndTime           uint32  `gorm:"end_time"`
	IsFreeShipping    uint8   `gorm:"is_free_shipping"`
	SpecialOfferId    uint32  `gorm:"special_offer_id"`
	Discount          float32 `gorm:"discount"`
	Title             string  `gorm:"title"`
	Keywords          string  `gorm:"keywords"`
	Descript          string  `gorm:"descript"`
	MiniDetail        string  `gorm:"mini_detail"`
	GroupCodeno       string  `gorm:"group_codeno"`
	Moq               int32   `gorm:"moq"`
	Mxoq              int32   `gorm:"mxoq"`
	IsBookable        int8    `gorm:"is_bookable"`
	B2bStatus         int8    `gorm:"b2b_status"`
	UserGroupIdValues string  `gorm:"user_group_id_values"`
	Volume            float32 `gorm:"volume"`
	SupplierRemark    string  `gorm:"supplier_remark"`
	Video             string  `gorm:"video"`
	IsInstock         int8    `gorm:"is_instock"`
	CreateDay         string  `gorm:"create_day"`
	BrandName         string  `gorm:"brand_name"`
	CategoryName      string  `gorm:"category_name"`
}

//	商品标签
type GoodsTag struct {
	TagId   uint32 `gorm:"tag_id"`
	TagName uint32 `gorm:"tag_name"`
	GoodsId uint32 `gorm:"goods_id"`
	StoreId uint32 `gorm:"store_id"`
}

//	商品推荐
type GoodsRecommend struct {
	RecId     uint32 `gorm:"rec_id"`
	RecIndex  uint32 `gorm:"rec_index"`
	RecName   string `gorm:"rec_name"`
	RecUpTime int32  `gorm:"rec_up_time"`
	GoodsId   uint32 `gorm:"goods_id"`
	StoreId   uint32 `gorm:"store_id"`
}

//	商品分类
type GoodsCategory struct {
	GoodsCategoryId   uint32 `gorm:"goods_category_id"`
	GoodsCategoryName string `gorm:"goods_category_name"`
}

//	商品附属分类
type GoodsSubCategory struct {
	GoodsId           uint32 `gorm:"goods_id"`
	StoreId           uint32 `gorm:"store_id"`
	GoodsCategoryId   uint32 `gorm:"goods_category_id"`
	GoodsCategoryName string `gorm:"goods_category_name"`
}

//	商品图片
type GoodsOtherImage struct {
	GoodsId uint32 `gorm:"goods_id"`
	StoreId uint32 `gorm:"store_id"`
	Image   string `gorm:"image"`
}

//	商品销量属性？
type GoodsSaleProperty struct {
	BaseName string `gorm:"base_name"`
	Image    string `gorm:"image"`
	GoodsId  uint32 `gorm:"goods_id"`
	StoreId  uint32 `gorm:"store_id"`
}

//	商品属性
type GoodsProperty struct {
	GoodsId    uint32 `gorm:"goods_id"`
	StoreId    uint32 `gorm:"store_id"`
	PropertyId int32  `gorm:"property_id"`
	ValueId    int32  `gorm:"value_id"`
	ValueName  string `gorm:"value_name"`
}

//	返回商品唯一ID
func UniqueId(storeId uint32, goodsId uint32) string {
	return fmt.Sprintf("%d-%d", storeId, goodsId)
}

func GetGoods(tableHash string, optionDatas []map[string]interface{}) (goodsLists []Godos) {
	goodsIds := make([]string, 0)
	for _, item := range optionDatas {
		goodsIds = append(goodsIds, item["goods_id"].(string))
	}
	goodsIds = utils.RemoveRepeat(goodsIds)
	//	查询goods
	sql := `SELECT
	CONCAT( CAST( g.store_id AS CHAR ), '-', CAST( g.id AS CHAR ) ) AS uniqueeid,
	g.*,
IF
	( g.stock_nums > 0, 1, IF ( g.is_bookable, 1, 0 ) ) AS is_instock,
	FROM_UNIXTIME( g.create_time, '%Y%m%d' ) AS create_day,
	b.base_name AS brand_name,
	c.base_name AS category_name 
FROM
	z_goods_` + tableHash + ` g
	LEFT JOIN z_brand AS b ON g.brand_id = b.id
	LEFT JOIN z_goods_category_` + tableHash + ` AS c ON g.category_id = c.id 
WHERE
    g.id IN(` + strings.Join(goodsIds, ",") + `) 
    AND g.store_id > 0
	AND g.STATUS != -1`

	DB.Raw(sql).Find(&goodsLists)
	return goodsLists
}

func GetGoodsTags(goodsIds []string, storeIds []string) map[string]GoodsTag {
	goodsIds = utils.RemoveRepeat(goodsIds)
	if len(goodsIds) < 1 {
		return nil
	}

	storeIds = utils.RemoveRepeat(storeIds)
	sql := `SELECT
	tag_id AS tag_id,
	base_name AS tag_name,
	r.goods_id,
	r.store_id
FROM
	z_goods_tag AS t
	LEFT JOIN z_goods_tag_rel AS r ON t.id = r.tag_id 
WHERE
    r.store_id in(` + strings.Join(storeIds, ",") + `) and r.goods_id in(` + strings.Join(goodsIds, ",") + `)`

	goodsTags := []GoodsTag{}
	DB.Raw(sql).Find(&goodsTags)

	tags := make(map[string]GoodsTag)
	for _, tag := range goodsTags {
		uniqueId := UniqueId(tag.StoreId, tag.GoodsId)
		tags[uniqueId] = tag
	}

	return tags
}

func GetGoodsRecommends(goodsIds []string, storeIds []string) map[string]GoodsRecommend {
	goodsIds = utils.RemoveRepeat(goodsIds)
	if len(goodsIds) < 1 {
		return nil
	}

	storeIds = utils.RemoveRepeat(storeIds)

	sql := `SELECT
	r.id AS rec_id,
	r.rec_index AS rec_index,
	r.base_name AS rec_name,
	rr.up_time AS rec_up_time,
	rr.goods_id,
	rr.store_id
FROM
	z_goods_recommend AS r
	LEFT JOIN z_goods_recommend_rel AS rr ON r.id = rr.goods_recommend_id 
WHERE
	rr.store_id in(` + strings.Join(storeIds, ",") + `) AND rr.goods_id IN (` + strings.Join(goodsIds, ",") + `)`

	goodsRecommends := []GoodsRecommend{}
	DB.Raw(sql).Find(&goodsRecommends)

	recommends := make(map[string]GoodsRecommend)
	for _, recommend := range goodsRecommends {
		uniqueId := UniqueId(recommend.StoreId, recommend.GoodsId)
		recommends[uniqueId] = recommend
	}

	return recommends
}

func GetGoodsCategories(tableHash string, categoryIds []string) map[uint32]GoodsCategory {
	categoryIds = utils.RemoveRepeat(categoryIds)
	if len(categoryIds) < 1 {
		return nil
	}
	sql := `SELECT
	id AS goods_category_id,
	base_name AS goods_category_name
FROM
	z_goods_category_` + tableHash + `
WHERE
	id IN ( ` + strings.Join(categoryIds, ",") + ` ) `

	goodsCategories := []GoodsCategory{}
	DB.Raw(sql).Find(&goodsCategories)

	categories := make(map[uint32]GoodsCategory)
	for _, category := range goodsCategories {
		categories[category.GoodsCategoryId] = category
	}

	return categories
}

func GetGoodsSubCategories(tableHash string, goodsIds []string, storeIds []string) map[string][]GoodsSubCategory {
	goodsIds = utils.RemoveRepeat(goodsIds)
	if len(goodsIds) < 1 {
		return nil
	}

	storeIds = utils.RemoveRepeat(storeIds)
	sql := `SELECT
    r.goods_id,
    r.store_id,
	c.id AS goods_category_id,
	c.base_name AS goods_category_name 
FROM
	z_goods_category_` + tableHash + ` c
	LEFT JOIN z_goods_category_rel_` + tableHash + ` r ON c.id = r.category_id 
WHERE
	r.store_id in(` + strings.Join(storeIds, ",") + `) 
	AND r.goods_id IN(` + strings.Join(goodsIds, ",") + `)`

	goodsSubCategories := []GoodsSubCategory{}
	DB.Raw(sql).Find(&goodsSubCategories)

	subCategories := make(map[string][]GoodsSubCategory)
	for _, category := range goodsSubCategories {
		uniqueId := UniqueId(category.StoreId, category.GoodsId)

		subCategories[uniqueId] = append(subCategories[uniqueId], category)
	}

	return subCategories
}

func GetGoodsOtherImages(tableHash string, goodsIds []string, storeIds []string) map[string]GoodsOtherImage {
	goodsIds = utils.RemoveRepeat(goodsIds)
	if len(goodsIds) < 1 {
		return nil
	}

	storeIds = utils.RemoveRepeat(storeIds)
	sql := `SELECT
	goods_id,
	store_id,
	image 
FROM
	z_image_` + tableHash + ` 
WHERE
    store_id IN ( ` + strings.Join(storeIds, ",") + ` ) 
	AND goods_id IN ( ` + strings.Join(goodsIds, ",") + ` ) 
	AND category = 'goods' 
	AND obj_id = 0 
ORDER BY
	listorder ASC`

	goodsOtherImages := []GoodsOtherImage{}
	DB.Raw(sql).Find(&goodsOtherImages)

	otherImages := make(map[string]GoodsOtherImage)
	for _, image := range goodsOtherImages {
		uniqueId := UniqueId(image.StoreId, image.GoodsId)
		otherImages[uniqueId] = image
	}

	return otherImages
}

func GetGoodsSaleProperties(tableHash string, goodsIds []string, storeIds []string) map[string]GoodsSaleProperty {
	goodsIds = utils.RemoveRepeat(goodsIds)
	if len(goodsIds) < 1 {
		return nil
	}

	storeIds = utils.RemoveRepeat(storeIds)
	sql := `SELECT
	a.base_name,
	a.image ,
	a.goods_id,
	b.store_id
FROM
	z_goods_sale_prop_` + tableHash + ` a
	LEFT JOIN z_goods_sale_prop_` + tableHash + ` b ON a.parent_id = b.id 
WHERE
    b.store_id IN ( ` + strings.Join(storeIds, ",") + ` ) 
	AND b.goods_id IN ( ` + strings.Join(goodsIds, ",") + ` ) 
	AND b.multi_image = 1
ORDER BY
	a.listorder ASC`

	goodsSaleProperties := []GoodsSaleProperty{}
	DB.Raw(sql).Find(&goodsSaleProperties)

	saleProperties := make(map[string]GoodsSaleProperty)
	for _, property := range goodsSaleProperties {
		uniqueId := UniqueId(property.StoreId, property.GoodsId)
		saleProperties[uniqueId] = property
	}

	return saleProperties
}

func GetGoodsProperties(tableHash string, goodsIds []string, storeIds []string) map[string]GoodsProperty {
	goodsIds = utils.RemoveRepeat(goodsIds)
	if len(goodsIds) < 1 {
		return nil
	}

	storeIds = utils.RemoveRepeat(storeIds)
	sql := `SELECT
	goods_id,
	store_id,
	property_id AS property_id,
	value_id AS value_id,
	value_name AS value_name 
FROM
	z_goods_property_rel_` + tableHash + ` 
WHERE
    store_id IN ( ` + strings.Join(storeIds, ",") + ` ) 
	AND goods_id IN ( ` + strings.Join(goodsIds, ",") + ` ) 
	AND value_id != 0`

	goodsProperties := []GoodsProperty{}
	DB.Raw(sql).Find(&goodsProperties)

	properties := make(map[string]GoodsProperty)
	for _, property := range goodsProperties {
		uniqueId := UniqueId(property.StoreId, property.GoodsId)
		properties[uniqueId] = property
	}

	return properties
}

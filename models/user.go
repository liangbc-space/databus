package models

import (
	"databus/utils"
	"errors"
	"github.com/jinzhu/gorm"
	"time"
)

type User struct {
	BaseModel
	Username      string `json:"username" redis:"username"`
	Mobile        string `gorm:"unique_index" json:"mobile" redis:"mobile"`
	Password      string `json:"password" redis:"password"`
	LastLoginTime int64  `json:"last_login_time" redis:"last_login_time"`
	Status        uint   `json:"status" redis:"status"`
}

func (user *User) BeforeCreate(scope *gorm.Scope) error {
	defer user.BaseModel.BeforeCreate(scope)

	user.Status = 1
	return nil
}

func Lists() string {
	return DB.NewScope(User{}).TableName()
}

func (user *User) Insert() error {
	return DB.Create(user).Error
}

func FindUserByMobile(mobile string) (User, error) {
	user := User{}

	DB = DB.Where(map[string]interface{}{"status": 1})
	err := DB.Where(map[string]interface{}{"mobile": mobile}).Find(&user).Error

	return user, err
}

func Login(mobile string, password string) (token string, err error) {
	user, err := FindUserByMobile(mobile)
	if err != nil {
		return token, err
	}
	if user == (User{}) {
		return token, errors.New("用户未找到，请先注册后重试")
	}

	//	密码验证
	if user.Password != utils.Md5(password) {
		return token, errors.New("账号或密码不正确")
	}

	//	生成jwt令牌
	token, err = utils.GenerateToken(user.Id, user.Username, user.Mobile)
	if err != nil {
		return token, err
	}

	//	缓存用户基本信息到redis hash表
	/*key := "login_info:" + strconv.Itoa(int(user.Id))
	if _, err := utils.RedisClient.Exec("hmset", key, user); err != nil {
		return token, err
	}
	if _, err := utils.RedisClient.Exec("expire", key, 600); err != nil {
		return token, err
	}*/

	//	更新用户登录时间
	DB.Model(&user).UpdateColumns(User{LastLoginTime: time.Now().Unix()})

	return token, nil
}

/**
登出	添加jwtToken黑名单
*/
func Logout(jwtToken string) (err error) {

	key := "jwt_token:blacklist"
	if _, err := utils.RedisClient.Exec("sadd", key, jwtToken); err != nil {
		return err
	}

	expTime := time.Now().Add(utils.JWTTokenExpireDuration).Unix()
	if _, err := utils.RedisClient.Exec("expireat", key, expTime); err != nil {
		return err
	}

	return nil
}

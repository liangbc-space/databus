package utils

import (
	"crypto/md5"
	"encoding/hex"
	"regexp"
)

//	手机正则验证
func ValidateMobile(mobile string) bool {
	reg := `^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))(\d{8})$`

	regx := regexp.MustCompile(reg)

	return regx.MatchString(mobile)
}

//	md5
func Md5(str string) string {
	hash := md5.New()
	return hex.EncodeToString(hash.Sum([]byte(str)))
}

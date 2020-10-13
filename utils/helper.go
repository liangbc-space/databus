package utils

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
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

//	数组|slice去重和空
func RemoveRepeat(list []string) []string {
	output := make([]string, 0)
	m := make(map[string]interface{})
	for _, item := range list {
		if _, ok := m[item]; !ok && item != "" {
			output = append(output, item)
			m[item] = true
		}
	}

	return output
}

//	检测元素在array|slice中存在
func In(list []interface{}, item interface{}) (isIn bool) {
	if len(list) < 1 {
		return false
	}

	for _, val := range list {
		if val == item {
			isIn = true
		} else {
			isIn = false
		}
	}

	return isIn
}

//	gzip压缩
func GzipEncode(input []byte) ([]byte, error) {
	buffer := new(bytes.Buffer)
	writer := gzip.NewWriter(buffer)
	defer writer.Close()

	if _, err := writer.Write(input); err != nil {
		return nil, err
	}
	if err := writer.Flush(); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

//	gzip解压
func GzipDecode(input []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(input))
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	return ioutil.ReadAll(reader)
}

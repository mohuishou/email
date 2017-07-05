package utils

import (
	"encoding/json"
	"time"
)

//APP 应用信息
type APP struct {
	Name   string `json:"name"`   //应用名称
	Time   int64  `json:"time"`   //应用签署日期
	Author string `json:"author"` //应用拥有者
	IP     string `json:"ip"`     //应用ip
}

//NewAppToken 新建应用秘钥
func NewAppToken(name, author, ip string, aes AES) (token string) {
	app := APP{Name: name, Time: time.Now().Unix(), Author: author, IP: ip}
	appkey, err := json.Marshal(app)
	if err != nil {
		panic(err)
	}
	return aes.Encrypt(string(appkey))
}

//DecryptToken token解密，获取应用信息
func DecryptToken(token string, aes AES) (app APP, err error) {
	appkey, err := aes.Decrypt(token)
	if err != nil {
		return app, err
	}
	err = json.Unmarshal(appkey, &app)
	return app, err
}

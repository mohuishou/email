package main

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

//new app token 新建应用秘钥
func newAppToken(name, author, ip string) (token string) {
	app := APP{Name: name, Time: time.Now().Unix(), Author: author, IP: ip}
	a, err := json.Marshal(app)
	if err != nil {
		panic(err)
	}
	return encrypt(string(a))
}

//token解密，获取应用信息
func decryptToken(token string) (app APP, err error) {
	a, err := decrypt(token)
	if err != nil {
		return app, err
	}
	err = json.Unmarshal(a, &app)
	return app, err
}

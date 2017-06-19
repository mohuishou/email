package main

import (
	gomail "gopkg.in/gomail.v2"
)

// send mall
func mail(to, title, content string) *gomail.Message {
	m := gomail.NewMessage()
	//发信人
	m.SetHeader("From", config.address)
	//收信
	m.SetHeader("To", to)
	//主题
	m.SetHeader("Subject", title)
	//内容
	m.SetBody("text/html", content)

	return m
}

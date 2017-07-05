package utils

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

//System 系统配置
type System struct {
	WorkerNumber int    `yaml:"worker_number"`  //并发数目
	MaxQueueSize int    `yaml:"max_queue_size"` //队列最大数目
	Delay        string `yaml:"delay"`          //延迟
	Key          string `yaml:"key"`            //key
}

//Email 邮件配置
type Email struct {
	Address  string //地址
	Server   string //服务器地址
	Password string //密码
	Port     int    //端口号
}

//Config 配置
type Config struct {
	System System  `yaml:"system"` //系统配置
	Emails []Email `yaml:"emails"` //邮件配置
}

//GetConfig 获取配置文件
//filepath 配置文件地址
func GetConfig(filepath string) Config {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic("配置文件读取错误！")
	}
	config := Config{}
	err = yaml.Unmarshal([]byte(data), &config)
	if err != nil {
		panic("配置文件读取错误！")
	}
	return config
}

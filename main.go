// Original code with Dispatcher
package main

import (
	_ "expvar"
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"

	"os"

	"github.com/gpmgo/gopm/modules/goconfig"
	"github.com/mohuishou/email/utils"
	"gopkg.in/gomail.v2"
)

func requestHandler(w http.ResponseWriter, r *http.Request, jobQueue chan Job) {
	// Make sure we can only be called with an HTTP POST request.
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	token := r.Header.Get("token")
	if token == "" {
		errorRetrun(w, "应用暂未授权", "")
		return
	}
	app, err := decryptToken(token)
	if err != nil {
		errorRetrun(w, "应用授权信息错误！", "")
		return
	}

	// Set name and validate value.
	title := r.FormValue("title")
	if title == "" {
		errorRetrun(w, "title 必须！", "")
		return
	}

	to := r.FormValue("address")
	if to == "" {
		errorRetrun(w, "地址 必须！", "")
		return
	}

	content := r.FormValue("content")
	if content == "" {
		errorRetrun(w, "content必须！", "")
		return
	}

	// 入队
	m := mail(to, title, content)
	job := Job{email: m, app: app.Name}
	jobQueue <- job

	// 添加成功
	successReturn(w, "已添加到后台队列！", "")
}

//Config 配置项目
type Config struct {
	workerNumber int
	delay        time.Duration
	maxQueueSize int
	address      string
	server       string
	port         int
	password     string
	key          string
}

//全局配置项
var config = Config{}

//配置项目初始化
func initConfig(configFile string) {
	c, err := goconfig.LoadConfigFile(configFile)
	if err != nil {
		log.Fatalln("配置文件加载错误！")
	}

	//系统配置
	config.workerNumber = c.MustInt("system", "worker_number", 5)
	config.maxQueueSize = c.MustInt("system", "max_queue_size", 100)
	delay, _ := c.GetValue("system", "delay")
	config.delay, _ = time.ParseDuration(delay)

	//email配置
	config.address, _ = c.GetValue("email", "address")
	config.server, _ = c.GetValue("email", "server")
	config.port = c.MustInt("email", "port", 465)
	config.password, _ = c.GetValue("email", "password")

	//密钥配置
	config.key, _ = c.GetValue("token", "key")
}

//发送邮件
var sendMailer *gomail.Dialer

func main() {
	//命令输入
	var (
		port       = flag.String("p", "8080", "The server port")
		configFile = flag.String("c", "config.ini", "the config file")
	)
	flag.Parse()

	//加载配置文件
	initConfig(*configFile)
	log.Println("配置文件加载成功！")

	//新建应用
	if len(os.Args) > 1 {
		if os.Args[1] == "new" {
			var (
				name   = flag.String("name", "email", "应用名称")
				ip     = flag.String("ip", "0.0.0.0", "ip 地址")
				author = flag.String("author", "mohuishou", "拥有者")
			)
			flag.Parse()
			token := utils.NewAppToken(*name, *author, *ip)
			fmt.Printf("the [%s] token:%s \n", *name, token)
			fmt.Println(utils.DecryptToken(token))
			return
		}
	}

	//发送邮件配置
	sendMailer = gomail.NewDialer(config.server, config.port, config.address, config.password)

	// 创建应用队列
	jobQueue := make(chan Job, config.maxQueueSize)

	// 队列分发
	dispatcher := NewDispatcher(jobQueue, config.workerNumber)
	dispatcher.run()

	// Start the HTTP handler.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		requestHandler(w, r, jobQueue)
	})
	log.Println("应用启动中，监听端口：", *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))

}

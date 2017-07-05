// Original code with Dispatcher
package main

import (
	_ "expvar"
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"

	"os"

	"github.com/mohuishou/email/utils"
	"gopkg.in/gomail.v2"
)

func requestHandler(w http.ResponseWriter, r *http.Request, jobQueue chan Job, aes utils.AES) {
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
	app, err := utils.DecryptToken(token, aes)
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

func main() {
	//命令输入
	var (
		port       = flag.String("p", "8080", "The server port")
		configFile = flag.String("c", "config.yaml", "the config file")
	)
	flag.Parse()

	config := utils.GetConfig(*configFile)

	aes := utils.NewAES(config.System.Key)

	//新建应用
	if len(os.Args) > 1 {
		if os.Args[1] == "new" {
			var (
				name   = flag.String("name", "email", "应用名称")
				ip     = flag.String("ip", "0.0.0.0", "ip 地址")
				author = flag.String("author", "mohuishou", "拥有者")
			)
			flag.Parse()
			token := utils.NewAppToken(*name, *author, *ip, aes)
			fmt.Printf("the [%s] token:%s \n", *name, token)
			return
		}
	}

	// 创建应用队列
	jobQueue := make(chan Job, config.System.MaxQueueSize)

	for _, e := range config.Emails {
		//发送邮件配置
		sendMailer := gomail.NewDialer(e.Server, e.Port, e.Address, e.Password)

		// 队列分发
		dispatcher := NewDispatcher(jobQueue, config.System.WorkerNumber, sendMailer)
		dispatcher.run()
	}

	// Start the HTTP handler.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		requestHandler(w, r, jobQueue, aes)
	})
	log.Println("应用启动中，监听端口：", *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))

}

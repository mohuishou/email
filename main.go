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
	"gopkg.in/gomail.v2"
)

// Job holds the attributes needed to perform unit of work.
type Job struct {
	app   string
	email *gomail.Message
}

// NewWorker creates takes a numeric id and a channel w/ worker pool.
func NewWorker(id int, workerPool chan chan Job) Worker {
	return Worker{
		id:         id,
		jobQueue:   make(chan Job),
		workerPool: workerPool,
		quitChan:   make(chan bool),
	}
}

//Worker worker
type Worker struct {
	id         int
	jobQueue   chan Job
	workerPool chan chan Job
	quitChan   chan bool
}

func (w Worker) start() {
	go func() {
		for {
			// Add my jobQueue to the worker pool.
			w.workerPool <- w.jobQueue

			select {
			case job := <-w.jobQueue:
				// Dispatcher has added a job to my jobQueue.
				fmt.Printf("worker%d: started %s, blocking for %f seconds\n", w.id, job.app, config.delay.Seconds())
				time.Sleep(config.delay)
				//开始发送
				if err := sendMailer.DialAndSend(job.email); err != nil {
					log.Printf("[Mailer][%s]Error:%s", job.app, err.Error())
				} else {
					log.Printf("[Mailer][%s]Success: 发送成功", job.app)
				}
			case <-w.quitChan:
				// We have been asked to stop.
				fmt.Printf("worker%d stopping\n", w.id)
				return
			}
		}
	}()
}

func (w Worker) stop() {
	go func() {
		w.quitChan <- true
	}()
}

// NewDispatcher creates, and returns a new Dispatcher object.
func NewDispatcher(jobQueue chan Job, maxWorkers int) *Dispatcher {
	workerPool := make(chan chan Job, maxWorkers)

	return &Dispatcher{
		jobQueue:   jobQueue,
		maxWorkers: maxWorkers,
		workerPool: workerPool,
	}
}

//Dispatcher 分发
type Dispatcher struct {
	workerPool chan chan Job
	maxWorkers int
	jobQueue   chan Job
}

func (d *Dispatcher) run() {
	for i := 0; i < d.maxWorkers; i++ {
		worker := NewWorker(i+1, d.workerPool)
		worker.start()
	}

	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		case job := <-d.jobQueue:
			go func() {
				fmt.Printf("fetching workerJobQueue for: %s\n", job.app)
				workerJobQueue := <-d.workerPool
				fmt.Printf("adding %s to workerJobQueue\n", job.app)
				workerJobQueue <- job
			}()
		}
	}
}

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

	// Create Job and push the work onto the jobQueue.
	m := mail(to, title, content)
	job := Job{email: m, app: app.Name}
	jobQueue <- job

	// Render success.
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

var config = Config{}

func initConfig(configFile string) {
	c, err := goconfig.LoadConfigFile(configFile)
	if err != nil {
		log.Fatalln("配置文件加载错误！")
	}
	config.workerNumber = c.MustInt("system", "worker_number", 5)
	config.maxQueueSize = c.MustInt("system", "max_queue_size", 100)
	delay, _ := c.GetValue("system", "delay")
	config.delay, _ = time.ParseDuration(delay)
	config.address, _ = c.GetValue("email", "address")
	config.server, _ = c.GetValue("email", "server")
	config.port = c.MustInt("email", "port", 465)
	config.password, _ = c.GetValue("email", "password")

	config.key, _ = c.GetValue("token", "key")
}

var sendMailer *gomail.Dialer

func main() {
	//命令输入
	var (
		port       = flag.String("p", "8080", "The server port")
		configFile = flag.String("c", "config.ini", "the config file")
	)
	flag.Parse()
	initConfig(*configFile)

	if len(os.Args) > 1 {
		if os.Args[1] == "new" {
			var (
				name   = flag.String("name", "email", "应用名称")
				ip     = flag.String("ip", "0.0.0.0", "ip 地址")
				author = flag.String("author", "mohuishou", "拥有者")
			)
			flag.Parse()
			token := newAppToken(*name, *author, *ip)
			fmt.Printf("the [%s] token:%s \n", *name, token)
			fmt.Println(decryptToken(token))
			return
		}
	}

	//加载配置文件
	sendMailer = gomail.NewDialer(config.server, config.port, config.address, config.password)

	// Create the job queue.
	jobQueue := make(chan Job, config.maxQueueSize)

	// Start the dispatcher.
	dispatcher := NewDispatcher(jobQueue, config.workerNumber)
	dispatcher.run()

	// Start the HTTP handler.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		requestHandler(w, r, jobQueue)
	})
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

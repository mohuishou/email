package main

import (
	"fmt"
	"log"
	"time"

	gomail "gopkg.in/gomail.v2"
)

// Job holds the attributes needed to perform unit of work.
type Job struct {
	app   string
	email Email
}

//Email 邮件
type Email struct {
	address string //地址
	title   string //标题
	content string //内容
}

//send 发送邮件
func (e Email) send(sender *gomail.Dialer) error {
	m := gomail.NewMessage()
	//发信人
	m.SetHeader("From", sender.Username)
	//收信
	m.SetHeader("To", e.address)
	//主题
	m.SetHeader("Subject", e.title)
	//内容
	m.SetBody("text/html", e.content)

	return sender.DialAndSend(m)
}

// NewWorker creates takes a numeric id and a channel w/ worker pool.
func NewWorker(id int, workerPool chan chan Job, sendMailer *gomail.Dialer, delay time.Duration) Worker {
	return Worker{
		id:         id,
		jobQueue:   make(chan Job),
		workerPool: workerPool,
		quitChan:   make(chan bool),
		sendMailer: sendMailer,
		delay:      delay,
	}
}

//Worker worker
type Worker struct {
	id         int
	jobQueue   chan Job
	workerPool chan chan Job
	quitChan   chan bool
	sendMailer *gomail.Dialer
	delay      time.Duration
}

func (w Worker) start() {
	go func() {
		for {
			// Add my jobQueue to the worker pool.
			w.workerPool <- w.jobQueue

			select {
			case job := <-w.jobQueue:
				time.Sleep(w.delay)
				//开始发送
				if err := job.email.send(w.sendMailer); err != nil {
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
func NewDispatcher(jobQueue chan Job, maxWorkers int, sendMailer *gomail.Dialer, delay time.Duration) *Dispatcher {
	workerPool := make(chan chan Job, maxWorkers)

	return &Dispatcher{
		jobQueue:   jobQueue,
		maxWorkers: maxWorkers,
		workerPool: workerPool,
		sendMailer: sendMailer,
		delay:      delay,
	}
}

//Dispatcher 分发
type Dispatcher struct {
	workerPool chan chan Job
	maxWorkers int
	jobQueue   chan Job
	sendMailer *gomail.Dialer
	delay      time.Duration
}

func (d *Dispatcher) run() {
	for i := 0; i < d.maxWorkers; i++ {
		worker := NewWorker(i+1, d.workerPool, d.sendMailer, d.delay)
		worker.start()
	}

	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		case job := <-d.jobQueue:
			go func() {
				workerJobQueue := <-d.workerPool
				workerJobQueue <- job
			}()
		}
	}
}

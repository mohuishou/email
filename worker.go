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
	email *gomail.Message
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
				if err := w.sendMailer.DialAndSend(job.email); err != nil {
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
func NewDispatcher(jobQueue chan Job, maxWorkers int, sendMailer *gomail.Dialer) *Dispatcher {
	workerPool := make(chan chan Job, maxWorkers)

	return &Dispatcher{
		jobQueue:   jobQueue,
		maxWorkers: maxWorkers,
		workerPool: workerPool,
		sendMailer: sendMailer,
	}
}

//Dispatcher 分发
type Dispatcher struct {
	workerPool chan chan Job
	maxWorkers int
	jobQueue   chan Job
	sendMailer *gomail.Dialer
}

func (d *Dispatcher) run() {
	for i := 0; i < d.maxWorkers; i++ {
		worker := NewWorker(i+1, d.workerPool, d.sendMailer)
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

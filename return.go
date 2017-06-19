package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type jsonData struct {
	Status int         `json:"status"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
}

func successReturn(w http.ResponseWriter, msg string, data interface{}) {
	jsonReturn(w, 1, msg, data)
}

func errorRetrun(w http.ResponseWriter, msg string, data interface{}) {
	jsonReturn(w, 0, msg, data)
}

func jsonReturn(w http.ResponseWriter, status int, msg string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	d := jsonData{
		Status: status,
		Msg:    msg,
		Data:   data,
	}
	dataByte, err := json.Marshal(d)
	if err != nil {
		fmt.Fprint(w, err)
	}
	w.Write(dataByte)
	return
}

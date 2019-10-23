package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"time"
)

// struct for http response
type EulerStruct struct {
	StartTime int64 `json: "startTime"`
	FetchTime int64 `json: "fetchTime"`
	Attempts  int32 `json: "attempts"`
}

// struct for writing out to log file
type ChanStruct struct {
	Ip              string `json: "ip"`
	LastFetchedTime int64  `json: "lastFetchTime"`
	StartTime       int64  `json: "startTime"`
}

// struct to capture unixtime from api
type ApiTime struct {
	UnixTime int64 `json: "unixtime"`
}

// function to retrieve IP calls
func GetIP() string {
	response, err := http.Get("https://api.ipify.org?format=text/plain")

	if err != nil {
		panic(err)
	}

	ip, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
		error.Error(err)
	}
	return string(ip)
}

func WorldApiBase(ip string) string {
	http, err := http.NewRequest("GET", "http://worldtimeapi.org/api/ip", nil)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	q := http.URL.Query()
	q.Add("ip", ip)
	http.URL.RawQuery = q.Encode()
	fmt.Printf(" %s <--- %s ", http.URL.String(), "successfully connected!")
	return http.URL.String()
}

func fetchTime(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		data, err := json.Marshal(EulerObject)
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		var vStruct ChanStruct
		vStruct.LastFetchedTime = EulerObject.FetchTime
		vStruct.StartTime = EulerObject.StartTime
		vStruct.Ip = req.RemoteAddr

		go writeToFile(chann)
		chann <- vStruct

		// fmt.Println("This is after the go statement")
		fmt.Fprintf(w, "%s", data)
	}
}

func writeToFile(ch <-chan ChanStruct) {
	select {
	case msg := <-ch:
		var file, err = os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}

		defer file.Close()
		res := msg.Ip + fmt.Sprintf("<%d>", msg.LastFetchedTime) + fmt.Sprintf("<%d>", msg.StartTime) + "\n"
		_, err = file.WriteString(res)
		if err != nil {
			log.Fatal(err)
			error.Error(err)
			return
		}

	}
	// fmt.Println("Written to file")
}

func EulerTime(ip string) {

	t := float64(math.E)
	for range time.Tick(time.Second * time.Duration(int64(t))) {
		EulerObject.Attempts++

		resp, err := http.Get(ip)
		if err != nil {
			log.Fatal(err)
		}
		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			log.Fatal(err)
		}
		var apiTimeObject ApiTime
		err = json.Unmarshal([]byte(body), &apiTimeObject)
		if err != nil {
			log.Fatal(err)
		}
		EulerObject.FetchTime = apiTimeObject.UnixTime
	}
}

//Global Variable
var EulerObject EulerStruct
var chann = make(chan ChanStruct)

func main() {

	// get public ip and create a url for worldapibase
	ip := GetIP()
	worldIpAdd := WorldApiBase(ip)
	EulerObject.StartTime = time.Now().Unix()
	_, err := os.Create("logs.txt")
	if err != nil {
		log.Fatal(err)
		fmt.Println("log.txt could not be created")
		return
	}

	go EulerTime(worldIpAdd)

	http.HandleFunc("/", fetchTime)
	http.ListenAndServe(":12345", nil)
}

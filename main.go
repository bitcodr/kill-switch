package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var initIP string

func init() {
	ip := getIP()
	os.Stdout.WriteString("kill-switch initilized by IP: " + ip + "\n")
	initIP = ip
}

func main() {
	for {
		ip := getIP()
		if ip != initIP {
			os.Stdout.WriteString("ip changed\n")
		} else {
			os.Stdout.WriteString("ip not chnaged\n")
		}

		time.Sleep(1 * time.Second)
	}
}

func getIP() string {
	url := "https://api.ipify.org?format=text"
	res, err := http.Get(url)
	if err != nil {
		os.Stderr.WriteString("cannot fetch public ip\n")
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		os.Stderr.WriteString("cannot read body\n")
	}
	return string(body)
}

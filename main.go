package main

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var initIP string

func init() {
	url := "https://api.ipify.org?format=text"
	res, err := http.Get(url)
	if err != nil {
		os.Stdout.WriteString("cannot fetch public ip\n")
		os.Exit(1)
	}
	defer res.Body.Close()

	parsedData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		os.Stderr.WriteString("cannot parse body\n")
		os.Exit(1)
	}

	ip := string(parsedData)

	os.Stdout.WriteString("kill-switch initilized by IP: " + ip + "\n")
	initIP = ip
}

func main() {
	ctx := context.Background()
	for {
		func() {
			ip := getIP(ctx)
			if ip != initIP {
				os.Stdout.WriteString("ip changed\n")
			} else {
				os.Stdout.WriteString("ip not chnaged\n")
			}

			time.Sleep(500 * time.Millisecond)
		}()
	}
}

func getIP(ctx context.Context) string {

	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	errorchan := make(chan error, 1)
	ipchan := make(chan string, 1)

	go func(ctx context.Context, errorchan chan<- error, ipchan chan<- string) {
		defer close(errorchan)
		defer close(ipchan)

		url := "https://api.ipify.org?format=text"

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			errorchan <- errors.New("cannot make request\n")
			return
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			errorchan <- errors.New("cannot perform request\n")
			return
		}

		if res.Body == nil {
			errorchan <- errors.New("response body is empty\n")
			return
		}

		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			errorchan <- errors.New("cannot read body\n")
			return
		}

		ipchan <- string(body)
	}(ctx, errorchan, ipchan)

	select {
	case <-ctx.Done():
		os.Stderr.WriteString("timeout\n")
	case err := <-errorchan:
		os.Stderr.WriteString(err.Error() + "\n")
	case ip := <-ipchan:
		return ip
	}

	return ""
}

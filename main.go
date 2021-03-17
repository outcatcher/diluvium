package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const errorMessage = "failed"

var (
	count int
	url   string
)

func init() {
	flag.StringVar(&url, "url", "", "URL to GET")
	flag.IntVar(&count, "count", 100, "Number of requests")
}

func request(address string) (string, error) {
	resp, err := http.Get(address)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func main() {
	flag.Parse()

	if url == "" {
		print("URL is missing, please provide the URL")
		os.Exit(2)
	}

	resChan := make(chan string, count)
	errChan := make(chan error, count)
	for i := 0; i < count; i++ {
		go func() {
			res, err := request(url)
			errChan <- err
			resChan <- res
		}()
	}

	stats := map[string]int{
		errorMessage: 0,
	}
	for i := 0; i < count; i++ {
		err := <-errChan
		res := <-resChan

		if err != nil {
			print(err.Error())
			stats[errorMessage] += 1
			continue
		}

		res = strings.Trim(res, "\n ")
		if _, ok := stats[res]; !ok {
			stats[res] = 0
		}
		stats[res] += 1
	}
	msg, _ := json.MarshalIndent(stats, "", "    ")
	fmt.Printf("\n=======================\nResults:\n%s", msg)
}

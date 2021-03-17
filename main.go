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

var (
	count int
	url   string
)

func init() {
	flag.StringVar(&url, "url", "", "URL to GET")
	flag.IntVar(&count, "count", 100, "Number of requests")
}

func request(address string) string {
	resp, err := http.Get(address)
	if err != nil {
		return fmt.Sprintf("failed: %s", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		print(err)
	}
	return string(body)
}

func main() {
	flag.Parse()

	if url == "" {
		print("URL is missing, please provide the URL")
		os.Exit(2)
	}

	resChan := make(chan string, count)
	for i := 0; i < count; i++ {
		go func() {
			resChan <- request(url)
		}()
	}

	stats := make(map[string]int, 2)
	for i := 0; i < count; i++ {
		res := <-resChan
		res = strings.Trim(res, "\n ")
		if _, ok := stats[res]; !ok {
			stats[res] = 0
		}
		stats[res] += 1
	}
	msg, _ := json.MarshalIndent(stats, "", "  ")
	fmt.Printf("%s", msg)
}

package main

import (
	"fmt"
	"github.com/kr/pretty"
	"net/http"
	"time"
)

type HttpResponse struct {
	url      string
	response *http.Response
	err      error
}

func asyncHttpGets(urls []string) []*HttpResponse {
	ch := make(chan *HttpResponse)
	responses := []*HttpResponse{}
	client := http.Client{}
	for _, url := range urls {
		go func(url string) {
			fmt.Printf("Fetching %s \n", url)
			resp, err := client.Get(url)
			ch <- &HttpResponse{url, resp, err}
			if err != nil && resp != nil && resp.StatusCode == http.StatusOK {
				resp.Body.Close()
			}
		}(url)
	}

	for {
		select {
		case r := <-ch:
			fmt.Printf("%s was fetched\n", r.url)
			if r.err != nil {
				fmt.Println("with an error", r.err)
			}
			responses = append(responses, r)
			if len(responses) == len(urls) {
				return responses
			}
		case <-time.After(50 * time.Millisecond):
			fmt.Printf(".")
		}
	}
	return responses
}

var urls = []string{
	"http://www.rubyconf.org/",
	"http://golang.org/",
	"http://matt.aimonetti.net/",
	"http://localhost:8080/",
}

func handler(w http.ResponseWriter, req *http.Request) {
	//	pretty.Println(r)
	req.ParseForm()
	pretty.Println(req.Form.Get("foo"))
	fmt.Fprintf(w, "Hi there, I love %s!", req.URL.Path[1:])
}

func busy() {
	for {
		fmt.Println("tick")
		time.Sleep(1000000000)
	}
}

func getResults() {
	results := asyncHttpGets(urls)
	for _, result := range results {
		if result != nil && result.response != nil {
			fmt.Printf("%s status: %s\n", result.url,
				result.response.Status)
		}
	}

}
func doServer() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func main() {
	go doServer()
	go getResults()
	busy()
}

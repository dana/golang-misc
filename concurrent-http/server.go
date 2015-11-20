package main

import (
	"fmt"
	"github.com/kr/pretty"
	"net/http"
)

func handler(w http.ResponseWriter, req *http.Request) {
	//	pretty.Println(r)
	req.ParseForm()
	pretty.Println(req.Form.Get("foo"))
	fmt.Fprintf(w, "Hi there, I love %s!", req.URL.Path[1:])
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

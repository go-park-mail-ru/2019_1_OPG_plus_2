package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type HelloWorld struct {
	Status int    `json:"status, int"`
	Text   string `json:"text, string"`
}

func NewHelloWorld(status int, text string) *HelloWorld {
	return &HelloWorld{Status: status, Text: text}
}

func runServer(addr string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/",
		func(res http.ResponseWriter, req *http.Request) {
			_, err := fmt.Fprintln(res, "<h1>Hello world!</h1>")
			if err != nil {
				panic(err)
			}
		})

	mux.HandleFunc("/api",
		func(res http.ResponseWriter, req *http.Request) {
			_, err := fmt.Fprintln(res, "<h1>CHOOSE API METHOD!</h1>")
			if err != nil {
				panic(err)
			}
		})

	mux.HandleFunc("/api/helloworld",
		func(res http.ResponseWriter, req *http.Request) {
			helloWorldResponse := NewHelloWorld(200, "HelloWorld!")
			result, _ := json.Marshal(helloWorldResponse)
			fmt.Println(string(result))
			_, _ = fmt.Fprintln(res, string(result))
		})

	server := http.Server{
		Addr:    addr,
		Handler: mux,
	}

	fmt.Println("starting server at", addr)
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func main() {
	runServer(":8080")
}

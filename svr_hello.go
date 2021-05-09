package main

import (
	"context"
	"fmt"
	"net/http"
)

var (
	hello_config struct {
		ListenAddr    string
	}
	hello_svr 	*SvrConfig
	hello_http    *http.Server
)

func init() {
	hello_svr = &SvrConfig{
		Name:   "hello",
		Config: &hello_config,
		Run:    hello_run,
	}
	Register(hello_svr)
}

func hello_run() error {
	fmt.Println("run hello begin")

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", hello)

	// server
	hello_http = &http.Server{
		Addr:    hello_config.ListenAddr,
		Handler: mux,
	}

	go func() {
		select {
		case <-hello_svr.ErrCtx.Done():
			fmt.Println("Shutdown hello")
			hello_http.Shutdown(context.Background())
			return
		}
	}()

	err := hello_http.ListenAndServe()
	fmt.Println("run hello end")
	return err
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	http.Error(w,"hello",200)
}
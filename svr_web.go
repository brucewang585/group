package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"path"
)

var (
	web_config struct {
		ListenAddr    string
		StaticPath    string
	}
	web_svr 	*SvrConfig
	web_http    *http.Server
)

func init() {
	web_svr = &SvrConfig{
		Name:   "web",
		Config: &web_config,
		Run:    web_run,
	}
	Register(web_svr)
}

func web_run() error {
	fmt.Println("run web begin")

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", website)

	// server
	web_http = &http.Server{
		Addr:    web_config.ListenAddr,
		Handler: mux,
	}

	go func() {
		select {
		case <-web_svr.ErrCtx.Done():
			fmt.Println("Shutdown web")
			web_http.Shutdown(context.Background())
			return
		}
	}()

	err := web_http.ListenAndServe()
	fmt.Println("run web end")
	return err
}

func website(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Path
	if filePath == "/" {
		filePath = "/index.html"
	}
	if mime := mime.TypeByExtension(path.Ext(filePath)); mime != "" {
		w.Header().Set("Content-Type", mime)
	}
	if f, err := ioutil.ReadFile(web_config.StaticPath + filePath); err == nil {
		if _, err = w.Write(f); err != nil {
			w.WriteHeader(505)
		}
	} else {
		w.Header().Set("Location", "/")
		w.WriteHeader(302)
	}
}

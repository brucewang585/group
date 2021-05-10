package websvr

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"path"

	"group/core"
)

type config struct {
	ListenAddr    string
	StaticPath    string
}

type WebSvr struct {
	name string
	cfg  *config
	ctx  context.Context
	http *http.Server
}

var (
	svr 	*WebSvr
)

func init() {
	svr = &WebSvr{
		name:   "web",
		cfg : &config{},
	}
	core.Register(svr)
}

func (s *WebSvr) Run() error {
	fmt.Println("run web begin")

	mux := http.NewServeMux()
	mux.HandleFunc("/", web)

	// server
	s.http = &http.Server{
		Addr:    s.cfg.ListenAddr,
		Handler: mux,
	}

	//接受关闭消息
	go func() {
		select {
		case <-s.ctx.Done():
			fmt.Println("Shutdown web")
			s.http.Shutdown(context.Background())
			return
		}
	}()

	err := s.http.ListenAndServe()
	fmt.Println("run web end")
	return err
}

func (h *WebSvr) Name() string {
	return h.name
}

func (h *WebSvr) SetConfig(cfg []byte) error {
	return json.Unmarshal(cfg, h.cfg)
}

func (h *WebSvr) SetCtx(ctx context.Context) {
	h.ctx = ctx
}

func web(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Path
	if filePath == "/" {
		filePath = "/index.html"
	}
	if mime := mime.TypeByExtension(path.Ext(filePath)); mime != "" {
		w.Header().Set("Content-Type", mime)
	}
	if f, err := ioutil.ReadFile(svr.cfg.StaticPath + filePath); err == nil {
		if _, err = w.Write(f); err != nil {
			w.WriteHeader(505)
		}
	} else {
		w.Header().Set("Location", "/")
		w.WriteHeader(302)
	}
}


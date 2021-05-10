package hellosvr

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"group/core"
)

type config struct {
	ListenAddr string
}

type HelloSvr struct {
	name string
	cfg  *config
	ctx  context.Context
	http *http.Server
}

var (
	svr *HelloSvr
)

func init() {
	svr = &HelloSvr{
		name: "hello",
		cfg:  &config{},
	}
	core.Register(svr)
}

func (s *HelloSvr) Run() error {
	fmt.Println("run hello begin")

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", hello)

	// server
	s.http = &http.Server{
		Addr:    s.cfg.ListenAddr,
		Handler: mux,
	}

	//接受关闭消息
	go func() {
		select {
		case <-s.ctx.Done():
			fmt.Println("Shutdown hello")
			s.http.Shutdown(context.Background())
			return
		}
	}()

	err := s.http.ListenAndServe()
	fmt.Println("run hello end")
	return err
}

func (h *HelloSvr) Name() string {
	return h.name
}

func (h *HelloSvr) SetConfig(cfg []byte) error {
	return json.Unmarshal(cfg, h.cfg)
}

func (h *HelloSvr) SetCtx(ctx context.Context) {
	h.ctx = ctx
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	http.Error(w, "hello", 200)
}

package main

import "golang.org/x/net/context"

var Svrs = make(map[string]*SvrConfig)

//SvrConfig 服务配置定义
type SvrConfig struct {
	Name      string                       //服务名称
	Config    interface{}                  //服务配置
	Run       func() error                 //服务启动函数
	ErrCtx    context.Context              //外部上下文
}

// RegisterSvr 注册服务
func Register(opt *SvrConfig) {
	Svrs[opt.Name] = opt
}

type BaseSvr interface {
	Run() error
}

//校验是否有协程已发生错误
func CheckGoroutineErr(cfg *SvrConfig) error {
	select {
	case <-cfg.ErrCtx.Done():
		return cfg.ErrCtx.Err()
	default:
		return nil
	}
}

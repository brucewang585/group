package core

import "golang.org/x/net/context"

var Svrs = make(map[string] BaseSvr)

//Svr 服务配置定义
type BaseSvr interface {
	Run() error

	Name() string
	SetConfig(cfg []byte) error
	SetCtx(ctx context.Context)
}

// RegisterSvr 注册服务
func Register(svr BaseSvr) {
	Svrs[svr.Name()] = svr
}



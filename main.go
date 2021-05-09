package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	xContext "golang.org/x/net/context"
	"golang.org/x/sync/errgroup"
)

var (
	group 		*errgroup.Group
	errCtx    	xContext.Context              //外部上下文
)

func main() {
	addr := flag.String("c", "", "config file")
	flag.Parse()
	if *addr == "" {
		Run(filepath.Join(filepath.Dir(os.Args[0]), "config.toml"))
	} else {
		Run(*addr)
	}
}

// Run
func Run(configFile string) (err error) {
	var (
		cfg         map[string]interface{}
		configData []byte

	)

	if configData, err = ioutil.ReadFile(configFile); err != nil {
		return
	}
	if _, err = toml.Decode(string(configData), &cfg);err != nil {
		return
	}

	group, errCtx = errgroup.WithContext(xContext.Background())

	//启动服务
	for name, config := range Svrs {
		if cfg, ok := cfg[name]; ok {
			b, _ := json.Marshal(cfg)
			if err = json.Unmarshal(b, config.Config); err != nil {
				continue
			}
			//填充
			config.ErrCtx = errCtx
		} else {
			continue
		}

		if config.Run != nil {
			group.Go(config.Run)
		}
	}

	group.Go(sign_monitor)

	group.Wait()
	fmt.Println("main exit")

	return
}

func sign_monitor() (err error){
	sigc := make(chan os.Signal, 1)
	defer close(sigc)

	//
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigc)

	//
	select {
	case <- sigc:
		fmt.Println("recv sign exit")
		return errors.New("sign exit")
	case <- errCtx.Done():
		fmt.Println("close sign")
		return nil
	}
}




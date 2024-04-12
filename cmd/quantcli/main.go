package main

import (
	"context"
	"log"
	"os"

	quantcfg "github.com/tensorchen/quant/config"

	"github.com/longportapp/openapi-go/config"
	"github.com/longportapp/openapi-go/trade"
	"gopkg.in/yaml.v3"
)

func main() {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalln(err)
	}

	var cfg quantcfg.Config
	if err = yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalln(err)
	}

	conf, err := config.New(config.WithConfigKey(cfg.LongBridge.AppKey, cfg.LongBridge.AppSecret, cfg.LongBridge.AccessToken))
	if err != nil {
		log.Fatalln(err)
	}

	tc, err := trade.NewFromCfg(conf)
	if err != nil {
		log.Fatalln(err)
	}

	rsp, err := tc.StockPositions(context.Background(), []string{"ARM.US"})
	if err != nil {
		log.Fatalln(err)
	}

	for _, sps := range rsp {
		for _, sp := range sps.Positions {
			log.Println(sp)
		}
	}
}

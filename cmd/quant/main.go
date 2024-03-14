package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/tensorchen/quant/config"
	"github.com/tensorchen/quant/handler"
	"github.com/tensorchen/quant/logger"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v3"
)

func main() {

	logger.NewLogger()

	data, err := os.ReadFile("config.yaml")
	if err != nil {
		logger.Logger().Fatal(err)
	}

	var cfg config.Config
	if err = yaml.Unmarshal(data, &cfg); err != nil {
		logger.Logger().Fatal(err)
	}

	h, err := handler.New(cfg)
	if err != nil {
		logger.Logger().Fatal(err)
	}

	r := mux.NewRouter()
	r.Handle("/tradingview", h)

	addr := fmt.Sprint(":" + cfg.Tquant.Port)
	logger.Logger().Fatal(http.ListenAndServe(addr, r))
}

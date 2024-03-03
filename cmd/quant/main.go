package main

import (
	"fmt"
	"github.com/tensorchen/quant/env"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/tensorchen/quant/handler"
	"github.com/tensorchen/quant/logger"
)

func main() {
	logger.NewLogger()

	h, err := handler.New()
	if err != nil {
		logger.Logger().Fatal(err)
	}

	r := mux.NewRouter()
	r.Handle("/tradingview", h)

	addr := fmt.Sprint(":" + os.Getenv(env.TquantPortKey))
	logger.Logger().Fatal(http.ListenAndServe(addr, r))
}

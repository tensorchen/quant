package main

import (
	"net/http"

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

	logger.Logger().Fatal(http.ListenAndServe(":16666", r))
}

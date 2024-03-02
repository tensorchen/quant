package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tensorchen/quant/entity"
	"github.com/tensorchen/quant/env"
	"github.com/tensorchen/quant/logger"
	"github.com/tensorchen/quant/notify"
	discordnotifier "github.com/tensorchen/quant/notify/discord"
	"github.com/tensorchen/quant/stock"
	"github.com/tensorchen/quant/stock/longbridge"
	"io"
	"net/http"
	"os"

	"github.com/disgoorg/disgo/discord"
)

type Response struct {
	Token string       `json:"token"`
	Trade entity.Trade `json:"trade"`
}

var _ http.Handler = (*Handler)(nil)

type Handler struct {
	notifier notify.Notifier
	stock    stock.Stock

	token string
}

func (h *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	var info entity.Information
	info.Color = 12397228
	defer func() {
		info.Title = "星火计划"
		if err := h.notifier.Notify(ctx, info); err != nil {
			logger.Logger().Error(err)
		}
	}()

	body, err := io.ReadAll(request.Body)
	if err != nil {
		handleError(&info, writer, err)
		return
	}

	var rsp Response
	if err = json.Unmarshal(body, &rsp); err != nil {
		handleError(&info, writer, errors.New(err.Error()+" : "+string(body)))
		return
	}

	if rsp.Token != h.token {
		err = errors.New(fmt.Sprintf("token [%s] 校验失败", rsp.Token))
		handleError(&info, writer, err)
		return
	}

	if err := h.stock.SubmitOrder(ctx, rsp.Trade); err != nil {
		handleError(&info, writer, err)
		return
	}

	info.Fields = []discord.EmbedField{
		{
			Name:   "股票",
			Value:  rsp.Trade.Ticker,
			Inline: newTrue(),
		},
		{
			Name:   "操作",
			Value:  rsp.Trade.Strategy.Order.Action,
			Inline: newTrue(),
		},
		{
			Name:   "数量",
			Value:  rsp.Trade.Strategy.Order.Contracts,
			Inline: newTrue(),
		},
	}

	info.Description = "✅ 交易执行成功" + ":scroll: "
}

func handleError(info *entity.Information, writer http.ResponseWriter, err error) {
	info.Description = "❌ " + err.Error()
	logger.Logger().Error(err)
	http.Error(writer, err.Error(), http.StatusInternalServerError)
}

func newTrue() *bool {
	t := true
	return &t
}

func New() (http.Handler, error) {

	//stockName := os.Getenv(env.TquantStockKey)
	//notifyName := os.Getenv(env.TquantNotifyKey)

	stock, err := longbridge.New()
	if err != nil {
		return nil, err
	}

	notifier, err := discordnotifier.New()
	if err != nil {
		return nil, err
	}

	return &Handler{
		stock:    stock,
		notifier: notifier,
		token:    os.Getenv(env.TquantTokenKey),
	}, nil
}

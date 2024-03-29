package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/tensorchen/quant/config"
	"github.com/tensorchen/quant/entity"
	"github.com/tensorchen/quant/logger"
	"github.com/tensorchen/quant/notify"
	discordnotifier "github.com/tensorchen/quant/notify/discord"
	"github.com/tensorchen/quant/stock"
	"github.com/tensorchen/quant/stock/longbridge"

	"github.com/disgoorg/disgo/discord"
)

type Request struct {
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
	var fields []discord.EmbedField

	defer func() {
		info.Title = "星火计划"
		info.Fields = fields
		if err := h.notifier.Notify(ctx, info); err != nil {
			logger.Logger().Error(err)
		}
	}()

	body, err := io.ReadAll(request.Body)
	if err != nil {
		handleError(&info, writer, err)
		return
	}

	logger.Logger().Info("交易请求开始: ", string(body))

	clearedBody, err := clearTokenValue(body)
	if err != nil {
		handleError(&info, writer, errors.New(err.Error()))
		return
	}

	info.Footer = &discord.EmbedFooter{
		Text: "请求参数：" + string(clearedBody),
	}

	var req Request
	if err = json.Unmarshal(body, &req); err != nil {
		handleError(&info, writer, errors.New(err.Error()+" : "+string(body)))
		return
	}

	if req.Token != h.token {
		err = errors.New(fmt.Sprintf("token [%s] 校验失败", req.Token))
		handleError(&info, writer, err)
		return
	}

	if err := h.stock.SubmitOrder(ctx, req.Trade); err != nil {
		handleError(&info, writer, err)
		return
	}

	logger.Logger().Info("交易请求结束: ", req.Trade)

	fields = []discord.EmbedField{
		{
			Name:   "股票",
			Value:  req.Trade.Ticker,
			Inline: newTrue(),
		},
		{
			Name:   "操作",
			Value:  req.Trade.Strategy.Order.Action,
			Inline: newTrue(),
		},
		{
			Name:   "数量",
			Value:  req.Trade.Strategy.Order.Contracts,
			Inline: newTrue(),
		},
		{
			Name:   "价格",
			Value:  req.Trade.Strategy.Order.Price,
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

func New(cfg config.Config) (http.Handler, error) {

	stock, err := longbridge.New(cfg.LongBridge.AppKey, cfg.LongBridge.AppSecret, cfg.LongBridge.AccessToken)
	if err != nil {
		return nil, err
	}

	notifier, err := discordnotifier.New(cfg.Discord.ID, cfg.Discord.Token)
	if err != nil {
		return nil, err
	}

	return &Handler{
		stock:    stock,
		notifier: notifier,
		token:    cfg.Tquant.Token,
	}, nil
}

func clearTokenValue(data []byte) ([]byte, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}

	if _, ok := m["token"]; ok {
		m["token"] = ""
	}

	return json.MarshalIndent(m, "", "\t")
}

package longbridge

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/tensorchen/quant/entity"
	"github.com/tensorchen/quant/stock"

	"github.com/longportapp/openapi-go/config"
	"github.com/longportapp/openapi-go/trade"
)

var _ stock.Stock = (*LongBridge)(nil)

type LongBridge struct {
	tc *trade.TradeContext
}

func (lb *LongBridge) SubmitOrder(ctx context.Context, tradeOrder entity.Trade) error {

	order := tradeOrder.Strategy.Order
	action := order.Action
	if action != "buy" && action != "sell" {
		return errors.New(fmt.Sprintf("交易行为 [%s] 不支持", action))
	}

	var side trade.OrderSide
	if action == "buy" {
		side = trade.OrderSideBuy
	}
	if action == "sell" {
		side = trade.OrderSideSell
	}

	contracts, err := strconv.Atoi(order.Contracts)
	if err != nil {
		return errors.New(fmt.Sprintf("contracts [%s] [%v]格式错误", order.Contracts, err))
	}

	submitOrder := &trade.SubmitOrder{
		Symbol:            tradeOrder.Ticker + ".US",
		OrderType:         trade.OrderTypeMO,
		Side:              side,
		SubmittedQuantity: uint64(contracts),
		TimeInForce:       trade.TimeTypeDay,
	}

	_, err = lb.tc.SubmitOrder(ctx, submitOrder)
	return err
}

func New(appKey, appSecret, accessToken string) (*LongBridge, error) {
	conf, err := config.New(config.WithConfigKey(appKey, appSecret, accessToken))
	if err != nil {
		return nil, err
	}

	tc, err := trade.NewFromCfg(conf)
	if err != nil {
		return nil, err
	}

	return &LongBridge{tc: tc}, nil
}

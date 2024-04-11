package longbridge

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/shopspring/decimal"
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

	marketPosition := tradeOrder.Strategy.MarketPosition
	if marketPosition != "long" && marketPosition != "short" && marketPosition != "flat" {
		return errors.New(fmt.Sprintf("marketPostion [%s] 不支持", marketPosition))
	}

	noop, err := lb.judgeNoop(ctx, tradeOrder)
	if err != nil {
		return errors.New(fmt.Sprintf("判断是否实际下单 [%v] 异常", err))
	}

	if noop {
		return errors.New("经判断无需要进行实际下单操作")
	}

	contracts, err := strconv.Atoi(order.Contracts)
	if err != nil {
		return errors.New(fmt.Sprintf("contracts [%s] [%v]格式错误", order.Contracts, err))
	}

	price, err := decimal.NewFromString(order.Price)
	if err != nil {
		return errors.New(fmt.Sprintf("contracts [%s] [%v]格式错误", order.Contracts, err))
	}
	price.StringFixed(2)

	var side trade.OrderSide
	if action == "buy" {
		side = trade.OrderSideBuy
	}
	if action == "sell" {
		side = trade.OrderSideSell
	}

	submitOrder := &trade.SubmitOrder{
		Symbol:            tradeOrder.Ticker + ".US",
		OrderType:         trade.OrderTypeLO,
		SubmittedPrice:    price,
		Side:              side,
		SubmittedQuantity: uint64(contracts),
		OutsideRTH:        trade.OutsideRTHAny,
		TimeInForce:       trade.TimeTypeDay,
	}

	_, err = lb.tc.SubmitOrder(ctx, submitOrder)
	return err
}

func (lb *LongBridge) judgeNoop(ctx context.Context, tradeOrder entity.Trade) (bool, error) {
	order := tradeOrder.Strategy.Order
	action := order.Action
	marketPosition := tradeOrder.Strategy.MarketPosition

	if marketPosition == "" || action == "" {
		return true, nil
	}

	if marketPosition == "flat" {
		stockPositionChannels, err := lb.tc.StockPositions(ctx, []string{tradeOrder.Ticker + ".US"})
		if err != nil {
			return true, err
		}

		if len(stockPositionChannels) == 0 {
			return true, errors.New(fmt.Sprintf("查询股票持仓 [%s] 信息不匹配", tradeOrder.Ticker+".US"))
		}

		stockInfo := stockPositionChannels[0]
		if len(stockInfo.Positions) == 0 {
			return true, errors.New(fmt.Sprintf("查询股票持仓 [%s] 信息不匹配", tradeOrder.Ticker+".US"))
		}

		quantity, err := strconv.Atoi(stockInfo.Positions[0].Quantity)
		if err != nil {
			return true, err
		}

		if action == "buy" && quantity >= 0 {
			return true, nil
		}
		if action == "sell" && quantity <= 0 {
			return true, nil
		}
	}

	return false, nil
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

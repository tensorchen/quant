package longbridge

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/tensorchen/quant/entity"
	"github.com/tensorchen/quant/logger"
	"github.com/tensorchen/quant/stock"

	"github.com/longportapp/openapi-go/config"
	"github.com/longportapp/openapi-go/trade"
	"github.com/shopspring/decimal"
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
		stockPositionChannels, err := lb.tc.StockPositions(ctx, []string{})
		if err != nil {
			return true, err
		}

		ticker := tradeOrder.Ticker + ".US"

		var quantity int
		for _, spc := range stockPositionChannels {
			for _, sp := range spc.Positions {
				if sp.Symbol == ticker {
					quantity, err = strconv.Atoi(sp.Quantity)
					if err != nil {
						return true, err
					}
					logger.Logger().Infof("查询持仓信息 [%s] 持有 [%d] [%+v]", ticker, quantity, sp)
				}
			}
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

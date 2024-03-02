package stock

import (
	"context"

	"github.com/tensorchen/quant/entity"
)

type Stock interface {
	SubmitOrder(ctx context.Context, tradeOrder entity.Trade) error
}

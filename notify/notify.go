package notify

import (
	"context"

	"github.com/tensorchen/quant/entity"
)

type Notifier interface {
	Notify(ctx context.Context, information entity.Information) error
}

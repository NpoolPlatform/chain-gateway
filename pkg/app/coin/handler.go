package appcoin

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Handler struct {
	ID                       *string
	AppID                    *string
	CoinTypeID               *string
	Name                     *string
	DisplayNames             []string
	Logo                     *string
	ForPay                   *bool
	WithdrawAutoReviewAmount *string
	MarketValue              *string
	SettlePercent            *uint32
	SettleTips               []string
	Setter                   *string
	ProductPage              *string
	DailyRewardAmount        *string
	Disabled                 *bool
	Display                  *bool
	DisplayIndex             *uint32
	MaxAmountPerWithdraw     *string
	Offset                   int32
	Limit                    int32
}

func NewHandler(ctx context.Context, options ...func(context.Context, *Handler) error) (*Handler, error) {
	handler := &Handler{}
	for _, opt := range options {
		if err := opt(ctx, handler); err != nil {
			return nil, err
		}
	}
	return handler, nil
}

func WithID(id *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			return nil
		}
		_, err := uuid.Parse(*id)
		if err != nil {
			return err
		}
		h.ID = id
		return nil
	}
}

func WithAppID(id *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			return nil
		}
		_, err := uuid.Parse(*id)
		if err != nil {
			return err
		}
		h.AppID = id
		return nil
	}
}

func WithCoinTypeID(id *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			return nil
		}
		_, err := uuid.Parse(*id)
		if err != nil {
			return err
		}
		h.CoinTypeID = id
		return nil
	}
}

func WithName(name *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if name == nil {
			return nil
		}
		if *name == "" {
			return fmt.Errorf("invalid coinname")
		}
		h.Name = name
		return nil
	}
}

func WithDisplayNames(names []string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.DisplayNames = names
		return nil
	}
}

func WithLogo(logo *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.Logo = logo
		return nil
	}
}

func WithForPay(forPay *bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.ForPay = forPay
		return nil
	}
}

func WithWithdrawAutoReviewAmount(amount *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if amount == nil {
			return nil
		}
		_, err := decimal.NewFromString(*amount)
		if err != nil {
			return err
		}
		h.WithdrawAutoReviewAmount = amount
		return nil
	}
}

func WithMarketValue(amount *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if amount == nil {
			return nil
		}
		_, err := decimal.NewFromString(*amount)
		if err != nil {
			return err
		}
		h.MarketValue = amount
		return nil
	}
}

func WithSettlePercent(percent *uint32) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if percent == nil {
			return nil
		}
		if *percent == 0 {
			return fmt.Errorf("invalid percent")
		}
		h.SettlePercent = percent
		return nil
	}
}

func WithSettleTips(tips []string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.SettleTips = tips
		return nil
	}
}

func WithSetter(id *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			return nil
		}
		_, err := uuid.Parse(*id)
		if err != nil {
			return err
		}
		h.Setter = id
		return nil
	}
}

func WithProductPage(page *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.ProductPage = page
		return nil
	}
}

func WithDisabled(disabled *bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.Disabled = disabled
		return nil
	}
}

func WithDailyRewardAmount(amount *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if amount == nil {
			return nil
		}
		_, err := decimal.NewFromString(*amount)
		if err != nil {
			return err
		}
		h.DailyRewardAmount = amount
		return nil
	}
}

func WithDisplay(display *bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.Display = display
		return nil
	}
}

func WithDisplayIndex(index *uint32) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.DisplayIndex = index
		return nil
	}
}

func WithMaxAmountPerWithdraw(amount *string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if amount == nil {
			return nil
		}
		_, err := decimal.NewFromString(*amount)
		if err != nil {
			return err
		}
		h.MaxAmountPerWithdraw = amount
		return nil
	}
}

func WithOffset(offset int32) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.Offset = offset
		return nil
	}
}

func WithLimit(limit int32) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.Limit = limit
		return nil
	}
}

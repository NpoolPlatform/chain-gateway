package coinfiat

import (
	"context"
	"fmt"

	constant "github.com/NpoolPlatform/chain-gateway/pkg/const"
	coinmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/coin"
	fiatmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/fiat"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"

	"github.com/google/uuid"
)

type Handler struct {
	ID          *uint32
	CoinTypeID  *string
	CoinTypeIDs []string
	FiatID      *string
	FeedType    *basetypes.CurrencyFeedType
	Offset      int32
	Limit       int32
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

func WithID(id *uint32, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.ID = id
		return nil
	}
}

func WithCoinTypeID(id *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			if must {
				return fmt.Errorf("invalid cointypeid")
			}
			return nil
		}
		_coin, err := coinmwcli.GetCoin(ctx, *id)
		if err != nil {
			return err
		}
		if _coin == nil {
			return fmt.Errorf("invalid coin")
		}
		h.CoinTypeID = id
		return nil
	}
}

func WithCoinTypeIDs(ids []string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		for _, id := range ids {
			if _, err := uuid.Parse(id); err != nil {
				return err
			}
		}
		h.CoinTypeIDs = ids
		return nil
	}
}

func WithFiatID(id *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			if must {
				return fmt.Errorf("invalid fiatid")
			}
			return nil
		}
		_fiat, err := fiatmwcli.GetFiat(ctx, *id)
		if err != nil {
			return err
		}
		if _fiat == nil {
			return fmt.Errorf("invalid fiat")
		}
		h.FiatID = id
		return nil
	}
}

func WithFeedType(feedType *basetypes.CurrencyFeedType, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if feedType == nil {
			if must {
				return fmt.Errorf("invalid feedtype")
			}
			return nil
		}
		switch *feedType {
		case basetypes.CurrencyFeedType_CoinGecko:
		case basetypes.CurrencyFeedType_CoinBase:
		case basetypes.CurrencyFeedType_StableUSDHardCode:
		default:
			return fmt.Errorf("invalid feedtype")
		}
		h.FeedType = feedType
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
		if limit == 0 {
			limit = constant.DefaultRowLimit
		}
		h.Limit = limit
		return nil
	}
}

package description

import (
	"context"
	"fmt"

	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	constant "github.com/NpoolPlatform/chain-gateway/pkg/const"
	coinmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/coin"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"

	"github.com/google/uuid"
)

type Handler struct {
	ID         *uint32
	EntID      *string
	AppID      *string
	CoinTypeID *string
	UsedFor    *basetypes.UsedFor
	Title      *string
	Message    *string
	Offset     int32
	Limit      int32
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
		if id == nil {
			if must {
				return fmt.Errorf("invalid id")
			}
			return nil
		}
		h.ID = id
		return nil
	}
}

func WithEntID(id *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			if must {
				return fmt.Errorf("invalid entid")
			}
			return nil
		}
		_, err := uuid.Parse(*id)
		if err != nil {
			return err
		}
		h.EntID = id
		return nil
	}
}

func WithAppID(id *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			if must {
				return fmt.Errorf("invalid appid")
			}
			return nil
		}
		_app, err := appmwcli.GetApp(ctx, *id)
		if err != nil {
			return err
		}
		if _app == nil {
			return fmt.Errorf("invalid app")
		}
		h.AppID = id
		return nil
	}
}

func WithCoinTypeID(id *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			if must {
				return fmt.Errorf("invalid coin")
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

func WithUsedFor(usedFor *basetypes.UsedFor, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if usedFor == nil {
			if must {
				return fmt.Errorf("invalid usedfor")
			}
			return nil
		}
		switch *usedFor {
		case basetypes.UsedFor_ProductPage:
		default:
			return fmt.Errorf("invalid usedfor")
		}
		h.UsedFor = usedFor
		return nil
	}
}

func WithTitle(title *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if title == nil {
			if must {
				return fmt.Errorf("invalid title")
			}
			return nil
		}
		if *title == "" {
			return fmt.Errorf("invalid title")
		}
		h.Title = title
		return nil
	}
}

func WithMessage(message *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if message == nil {
			if must {
				return fmt.Errorf("invalid message")
			}
			return nil
		}
		if *message == "" {
			return fmt.Errorf("invalid message")
		}
		h.Message = message
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

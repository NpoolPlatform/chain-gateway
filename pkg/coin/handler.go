package coin

import (
	"context"
	"fmt"

	constant "github.com/NpoolPlatform/chain-gateway/pkg/const"
	coinmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/coin"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Handler struct {
	ID                          *uint32
	EntID                       *string
	Name                        *string
	Logo                        *string
	Presale                     *bool
	Unit                        *string
	ENV                         *string
	ReservedAmount              *string
	ForPay                      *bool
	HomePage                    *string
	Specs                       *string
	FeeCoinTypeID               *string
	WithdrawFeeByStableUSD      *bool
	WithdrawFeeAmount           *string
	CollectFeeAmount            *string
	HotWalletFeeAmount          *string
	LowFeeAmount                *string
	HotLowFeeAmount             *string
	HotWalletAccountAmount      *string
	PaymentAccountCollectAmount *string
	Disabled                    *bool
	StableUSD                   *bool
	LeastTransferAmount         *string
	NeedMemo                    *bool
	RefreshCurrency             *bool
	CheckNewAddressBalance      *bool
	Offset                      int32
	Limit                       int32
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

func WithName(name *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if name == nil {
			if must {
				return fmt.Errorf("invalid coinname")
			}
			return nil
		}
		if *name == "" {
			return fmt.Errorf("invalid coinname")
		}
		h.Name = name
		return nil
	}
}

func WithLogo(logo *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.Logo = logo
		return nil
	}
}

func WithPresale(presale *bool, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.Presale = presale
		return nil
	}
}

func WithUnit(unit *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if unit == nil {
			if must {
				return fmt.Errorf("invalid coinunit")
			}
			return nil
		}
		if *unit == "" {
			return fmt.Errorf("invalid coinunit")
		}
		h.Unit = unit
		return nil
	}
}

func WithENV(env *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if env == nil {
			if must {
				return fmt.Errorf("invalid coinenv")
			}
			return nil
		}
		switch *env {
		case "main":
		case "test":
		case "local":
		default:
			return fmt.Errorf("invalid coinenv")
		}
		h.ENV = env
		return nil
	}
}

func WithReservedAmount(amount *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if amount == nil {
			if must {
				return fmt.Errorf("invalid reservedamount")
			}
			return nil
		}
		_, err := decimal.NewFromString(*amount)
		if err != nil {
			return err
		}
		h.ReservedAmount = amount
		return nil
	}
}

func WithForPay(forPay *bool, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.ForPay = forPay
		return nil
	}
}

func WithHomePage(homePage *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.HomePage = homePage
		return nil
	}
}

func WithSpecs(specs *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.Specs = specs
		return nil
	}
}

func WithFeeCoinTypeID(id *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if id == nil {
			if must {
				return fmt.Errorf("invalid feecointypeid")
			}
			return nil
		}
		_coin, err := coinmwcli.GetCoin(ctx, *id)
		if err != nil {
			return err
		}
		if _coin == nil {
			return fmt.Errorf("invalid feecoin")
		}
		h.FeeCoinTypeID = id
		return nil
	}
}

func WithWithdrawFeeByStableUSD(stable *bool, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.WithdrawFeeByStableUSD = stable
		return nil
	}
}

func WithWithdrawFeeAmount(amount *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if amount == nil {
			if must {
				return fmt.Errorf("invalid withdrawfeeamount")
			}
			return nil
		}
		_, err := decimal.NewFromString(*amount)
		if err != nil {
			return err
		}
		h.WithdrawFeeAmount = amount
		return nil
	}
}

func WithCollectFeeAmount(amount *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if amount == nil {
			if must {
				return fmt.Errorf("invalid collectfeeamount")
			}
			return nil
		}
		_, err := decimal.NewFromString(*amount)
		if err != nil {
			return err
		}
		h.CollectFeeAmount = amount
		return nil
	}
}

func WithHotWalletFeeAmount(amount *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if amount == nil {
			if must {
				return fmt.Errorf("invalid hotwalletfeeamount")
			}
			return nil
		}
		_, err := decimal.NewFromString(*amount)
		if err != nil {
			return err
		}
		h.HotWalletFeeAmount = amount
		return nil
	}
}

func WithLowFeeAmount(amount *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if amount == nil {
			if must {
				return fmt.Errorf("invalid lowfeeamount")
			}
			return nil
		}
		_, err := decimal.NewFromString(*amount)
		if err != nil {
			return err
		}
		h.LowFeeAmount = amount
		return nil
	}
}

func WithHotLowFeeAmount(amount *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if amount == nil {
			if must {
				return fmt.Errorf("invalid hotlowfeeamount")
			}
			return nil
		}
		_, err := decimal.NewFromString(*amount)
		if err != nil {
			return err
		}
		h.HotLowFeeAmount = amount
		return nil
	}
}

func WithHotWalletAccountAmount(amount *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if amount == nil {
			if must {
				return fmt.Errorf("invalid hotwalletaccountamount")
			}
			return nil
		}
		_, err := decimal.NewFromString(*amount)
		if err != nil {
			return err
		}
		h.HotWalletAccountAmount = amount
		return nil
	}
}

func WithPaymentAccountCollectAmount(amount *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if amount == nil {
			if must {
				return fmt.Errorf("invalid paymentaccountcollectamount")
			}
			return nil
		}
		_, err := decimal.NewFromString(*amount)
		if err != nil {
			return err
		}
		h.PaymentAccountCollectAmount = amount
		return nil
	}
}

func WithDisabled(disabled *bool, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.Disabled = disabled
		return nil
	}
}

func WithStableUSD(stable *bool, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.StableUSD = stable
		return nil
	}
}

func WithLeastTransferAmount(amount *string, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if amount == nil {
			if must {
				return fmt.Errorf("invalid leasttransferamount")
			}
			return nil
		}
		_, err := decimal.NewFromString(*amount)
		if err != nil {
			return err
		}
		h.LeastTransferAmount = amount
		return nil
	}
}

func WithNeedMemo(needMemo *bool, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.NeedMemo = needMemo
		return nil
	}
}

func WithRefreshCurrency(refresh *bool, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.RefreshCurrency = refresh
		return nil
	}
}

func WithCheckNewAddressBalance(check *bool, must bool) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.CheckNewAddressBalance = check
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

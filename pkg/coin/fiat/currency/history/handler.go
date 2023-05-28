package currencyhistory

import (
	"context"

	"github.com/google/uuid"
)

type Handler struct {
	CoinTypeIDs []string
	StartAt     *uint32
	EndAt       *uint32
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

func WithCoinTypeIDs(ids []string) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		if len(ids) == 0 {
			return nil
		}
		for _, id := range ids {
			if _, err := uuid.Parse(id); err != nil {
				return err
			}
		}
		h.CoinTypeIDs = ids
		return nil
	}
}

func WithStartAt(startAt *uint32) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.StartAt = startAt
		return nil
	}
}

func WithEndAt(endAt *uint32) func(context.Context, *Handler) error {
	return func(ctx context.Context, h *Handler) error {
		h.EndAt = endAt
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
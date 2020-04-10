package service

import (
	"context"
	"errors"
)

var (
	ErrBillNotFound    = errors.New("bill not found")
	ErrCmdRepository   = errors.New("unable to command repository")
	ErrQueryRepository = errors.New("unable to query repository")
)

// Service describes the Order service.
type Service interface {
	Create(ctx context.Context, bill Bill) (string, error)
	GetByID(ctx context.Context, id string) (Bill, error)
}

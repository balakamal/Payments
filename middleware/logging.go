package middleware

import (
	"context"
	"github.com/go-kit/kit/log/level"
	"time"

	"github.com/go-kit/kit/log"

	"kkagitala/go-rest-api/service"
)

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next service.Service) service.Service {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   service.Service
	logger log.Logger
}

func (mw loggingMiddleware) Create(ctx context.Context, order service.Bill) (id string, err error) {
	defer func(begin time.Time) {
		level.Debug(mw.logger).Log("method", "Create", "CustomerID", order.CampaignID, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.Create(ctx, order)
}

func (mw loggingMiddleware) GetByID(ctx context.Context, id string) (order service.Bill, err error) {
	defer func(begin time.Time) {
		level.Debug(mw.logger).Log("method", "GetByID", "OrderID", id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetByID(ctx, id)
}

package implementation

import (
	"context"
	"database/sql"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/gofrs/uuid"

	billsvc "kkagitala/go-rest-api/service"
)

// service implements the Order Service
type service struct {
	repository billsvc.Repository
	logger     log.Logger
	seq        int32
}

// NewService creates and returns a new Order service instance
func NewService(rep billsvc.Repository, logger log.Logger) billsvc.Service {
	return &service{
		repository: rep,
		logger:     logger,
		seq:        int32(1),
	}
}

// Create makes an order
func (s *service) Create(ctx context.Context, bill billsvc.Bill) (string, error) {
	logger := log.With(s.logger, "method", "Create")
	uuid, _ := uuid.NewV4()
	id := uuid.String()
	bill.BillID = s.seq
	s.seq++
	level.Info(logger).Log("Call coming from client")
	if err := s.repository.CreateBill(ctx, bill); err != nil {
		level.Error(logger).Log("err", err)
		return "", billsvc.ErrCmdRepository
	}
	return id, nil
}

// GetByID returns an order given by id
func (s *service) GetByID(ctx context.Context, id string) (billsvc.Bill, error) {
	logger := log.With(s.logger, "method", "GetByID")
	order, err := s.repository.GetBillByID(ctx, id)
	if err != nil {
		level.Error(logger).Log("err", err)
		if err == sql.ErrNoRows {
			return order, billsvc.ErrBillNotFound
		}
		return order, billsvc.ErrQueryRepository
	}
	return order, nil
}

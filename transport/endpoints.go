package transport

import (
	"context"
	"kkagitala/go-rest-api/transport/pb"

	"github.com/go-kit/kit/endpoint"

	"kkagitala/go-rest-api/service"
)

// Endpoints holds all Go kit endpoints for the Order service.
type Endpoints struct {
	Create       endpoint.Endpoint
	GetByID      endpoint.Endpoint
	ChangeStatus endpoint.Endpoint
}

// MakeEndpoints initializes all Go kit endpoints for the Order service.
func MakeEndpoints(s service.Service) Endpoints {
	return Endpoints{
		Create:  makeCreateEndpoint(s),
		GetByID: makeGetByIDEndpoint(s),
	}
}

func makeCreateEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*pb.CreateRequest) // type assertion
		s.Create(ctx, service.Bill{
			BillID:      req.BillId,
			UserID:      req.UserId,
			CampaignID:  req.CampaignId,
			PledgeID:    req.PledgeId,
			VatChargeID: 0,
			AmountCents: req.AmountCents,
			Status:      req.Status,
		})
		return pb.CreateResponse{Id: 200}, nil
	}
}

func makeGetByIDEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetByIDRequest)
		orderRes, err := s.GetByID(ctx, req.ID)
		return GetByIDResponse{Bill: orderRes, Err: err}, nil
	}
}

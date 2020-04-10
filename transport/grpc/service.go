package grpc

import (
	"context"
	"github.com/go-kit/kit/log"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	oldcontext "golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"kkagitala/go-rest-api/transport"
	"kkagitala/go-rest-api/transport/pb"
)

// grpc transport service for Account service.
type grpcServer struct {
	createBill kitgrpc.Handler
	logger     log.Logger
}

// NewGRPCServer returns a new gRPC service for the provided Go kit endpoints
func NewGRPCServer(
	endpoints transport.Endpoints, options []kitgrpc.ServerOption,
	logger log.Logger,
) pb.SubscriptionServer {
	errorLogger := kitgrpc.ServerErrorLogger(logger)
	options = append(options, errorLogger)

	return &grpcServer{
		createBill: kitgrpc.NewServer(
			endpoints.Create, decodeCreateCustomerRequest, encodeCreateCustomerResponse, options...,
		),
		logger: logger,
	}
}

// Generate glues the gRPC method to the Go kit service method
func (s *grpcServer) Create(ctx oldcontext.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	_, rep, err := s.createBill.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.CreateResponse), nil
}

// decodeCreateCustomerRequest decodes the incoming grpc payload to our go kit payload
func decodeCreateCustomerRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.CreateRequest) // Type assertion
	return req, nil
}

// encodeCreateCustomerResponse encodes the outgoing go kit payload to the grpc payload
func encodeCreateCustomerResponse(_ context.Context, response interface{}) (interface{}, error) {
	//res := response.(pb.CreateResponse)
	err := getError(nil)
	if err == nil {
		return &pb.CreateResponse{
			Id: 1000,
		}, nil
	}
	return nil, err
}

func getError(err error) error {
	switch err {
	case nil:
		return nil
	default:
		return status.Error(codes.Unknown, err.Error())
	}
}

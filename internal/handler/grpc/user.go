package grpc

import (
	"context"

	pb "github.com/juanpblasi/go-template/api/proto/v1"
	"github.com/juanpblasi/go-template/internal/service"
	"github.com/juanpblasi/go-template/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type userGrpcHandler struct {
	pb.UnimplementedUserServiceServer
	svc service.UserService
}

func NewUserGrpcHandler(svc service.UserService) pb.UserServiceServer {
	return &userGrpcHandler{svc: svc}
}

func mapErrorToGrpc(err error) error {
	if errors.IsType(err, errors.ErrNotFound) {
		return status.Error(codes.NotFound, err.Error())
	}
	if errors.IsType(err, errors.ErrInvalidRequest) {
		return status.Error(codes.InvalidArgument, err.Error())
	}
	return status.Error(codes.Internal, "Internal Server Error")
}

func (h *userGrpcHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err := h.svc.GetUser(ctx, req.GetId())
	if err != nil {
		return nil, mapErrorToGrpc(err)
	}

	return &pb.GetUserResponse{
		User: &pb.User{
			Id:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
	}, nil
}

func (h *userGrpcHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	user, err := h.svc.CreateUser(ctx, req.GetName(), req.GetEmail())
	if err != nil {
		return nil, mapErrorToGrpc(err)
	}

	return &pb.CreateUserResponse{
		User: &pb.User{
			Id:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
	}, nil
}

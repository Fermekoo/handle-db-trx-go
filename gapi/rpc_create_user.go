package gapi

import (
	"context"

	db "github.com/Fermekoo/handle-db-tx-go/db/sqlc"
	"github.com/Fermekoo/handle-db-tx-go/pb"
	"github.com/Fermekoo/handle-db-tx-go/utils"
	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) Register(ctx context.Context, request *pb.RegisterRequest) (*pb.UserResponse, error) {
	hashed_password, err := utils.HashPassword(request.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password %s", err)
	}

	arg := db.CreateUserParams{
		Fullname: request.GetFullname(),
		Username: request.GetFullname(),
		Email:    request.GetEmail(),
		Password: hashed_password,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "username already exists %s", err)
			}
		}
	}

	result := &pb.UserResponse{
		User: convertUser(user),
	}

	return result, nil
}

package gapi

import (
	"context"

	db "github.com/Fermekoo/handle-db-tx-go/db/sqlc"
	"github.com/Fermekoo/handle-db-tx-go/pb"
	"github.com/Fermekoo/handle-db-tx-go/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := server.store.GetUserByEmail(ctx, request.GetEmail())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "user not found")
	}

	err = utils.CheckPassword(request.GetPassword(), user.Password)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid email or password")
	}

	token, acces_payload, err := server.tokenMaker.CreateToken(user.ID, server.config.AccessTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create token %s", err)
	}

	refresh_token, refresh_payload, err := server.tokenMaker.CreateToken(user.ID, server.config.RefreshTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create refresh token %s", err)
	}

	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refresh_payload.ID,
		UserID:       refresh_payload.UserID,
		RefreshToken: refresh_token,
		UserAgent:    "",
		ClientIp:     "",
		IsBlocked:    false,
		ExpiresAt:    refresh_payload.ExpiredAt,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create session %s", err)
	}

	result := &pb.LoginResponse{
		User:                  convertUser(user),
		SessionId:             session.ID.String(),
		AccessToken:           token,
		RefreshToken:          refresh_token,
		AccessTokenExpiresAt:  timestamppb.New(acces_payload.ExpiredAt),
		RefreshTokenExpiresAt: timestamppb.New(refresh_payload.ExpiredAt),
	}

	return result, nil
}

package services

import (
	"context"
	"fmt"

	"github.com/sonnyvictok/miniapp_taptoearn/internal/pb"
	worker "github.com/sonnyvictok/miniapp_taptoearn/internal/workers"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ServerGapi) ClickToEarn(ctx context.Context, req *pb.ClickToEarnRequest) (*pb.ClickToEarnResponse, error) {
	authPayload, err := s.authorizeUser(ctx)
	if err != nil {
		return nil, err
	}
	userTelegramId, err := s.repo.FindByTelegramID(authPayload.ID)
	if err != nil {
		return nil, err
	}
	if authPayload.ID != userTelegramId.TelegramID {
		return nil, status.Errorf(codes.Unauthenticated, "user not found")
	}

	newScore := userTelegramId.Score + 1
	// err = s.repo.UpdateScore(userTelegramId.TelegramID, newScore)
	// if err != nil {
	// 	return nil, err
	// }

	err = s.taskDistributor.DistributeTaskHandleTapToEarn(context.Background(), &worker.PayloadHandleTapToEarn{
		TelegramID: userTelegramId.TelegramID,
		Score:      newScore,
	})
	if err != nil {
		return nil, err
	}

	rsp := &pb.ClickToEarnResponse{
		UserId:   int32(authPayload.ID),
		UserName: authPayload.Username,
		Score:    int32(newScore),
	}
	return rsp, nil
}

func (s *ServerGapi) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	authPayload, err := s.authorizeUser(ctx)
	fmt.Println("authPayload:", authPayload)
	if err != nil {
		return nil, err
	}

	userTelegramId, err := s.repo.FindByTelegramID(authPayload.ID)
	fmt.Println("userTelegramId:", userTelegramId)
	if err != nil {
		return nil, err
	}
	if authPayload.ID != userTelegramId.TelegramID {
		return nil, status.Errorf(codes.Unauthenticated, "user not found")
	}

	rsp := &pb.GetUserResponse{
		UserId:   int32(authPayload.ID),
		UserName: authPayload.Username,
		Score:    int32(userTelegramId.Score),
	}
	return rsp, nil
}

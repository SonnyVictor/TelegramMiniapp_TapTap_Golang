package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/sonnyvictok/miniapp_taptoearn/internal/domain"
	initdata "github.com/telegram-mini-apps/init-data-golang"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	authorizationHeader = "authorization"
	authorizationTMA    = "tma"
)

func (s *ServerGapi) authorizeUser(ctx context.Context) (*domain.UserTelegramPayload, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "authorization header is required")
	}

	headerTMA := md.Get(authorizationTMA)
	if len(headerTMA) == 0 {
		return nil, status.Error(codes.Unauthenticated, "Get Header TMA Error")
	}

	authHeader := headerTMA[0]
	fields := strings.Fields(authHeader)

	initDataTlg := fields[0]

	initData, err := initdata.Parse(initDataTlg)
	if err != nil {
		return nil, fmt.Errorf("invalid init data: %v", err)
	}

	userData := initData.User
	payload := &domain.UserTelegramPayload{
		ID:        userData.ID,
		Username:  userData.Username,
		FirstName: userData.FirstName,
		LastName:  userData.LastName,
	}
	return payload, nil
}

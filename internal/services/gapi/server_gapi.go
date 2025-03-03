package services

import (
	"github.com/sonnyvictok/miniapp_taptoearn/internal/pb"
	"github.com/sonnyvictok/miniapp_taptoearn/internal/repositories"
	worker "github.com/sonnyvictok/miniapp_taptoearn/internal/workers"
)

type ServerGapi struct {
	pb.UnimplementedTapServiceServer
	repo            *repositories.UserRepository
	taskDistributor worker.TaskDistributor
}

func NewServerGapi(repo *repositories.UserRepository, taskDistributor worker.TaskDistributor) (*ServerGapi, error) {
	return &ServerGapi{
		repo:            repo,
		taskDistributor: taskDistributor,
	}, nil
}

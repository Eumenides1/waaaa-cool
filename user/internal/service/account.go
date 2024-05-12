package service

import (
	"common/logs"
	"context"
	"core/repo"
	"user/pb"
)

// AccountService 账号相关服务
type AccountService struct {
	pb.UnimplementedUserServiceServer
}

func NewAccountService(manager *repo.Manager) *AccountService {
	return &AccountService{}
}

func (a *AccountService) Register(ctx context.Context, params *pb.RegisterParams) (*pb.RegisterResponse, error) {
	// 注册的业务逻辑
	logs.Info("register server called.....")
	return &pb.RegisterResponse{
		Uid: "10000",
	}, nil
}

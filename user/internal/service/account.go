package service

import (
	"common/biz"
	"common/logs"
	"context"
	"core/dao"
	"core/models/entity"
	"core/models/requests"
	"core/repo"
	"framework/waError"
	"time"
	"user/pb"
)

// AccountService 账号相关服务
type AccountService struct {
	accountDao *dao.AccountDao
	redisDao   *dao.RedisDao
	pb.UnimplementedUserServiceServer
}

func NewAccountService(manager *repo.Manager) *AccountService {
	return &AccountService{
		accountDao: dao.NewAccountDao(manager),
		redisDao:   dao.NewRedisDao(manager),
	}
}

func (a *AccountService) Register(ctx context.Context, req *pb.RegisterParams) (*pb.RegisterResponse, error) {
	// 注册的业务逻辑
	// 1. 封装一个Account 结构 (mongo -> 分布式ID)
	logs.Info("register server called.....")

	if req.LoginPlatform == requests.WeiXin {
		ac, err := a.wxRegister(req)
		if err != nil {
			return &pb.RegisterResponse{}, waError.GrpcError(err)
		}
		return &pb.RegisterResponse{
			Uid: ac.Uid,
		}, nil
	}
	return &pb.RegisterResponse{}, nil

}

func (a *AccountService) wxRegister(req *pb.RegisterParams) (*entity.Account, *waError.Error) {
	ac := &entity.Account{
		WxAccount:  req.Account,
		CreateTime: time.Now(),
	}
	// 3. 生成唯一识别ID （Redis 自增）
	uid, err := a.redisDao.NextAccountId()
	if err != nil {
		return ac, biz.SqlError
	}
	ac.Uid = uid
	err = a.accountDao.SaveAccount(context.TODO(), ac)
	if err != nil {
		return ac, biz.SqlError
	}
	return ac, nil
}

package service

import (
	pb "banana/api/account-center/service/v1"
	"banana/app/account-center/service/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewAccountCenterService)

type AccountCenterService struct {
	pb.UnimplementedAccountCenterServer

	ac *biz.AccountCenterCase

	log *log.Helper
}
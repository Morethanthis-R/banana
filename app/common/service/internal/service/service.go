package service

import (
	pb "banana/api/common/service/v1"
	"banana/app/common/service/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewCommonService)

type CommonService struct {
	pb.UnimplementedCommonServer

	cm *biz.CommonCase

	log *log.Helper
}


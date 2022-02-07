package service

import (
	pb "banana/api/transfer/service/v1"
	"banana/app/transfer/service/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)


var ProviderSet = wire.NewSet(NewTransferService)

type TransferService struct {
	pb.UnimplementedTransferServer

	tf *biz.TransferCase

	log *log.Helper
}
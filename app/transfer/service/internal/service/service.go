package service

import (
	pb "banana/api/transfer/service/v1"
	"banana/app/transfer/service/internal/biz"
	"banana/app/transfer/service/internal/data"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)


var ProviderSet = wire.NewSet(NewTransferService,NewMqService)

type TransferService struct {
	pb.UnimplementedTransferServer

	tf *biz.TransferCase

	log *log.Helper
}

type MqService struct {
	mq *data.RabbitMQ
	db *data.Data
}
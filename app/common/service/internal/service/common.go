package service

import (
	"banana/app/common/service/internal/biz"
	"context"
	"github.com/go-kratos/kratos/v2/log"

	pb "banana/api/common/service/v1"
)

func NewCommonService(cm *biz.CommonCase, logger log.Logger) *CommonService {
	return &CommonService{
		cm:cm,
		log:log.NewHelper(logger),
	}
}

func (s *CommonService) CreateNotify(ctx context.Context, req *pb.ReqCreateNotify) (*pb.RespCreateNotify, error) {
	return s.cm.CreateNotify(ctx,req)
}
func (s *CommonService) DeleteNotify(ctx context.Context, req *pb.ReqDeleteNotify) (*pb.RespDeleteNotify, error) {
	return s.cm.DeleteNotify(ctx,req)
}
func (s *CommonService) GetNotifyList(ctx context.Context, req *pb.ReqGetNotifyList) (*pb.RespGetNotifyList, error) {
	return s.cm.GetNotifyList(ctx,req)
}
func (s *CommonService) GetNotifyObject(ctx context.Context, req *pb.ReqGetNotifyObject) (*pb.RespGetNotifyObject, error) {
	return s.cm.GetNotify(ctx,req)
}
func (s *CommonService) CreateNotifyType(ctx context.Context, req *pb.ReqCreateNotifyType) (*pb.RespCreateNotifyType, error) {
	return s.cm.CreateNType(ctx,req)
}
func (s *CommonService) UpdateNotifyType(ctx context.Context, req *pb.ReqUpdateNotifyType) (*pb.RespUpdateNotifyType, error) {
	return s.cm.UpdateNType(ctx,req)
}
func (s *CommonService) DeleteNotifyType(ctx context.Context, req *pb.ReqDeleteNotifyType) (*pb.RespDeleteNotifyType, error) {
	return s.cm.DeleteNType(ctx,req)
}
func (s *CommonService) GetNotifyTypeList(ctx context.Context, req *pb.ReqGetNotifyTypeList) (*pb.RespGetNotifyTypeList, error) {
	return s.cm.GetNTypeList(ctx,req)
}
func (s *CommonService) CreateAdv(ctx context.Context,req *pb.ReqCreateAdv)(*pb.RespCreateAdv,error){
	return s.cm.CreateAdv(ctx,req)
}

func (s *CommonService) DeleteAdv(ctx context.Context,req *pb.ReqDeleteAdv)(*pb.RespDeleteAdv,error){
	return s.cm.DeleteAdv(ctx,req)
}
func (s *CommonService) UpdateAdv(ctx context.Context,req *pb.ReqUpdateAdv)(*pb.RespUpdateAdv,error){
	return s.cm.UpdateAdv(ctx,req)
}

func (s *CommonService) GetAdvList(ctx context.Context,req *pb.ReqGetAdvList)(*pb.RespGetAdvList,error){
	return s.cm.GetAdvList(ctx,req)
}
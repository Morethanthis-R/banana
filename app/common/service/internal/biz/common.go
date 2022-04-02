package biz

import (
	pb "banana/api/common/service/v1"
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type TNotify struct {
	ID           int    `gorm:"primary_key" json:"id"`
	Name         string `gorm:"type:varchar(20) not null" json:"name"`
	Body         string `gorm:"type:text not null" json:"body""`
	IsGlobal     int8   `gorm:"default:0" json:"is_global"`
	Type         string `gorm:"type:varchar(256) not null" json:"type"`
	SendTime     int64  `json:"send_time"`
	Status       int8   `gorm:"type:tinyint(5);default:1" json:"status"`
	CreatedAt    int64  `gorm:"autoCreateAt" json:"created_at"`
	UpdatedAt    int64  `gorm:"autoUpdateAt" json:"updated_at"`
}
type TUserNotify struct {
	ID           int    `gorm:"primary_key" json:"id"`
	Nid         int     `gorm:"type:tinyint(5)" json:"nid"`
	Uid         int     `gorm:"type:tinyint(5)" json:"uid"`
	CreatedAt    int64  `gorm:"autoCreateAt" json:"created_at"`
	UpdatedAt    int64  `gorm:"autoUpdateAt" json:"updated_at"`
}

type TType struct {
	ID           int    `gorm:"primary_key" json:"id"`
	Name         string `gorm:"type:varchar(20) not null" json:"name"`
	Describe     string `gorm:"type:text not null" json:"describe""`
	CreatedAt    int64  `gorm:"autoCreateAt" json:"created_at"`
	UpdatedAt    int64  `gorm:"autoUpdateAt" json:"updated_at"`
}

type TAdv struct {
	ID           int    `gorm:"primary_key" json:"id"`
	Name         string `gorm:"type:varchar(20) not null" json:"name"`
	Describe     int8   `gorm:"type:tinyint(5) not null" json:"describe""`
	LinkUrl      string `gorm:"type:text not null" json:"link_url"`
	ImageUrl     string `gorm:"type:text not null" json:"image_url"`
	Status       int8   `gorm:"type:tinyint(5) not null" json:"status"`
	CreatedAt    int64  `gorm:"autoCreateAt" json:"created_at"`
	UpdatedAt    int64  `gorm:"autoUpdateAt" json:"updated_at"`
}
type CommonRepo interface {
	CreateNotify(ctx context.Context, rq *pb.ReqCreateNotify) (*pb.RespCreateNotify,error)
	DeleteNotify(ctx context.Context, rq *pb.ReqDeleteNotify) (*pb.RespDeleteNotify,error)
	GetNotifyList(ctx context.Context,rq *pb.ReqGetNotifyList) (*pb.RespGetNotifyList,error)
	GetNotify(ctx context.Context,rq *pb.ReqGetNotifyObject) (*pb.RespGetNotifyObject,error)
	CreateNType(ctx context.Context,rq *pb.ReqCreateNotifyType) (*pb.RespCreateNotifyType,error)
	UpdateNType(ctx context.Context,rq *pb.ReqUpdateNotifyType) (*pb.RespUpdateNotifyType,error)
	DeleteNType(ctx context.Context,rq *pb.ReqDeleteNotifyType) (*pb.RespDeleteNotifyType,error)
	GetNTypeList(ctx context.Context,rq *pb.ReqGetNotifyTypeList) (*pb.RespGetNotifyTypeList,error)
	CreateAdv(ctx context.Context,rq *pb.ReqCreateAdv) (*pb.RespCreateAdv,error)
	UpdateAdv(ctx context.Context,rq *pb.ReqUpdateAdv) (*pb.RespUpdateAdv,error)
	GetAdvList(ctx context.Context,rq *pb.ReqGetAdvList) (*pb.RespGetAdvList,error)
	DeleteAdv(ctx context.Context,rq *pb.ReqDeleteAdv) (*pb.RespDeleteAdv,error)
}

type CommonCase struct {
	repo CommonRepo
	log *log.Helper
}

func NewCommonCase(repo CommonRepo,logger log.Logger) *CommonCase{
	return &CommonCase{repo: repo,log: log.NewHelper(log.With(logger,"module","common/case"))}
}

func (cm *CommonCase) CreateNotify(ctx context.Context, rq *pb.ReqCreateNotify) (*pb.RespCreateNotify,error){
	return cm.repo.CreateNotify(ctx,rq)
}

func (cm *CommonCase) DeleteNotify(ctx context.Context, rq *pb.ReqDeleteNotify) (*pb.RespDeleteNotify,error){
	return cm.repo.DeleteNotify(ctx,rq)
}

func (cm *CommonCase) GetNotifyList(ctx context.Context,rq *pb.ReqGetNotifyList) (*pb.RespGetNotifyList,error){
	return cm.repo.GetNotifyList(ctx,rq)
}
func (cm *CommonCase) GetNotify(ctx context.Context,rq *pb.ReqGetNotifyObject) (*pb.RespGetNotifyObject,error){
	return cm.repo.GetNotify(ctx,rq)
}
func (cm *CommonCase) CreateNType(ctx context.Context,rq *pb.ReqCreateNotifyType) (*pb.RespCreateNotifyType,error){
	return cm.repo.CreateNType(ctx,rq)
}
func (cm *CommonCase) UpdateNType(ctx context.Context,rq *pb.ReqUpdateNotifyType) (*pb.RespUpdateNotifyType,error){
	return cm.repo.UpdateNType(ctx,rq)
}
func (cm *CommonCase)  DeleteNType(ctx context.Context,rq *pb.ReqDeleteNotifyType) (*pb.RespDeleteNotifyType,error){
	return cm.repo.DeleteNType(ctx,rq)
}
func (cm *CommonCase) GetNTypeList(ctx context.Context,rq *pb.ReqGetNotifyTypeList) (*pb.RespGetNotifyTypeList,error){
	return cm.repo.GetNTypeList(ctx,rq)
}

func (cm *CommonCase) CreateAdv(ctx context.Context,rq *pb.ReqCreateAdv) (*pb.RespCreateAdv,error){
	return cm.repo.CreateAdv(ctx,rq)
}
func (cm *CommonCase) DeleteAdv(ctx context.Context,rq *pb.ReqDeleteAdv) (*pb.RespDeleteAdv,error){
	return cm.repo.DeleteAdv(ctx,rq)
}
func (cm *CommonCase) UpdateAdv(ctx context.Context,rq *pb.ReqUpdateAdv) (*pb.RespUpdateAdv,error){
	return cm.repo.UpdateAdv(ctx,rq)
}

func (cm *CommonCase) GetAdvList(ctx context.Context,rq *pb.ReqGetAdvList) (*pb.RespGetAdvList,error){
	return cm.repo.GetAdvList(ctx,rq)
}
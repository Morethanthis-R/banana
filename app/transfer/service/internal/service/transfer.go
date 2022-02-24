package service

import (
	"banana/app/transfer/service/internal/biz"
	"context"
	"github.com/go-kratos/kratos/v2/log"

	pb "banana/api/transfer/service/v1"
)

func NewTransferService(tf *biz.TransferCase, logger log.Logger) *TransferService {
	return &TransferService{
		tf:  tf,
		log: log.NewHelper(logger),
	}
}

func (s *TransferService) UploadEntry(ctx context.Context, req *pb.ReqUpload) (*pb.RespUpload, error) {
	res, err := s.tf.UploadEntry(ctx, req)
	return res, err
}
func (s *TransferService) DownLoadEntry(ctx context.Context, req *pb.ReqDownload) (*pb.RespDownload, error) {
	return s.tf.DownloadEntry(ctx, req)
}
func (s *TransferService) UploadStatic(ctx context.Context, req *pb.ReqStatic) (*pb.RespStatic, error) {
	return s.tf.UploadStatic(ctx, req)
}

func (s *TransferService) GetUserFileTree(ctx context.Context, req *pb.ReqGetUserFileTree) (*pb.RespGetUserFileTree, error) {
	return s.tf.GetUserFileTree(ctx, req)
}
func (s *TransferService) GetUserTrashList(ctx context.Context,req *pb.ReqGetUserTrashBin) (*pb.RespGetUserTrashBin,error){
	return s.tf.GetTrashBin(ctx,req)
}
func (s *TransferService) DeleteFile(ctx context.Context, req *pb.ReqDeleteFile) (*pb.RespDelete, error) {
	return s.tf.DeleteFile(ctx, req)
}
func (s *TransferService) DeleteDir(ctx context.Context, req *pb.ReqDeleteDir) (*pb.RespDelete, error) {
	return s.tf.DeleteDir(ctx, req)
}
func (s *TransferService) ShareFile(ctx context.Context, req *pb.ReqShareFileStr) (*pb.RespShareFileStr, error) {
	return s.tf.ShareFile(ctx, req)
}

func (s *TransferService) PreviewFile(ctx context.Context, req *pb.ReqPreviewFile) (*pb.RespPreviewFile, error) {
	return s.tf.PreviewFile(ctx, req)
}

func (s *TransferService) FileCensus(ctx context.Context, req *pb.ReqFileCensus) (*pb.RespFileCensus, error) {
	return s.tf.FileCensus(ctx, req)
}
func (s *TransferService) SearchFile(ctx context.Context, req *pb.ReqSearchFile) (*pb.RespSearchFile, error) {
	return s.tf.SearchFile(ctx, req)
}
func (s *TransferService) CleanTrashFile(ctx context.Context, req *pb.ReqCleanTrashFile) (*pb.RespCleanTrash, error) {
	return s.tf.CleanTrashFile(ctx, req)
}
func (s *TransferService) CleanTrashDir(ctx context.Context, req *pb.ReqCleanTrashDir) (*pb.RespCleanTrash, error) {
	return s.tf.CleanTrashDir(ctx,req)
}
func (s *TransferService) WithDrawFile(ctx context.Context,req *pb.ReqWithDrawFile) (*pb.RespWithDraw,error) {
	return s.tf.WithDrawFile(ctx,req)
}
func (s *TransferService) WithDrawDir(ctx context.Context,req *pb.ReqWithDrawDir) (*pb.RespWithDraw,error){
	return s.tf.WithDrawDir(ctx,req)
}
func (s *TransferService) CreateDir(ctx context.Context,req *pb.ReqCreateDir)(*pb.RespCreateDir,error){
	return s.tf.CreateDir(ctx,req)
}

func (s *TransferService) GuestUpload(ctx context.Context,req *pb.ReqGuestUpload)(*pb.RespGuestUpload,error){
	return s.tf.GuestUpload(ctx,req)
}
func (s *TransferService) GetCodeDownload(ctx context.Context,req *pb.ReqGetCodeDownLoad)(*pb.RespGetCOdeDownload,error){
	return s.tf.GetCodeDownload(ctx,req)
}
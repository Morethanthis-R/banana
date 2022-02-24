package biz

import (
	pb "banana/api/transfer/service/v1"
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type File struct {
	ID            int    `gorm:"primary_key" json:"id"`
	DirectoryId   int    `gorm:"index" json:"directory_id"`
	FileName      string `json:"file_name"`
	FileStatus    int8   `gorm:"type:tinyint(5);default:1" json:"file_status"` //1存在 0已删除(彻底删除) 2删除(垃圾桶)
	FileHash      string `gorm:"type:text" json:"file_hash"`
	FilePath      string `gorm:"type:text" json:"file_path"` //minio用
	FileSize      int64  `json:"file_size"`
	FileStr       string `gorm:"type:text ;default:''" json:"file_str"` //文件目录用
	Attribute     int8   `json:"attribute"`
	ContentType   string `gorm:"type:varchar(256)" json:"content_type"`
	Suffix        string `gorm:"type:varchar(10)" json:"suffix"`
	DownloadCount int    `gorm:"default:0" json:"download_count"`
	CreatedAt     int64  `gorm:"autoCreateAt" json:"created_at"`
	UpdatedAt     int64  `gorm:"autoUpdateAt" json:"updated_at"`
}

type UserFile struct {
	ID          int    `gorm:"primary_key" json:"id"`
	FileId      int    `gorm:"index" json:"file_id"`
	UserId      int    `gorm:"index" json:"user_id"`
	DirectoryId int    `gorm:"index" json:"directory_id"`
	UserNum     string `gorm:"type:varchar(10) not null " json:"user_num"`
	CreatedAt   int64  `gorm:"autoCreateAt" json:"created_at"`
	UpdatedAt   int64  `gorm:"autoUpdateAt" json:"updated_at"`
}

type UserDirectory struct {
	ID        int    `gorm:"primary_key" json:"id"`
	FatherId  int    `gorm:"index" json:"father_id"` //根结点默认0  用户专属根目录为 fatherid==0，userid==globaluid
	UserId    int    `gorm:"index" json:"user_id"`
	Size      int64  `json:"size"`
	PathStr   string `gorm:"type:varchar(128)" json:"path_str"`
	Name      string `gorm:"type:varchar(128)" json:"name"`
	DirStatus int8   `gorm:"type:tinyint(5);default:1" json:"dir_status""` //1存在 0已删除(彻底删除) 2删除(垃圾桶)
	Key       string `gorm:"default:''" json:"key"`
	CreatedAt int64  `gorm:"autoCreateAt" json:"created_at"`
	UpdatedAt int64  `gorm:"autoUpdateAt" json:"updated_at"`
}

type TransferRepo interface {
	UploadEntry(ctx context.Context, u *pb.ReqUpload) (*pb.RespUpload, error)
	DownloadEntry(ctx context.Context, u *pb.ReqDownload) (*pb.RespDownload, error)
	UploadStatic(ctx context.Context, u *pb.ReqStatic) (*pb.RespStatic, error)
	GetUserFileTree(ctx context.Context, u *pb.ReqGetUserFileTree) (*pb.RespGetUserFileTree, error)
	GetTrashBin(ctx context.Context, u *pb.ReqGetUserTrashBin) (*pb.RespGetUserTrashBin, error)
	DeleteFile(ctx context.Context, u *pb.ReqDeleteFile) (*pb.RespDelete, error)
	DeleteDir(ctx context.Context, u *pb.ReqDeleteDir) (*pb.RespDelete, error)
	CleanTrashFile(ctx context.Context, u *pb.ReqCleanTrashFile) (*pb.RespCleanTrash, error)
	CleanTrashDir(ctx context.Context, u *pb.ReqCleanTrashDir) (*pb.RespCleanTrash, error)
	ShareFile(ctx context.Context, u *pb.ReqShareFileStr) (*pb.RespShareFileStr, error)
	PreviewFile(ctx context.Context, u *pb.ReqPreviewFile) (*pb.RespPreviewFile, error)
	FileCensus(ctx context.Context, u *pb.ReqFileCensus) (*pb.RespFileCensus, error)
	SearchFile(ctx context.Context, u *pb.ReqSearchFile) (*pb.RespSearchFile, error)
	WithDrawFile(ctx context.Context,u *pb.ReqWithDrawFile)(*pb.RespWithDraw,error)
	WithDrawDir(ctx context.Context,u *pb.ReqWithDrawDir) (*pb.RespWithDraw,error)
	CreateDir(ctx context.Context,u *pb.ReqCreateDir) (*pb.RespCreateDir,error)
	GuestUpload(ctx context.Context,u *pb.ReqGuestUpload) (*pb.RespGuestUpload,error)
	GetCodeDownload(ctx context.Context,u *pb.ReqGetCodeDownLoad) (*pb.RespGetCOdeDownload,error)
}

type TransferCase struct {
	repo TransferRepo
	log  *log.Helper
}

func NewTransferCase(repo TransferRepo, logger log.Logger) *TransferCase {
	return &TransferCase{repo: repo, log: log.NewHelper(log.With(logger, "module", "transfer/case"))}
}

func (tf *TransferCase) UploadEntry(ctx context.Context, t *pb.ReqUpload) (*pb.RespUpload, error) {
	return tf.repo.UploadEntry(ctx, t)
}

func (tf *TransferCase) DownloadEntry(ctx context.Context, t *pb.ReqDownload) (*pb.RespDownload, error) {
	return tf.repo.DownloadEntry(ctx, t)
}

func (tf *TransferCase) UploadStatic(ctx context.Context, t *pb.ReqStatic) (*pb.RespStatic, error) {
	return tf.repo.UploadStatic(ctx, t)
}

func (tf *TransferCase) GetUserFileTree(ctx context.Context, t *pb.ReqGetUserFileTree) (*pb.RespGetUserFileTree, error) {
	return tf.repo.GetUserFileTree(ctx, t)
}
func (tf *TransferCase) GetTrashBin(ctx context.Context, t *pb.ReqGetUserTrashBin) (*pb.RespGetUserTrashBin, error) {
	return tf.repo.GetTrashBin(ctx,t)
}

func (tf *TransferCase) DeleteFile(ctx context.Context, t *pb.ReqDeleteFile) (*pb.RespDelete, error) {
	return tf.repo.DeleteFile(ctx, t)
}
func (tf *TransferCase) DeleteDir(ctx context.Context, t *pb.ReqDeleteDir) (*pb.RespDelete, error) {
	return tf.repo.DeleteDir(ctx,t)
}

func (tf *TransferCase) ShareFile(ctx context.Context, t *pb.ReqShareFileStr) (*pb.RespShareFileStr, error) {
	return tf.repo.ShareFile(ctx, t)
}

func (tf *TransferCase) PreviewFile(ctx context.Context, t *pb.ReqPreviewFile) (*pb.RespPreviewFile, error) {
	return tf.repo.PreviewFile(ctx, t)
}
func (tf *TransferCase) FileCensus(ctx context.Context, t *pb.ReqFileCensus) (*pb.RespFileCensus, error) {
	return tf.repo.FileCensus(ctx, t)
}
func (tf *TransferCase) SearchFile(ctx context.Context, t *pb.ReqSearchFile) (*pb.RespSearchFile, error) {
	return tf.repo.SearchFile(ctx, t)
}
func (tf *TransferCase) CleanTrashFile(ctx context.Context, t *pb.ReqCleanTrashFile) (*pb.RespCleanTrash, error){
	return tf.repo.CleanTrashFile(ctx,t)
}
func (tf *TransferCase) CleanTrashDir (ctx context.Context, t *pb.ReqCleanTrashDir) (*pb.RespCleanTrash, error){
	return tf.repo.CleanTrashDir(ctx,t)
}

func (tf *TransferCase) WithDrawFile(ctx context.Context,t *pb.ReqWithDrawFile) (*pb.RespWithDraw,error){
	return tf.repo.WithDrawFile(ctx,t)
}

func (tf *TransferCase) WithDrawDir(ctx context.Context,t *pb.ReqWithDrawDir) (*pb.RespWithDraw,error){
	return tf.repo.WithDrawDir(ctx,t)
}

func (tf *TransferCase) CreateDir(ctx context.Context,t *pb.ReqCreateDir) (*pb.RespCreateDir,error){
	return tf.repo.CreateDir(ctx,t)
}

func (tf *TransferCase) GuestUpload(ctx context.Context,t *pb.ReqGuestUpload) (*pb.RespGuestUpload,error){
	return tf.repo.GuestUpload(ctx,t)
}

func (tf *TransferCase) GetCodeDownload(ctx context.Context,t *pb.ReqGetCodeDownLoad) (*pb.RespGetCOdeDownload,error){
	return tf.repo.GetCodeDownload(ctx,t)
}
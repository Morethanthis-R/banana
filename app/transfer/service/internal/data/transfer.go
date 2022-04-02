package data

import (
	pb "banana/api/transfer/service/v1"
	"banana/app/transfer/service/internal/biz"
	"banana/pkg/ecode"
	"banana/pkg/middleware"
	"banana/pkg/util"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/minio/minio-go/v7"
	"io"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"
)

var _ biz.TransferRepo = (*transferRepo)(nil)

type transferRepo struct {
	mq   *RabbitMQ
	data *Data
	log  *log.Helper
}

const (
	PUBLIC  = "peach-public"
	PRIVATE = "peach-private"
	DELETE  = "peach-delete"
	STATIC  = "peach-static"

	EXIST = 1
	DEL   = 2
	RUIN  = 3

	DEFAULT = 0
	ASC     = 1
	DESC    = 2

	NAME = 1
	TIME = 2
	SIZE = 3
)

func NewTransferRepo(data *Data, logger log.Logger, mq *RabbitMQ) biz.TransferRepo {
	return &transferRepo{
		mq : mq,
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "data/tf")),
	}
}
type Produce struct {
	msgContent string
}
// 实现发送者
func (t *Produce) MsgContent() string {
	return t.msgContent
}

func (t *transferRepo) GuestUpload(ctx context.Context, req *pb.ReqGuestUpload) (*pb.RespGuestUpload, error) {

	//id := ctx.Value("x-md-global-uid").(int)
	res := &pb.RespGuestUpload{}
	claims := ctx.Value("claims").(*middleware.Claims)
	var err error
	//fileinfo := minio.UploadInfo{}
	filetype := util.String2StringArrWithSeparate(req.File.Filename, ".", true)
	finalName := fmt.Sprintf("%s.%s", req.File.FileHash, filetype[len(filetype)-1])

	//快传检索
	checkFile := &biz.File{}
	if req.File.FileHash != "" {
		err = t.data.Db.Model(&biz.File{}).Where("file_hash = ? and file_status = ?", req.File.FileHash,EXIST).First(checkFile).Error
		if err != nil && err.Error() != "record not found" {
			t.log.Error(err)
			return nil, ecode.MYSQL_ERR
		}

		if checkFile.ID != 0 {
			userFile := &biz.UserFile{}
			err = t.data.Db.Model(&biz.UserFile{}).Where("file_id = ?", checkFile.ID).Last(userFile).Error
			if err != nil && err.Error() != "record not found" {
				t.log.Error(err)
				return nil, ecode.MYSQL_ERR
			}
			if userFile.UserId == 0{
				res.Fid = int32(checkFile.ID)
				return res,nil
			}
			if userFile.ID != 0 {
				client := t.data.Minio_internal
				bucket := PUBLIC
				if userFile.UserId != 0 {
					bucket = PRIVATE
				}
				copyPath := fmt.Sprintf("/%s", finalName)
				if bucket != PUBLIC {
					src := minio.CopySrcOptions{
						Bucket: bucket,
						Object: checkFile.FilePath,
					}
					dst := minio.CopyDestOptions{
						Bucket: PUBLIC,
						Object: copyPath,
					}

					_, err := client.CopyObject(context.TODO(), dst, src)
					if err != nil {
						return nil, ecode.EXTERNAL_API_FAIL.SetMessage(err.Error())
					}
				} else {
					res.Fid = int32(checkFile.ID)
					return res, nil
				}

				newFile := &biz.File{
					FileName:    req.File.Filename,
					FileStatus:  EXIST,
					FileHash:    req.File.FileHash,
					FilePath:    copyPath,
					FileSize:    checkFile.FileSize,
					FileStr:     fmt.Sprintf("%s", req.File.Filename),
					Attribute:   2,
					ContentType: req.File.ContentType,
					Suffix:      filetype[len(filetype)-1],
				}
				err = t.data.Db.Model(&biz.File{}).Create(newFile).Error
				if err != nil {
					t.log.Error(err)
					return nil, ecode.MYSQL_ERR.SetMessage(err.Error())
				}
				newUFile := &biz.UserFile{
					FileId:  newFile.ID,
					UserNum: claims.UserNum,
				}
				err = t.data.Db.Model(&biz.UserFile{}).Create(newUFile).Error
				if err != nil {
					t.log.Error(err)
					return nil, ecode.MYSQL_ERR.SetMessage(err.Error())
				}
				res.Fid = int32(newFile.ID)
				return res, nil
			}
		}

	}
	//var minioUpload = func(bucket, objectName string) (minio.UploadInfo, error) {
	//	client := t.data.Minio_internal
	//	fileinfo, err = client.FPutObject(ctx, bucket, objectName, req.File.Filename, minio.PutObjectOptions{ContentType: req.File.ContentType})
	//	if err != nil {
	//		t.log.Error(err)
	//		return fileinfo, err
	//	}
	//	return fileinfo, err
	//}
	//fileinfo, err = minioUpload(PUBLIC, finalName)
	//if err != nil {
	//	t.log.Error(err)
	//	return res, ecode.New(500).SetMessage("minio客户端错误")
	//}

	attribute := 2
	file := &biz.File{
		FileName:    req.File.Filename,
		FileHash:    req.File.FileHash,
		FilePath:    finalName,
		FileSize:    0,//fileinfo.Size,
		FileStr:     "",
		Attribute:   int8(attribute),
		Suffix:      filetype[len(filetype)-1],
		ContentType: req.File.ContentType,
	}
	err = t.data.Db.Model(&biz.File{}).WithContext(ctx).Create(file).Error
	if err != nil {
		t.log.Error(err)
		return res, ecode.MYSQL_ERR
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		obj := MessageObject{
			Fid:         file.ID,
			FileName:    file.FileName,
			FileHash:    file.FileHash,
			FileStr:     file.FileStr,
			FilePath:    file.FilePath,
			ContentType: file.ContentType,
			Bucket:      PUBLIC,
		}
		msg ,_:= json.Marshal(obj)
		produce := &Produce{msgContent: string(msg)}
		t.mq.RegisterProducer(produce)
		wg.Done()
	}()
	wg.Wait()


	userFile := &biz.UserFile{
		FileId:  file.ID,
		UserNum: claims.UserNum,
		UserId:  claims.UserId,
	}
	err = t.data.Db.Model(&biz.UserFile{}).WithContext(ctx).Create(userFile).Error
	if err != nil {
		t.log.Error(err)
		return res, ecode.MYSQL_ERR
	}
	res.Fid = int32(file.ID)
	return res, nil
}

func (t *transferRepo) GetCodeDownload(ctx context.Context, req *pb.ReqGetCodeDownLoad) (*pb.RespGetCOdeDownload, error) {
	res := &pb.RespGetCOdeDownload{}
	var err error
	url, err := t.data.cache.Get(context.TODO(), req.GetCode).Result()
	if err != nil {
		return res, ecode.REDIS_ERR
	}
	res.DownloadStr = url
	shareHis := &biz.ShareHistory{}
	err = t.data.Db.Model(&biz.ShareHistory{}).Where("get_code = ?", req.GetCode).First(shareHis).Error
	if err != nil {
		return res, ecode.MYSQL_ERR.SetMessage(err.Error())
	}
	res.Title = shareHis.Title
	res.Describe = shareHis.Describe
	res.CreateTime = shareHis.CreatedAt
	res.ExpireTime = shareHis.ExpireTime
	res.FileName = shareHis.FileName
	res.FileSize = util.FormatFileSize(shareHis.FileSize)
	return res, nil
}

func (t *transferRepo) UploadEntry(ctx context.Context, req *pb.ReqUpload) (*pb.RespUpload, error) {
	id := ctx.Value("x-md-global-uid").(int)
	claims := ctx.Value("claims").(*middleware.Claims)
	var err error
	res := &pb.RespUpload{}
	filetype := util.String2StringArrWithSeparate(req.File.Filename, ".", true)
	finalName := fmt.Sprintf("%s.%s", req.File.FileHash, filetype[len(filetype)-1])

	//首先找到用户根目录
	directory := &biz.UserDirectory{}
	if id != 0 {
		err = t.data.Db.Model(&biz.UserDirectory{}).Where("user_id =? and father_id=?", id, 0).First(directory).Error
		if err != nil && err.Error() != "record not found" {
			t.log.Error(err)
			return nil, ecode.MYSQL_ERR
		}
		//没有则创建一个
		if directory.ID == 0 {
			directory.Name = claims.UserNum
			directory.UserId = id
			directory.FatherId = 0
			directory.PathStr = fmt.Sprintf("%s/", claims.UserNum)
			directory.Key = util.GetRandomDirString(40)
			err = t.data.Db.Model(&biz.UserDirectory{}).Create(directory).Error
			if err != nil {
				t.log.Error(err)
				return nil, ecode.MYSQL_ERR
			}
			directory.PathTree = fmt.Sprintf("%d/",directory.ID)
			err = t.data.Db.Model(&biz.UserDirectory{}).Where("id = ?",directory.ID).Save(directory).Error
			if err != nil {
				t.log.Error(err)
				return nil, ecode.MYSQL_ERR
			}
		}
	}
	//快传检索
	checkFile := &biz.File{}
	if req.File.FileHash != "" {
		err = t.data.Db.Model(&biz.File{}).Where("file_hash = ? and file_status = ?", req.File.FileHash,EXIST).Last(checkFile).Error
		if err != nil && err.Error() != "record not found" {
			t.log.Error(err)
			return nil, ecode.MYSQL_ERR
		}
		if checkFile.FileHash != "" {
			dir := &biz.UserDirectory{}
			if req.Did != 0 {
				err = t.data.Db.Model(&biz.UserDirectory{}).Where("id = ? and user_id =?", req.Did, id).First(dir).Error
				if err != nil {
					t.log.Error(err)
					return nil, ecode.MYSQL_ERR
				}
			} else {
				dir = directory
			}
			samePathFile := &biz.File{}
			err = t.data.Db.Model(&biz.File{}).Where("file_path = ? and file_status = ?", fmt.Sprintf("%s%s", dir.PathStr, finalName),EXIST).First(&samePathFile).Error
			if err != nil && err.Error() != "record not found" {
				t.log.Error(err)
				return nil, ecode.MYSQL_ERR
			}
			if samePathFile.ID != 0 {
				return nil, ecode.OK.SetMessage("上传文件与该路径下文件内容相同，本次请求已忽略")
			}
		}
		if checkFile.ID != 0 {
			userFile := &biz.UserFile{}
			err = t.data.Db.Model(&biz.UserFile{}).Where("file_id = ?", checkFile.ID).Last(userFile).Error
			if err != nil && err.Error() != "record not found" {
				t.log.Error(err)
				return nil, ecode.MYSQL_ERR
			}
			if userFile.ID != 0 {
				dir := &biz.UserDirectory{}
				if req.Did != 0 {
					err = t.data.Db.Model(&biz.UserDirectory{}).Where("id = ? and user_id = ?", req.Did, id).Find(dir).Error
					if err != nil {
						t.log.Error(err)
						return nil, ecode.MYSQL_ERR.SetMessage(err.Error())
					}
				} else {
					dir = directory
				}

				copyPath := fmt.Sprintf("%s%s", dir.PathStr, finalName)
				client := t.data.Minio_internal
				bucket := PUBLIC

				if userFile.UserId != 0 {
					bucket = PRIVATE
				}
				src := minio.CopySrcOptions{
					Bucket: bucket,
					Object: checkFile.FilePath,
				}
				dst := minio.CopyDestOptions{
					Bucket: PRIVATE,
					Object: copyPath,
				}

				_, err := client.CopyObject(context.TODO(), dst, src)
				if err != nil {
					return nil, ecode.EXTERNAL_API_FAIL.SetMessage(err.Error())
				}

				newFile := &biz.File{
					DirectoryId: dir.ID,
					FileName:    req.File.Filename,
					FileStatus:  EXIST,
					FileHash:    req.File.FileHash,
					FilePath:    copyPath,
					FileSize:    checkFile.FileSize,
					FileStr:     fmt.Sprintf("%s%s", dir.PathStr, req.File.Filename),
					Attribute:   1,
					ContentType: req.File.ContentType,
					Suffix:      filetype[len(filetype)-1],
				}
				err = t.data.Db.Model(&biz.File{}).Create(newFile).Error
				if err != nil {
					t.log.Error(err)
					return nil, ecode.MYSQL_ERR.SetMessage(err.Error())
				}
				newUFile := &biz.UserFile{
					FileId:      newFile.ID,
					UserId:      id,
					DirectoryId: dir.ID,
					UserNum:     claims.UserNum,
				}
				err = t.data.Db.Model(&biz.UserFile{}).Create(newUFile).Error
				if err != nil {
					t.log.Error(err)
					return nil, ecode.MYSQL_ERR.SetMessage(err.Error())
				}

				dirs := []*biz.UserDirectory{}
				if dir.Name == claims.UserNum {
					err = t.data.Db.Model(&biz.UserDirectory{}).Where("user_id = ? and father_id = ?", id, 0).Find(&dirs).Error
					if err != nil {
						t.log.Error(err)
						return nil, ecode.MYSQL_ERR
					}
				} else {
					fatherId := dir.FatherId
					for fatherId != 0 {
						temp := &biz.UserDirectory{}
						err = t.data.Db.Model(&biz.UserDirectory{}).Where("id = ?", fatherId).First(temp).Error
						if err != nil {
							t.log.Error(err)
							return nil, ecode.MYSQL_ERR
						}
						dirs = append(dirs, temp)
						fatherId = temp.FatherId
						temp = &biz.UserDirectory{}
					}
					dirs = append(dirs, dir)
				}
				for _, v := range dirs {
					v.Size += checkFile.FileSize
				}
				err = t.data.Db.Model(&biz.UserDirectory{}).Save(&dirs).Error
				if err != nil {
					t.log.Error(err)
					return nil, ecode.MYSQL_ERR
				}
				res.Fid = int32(newFile.ID)
				os.Remove(req.File.Filename)
				return res, nil
			}
		}
	}

	//filepath := fmt.Sprintf("./%s", req.File.Filename)
	//fileinfo := minio.UploadInfo{}
	//var minioUpload = func(bucket, objectName string) (minio.UploadInfo, error) {
	//	client := t.data.Minio_internal
	//	fileinfo, err = client.FPutObject(ctx, bucket, objectName, filepath, minio.PutObjectOptions{ContentType: req.File.ContentType})
	//	if err != nil {
	//		t.log.Error(err)
	//		return fileinfo, err
	//	}
	//	return fileinfo, err
	//}

	//私人文件入库
	locDir := &biz.UserDirectory{}
	if req.Did == 0 {
		locDir = directory
	} else {
		err = t.data.Db.Model(&biz.UserDirectory{}).Where("id = ? and user_id = ?", req.Did, id).First(&locDir).Error
		if err != nil {
			t.log.Error(err)
			return nil, ecode.MYSQL_ERR
		}
	}
	finalName = fmt.Sprintf("%s%s", locDir.PathStr, finalName)
	//fileinfo, err = minioUpload(PRIVATE, finalName)
	//if err != nil {
	//	t.log.Error(err)
	//	return nil, ecode.New(500).SetMessage("minio客户端错误")
	//}
	file := &biz.File{
		FileName:    req.File.Filename,
		FileHash:    req.File.FileHash,
		FilePath:    finalName,
		FileSize:    req.File.Filesize,
		FileStr:     fmt.Sprintf("%s%s", locDir.PathStr, req.File.Filename),
		Attribute:   1,
		Suffix:      filetype[len(filetype)-1],
		ContentType: req.File.ContentType,
		DirectoryId: locDir.ID,
	}
	err = t.data.Db.Model(&biz.File{}).Create(file).Error
	if err != nil {
		t.log.Error(err)
		return nil, ecode.MYSQL_ERR
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		obj := MessageObject{
			Fid:         file.ID,
			FileName:    file.FileName,
			FileHash:    file.FileHash,
			FileStr:     file.FileStr,
			FilePath:    file.FilePath,
			ContentType: file.ContentType,
			Bucket:      PRIVATE,
		}
		msg ,_:= json.Marshal(obj)
		produce := &Produce{msgContent: string(msg)}
		t.mq.RegisterProducer(produce)
		wg.Done()
	}()
	wg.Wait()
	//文件存量统计
	dirs := []*biz.UserDirectory{}
	if req.Did == 0 {
		dirs = append(dirs, directory)
	} else {
		fatherId := locDir.FatherId
		for fatherId != 0 {
			temp := &biz.UserDirectory{}
			err = t.data.Db.Model(&biz.UserDirectory{}).Where("id = ?", fatherId).First(temp).Error
			if err != nil {
				t.log.Error(err)
				return nil, ecode.MYSQL_ERR
			}
			dirs = append(dirs, temp)
			fatherId = temp.FatherId
			temp = &biz.UserDirectory{}
		}
		dirs = append(dirs, locDir)
	}
	for _, v := range dirs {
		v.Size += req.File.Filesize
	}
	err = t.data.Db.Model(&biz.UserDirectory{}).Save(&dirs).Error
	if err != nil {
		t.log.Error(err)
		return nil, ecode.MYSQL_ERR
	}
	userFile := &biz.UserFile{
		FileId:      file.ID,
		UserNum:     claims.UserNum,
		UserId:      claims.UserId,
		DirectoryId: locDir.ID,
	}
	err = t.data.Db.Model(&biz.UserFile{}).WithContext(ctx).Create(userFile).Error
	if err != nil {
		t.log.Error(err)
		return nil, ecode.MYSQL_ERR
	}
	res.Fid = int32(file.ID)
	return res, nil
}

func (t *transferRepo) DownloadEntry(ctx context.Context, req *pb.ReqDownload) (*pb.RespDownload, error) {
	res := &pb.RespDownload{}
	var err error
	//id := ctx.Value("x-md-global-uid").(int)
	//claims:=ctx.Value("claims").(*middleware.Claims)
	userFile := &biz.UserFile{}
	err = t.data.Db.Model(&biz.UserFile{}).Where("file_id = ?", req.Fid).First(userFile).Error
	if err != nil {
		return res, ecode.MYSQL_ERR
	}
	File := &biz.File{}
	err = t.data.Db.Model(&biz.File{}).Where("id = ?", req.Fid).First(File).Error
	if err != nil {
		return res, ecode.MYSQL_ERR
	}
	bucketName := PUBLIC
	if File.Attribute == 1 {
		bucketName = PRIVATE
	}

	client := t.data.minio_online
	object, err := client.GetObject(ctx, bucketName, File.FilePath, minio.GetObjectOptions{})
	if err != nil {
		return res, ecode.New(500).SetMessage("minio客户端错误")
	}
	localfile, err := os.Create(File.FileName)
	if err != nil {
		return res, ecode.New(500).SetMessage("创建文件错误")
		//return
	}
	if _, err = io.Copy(localfile, object); err != nil {
		return res, ecode.New(500).SetMessage("创建文件错误")
	}

	File.DownloadCount += 1
	err = t.data.Db.Model(&biz.File{}).Where("id = ?", File.ID).Save(File).Error
	if err != nil {
		return res, ecode.MYSQL_ERR
	}
	res.Message = "success"
	res.Status = true
	res.Filename = File.FileName
	res.Type = File.ContentType
	return res, nil
}

func (t *transferRepo) UploadStatic(ctx context.Context, req *pb.ReqStatic) (*pb.RespStatic, error) {
	res := &pb.RespStatic{}
	filepath := fmt.Sprintf("./%s", req.Filename)
	objectName := fmt.Sprintf("/%s", req.Filename)
	client := t.data.Minio_internal
	fileinfo, err := client.FPutObject(ctx, STATIC, objectName, filepath, minio.PutObjectOptions{ContentType: req.ContentType})
	fmt.Println(fileinfo.Size)
	if err != nil {
		t.log.Error(err)
		return nil, ecode.New(500).SetMessage("minio客户端错误")
	}
	//
	//timeDur := time.Second * 24 * 60 * 60
	//reqParams := make(url.Values)
	//content := fmt.Sprintf("attachment; filename=\"%s\"",req.Filename)
	//reqParams.Set("response-content-disposition", content)
	//presignedUrl,err := client.PresignedGetObject(ctx,STATIC,objectName,timeDur,reqParams)
	//if err != nil {
	//	t.log.Error(err)
	//	return nil,ecode.New(500).SetMessage("minio客户端错误")
	//}
	//fmt.Println(presignedUrl.String())
	res.FileAddress = fmt.Sprintf("http://47.107.95.82:8000/peach-static/%s", req.Filename)
	return res, nil
}

func (t *transferRepo) GetUserFileTree(ctx context.Context, req *pb.ReqGetUserFileTree) (*pb.RespGetUserFileTree, error) {
	id := ctx.Value("x-md-global-uid").(int)
	claims := ctx.Value("claims").(*middleware.Claims)
	var err error
	res := &pb.RespGetUserFileTree{}
	root := &biz.UserDirectory{}
	childDir := []*biz.UserDirectory{}
	sortName := req.SortObject
	sortType := req.SortType
	keyWords := req.Keywords
	nameNId := []*pb.DirFileNameAndId{}
	father := &pb.DirFileNameAndId{Name: "根目录",Did: 0}
	nameNId = append(nameNId,father)
	if req.DirectoryId == 0 {
		err = t.data.Db.Model(&biz.UserDirectory{}).Where("user_id = ? and name = ?", id, claims.UserNum).First(root).Error
		if err != nil {
			return res, ecode.MYSQL_ERR.SetMessage(err.Error())
		}
	} else {
		err = t.data.Db.Model(&biz.UserDirectory{}).Where("id = ? and dir_status = ?", req.DirectoryId, EXIST).First(root).Error
		if err != nil {
			return res, ecode.MYSQL_ERR.SetMessage(err.Error())
		}
		pathTree := root.PathTree
		pathStr := root.PathStr
		tree:=util.String2StringArrWithSeparate(pathTree,"/",true)
		str := util.String2StringArrWithSeparate(pathStr,"/",true)
		for k,v := range str{
			if v== claims.UserNum{
				continue
			}
			did,_ := strconv.Atoi(tree[k])
			temp := &pb.DirFileNameAndId{Name: v,Did: int32(did)}
			nameNId = append(nameNId,temp)
		}
	}
	files := []*biz.File{}
	if root.ID != 0 {
		if keyWords == "" {
			err = t.data.Db.Model(&biz.File{}).Where("directory_id = ? and file_status = ?", root.ID, EXIST).Find(&files).Error
			if err != nil {
				return res, ecode.MYSQL_ERR.SetMessage(err.Error())
			}
			err = t.data.Db.Model(&biz.UserDirectory{}).Where("father_id = ? and dir_status = ?", root.ID, EXIST).Find(&childDir).Error
			if err != nil && err.Error() != "record not found" {
				return res, ecode.MYSQL_ERR.SetMessage(err.Error())
			}
		} else {
			err = t.data.Db.Model(&biz.File{}).Where("directory_id = ? and file_name like ? and file_status = ?", root.ID, "%"+keyWords+"%", EXIST).Find(&files).Error
			if err != nil {
				return res, ecode.MYSQL_ERR.SetMessage(err.Error())
			}
			err = t.data.Db.Model(&biz.UserDirectory{}).Where("father_id = ? and name like ? and dir_status = ?", root.ID, "%"+keyWords+"%", EXIST).Find(&childDir).Error
			if err != nil && err.Error() != "record not found" {
				return res, ecode.MYSQL_ERR.SetMessage(err.Error())
			}
		}
	}

	//排序两参同传
	if sortName != 0 && sortType != 0 {
		switch sortName {
		case NAME:
			SortFile(files, func(p, q *biz.File) bool {
				if sortType == DESC {
					return len(p.FileName) > len(q.FileName)
				}
				return len(p.FileName) < len(q.FileName)
			})
			SortDir(childDir, func(p, q *biz.UserDirectory) bool {
				if sortType == DESC {
					return len(p.Name) > len(q.Name)
				}
				return len(p.Name) < len(q.Name)
			})
		case TIME:
			SortFile(files, func(p, q *biz.File) bool {
				if sortType == DESC {
					return p.UpdatedAt > q.UpdatedAt
				}
				return p.UpdatedAt < q.UpdatedAt
			})
			SortDir(childDir, func(p, q *biz.UserDirectory) bool {
				if sortType == DESC {
					return p.UpdatedAt > q.UpdatedAt
				}
				return p.UpdatedAt < q.UpdatedAt
			})
		case SIZE:
			SortFile(files, func(p, q *biz.File) bool {
				if sortType == DESC {
					return p.FileSize > q.FileSize
				}
				return p.FileSize < q.FileSize
			})
			SortDir(childDir, func(p, q *biz.UserDirectory) bool {
				if sortType == DESC {
					return p.Size > q.Size
				}
				return p.Size < q.Size
			})
		}
	}

	FileMeta := []*pb.FileMetaObject{}
	fids := []int32{}
	for _, v := range files {
		meta := &pb.FileMetaObject{
			Fid:          int32(v.ID),
			Size:         util.FormatFileSize(v.FileSize),
			FileName:     v.FileName,
			FileType:     v.Suffix,
			LastModified: v.UpdatedAt,
			Key:          v.FileHash,
		}
		FileMeta = append(FileMeta, meta)
		fids = append(fids, int32(v.ID))
	}

	DirMeta := []*pb.DirMetaObject{}
	dids := []int32{}
	for _, v := range childDir {
		meta := &pb.DirMetaObject{
			Did:          int32(v.ID),
			Size:         util.FormatFileSize(v.Size),
			DirName:      v.Name,
			LastModified: v.UpdatedAt,
			Key:          v.Key,
		}
		DirMeta = append(DirMeta, meta)
		dids = append(dids, int32(v.ID))
	}
	res.UserId = int32(id)
	res.Total = int32(len(FileMeta) + len(DirMeta))
	res.FileObject = FileMeta
	res.DirObject = DirMeta
	res.Fids = fids
	res.Dids = dids
	res.LocId = req.DirectoryId
	res.DirNameId = nameNId
	return res, nil
}

func (t *transferRepo) GetTrashBin(ctx context.Context, req *pb.ReqGetUserTrashBin) (*pb.RespGetUserTrashBin, error) {
	id := ctx.Value("x-md-global-uid").(int)
	//claims := ctx.Value("claims").(*middleware.Claims)
	var err error
	res := &pb.RespGetUserTrashBin{}
	File := []*biz.File{}
	Dir := []*biz.UserDirectory{}
	sortName := req.SortObject
	sortType := req.SortType
	err = t.data.Db.Model(&biz.UserDirectory{}).Where("user_id = ? and dir_status = ?", id, DEL).Find(&Dir).Error
	if err != nil && err.Error() != "record not found" {
		return res, ecode.MYSQL_ERR.SetMessage(err.Error())
	}
	dirMap := make(map[int]int)
	for _, v := range Dir {
		dirMap[v.ID] = v.FatherId
	}
	//newDir:=Dir
doit:
	for k, v := range Dir {
		if _, exist := dirMap[v.FatherId]; exist {
			Dir = append(Dir[:k], Dir[k+1:]...)
			goto doit
		}
	}

	err = t.data.Db.Raw("select * from d_storage.files f where file_status  = 2 and id in "+
		"(select id from d_storage.user_files uf where uf.user_id = ?) and directory_id not in "+
		"(select id from d_storage.user_directories ud where ud.dir_status =2 and ud.user_id = ? );",
		id, id).Scan(&File).Error
	if err != nil {
		return res, ecode.MYSQL_ERR.SetMessage(err.Error())
	}
	//排序两参同传
	if sortName != 0 && sortType != 0 {
		switch sortName {
		case NAME:
			SortFile(File, func(p, q *biz.File) bool {
				if sortType == DESC {
					return p.FileName > q.FileName
				}
				return p.FileName < q.FileName
			})
			SortDir(Dir, func(p, q *biz.UserDirectory) bool {
				if sortType == DESC {
					return p.Name > q.Name
				}
				return p.Name < q.Name
			})
		case TIME:
			SortFile(File, func(p, q *biz.File) bool {
				if sortType == DESC {
					return p.UpdatedAt > q.UpdatedAt
				}
				return p.UpdatedAt < q.UpdatedAt
			})
			SortDir(Dir, func(p, q *biz.UserDirectory) bool {
				if sortType == DESC {
					return p.UpdatedAt > q.UpdatedAt
				}
				return p.UpdatedAt < q.UpdatedAt
			})
		case SIZE:
			SortFile(File, func(p, q *biz.File) bool {
				if sortType == DESC {
					return p.FileSize > q.FileSize
				}
				return p.FileSize < q.FileSize
			})
			SortDir(Dir, func(p, q *biz.UserDirectory) bool {
				if sortType == DESC {
					return p.Size > q.Size
				}
				return p.Size < q.Size
			})
		}
	}

	FileMeta := []*pb.FileMetaObject{}
	fids := []int32{}
	for _, v := range File {
		meta := &pb.FileMetaObject{
			Fid:          int32(v.ID),
			Size:         util.FormatFileSize(v.FileSize),
			FileName:     v.FileName,
			FileType:     v.Suffix,
			LastModified: v.UpdatedAt,
			Key:          v.FileHash,
		}
		FileMeta = append(FileMeta, meta)
		fids = append(fids, int32(v.ID))
	}

	DirMeta := []*pb.DirMetaObject{}
	dids := []int32{}
	for _, v := range Dir {
		meta := &pb.DirMetaObject{
			Did:          int32(v.ID),
			Size:         util.FormatFileSize(v.Size),
			DirName:      v.Name,
			LastModified: v.UpdatedAt,
			Key:          v.Key,
		}
		DirMeta = append(DirMeta, meta)
		dids = append(dids, int32(v.ID))
	}
	res.UserId = int32(id)
	res.Total = int32(len(FileMeta) + len(DirMeta))
	res.FileObject = FileMeta
	res.DirObject = DirMeta
	res.Fids = fids
	res.Dids = dids
	return res, nil
}

func (t *transferRepo) DeleteFile(ctx context.Context, req *pb.ReqDeleteFile) (*pb.RespDelete, error) {
	res := &pb.RespDelete{}
	id := ctx.Value("x-md-global-uid").(int)
	var err error
	dids := []int{}
	files := []*biz.File{}
	if len(req.Fid) != 0 {
		err = t.data.Db.Model(&biz.File{}).Where("id in (?)  and file_status = ?", req.Fid, EXIST).Find(&files).Error
		if err != nil {
			return res, ecode.MYSQL_ERR.SetMessage(err.Error())
		}
	}

	client := t.data.Minio_internal
	for _, v := range files {
		src := minio.CopySrcOptions{
			Bucket: PRIVATE,
			Object: v.FilePath,
		}
		dst := minio.CopyDestOptions{
			Bucket: DELETE,
			Object: v.FilePath,
		}
		_, err = client.CopyObject(ctx, dst, src)
		if err != nil {
			return res, ecode.New(500).SetMessage(err.Error())
		}
		ropt := minio.RemoveObjectOptions{
			ForceDelete: true,
		}
		err = client.RemoveObject(ctx, PRIVATE, v.FilePath, ropt)
		if err != nil {
			return res, ecode.New(500).SetMessage(err.Error())
		}
		v.FileStatus = DEL
		dids = append(dids, v.DirectoryId)
	}
	dids = util.UniqueIntArr(dids)
	//没办法了
	dirs := []*biz.UserDirectory{}
	//找到当前目录的文件夹
	err = t.data.Db.Model(&biz.UserDirectory{}).Where("id in (?) and user_id =?", dids, id).Find(&dirs).Error
	if err != nil {
		return res, ecode.MYSQL_ERR
	}
	dirMap := make(map[int][]*biz.File)
	for _, v := range dirs {
		temp := []*biz.File{}
		dirMap[v.ID] = temp
	}

	for _, v := range files {
		if _, exist := dirMap[v.DirectoryId]; exist {
			dirMap[v.DirectoryId] = append(dirMap[v.DirectoryId], v)
		}
	}

	sizeMap := make(map[int]int64)
	for _, v := range dirs {
		countFile := dirMap[v.ID]
		var sizeCount int64
		for _, file := range countFile {
			sizeCount += file.FileSize
		}
		if sizeCount > 0 {
			v.Size -= sizeCount
		}
		sizeMap[v.ID] = sizeCount
		sizeCount = 0
	}
	tx := t.data.Db.Begin()
	err = tx.Model(&biz.UserDirectory{}).Save(&dirs).Error
	if err != nil {
		tx.Rollback()
		return res, ecode.MYSQL_ERR.SetMessage(err.Error())
	}

	err = tx.Model(&biz.File{}).Save(files).Error
	if err != nil {
		tx.Rollback()
		return res, ecode.MYSQL_ERR.SetMessage(err.Error())
	}
	tx.Commit()

	pathMap := make(map[int][]string)
	for _, v := range dirs {
		pathArr := util.String2StringArrWithSeparate(v.PathTree, "/", true)
		pathMap[v.ID] = pathArr
	}

	for k, v := range pathMap {
		err = t.data.Db.Raw("update d_storage.user_directories ud set ud.size = ud.size - ? where id in(?) and user_id = ?",
			sizeMap[k], v, id).Error
		if err != nil {
			return res, ecode.MYSQL_ERR.SetMessage(err.Error())
		}
	}
	res.Status = true
	res.Message = "success"
	return res, err
}
func (t *transferRepo) DeleteDir(ctx context.Context, req *pb.ReqDeleteDir) (*pb.RespDelete, error) {
	id := ctx.Value("x-md-global-uid").(int)
	res := &pb.RespDelete{}
	var err error
	dirs := []*biz.UserDirectory{}

	//找到文件夹下的子目录
	dirMap := make(map[int][]*biz.UserDirectory)
	if len(req.Did) != 0 {
		err = t.data.Db.Model(&biz.UserDirectory{}).Where("id in (?) and user_id= ? and dir_status = ?", req.Did, id, EXIST).Find(&dirs).Error
		if err != nil {
			return res, ecode.MYSQL_ERR.SetMessage(err.Error())
		}
	}
	dirSizeMap := make(map[int]int64)
	for _, v := range dirs {
		dirSizeMap[v.ID] = v.Size
	}
	childDir := []*biz.UserDirectory{}

	for _, v := range dirs {

		err = t.data.Db.Raw(
			"select * from  user_directories a "+
				"left join(select path_tree from user_directories d where d.path_tree like ? and user_id = ?)b "+
				"on a.path_tree=b.path_tree where b.path_tree is not null ", "%"+strconv.Itoa(v.ID)+"/%",id).Scan(&childDir).Error

		if err != nil && err.Error() != "record not found" {
			return res, ecode.MYSQL_ERR.SetMessage(err.Error())
		}
		if len(childDir) != 0 {
			dirMap[v.ID] = childDir
			childDir = []*biz.UserDirectory{}
		}
	}

	//找到目录下的文件
	fileMap := make(map[int][]*biz.File)
	saveDir := []*biz.UserDirectory{}
	for _, dirs := range dirMap {
		for _, dir := range dirs {
			files := []*biz.File{}
			err = t.data.Db.Model(&biz.File{}).Where("directory_id = ? and file_status = ?", dir.ID, EXIST).Find(&files).Error
			if err != nil && err.Error() != "record not found" {
				return res, ecode.MYSQL_ERR.SetMessage(err.Error())
			}
			dir.DirStatus = DEL
			saveDir = append(saveDir, dir)
			if len(files) != 0 {
				fileMap[dir.ID] = files
				files = []*biz.File{}
			}
		}
	}
	saveFile := []*biz.File{}
	client := t.data.Minio_internal
	for _, files := range fileMap {
		for _, file := range files {
			src := minio.CopySrcOptions{
				Bucket: PRIVATE,
				Object: file.FilePath,
			}
			dst := minio.CopyDestOptions{
				Bucket: DELETE,
				Object: file.FilePath,
			}
			_, err = client.CopyObject(ctx, dst, src)
			if err != nil {
				return res, ecode.New(500).SetMessage(err.Error())
			}
			ropt := minio.RemoveObjectOptions{
				ForceDelete: true,
			}
			err = client.RemoveObject(ctx, PRIVATE, file.FilePath, ropt)
			if err != nil {
				return res, ecode.New(500).SetMessage(err.Error())
			}
			file.FileStatus = DEL
			saveFile = append(saveFile, file)
		}
	}
	tx := t.data.Db.Begin()
	if len(saveFile) != 0 {
		err = tx.Model(&biz.File{}).Save(&saveFile).Error
		if err != nil {
			tx.Rollback()
			return res, ecode.MYSQL_ERR.SetMessage(err.Error())
		}
	}
	err = tx.Model(&biz.UserDirectory{}).Save(&saveDir).Error
	if err != nil {
		tx.Rollback()
		return res, ecode.MYSQL_ERR.SetMessage(err.Error())
	}
	tx.Commit()

	pathMap := make(map[int][]string)
	for _, v := range dirs {
		pathArr := util.String2StringArrWithSeparate(v.PathTree, "/", true)
		//for _,value := range pathArr{
		//
		//}
		pathMap[v.ID] = pathArr
	}

	for k, v := range pathMap {
		err = t.data.Db.Raw("update d_storage.user_directories ud set ud.size = ud.size - ? where name in(?) and user_id = ?",
			dirSizeMap[k], v, id).Error
		if err != nil {
			return res, ecode.MYSQL_ERR.SetMessage(err.Error())
		}
	}


	res.Status = true
	res.Message = "success"
	return res, err
}
func (t *transferRepo) WithDrawFile(ctx context.Context, req *pb.ReqWithDrawFile) (*pb.RespWithDraw, error) {
	id := ctx.Value("x-md-global-uid").(int)
	res := &pb.RespWithDraw{}
	var err error
	dids := []int{}
	files := []*biz.File{}
	if len(req.Fid) != 0 {
		err = t.data.Db.Model(&biz.File{}).Where("id in (?)  and file_status = ?", req.Fid, DEL).Find(&files).Error
		if err != nil {
			return res, ecode.MYSQL_ERR.SetMessage(err.Error())
		}
	}
	client := t.data.Minio_internal
	for _, v := range files {
		src := minio.CopySrcOptions{
			Bucket: DELETE,
			Object: v.FilePath,
		}
		dst := minio.CopyDestOptions{
			Bucket: PRIVATE,
			Object: v.FilePath,
		}
		_, err = client.CopyObject(ctx, dst, src)
		if err != nil {
			return res, ecode.New(500).SetMessage(err.Error())
		}
		ropt := minio.RemoveObjectOptions{
			ForceDelete: true,
		}
		err = client.RemoveObject(ctx, DELETE, v.FilePath, ropt)
		if err != nil {
			return res, ecode.New(500).SetMessage(err.Error())
		}
		v.FileStatus = EXIST
		dids = append(dids, v.DirectoryId)
	}
	dids = util.UniqueIntArr(dids)
	dirs := []*biz.UserDirectory{}
	err = t.data.Db.Model(&biz.UserDirectory{}).Where("id in (?) and user_id =?", dids, id).Find(&dirs).Error
	if err != nil {
		return res, ecode.MYSQL_ERR
	}
	dirMap := make(map[int][]*biz.File)
	for _, v := range dirs {
		temp := []*biz.File{}
		dirMap[v.ID] = temp
	}
	for _, v := range files {
		if _, exist := dirMap[v.DirectoryId]; exist {
			dirMap[v.DirectoryId] = append(dirMap[v.DirectoryId], v)
		}
	}

	sizeMap := make(map[int]int64)
	for _, v := range dirs {
		countFile := dirMap[v.ID]
		var sizeCount int64
		for _, file := range countFile {
			sizeCount += file.FileSize
		}
		if sizeCount > 0 {
			v.Size += sizeCount
		}
		sizeMap[v.ID] = sizeCount
		sizeCount = 0
	}

	tx := t.data.Db.Begin()
	err = tx.Model(&biz.UserDirectory{}).Save(&dirs).Error
	if err != nil {
		tx.Rollback()
		return res, ecode.MYSQL_ERR.SetMessage(err.Error())
	}
	err = tx.Model(&biz.File{}).Save(files).Error
	if err != nil {
		return res, ecode.MYSQL_ERR.SetMessage(err.Error())
	}
	tx.Commit()

	pathMap := make(map[int][]string)
	for _, v := range dirs {
		pathArr := util.String2StringArrWithSeparate(v.PathTree, "/", true)
		pathMap[v.ID] = pathArr
	}

	for k, v := range pathMap {
		err = t.data.Db.Raw("update d_storage.user_directories ud set ud.size = ud.size + ? where id in(?) and user_id = ?",
			sizeMap[k], v, id).Error
		if err != nil {
			return res, ecode.MYSQL_ERR.SetMessage(err.Error())
		}
	}
	res.Status = true
	res.Message = "success"
	return res, err

}
func (t *transferRepo) WithDrawDir(ctx context.Context, req *pb.ReqWithDrawDir) (*pb.RespWithDraw, error) {
	id := ctx.Value("x-md-global-uid").(int)
	res := &pb.RespWithDraw{}
	var err error
	dirs := []*biz.UserDirectory{}
	//找到文件夹下的子目录
	dirMap := make(map[int][]*biz.UserDirectory)
	if len(req.Did) != 0 {
		err = t.data.Db.Model(&biz.UserDirectory{}).Where("id in (?) and user_id= ? and dir_status = ?", req.Did, id, DEL).Find(&dirs).Error
		if err != nil {
			return res, ecode.MYSQL_ERR.SetMessage(err.Error())
		}
	}
	dirSizeMap := make(map[int]int64)
	for _, v := range dirs {
		dirSizeMap[v.ID] = v.Size
	}
	childDir := []*biz.UserDirectory{}

	for _, v := range dirs {
		err = t.data.Db.Raw(
			"select * from  user_directories a "+
				"left join(select path_str from user_directories d where d.path_tree like ? and d.father_id =?)b "+
				"on a.path_tree=b.path_tree where b.path_str is not null", "%"+strconv.Itoa(v.ID)+"/%", v.FatherId).Scan(&childDir).Error
		if err != nil && err.Error() != "record not found" {
			return res, ecode.MYSQL_ERR.SetMessage(err.Error())
		}
		if len(childDir) != 0 {
			dirMap[v.ID] = childDir
			childDir = []*biz.UserDirectory{}
		}
	}

	//找到目录下的文件
	fileMap := make(map[int][]*biz.File)
	saveDir := []*biz.UserDirectory{}
	for _, dirs := range dirMap {
		for _, dir := range dirs {
			files := []*biz.File{}
			err = t.data.Db.Model(&biz.File{}).Where("directory_id = ? and file_status = ?", dir.ID, DEL).Find(&files).Error
			if err != nil && err.Error() != "record not found" {
				return res, ecode.MYSQL_ERR.SetMessage(err.Error())
			}
			dir.DirStatus = EXIST
			saveDir = append(saveDir, dir)
			if len(files) != 0 {
				fileMap[dir.ID] = files
				files = []*biz.File{}
			}
		}
	}
	saveFile := []*biz.File{}

	client := t.data.Minio_internal
	for _, files := range fileMap {
		for _, file := range files {
			src := minio.CopySrcOptions{
				Bucket: DELETE,
				Object: file.FilePath,
			}
			dst := minio.CopyDestOptions{
				Bucket: PRIVATE,
				Object: file.FilePath,
			}
			_, err = client.CopyObject(ctx, dst, src)
			if err != nil {
				return res, ecode.New(500).SetMessage(err.Error())
			}
			ropt := minio.RemoveObjectOptions{
				ForceDelete: true,
			}
			err = client.RemoveObject(ctx, DELETE, file.FilePath, ropt)
			if err != nil {
				return res, ecode.New(500).SetMessage(err.Error())
			}
			file.FileStatus = EXIST
			//dirSizeMap[file.DirectoryId] += file.FileSize
			saveFile = append(saveFile, file)
		}
	}

	tx := t.data.Db.Begin()
	if len(saveFile) != 0 {
		err = tx.Model(&biz.File{}).Save(&saveFile).Error
		if err != nil {
			tx.Rollback()
			return res, ecode.MYSQL_ERR.SetMessage(err.Error())
		}
	}

	err = tx.Model(&biz.UserDirectory{}).Save(&saveDir).Error
	if err != nil {
		tx.Rollback()
		return res, ecode.MYSQL_ERR.SetMessage(err.Error())
	}
	tx.Commit()

	pathMap := make(map[int][]string)
	for _, v := range dirs {
		pathArr := util.String2StringArrWithSeparate(v.PathStr, "/", true)
		pathMap[v.ID] = pathArr
	}

	for k, v := range pathMap {
		err = t.data.Db.Raw("update d_storage.user_directories ud set ud.size = ud.size + ? where id in(?) and user_id = ?",
			dirSizeMap[k], v, id).Error
		if err != nil {
			return res, ecode.MYSQL_ERR.SetMessage(err.Error())
		}
	}

	res.Status = true
	res.Message = "success"
	return res, err
}
func (t *transferRepo) CleanTrashFile(ctx context.Context, req *pb.ReqCleanTrashFile) (*pb.RespCleanTrash, error) {
	res := &pb.RespCleanTrash{}
	var err error
	//id := ctx.Value("x-md-global-uid").(int)

	files := []*biz.File{}
	if len(req.Fid) != 0 {
		err = t.data.Db.Model(&biz.File{}).Where("id in (?) and file_status = ?", req.Fid, DEL).Find(&files).Error
		if err != nil {
			return res, ecode.MYSQL_ERR.SetMessage(err.Error())
		}
	}

	client := t.data.Minio_internal
	//var sizeCut int64
	for _, v := range files {
		ropt := minio.RemoveObjectOptions{
			ForceDelete: true,
		}
		err = client.RemoveObject(ctx, DELETE, v.FilePath, ropt)
		if err != nil {
			return res, ecode.New(500).SetMessage(err.Error())
		}
		//sizeCut += v.FileSize
		v.FileStatus = RUIN
	}

	//tx := t.data.db
	err = t.data.Db.Model(&biz.File{}).Save(&files).Error
	if err != nil {
		return res, ecode.MYSQL_ERR.SetMessage(err.Error())
	}
	//tx.Commit()
	res.Status = true
	res.Message = "success"
	return res, err
}
func (t *transferRepo) CleanTrashDir(ctx context.Context, req *pb.ReqCleanTrashDir) (*pb.RespCleanTrash, error) {
	id := ctx.Value("x-md-global-uid").(int)
	res := &pb.RespCleanTrash{}
	var err error
	dirs := []*biz.UserDirectory{}
	//找到文件夹下的子目录
	dirMap := make(map[int][]*biz.UserDirectory)
	if len(req.Did) != 0 {
		err = t.data.Db.Model(&biz.UserDirectory{}).Where("id in (?) user_id= ? and dir_status = ?", req.Did, id, DEL).Find(&dirs).Error
		if err != nil {
			return res, ecode.MYSQL_ERR.SetMessage(err.Error())
		}
	}
	childDir := []*biz.UserDirectory{}

	for _, v := range dirs {
		err = t.data.Db.Raw(
			"select * from  user_directories a "+
				"left join(select path_str from user_directories d where d.path_tree like ?)b "+
				"on a.path_tree=b.path_tree where b.path_str is not null", "%"+strconv.Itoa(v.ID)+"/%").Scan(&childDir).Error
		if err != nil && err.Error() != "record not found" {
			return res, ecode.MYSQL_ERR.SetMessage(err.Error())
		}
		dirMap[v.ID] = childDir
		childDir = []*biz.UserDirectory{}
	}

	//找到目录下的文件
	fileMap := make(map[int][]*biz.File)
	saveDir := []*biz.UserDirectory{}
	for _, dirs := range dirMap {
		for _, dir := range dirs {
			files := []*biz.File{}
			err = t.data.Db.Model(&biz.File{}).Where("directory_id = ? and file_status = ?", dir.ID, DEL).Find(&files).Error
			if err != nil && err.Error() != "record not found" {
				return res, ecode.MYSQL_ERR.SetMessage(err.Error())
			}
			if len(files) != 0 {
				fileMap[dir.ID] = files
			}
			dir.DirStatus = RUIN
			saveDir = append(saveDir, dir)
			files = []*biz.File{}
		}
	}
	saveFile := []*biz.File{}
	client := t.data.Minio_internal
	for _, files := range fileMap {
		for _, file := range files {
			ropt := minio.RemoveObjectOptions{
				ForceDelete: true,
			}
			err = client.RemoveObject(ctx, DELETE, file.FilePath, ropt)
			if err != nil {
				return res, ecode.New(500).SetMessage(err.Error())
			}
			file.FileStatus = RUIN
			saveFile = append(saveFile, file)
		}
	}

	tx := t.data.Db.Begin()
	err = tx.Model(&biz.File{}).Save(&saveFile).Error
	if err != nil {
		tx.Rollback()
		return res, ecode.MYSQL_ERR.SetMessage(err.Error())
	}

	err = tx.Model(&biz.UserDirectory{}).Save(&saveDir).Error
	if err != nil {
		tx.Rollback()
		return res, ecode.MYSQL_ERR.SetMessage(err.Error())
	}
	tx.Commit()
	res.Status = true
	res.Message = "success"
	return res, err
}

func (t *transferRepo) ShareFile(ctx context.Context, req *pb.ReqShareFileStr) (*pb.RespShareFileStr, error) {
	id := ctx.Value("x-md-global-uid").(int)
	res := &pb.RespShareFileStr{}
	//claims := ctx.Value("claims").(*middleware.Claims)
	var err error
	userfile := &biz.UserFile{}
	err = t.data.Db.Model(&biz.UserFile{}).Where("user_id = ? and file_id = ?", id, req.Fid).Find(userfile).Error
	if err != nil {
		return res, ecode.MYSQL_ERR.SetMessage(err.Error())
	}
	file := &biz.File{}
	err = t.data.Db.Model(&biz.File{}).Where("id = ?", userfile.FileId).Find(file).Error
	if err != nil {
		return res, ecode.MYSQL_ERR.SetMessage(err.Error())
	}
	timeDur := time.Minute * 24 * 60 * 7
	if req.ExpireTime != 0 {
		timeDur = time.Second * time.Duration(req.ExpireTime)
	}
	reqParams := make(url.Values)
	content := fmt.Sprintf("attachment; filename=\"%s\"", file.FileName)
	reqParams.Set("response-content-disposition", content)
	bucket := PUBLIC
	if file.Attribute == 1 {
		bucket = PRIVATE
	}
	presignedUrl, err := t.data.minio_online.PresignedGetObject(ctx, bucket, file.FilePath, timeDur, reqParams)
	if err != nil {
		t.log.Error(err)
		return res, ecode.EXTERNAL_API_FAIL.SetMessage("minio客户端错误")
	}
	getCode := util.GetRandomInt(6)
	err = t.data.cache.Set(context.TODO(), getCode, presignedUrl.String(), timeDur).Err()
	if err != nil {
		return res, ecode.EXTERNAL_API_FAIL.SetMessage("minio客户端错误")
	}
	res.GetCode = getCode
	res.Fid = int32(file.ID)
	newHistory := &biz.ShareHistory{
		Fid:        file.ID,
		Describe:   req.Describe,
		ExpireTime: timeDur.Nanoseconds() / 1e9,
		Title:      req.Title,
		GetCode:    getCode,
		FileName:   file.FileName,
		FileSize:   file.FileSize,
	}
	err = t.data.Db.Model(&biz.ShareHistory{}).Create(newHistory).Error
	if err != nil {
		return res, ecode.MYSQL_ERR.SetMessage(err.Error())
	}
	return res, nil

}
func (t *transferRepo) PreviewFile(ctx context.Context, req *pb.ReqPreviewFile) (*pb.RespPreviewFile, error) {
	res := &pb.RespPreviewFile{}
	var err error
	id := ctx.Value("x-md-global-uid").(int)
	if id == 0 {

	}
	userFile := &biz.UserFile{}
	err = t.data.Db.Model(&biz.UserFile{}).Where("user_id = ? and file_id = ?", id, req.Fid).First(userFile).Error
	if err != nil {
		return res, ecode.MYSQL_ERR
	}
	if userFile == nil {
		return res, ecode.MYSQL_ERR.SetMessage("找不到此文件")
	}
	File := &biz.File{}
	err = t.data.Db.Model(&biz.File{}).Where("id = ?", req.Fid).First(File).Error
	if err != nil {
		return res, ecode.MYSQL_ERR
	}
	client := t.data.minio_online
	timeDur := time.Second * 24 * 60 * 60
	reqParams := make(url.Values)
	content := fmt.Sprintf("inline; filename=\"%s\"", File.FileName)
	reqParams.Set("response-content-disposition", content)
	bucket := PUBLIC
	if File.Attribute == 1 {
		bucket = PRIVATE
	}
	presignedUrl, err := client.PresignedGetObject(ctx, bucket, File.FilePath, timeDur, reqParams)
	if err != nil {
		t.log.Error(err)
		return nil, ecode.New(500).SetMessage("minio客户端错误")
	}
	res.PreviewStr = presignedUrl.String()
	return res, nil
}

func (t *transferRepo) FileCensus(ctx context.Context, req *pb.ReqFileCensus) (*pb.RespFileCensus, error) {
	id := ctx.Value("x-md-global-uid").(int)
	var err error
	res := &pb.RespFileCensus{}
	dirIds := []int{}
	err = t.data.Db.Model(&biz.UserDirectory{}).Where("user_id = ? and dir_status = ?", id, EXIST).Pluck("id", &dirIds).Error
	if err != nil && err.Error() != "record not found" {
		return res, ecode.MYSQL_ERR.SetMessage(err.Error())
	}
	files := []*biz.File{}
	err = t.data.Db.Model(&biz.File{}).
		Where("directory_id in (?) and file_status = ?", dirIds, EXIST).
		Order("download_count desc").
		Find(&files).Error
	if err != nil && err.Error() != "record not found" {
		return res, ecode.MYSQL_ERR.SetMessage(err.Error())
	}
	var total int64
	//err = t.data.db.Model(&biz.File{}).Select("sum(file_size)").Error
	topTen := &pb.DownloadCensus{}
	name := []string{}
	value := []int32{}
	if len(files) <= 10 {
		for _, v := range files {
			name = append(name, v.FileName)
			value = append(value, int32(v.DownloadCount))
		}
	} else {
		temp := files[:10]
		for _, v := range temp {
			name = append(name, v.FileName)
			value = append(value, int32(v.DownloadCount))
		}
	}
	topTen.Name = name
	topTen.Value = value
	countMap := make(map[string]int)
	for _, v := range files {
		total += v.FileSize
		if _, exist := TypeMap[v.Suffix]; exist {
			countMap[TypeMap[v.Suffix]] += 1
		} else {
			countMap["other"] += 1
		}
	}
	last := float32(total) / float32(1024*1024*1024*10)
	err = t.data.Db.Model(&biz.UserDirectory{}).Where("user_id  = ? and father_id = 0", id).Update("size", total).Error
	if err != nil {
		return res, ecode.MYSQL_ERR
	}
	usage := &pb.Usage{
		UseStr: fmt.Sprintf("%s/%s", util.FormatFileSize(total), util.FormatFileSize(1024*1024*1024*10)),
		Used:   last,
	}
	ratios := []*pb.FileRatio{}
	for k, v := range countMap {
		ratio := &pb.FileRatio{
			Name:  k,
			Value: int32(v),
		}
		ratios = append(ratios, ratio)
	}
	res.FileRatio = ratios
	res.Usage = usage
	res.TopTen = topTen
	return res, nil

}

func (t *transferRepo) SearchFile(ctx context.Context, req *pb.ReqSearchFile) (*pb.RespSearchFile, error) {
	return nil, nil
}

func (t *transferRepo) CreateDir(ctx context.Context, req *pb.ReqCreateDir) (*pb.RespCreateDir, error) {
	res := &pb.RespCreateDir{}
	id := ctx.Value("x-md-global-uid").(int)
	claims := ctx.Value("claims").(*middleware.Claims)
	var err error
	if req.DirName == "" {
		return res, ecode.INVALID_PARAM.SetMessage("文件夹名不能为空")
	}
	directory := &biz.UserDirectory{}
	if id != 0 {
		err = t.data.Db.Model(&biz.UserDirectory{}).Where("user_id =? and father_id=?", id, 0).First(directory).Error
		if err != nil && err.Error() != "record not found" {
			t.log.Error(err)
			return nil, ecode.MYSQL_ERR
		}
		//没有则创建一个
		if directory.ID == 0 {
			directory.Name = claims.UserNum
			directory.UserId = id
			directory.FatherId = 0
			directory.PathStr = fmt.Sprintf("%s/", claims.UserNum)
			directory.Key = util.GetRandomDirString(40)
			err = t.data.Db.Model(&biz.UserDirectory{}).Create(directory).Error
			if err != nil {
				t.log.Error(err)
				return nil, ecode.MYSQL_ERR
			}
			directory.PathTree = fmt.Sprintf("%d/",directory.ID)
			err = t.data.Db.Model(&biz.UserDirectory{}).Where("id = ?",directory.ID).Save(directory).Error
			if err != nil {
				t.log.Error(err)
				return nil, ecode.MYSQL_ERR
			}
		}
	}
	existCheck := &biz.UserDirectory{}
	err = t.data.Db.Model(&biz.UserDirectory{}).Where("name = ? and user_id = ?", req.DirName, id).First(existCheck).Error
	if existCheck.ID != 0 {
		dir := &biz.UserDirectory{}
		if req.LocDid == 0 {
			err = t.data.Db.Model(&biz.UserDirectory{}).Where("user_id = ? and father_id = ?", id, 0).First(dir).Error
			if err != nil {
				return res, ecode.MYSQL_ERR
			}
			if existCheck.FatherId == dir.ID{
				return res, ecode.INVALID_PARAM.SetMessage("新建文件夹与当前目录下重名")
			}
		} else {
			err = t.data.Db.Model(&biz.UserDirectory{}).Where("id = ?", req.LocDid).First(dir).Error
			if err != nil {
				return res, ecode.MYSQL_ERR
			}
			if existCheck.FatherId == dir.FatherId {
				return res, ecode.INVALID_PARAM.SetMessage("新建文件夹与当前目录下重名")
			}
		}

	}
	dir := &biz.UserDirectory{}
	if req.LocDid == 0 {
		dir = directory
	} else {
		err = t.data.Db.Model(&biz.UserDirectory{}).Where("user_id = ? and id = ?", id, req.LocDid).First(dir).Error
		if err != nil {
			return res, ecode.MYSQL_ERR
		}
	}

	newDir := &biz.UserDirectory{
		FatherId:  dir.ID,
		UserId:    id,
		PathStr:   fmt.Sprintf("%s%s/", dir.PathStr, req.DirName),
		Name:      req.DirName,
		DirStatus: EXIST,
		Key:       util.GetRandomDirString(40),
	}

	err = t.data.Db.Model(&biz.UserDirectory{}).Create(newDir).Error
	if err != nil {
		return res, ecode.MYSQL_ERR
	}
	res.Did = int32(newDir.ID)
	err = t.data.Db.Model(&biz.UserDirectory{}).Where("id = ?",newDir.ID).Update("path_tree",fmt.Sprintf("%s%d/",dir.PathTree,newDir.ID)).Error
	if err != nil {
		return res, ecode.MYSQL_ERR
	}
	return res, nil
}

package data

import (
	pb "banana/api/transfer/service/v1"
	"banana/app/transfer/service/internal/biz"
	"banana/pkg/ecode"
	"banana/pkg/middleware"
	"banana/pkg/util"
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/minio/minio-go/v7"
	"io"
	"net/url"
	"os"
	"time"
)

var _ biz.TransferRepo = (*transferRepo)(nil)

type transferRepo struct {
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
	RUIN  = 0

	DEFAULT = 0
	ASC     = 1
	DESC    = 2

	NAME = 1
	TIME = 2
	SIZE = 3
)

func NewTransferRepo(data *Data, logger log.Logger) biz.TransferRepo {
	return &transferRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "data/tf")),
	}
}

func (t *transferRepo) UploadEntry(ctx context.Context, req *pb.ReqUpload) (*pb.RespUpload, error) {
	id := ctx.Value("x-md-global-uid").(int)
	claims := ctx.Value("claims").(*middleware.Claims)
	var err error

	//首先找到用户根目录
	directory := &biz.UserDirectory{}
	childDir := &biz.UserDirectory{}
	objectPath := ""
	filetype := util.String2StringArrWithSeparate(req.File.Filename, ".", true)
	finalName := fmt.Sprintf("%s.%s", req.File.FileHash, filetype[len(filetype)-1])
	finalStr := ""
	localDid := 0
	if id != 0 {
		err = t.data.db.Model(&biz.UserDirectory{}).Where("user_id =? and father_id=?", id, 0).First(directory).Error
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
			err = t.data.db.Model(&biz.UserDirectory{}).Create(directory).Error
			if err != nil {
				t.log.Error(err)
				return nil, ecode.MYSQL_ERR
			}
		}
		localDid = directory.ID

		//处理子路径
		dir := req.Directory
		if dir != "" {
			dirs := util.String2StringArrWithSeparate(dir, "/", true)
			childName := dirs[len(dirs)-1]
			if len(dirs) == 1 { //如果是根目录下的文件夹
				err = t.data.db.Model(&biz.UserDirectory{}).Where("name = ? and user_id = ?", childName, id).First(childDir).Error
				if err != nil && err.Error() != "record not found" {
					t.log.Error(err)
					return nil, ecode.MYSQL_ERR
				}
				if childDir.ID == 0 {
					childDir.UserId = id
					childDir.FatherId = directory.ID
					childDir.Name = childName
					childDir.PathStr = fmt.Sprintf("%s%s/", directory.PathStr, childName)

					err = t.data.db.Model(&biz.UserDirectory{}).Create(childDir).Error
					if err != nil {
						t.log.Error(err)
						return nil, ecode.MYSQL_ERR
					}
				}

			} else {
				err = t.data.db.Model(&biz.UserDirectory{}).Where("name = ? and user_id = ?", childName, id).First(childDir).Error
				if err != nil && err.Error() != "record not found" {
					t.log.Error(err)
					return nil, ecode.MYSQL_ERR
				}
				if childDir.ID == 0 {
					lastDir := &biz.UserDirectory{}
					err = t.data.db.Model(&biz.UserDirectory{}).Where("name = ? and user_id = ?", dirs[len(dirs)-2], id).First(lastDir).Error
					if err != nil && err.Error() != "record not found" {
						t.log.Error(err)
						return nil, ecode.MYSQL_ERR
					}
					if lastDir.ID == 0 {
						return nil, ecode.New(500).SetMessage("路径不存在")
					}
					childDir.UserId = id
					childDir.FatherId = lastDir.ID
					childDir.Name = childName
					childDir.PathStr = fmt.Sprintf("%s%s/", lastDir.PathStr, childName)
					err = t.data.db.Model(&biz.UserDirectory{}).Create(&childDir).Error
					if err != nil {
						t.log.Error(err)
						return nil, ecode.MYSQL_ERR
					}
				}
			}
		}
		if childDir.Name != "" {
			objectPath = childDir.PathStr
			localDid = childDir.ID
		} else {
			objectPath = directory.PathStr
		}
		objectName := fmt.Sprintf("%s.%s", req.File.FileHash, filetype[len(filetype)-1])
		finalName = fmt.Sprintf("%s%s", objectPath, objectName)
		finalStr = fmt.Sprintf("%s%s", objectPath, req.File.Filename)
	}

	filepath := fmt.Sprintf("./%s", req.File.Filename)
	fileinfo := minio.UploadInfo{}
	attribute := 0
	var minioUpload = func(bucket, objectName string) (minio.UploadInfo, error) {
		client := t.data.minio
		fileinfo, err = client.FPutObject(ctx, bucket, objectName, filepath, minio.PutObjectOptions{ContentType: req.File.ContentType})
		if err != nil {
			t.log.Error(err)
			return fileinfo, err
		}
		return fileinfo, err
	}
	fid := 0
	if id == 0 { //游客用户入库
		fileinfo, err = minioUpload(PUBLIC, finalName)
		if err != nil {
			t.log.Error(err)
			return nil, ecode.New(500).SetMessage("minio客户端错误")
		}
		attribute = 2
		file := &biz.File{
			FileName:    req.File.Filename,
			FileHash:    req.File.FileHash,
			FilePath:    finalName,
			FileSize:    fileinfo.Size,
			FileStr:     "",
			Attribute:   int8(attribute),
			Suffix:      filetype[len(filetype)-1],
			ContentType: req.File.ContentType,
		}
		err = t.data.db.Model(&biz.File{}).WithContext(ctx).Create(file).Error
		if err != nil {
			t.log.Error(err)
			return nil, ecode.MYSQL_ERR
		}

		userFile := &biz.UserFile{
			FileId:  file.ID,
			UserNum: claims.UserNum,
			UserId:  claims.UserId,
		}
		err = t.data.db.Model(&biz.UserFile{}).WithContext(ctx).Create(userFile).Error
		if err != nil {
			t.log.Error(err)
			return nil, ecode.MYSQL_ERR
		}
		fid = file.ID
	} else { //私人文件入库
		fileinfo, err = minioUpload(PRIVATE, finalName)
		if err != nil {
			t.log.Error(err)
			return nil, ecode.New(500).SetMessage("minio客户端错误")
		}
		attribute = 1

		file := &biz.File{
			FileName:    req.File.Filename,
			FileHash:    req.File.FileHash,
			FilePath:    finalName,
			FileSize:    fileinfo.Size,
			FileStr:     finalStr,
			Attribute:   int8(attribute),
			Suffix:      filetype[len(filetype)-1],
			ContentType: req.File.ContentType,
			DirectoryId: localDid,
		}
		err = t.data.db.Model(&biz.File{}).WithContext(ctx).Create(file).Error
		if err != nil {
			t.log.Error(err)
			return nil, ecode.MYSQL_ERR
		}
		//文件存量统计
		dirs := []*biz.UserDirectory{}
		err = t.data.db.Model(&biz.UserDirectory{}).Where("path_str like ? and user_id = ?", "%"+directory.Name+"%", id).Find(&dirs).Error
		if err != nil {
			t.log.Error(err)
			return nil, ecode.MYSQL_ERR
		}
		for _, v := range dirs {
			v.Size += fileinfo.Size
		}
		err = t.data.db.Model(&biz.UserDirectory{}).Save(&dirs).Error
		if err != nil {
			t.log.Error(err)
			return nil, ecode.MYSQL_ERR
		}
		did := 0
		if childDir.ID == 0 {
			did = directory.ID
		} else {
			did = childDir.ID
		}
		userFile := &biz.UserFile{
			FileId:      file.ID,
			UserNum:     claims.UserNum,
			UserId:      claims.UserId,
			DirectoryId: did,
		}
		err = t.data.db.Model(&biz.UserFile{}).WithContext(ctx).Create(userFile).Error
		if err != nil {
			t.log.Error(err)
			return nil, ecode.MYSQL_ERR
		}
		fid = file.ID
	}
	res := &pb.RespUpload{Fid: int32(fid)}
	return res, nil
}

func (t *transferRepo) DownloadEntry(ctx context.Context, req *pb.ReqDownload) (*pb.RespDownload, error) {
	res := &pb.RespDownload{}
	var err error
	//id := ctx.Value("x-md-global-uid").(int)
	//claims:=ctx.Value("claims").(*middleware.Claims)
	userFile := &biz.UserFile{}
	err = t.data.db.Model(&biz.UserFile{}).Where("file_id = ?", req.Fid).First(userFile).Error
	if err != nil {
		return res, ecode.MYSQL_ERR
	}
	File := &biz.File{}
	err = t.data.db.Model(&biz.File{}).Where("id = ?", req.Fid).First(File).Error
	if err != nil {
		return res, ecode.MYSQL_ERR
	}
	//attribute = 1 私人
	bucketName := PUBLIC
	if File.Attribute == 1 {
		bucketName = PRIVATE
	}

	client := t.data.minio

	object, err := client.GetObject(ctx, bucketName, File.FilePath, minio.GetObjectOptions{})
	if err != nil {
		return res, ecode.New(500).SetMessage("minio客户端错误")
	}
	localfile, err := os.Create(File.FileName)
	if err != nil {
		return res, ecode.New(500).SetMessage("创建文件错误")
	}
	if _, err = io.Copy(localfile, object); err != nil {
		return res, ecode.New(500).SetMessage("创建文件错误")
	}

	File.DownloadCount += 1
	err = t.data.db.Model(&biz.File{}).Where("id = ?", File.ID).Save(File).Error
	if err != nil {
		return res, ecode.MYSQL_ERR
	}
	res.Message = "success"
	res.Status = true
	res.Filename = File.FileName
	return res, nil
}

func (t *transferRepo) UploadStatic(ctx context.Context, req *pb.ReqStatic) (*pb.RespStatic, error) {
	res := &pb.RespStatic{}
	filepath := fmt.Sprintf("./%s", req.Filename)
	objectName := fmt.Sprintf("/%s", req.Filename)
	client := t.data.minio
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
	res.FileAddress = fmt.Sprintf("http://47.107.95.82:9000/peach-static/%s", req.Filename)
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
	if req.DirectoryId == 0 {
		err = t.data.db.Model(&biz.UserDirectory{}).Where("user_id = ? and name = ?", id, claims.UserNum).First(root).Error
		if err != nil {
			return res, ecode.MYSQL_ERR.SetMessage(err.Error())
		}
	} else {
		err = t.data.db.Model(&biz.UserDirectory{}).Where("id = ?", req.DirectoryId).First(root).Error
		if err != nil {
			return res, ecode.MYSQL_ERR.SetMessage(err.Error())
		}
	}
	files := []*biz.File{}
	if root.ID != 0 {
		if keyWords == "" {
			err = t.data.db.Model(&biz.File{}).Where("directory_id = ? and file_status = ?", root.ID, EXIST).Find(&files).Error
			if err != nil {
				return res, ecode.MYSQL_ERR.SetMessage(err.Error())
			}
			err = t.data.db.Model(&biz.UserDirectory{}).Where("father_id = ? and dir_status = ?", root.ID, EXIST).Find(&childDir).Error
			if err != nil && err.Error() != "record not found" {
				return res, ecode.MYSQL_ERR.SetMessage(err.Error())
			}
		} else {
			err = t.data.db.Model(&biz.File{}).Where("directory_id = ? and file_name like ? and file_status = ?", root.ID, "%"+keyWords+"%", EXIST).Find(&files).Error
			if err != nil {
				return res, ecode.MYSQL_ERR.SetMessage(err.Error())
			}
			err = t.data.db.Model(&biz.UserDirectory{}).Where("father_id = ? and name like ? and dir_status = ?", root.ID, "%"+keyWords+"%", EXIST).Find(&childDir).Error
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
	for _, v := range files {
		meta := &pb.FileMetaObject{
			Fid:          int32(v.ID),
			Size:         util.FormatFileSize(v.FileSize),
			FileName:     v.FileName,
			FileType:     v.Suffix,
			LastModified: v.UpdatedAt,
		}
		FileMeta = append(FileMeta, meta)
	}

	DirMeta := []*pb.DirMetaObject{}
	for _, v := range childDir {
		meta := &pb.DirMetaObject{
			Did:          int32(v.ID),
			Size:         util.FormatFileSize(v.Size),
			DirName:      v.Name,
			LastModified: v.UpdatedAt,
		}
		DirMeta = append(DirMeta, meta)
	}
	res.UserId = int32(id)
	res.Total = int32(len(FileMeta) + len(DirMeta))
	res.FileObject = FileMeta
	res.DirObject = DirMeta

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
	err = t.data.db.Model(&biz.UserDirectory{}).Where("user_id = ? and dir_status = ?", id, DEL).Find(&File).Error
	if err != nil && err.Error() != "record not found" {
		return res, ecode.MYSQL_ERR.SetMessage(err.Error())
	}
	err = t.data.db.Model(&biz.UserDirectory{}).Where("user_id = ? and file_status = ?",id,DEL).Find(&Dir).Error
	if err != nil && err.Error() != "record not found" {
		return res, ecode.MYSQL_ERR.SetMessage(err.Error())
	}

	//排序两参同传
	if sortName != 0 && sortType != 0 {
		switch sortName {
		case NAME:
			SortFile(File, func(p, q *biz.File) bool {
				if sortType == ASC {
					return p.FileName > q.FileName
				}
				return p.FileName < q.FileName
			})
			SortDir(Dir, func(p, q *biz.UserDirectory) bool {
				if sortType == ASC {
					return p.Name > q.Name
				}
				return p.Name < q.Name
			})
		case TIME:
			SortFile(File, func(p, q *biz.File) bool {
				if sortType == ASC {
					return p.UpdatedAt > q.UpdatedAt
				}
				return p.UpdatedAt < q.UpdatedAt
			})
			SortDir(Dir, func(p, q *biz.UserDirectory) bool {
				if sortType == ASC {
					return p.UpdatedAt > q.UpdatedAt
				}
				return p.UpdatedAt < q.UpdatedAt
			})
		case SIZE:
			SortFile(File, func(p, q *biz.File) bool {
				if sortType == ASC {
					return p.FileSize > q.FileSize
				}
				return p.FileSize < q.FileSize
			})
			SortDir(Dir, func(p, q *biz.UserDirectory) bool {
				if sortType == ASC {
					return p.Size > q.Size
				}
				return p.Size < q.Size
			})
		}
	}

	FileMeta := []*pb.FileMetaObject{}
	for _, v := range File {
		meta := &pb.FileMetaObject{
			Fid:          int32(v.ID),
			Size:         util.FormatFileSize(v.FileSize),
			FileName:     v.FileName,
			FileType:     v.Suffix,
			LastModified: v.UpdatedAt,
		}
		FileMeta = append(FileMeta, meta)
	}

	DirMeta := []*pb.DirMetaObject{}
	for _, v := range Dir {
		meta := &pb.DirMetaObject{
			Did:          int32(v.ID),
			Size:         util.FormatFileSize(v.Size),
			DirName:      v.Name,
			LastModified: v.UpdatedAt,
		}
		DirMeta = append(DirMeta, meta)
	}
	res.UserId = int32(id)
	res.Total = int32(len(FileMeta) + len(DirMeta))
	res.FileObject = FileMeta
	res.DirObject = DirMeta

	return res, nil
}

func (t *transferRepo) DeleteFile(ctx context.Context, req *pb.ReqDeleteFile) (*pb.RespDelete, error) {
	res := &pb.RespDelete{}
	//id := ctx.Value("x-md-global-uid").(int)
	///claims := ctx.Value("claims").(*middleware.Claims)
	var err error
	files := []*biz.File{}
	if len(req.Fid) != 0 {
		err = t.data.db.Model(&biz.File{}).Where("id in (?)  and file_status = ?", req.Fid, EXIST).Find(&files).Error
		if err != nil {
			return res, ecode.MYSQL_ERR.SetMessage(err.Error())
		}
	}
	client := t.data.minio
	var sizeCut int64
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
		sizeCut += v.FileSize
		v.FileStatus = DEL
	}
	tx := t.data.db.Begin()
	dirs := []*biz.UserDirectory{}
	//没办法了
	father := &biz.UserDirectory{}
	//找到当前目录的文件夹
	err = tx.Model(&biz.UserDirectory{}).Where("id = ?", files[0].DirectoryId).First(father).Error
	fatherId := father.FatherId
	dirs = append(dirs, father)
	for fatherId != 0 {
		father = &biz.UserDirectory{}
		err = t.data.db.Model(&biz.UserDirectory{}).Where("id = ?", fatherId).First(father).Error
		if err != nil {
			tx.Rollback()
			return res, ecode.MYSQL_ERR.SetMessage(err.Error())
		}
		fatherId = father.FatherId
		dirs = append(dirs, father)
	}

	for _, v := range dirs {
		v.Size -= sizeCut
	}
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
		err = t.data.db.Model(&biz.UserDirectory{}).Where("id in (?) user_id= ? and dir_status = ?", req.Did, id, EXIST).Find(&dirs).Error
		if err != nil {
			return res, ecode.MYSQL_ERR.SetMessage(err.Error())
		}
	}
	childDir := []*biz.UserDirectory{}

	for _, v := range dirs {
		err = t.data.db.Raw(
			"select * from  user_directory a "+
				"left join(select path_str from user_directory d where d.path_str like ?)b "+
				"on a.path_str=b.path_str where b.path_str is not null", "%"+v.Name+"/%").Scan(&childDir).Error
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
			err = t.data.db.Model(&biz.File{}).Where("directory_id = ? and file_status = ?", dir.ID, EXIST).Find(&files).Error
			if err != nil && err.Error() != "record not found" {
				return res, ecode.MYSQL_ERR.SetMessage(err.Error())
			}
			fileMap[dir.ID] = files
			dir.DirStatus = DEL
			saveDir = append(saveDir, dir)
			files = []*biz.File{}
		}
	}
	saveFile := []*biz.File{}
	client := t.data.minio
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

	tx := t.data.db.Begin()
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
func (t *transferRepo) WithDrawFile(ctx context.Context, req *pb.ReqWithDrawFile) (*pb.RespWithDraw, error) {
	//id := ctx.Value("x-md-global-uid").(int)
	res := &pb.RespWithDraw{}
	var err error

	files := []*biz.File{}
	if len(req.Fid) != 0 {
		err = t.data.db.Model(&biz.File{}).Where("id in (?)  and file_status = ?", req.Fid, DEL).Find(&files).Error
		if err != nil {
			return res, ecode.MYSQL_ERR.SetMessage(err.Error())
		}
	}

	client := t.data.minio
	var sizeCut int64
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
		sizeCut += v.FileSize
		v.FileStatus = EXIST
	}

	tx := t.data.db
	err = tx.Model(&biz.File{}).Save(files).Error
	if err != nil {
		return res, ecode.MYSQL_ERR.SetMessage(err.Error())
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
		err = t.data.db.Model(&biz.UserDirectory{}).Where("id in (?) user_id= ? and dir_status = ?", req.Did, id, DEL).Find(&dirs).Error
		if err != nil {
			return res, ecode.MYSQL_ERR.SetMessage(err.Error())
		}
	}
	childDir := []*biz.UserDirectory{}

	for _, v := range dirs {
		err = t.data.db.Raw(
			"select * from  user_directory a "+
				"left join(select path_str from user_directory d where d.path_str like ?)b "+
				"on a.path_str=b.path_str where b.path_str is not null", "%"+v.Name+"/%").Scan(&childDir).Error
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
			err = t.data.db.Model(&biz.File{}).Where("directory_id = ? and file_status = ?", dir.ID, DEL).Find(&files).Error
			if err != nil && err.Error() != "record not found" {
				return res, ecode.MYSQL_ERR.SetMessage(err.Error())
			}
			fileMap[dir.ID] = files
			dir.DirStatus = EXIST
			saveDir = append(saveDir, dir)
			files = []*biz.File{}
		}
	}
	saveFile := []*biz.File{}
	client := t.data.minio
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
			file.FileStatus = DEL
			saveFile = append(saveFile, file)
		}
	}

	tx := t.data.db.Begin()
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
func (t *transferRepo) CleanTrashFile(ctx context.Context, req *pb.ReqCleanTrashFile) (*pb.RespCleanTrash, error) {
	res := &pb.RespCleanTrash{}
	var err error
	//id := ctx.Value("x-md-global-uid").(int)

	files := []*biz.File{}
	if len(req.Fid) != 0 {
		err = t.data.db.Model(&biz.File{}).Where("directory_id in (?) and file_status = ?", req.Fid, DEL).Find(&files).Error
		if err != nil {
			return res, ecode.MYSQL_ERR.SetMessage(err.Error())
		}
	}

	client := t.data.minio
	var sizeCut int64
	for _, v := range files {
		ropt := minio.RemoveObjectOptions{
			ForceDelete: true,
		}
		err = client.RemoveObject(ctx, DELETE, v.FilePath, ropt)
		if err != nil {
			return res, ecode.New(500).SetMessage(err.Error())
		}
		sizeCut += v.FileSize
		v.FileStatus = RUIN
	}

	tx := t.data.db
	err = tx.Model(&biz.File{}).Save(&files).Error
	if err != nil {
		return res, ecode.MYSQL_ERR.SetMessage(err.Error())
	}

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
		err = t.data.db.Model(&biz.UserDirectory{}).Where("id in (?) user_id= ? and dir_status = ?", req.Did, id, DEL).Find(&dirs).Error
		if err != nil {
			return res, ecode.MYSQL_ERR.SetMessage(err.Error())
		}
	}
	childDir := []*biz.UserDirectory{}

	for _, v := range dirs {
		err = t.data.db.Raw(
			"select * from  user_directory a "+
				"left join(select path_str from user_directory d where d.path_str like ?)b "+
				"on a.path_str=b.path_str where b.path_str is not null", "%"+v.Name+"/%").Scan(&childDir).Error
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
			err = t.data.db.Model(&biz.File{}).Where("directory_id = ? and file_status = ?", dir.ID, DEL).Find(&files).Error
			if err != nil && err.Error() != "record not found" {
				return res, ecode.MYSQL_ERR.SetMessage(err.Error())
			}
			fileMap[dir.ID] = files
			dir.DirStatus = RUIN
			saveDir = append(saveDir, dir)
			files = []*biz.File{}
		}
	}
	saveFile := []*biz.File{}
	client := t.data.minio
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

	tx := t.data.db.Begin()
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
	return nil, nil
}
func (t *transferRepo) PreviewFile(ctx context.Context, req *pb.ReqPreviewFile) (*pb.RespPreviewFile, error) {
	res := &pb.RespPreviewFile{}
	var err error
	id := ctx.Value("x-md-global-uid").(int)
	if id == 0 {

	}
	userFile := &biz.UserFile{}
	err = t.data.db.Model(&biz.UserFile{}).Where("user_id = ? and file_id = ?", id, req.Fid).First(userFile).Error
	if err != nil {
		return res, ecode.MYSQL_ERR
	}
	if userFile == nil {
		return res, ecode.MYSQL_ERR.SetMessage("找不到此文件")
	}
	File := &biz.File{}
	err = t.data.db.Model(&biz.File{}).Where("id = ?", req.Fid).First(File).Error
	if err != nil {
		return res, ecode.MYSQL_ERR
	}
	client := t.data.minio
	//
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
	err = t.data.db.Model(&biz.UserDirectory{}).Where("user_id = ?",id).Pluck("id",&dirIds).Error
	if err != nil && err.Error()!="record not found"{
		return res, ecode.MYSQL_ERR.SetMessage(err.Error())
	}
	files := []*biz.File{}
	err = t.data.db.Model(&biz.File{}).
		Where("directory_id in (?)",dirIds).
		Order("download_count desc").
		Find(&files).Error
	if err != nil && err.Error()!="record not found"{
		return res, ecode.MYSQL_ERR.SetMessage(err.Error())
	}
	var total int64
	//err = t.data.db.Model(&biz.File{}).Select("sum(file_size)").Error
	topTen := []*pb.DownloadCensus{}
	if len(files) <=10{
		for _,v:=range files{
			download := &pb.DownloadCensus{
				FileName: v.FileName,
				Count:    int32(v.DownloadCount),
			}
			topTen = append(topTen,download)
		}
	} else {
		temp := files[:10]
		for _,v:=range temp{
			download := &pb.DownloadCensus{
				FileName: v.FileName,
				Count:    int32(v.DownloadCount),
			}
			topTen = append(topTen,download)
		}
	}
	countMap := make(map[string]int)
	for _,v := range files {
		total += v.FileSize
		if _,exist := TypeMap[v.Suffix];exist{
			countMap[TypeMap[v.Suffix]] +=1
		} else {
			countMap["other"] +=1
		}
	}
	last := float32(total)/float32(1024*1024*1024*10)
	usage := &pb.Usage{
		UseStr: fmt.Sprintf("%s/%s",util.FormatFileSize(total),util.FormatFileSize(1024*1024*1024*10)),
		Used:   last,
	}
	audio := countMap["audio"]
	video := countMap["video"]
	image := countMap["img"]
	doc := countMap["doc"]
	other := countMap["other"]
	ratio := &pb.FileRatio{
		Audio: int32(audio),
		Video: int32(video),
		Image: int32(image),
		Doc:  int32(doc),
		Other: int32(other),
	}

	res.FileRatio = ratio
	res.Usage = usage
	res.TopTen = topTen
	return res,nil

}

func (t *transferRepo) SearchFile(ctx context.Context, req *pb.ReqSearchFile) (*pb.RespSearchFile, error) {
	return nil, nil
}

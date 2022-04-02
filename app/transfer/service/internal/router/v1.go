package router

import (
	pb "banana/api/transfer/service/v1"
	"banana/app/transfer/service/internal/service"
	"banana/pkg/middleware"
	"banana/pkg/response"
	"banana/pkg/util"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"io"
	"os"
	"strconv"
	"strings"
)

var transferService *service.TransferService

func apiV1(group gin.IRoutes, tf *service.TransferService) {
	transferService = tf
	group.GET("/download", DownloadHandler)
	group.GET("/code-download",GetCodeDownload)
	group.Use(middleware.JWTAuth())
	group.GET("/file-list", GetUserFileList)
	group.POST("/del-file", DeleteFile)
	group.POST("/share", ShareFile)
	group.GET("/preview", PreviewFile)
	group.GET("/census", FileCensus)
	group.GET("/trash-list",GetUserTrashList)
	group.POST("/del-dir",DeleteDirs)
	group.POST("/withdraw-file",WithDrawFile)
	group.POST("/withdraw-dir",WithDrawDir)
	group.POST("/add-dir",CreateDir)
	group.POST("/clean-file",CleanFiles)
	group.POST("/clean-dir",CleanDirs)

}

/// 解析多个文件上传中，每个具体的文件的信息
type FileHeader struct {
	ContentDisposition string
	Name               string
	FileName           string ///< 文件名
	ContentType        string
	ContentLength      int64
}

/// 解析描述文件信息的头部
/// @return FileHeader 文件名等信息的结构体
/// @return bool 解析成功还是失败
func ParseFileHeader(h []byte) (FileHeader, bool) {
	arr := bytes.Split(h, []byte("\r\n"))
	var out_header FileHeader
	out_header.ContentLength = -1
	const (
		CONTENT_DISPOSITION = "Content-Disposition: "
		NAME                = "name=\""
		FILENAME            = "filename=\""
		CONTENT_TYPE        = "Content-Type: "
		CONTENT_LENGTH      = "Content-Length: "
	)
	for _, item := range arr {
		if bytes.HasPrefix(item, []byte(CONTENT_DISPOSITION)) {
			l := len(CONTENT_DISPOSITION)
			arr1 := bytes.Split(item[l:], []byte("; "))
			out_header.ContentDisposition = string(arr1[0])
			if bytes.HasPrefix(arr1[1], []byte(NAME)) {
				out_header.Name = string(arr1[1][len(NAME) : len(arr1[1])-1])
			}
			l = len(arr1[2])
			if bytes.HasPrefix(arr1[2], []byte(FILENAME)) && arr1[2][l-1] == 0x22 {
				out_header.FileName = string(arr1[2][len(FILENAME) : l-1])
			}
		} else if bytes.HasPrefix(item, []byte(CONTENT_TYPE)) {
			l := len(CONTENT_TYPE)
			out_header.ContentType = string(item[l:])
		} else if bytes.HasPrefix(item, []byte(CONTENT_LENGTH)) {
			l := len(CONTENT_LENGTH)
			s := string(item[l:])
			contentLength, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				return out_header, false
			} else {
				out_header.ContentLength = contentLength
			}
		} else {
			util.Println("unknown:%s\n", string(item))
		}
	}
	if len(out_header.FileName) == 0 {
		return out_header, false
	}
	return out_header, true
}

/// 从流中一直读到文件的末位
/// @return []byte 没有写到文件且又属于下一个文件的数据
/// @return bool 是否已经读到流的末位了
/// @return error 是否发生错误
func ReadToBoundary(boundary []byte, stream io.ReadCloser, target io.WriteCloser) ([]byte, bool, error) {
	readData := make([]byte, 1024*8)
	read_data_len := 0
	buf := make([]byte, 1024*4)
	b_len := len(boundary)
	reach_end := false
	for !reach_end {
		read_len, err := stream.Read(buf)
		if err != nil {
			if err != io.EOF && read_len <= 0 {
				return nil, true, err
			}
			reach_end = true
		}
		//todo: 下面这一句很蠢，值得优化
		copy(readData[read_data_len:], buf[:read_len]) //追加到另一块buffer，仅仅只是为了搜索方便
		read_data_len += read_len
		if read_data_len < b_len+4 {
			continue
		}
		loc := bytes.Index(readData[:read_data_len], boundary)
		if loc >= 0 {
			//找到了结束位置
			target.Write(readData[:loc-4])
			return readData[loc:read_data_len], reach_end, nil
		}

		target.Write(readData[:read_data_len-b_len-4])
		copy(readData[0:], readData[read_data_len-b_len-4:])
		read_data_len = b_len + 4
	}
	target.Write(readData[:read_data_len])
	return nil, reach_end, nil
}

/// 解析表单的头部
/// @param readData 已经从流中读到的数据
/// @param readTotal 已经从流中读到的数据长度
/// @param boundary 表单的分割字符串
/// @param stream 输入流
/// @return FileHeader 文件名等信息头
///			[]byte 已经从流中读到的部分
///			error 是否发生错误
func ParseFromHead(readData []byte, readTotal int, boundary []byte, stream io.ReadCloser) (FileHeader, []byte, error) {
	buf := make([]byte, 1024*4)
	found_boundary := false
	boundary_loc := -1
	var fileHeader FileHeader
	for {
		read_len, err := stream.Read(buf)
		if err != nil {
			if err != io.EOF {
				return fileHeader, nil, err
			}
			break
		}
		if readTotal+read_len > cap(readData) {
			return fileHeader, nil, fmt.Errorf("not found boundary")
		}
		copy(readData[readTotal:], buf[:read_len])
		readTotal += read_len
		if !found_boundary {
			boundary_loc = bytes.Index(readData[:readTotal], boundary)
			if -1 == boundary_loc {
				continue
			}
			found_boundary = true
		}
		start_loc := boundary_loc + len(boundary)
		file_head_loc := bytes.Index(readData[start_loc:readTotal], []byte("\r\n\r\n"))
		if -1 == file_head_loc {
			continue
		}
		file_head_loc += start_loc
		ret := false
		fileHeader, ret = ParseFileHeader(readData[start_loc:file_head_loc])
		if !ret {
			return fileHeader, nil, fmt.Errorf("ParseFileHeader fail:%s", string(readData[start_loc:file_head_loc]))
		}
		return fileHeader, readData[file_head_loc+4 : readTotal], nil
	}
	return fileHeader, nil, fmt.Errorf("reach to sream EOF")
}

func GuestUpload (c *gin.Context) {
	var contentLength int64
	contentLength = c.Request.ContentLength
	if contentLength <= 0 || contentLength > 1024*1024*1024*2 {
		util.Println("contentLength error\n")
		response.NewErrWithCodeAndMsg(c, 200, "contentLength error")
		return
	}
	contentType_, has_key := c.Request.Header["Content-Type"]
	if !has_key {
		util.Println("Content-Type error\n")
		response.NewErrWithCodeAndMsg(c, 200, "Content-Type error")
		return
	}
	if len(contentType_) != 1 {
		util.Println("Content-Type count error\\n")
		response.NewErrWithCodeAndMsg(c, 200, "Content-Type count error")
		return
	}
	contentType := contentType_[0]
	const BOUNDARY string = "; boundary="
	loc := strings.Index(contentType, BOUNDARY)
	if -1 == loc {
		util.Println("Content-Type error, no boundary\n")
		response.NewErrWithCodeAndMsg(c, 200, "Content-Type error, no boundary")
		return
	}
	boundary := []byte(contentType[(loc + len(BOUNDARY)):])
	//
	readData := make([]byte, 1024*12)
	var readTotal int = 0
	fileMap := make(map[string]string)
	for {
		fileHeader, fileData, err := ParseFromHead(readData, readTotal, append(boundary, []byte("\r\n")...), c.Request.Body)
		if err != nil {
			util.Println(err)
			response.NewErrWithCodeAndMsg(c, 200, err.Error())
			return
		}
		f, err := os.Create(fileHeader.FileName)
		if err != nil {
			util.Println(err)
			response.NewErrWithCodeAndMsg(c, 200, err.Error())
			return
		}
		f.Write(fileData)
		fileMap[fileHeader.FileName] = fileHeader.ContentType
		fileData = nil
		//需要反复搜索boundary
		temp_data, reach_end, err := ReadToBoundary(boundary, c.Request.Body, f)
		f.Close()

		if err != nil {
			util.Println(err)
			response.NewErrWithCodeAndMsg(c, 200, err.Error())
			return
		}
		if reach_end {
			break
		} else {
			copy(readData[0:], temp_data)
			readTotal = len(temp_data)
			continue
		}
	}



	fid := []int32{}
	for k, v := range fileMap {
		f,err := os.Open(k)
		f.Seek(0, 0)
		filehash := util.FileSha1(f)
		f.Seek(0, 0)
		req := &pb.ReqGuestUpload{
			File: &pb.File{
				Filename: k,
				FileHash: filehash,
				ContentType: v,
			},
		}
		res, err := transferService.GuestUpload(c, req)
		if err != nil {

			response.NewErrWithCodeAndMsg(c, 200, err.Error())
			return
		}
		fid = append(fid,res.Fid)

	}
	//
	response.NewSuccess(c, gin.H{
		"message": "success",
		"fids": fid,
	})
}
func UploadHandler(c *gin.Context) {
	var contentLength int64
	contentLength = c.Request.ContentLength
	if contentLength <= 0 || contentLength > 1024*1024*1024*2 {
		util.Println("contentLength error\n")
		response.NewErrWithCodeAndMsg(c, 200, "contentLength error")
		return
	}
	contentType_, has_key := c.Request.Header["Content-Type"]
	if !has_key {
		util.Println("Content-Type error\n")
		response.NewErrWithCodeAndMsg(c, 200, "Content-Type error")
		return
	}
	if len(contentType_) != 1 {
		util.Println("Content-Type count error\\n")
		response.NewErrWithCodeAndMsg(c, 200, "Content-Type count error")
		return
	}
	contentType := contentType_[0]
	const BOUNDARY string = "; boundary="
	loc := strings.Index(contentType, BOUNDARY)
	if -1 == loc {
		util.Println("Content-Type error, no boundary\n")
		response.NewErrWithCodeAndMsg(c, 200, "Content-Type error, no boundary")
		return
	}
	boundary := []byte(contentType[(loc + len(BOUNDARY)):])
	//
	readData := make([]byte, 1024*12)
	var readTotal int = 0
	fileMap := make(map[string]string)

	fileIds := []int{}
	for {
		fileHeader, fileData, err := ParseFromHead(readData, readTotal, append(boundary, []byte("\r\n")...), c.Request.Body)
		if err != nil {
			util.Println(err)
			response.NewErrWithCodeAndMsg(c, 200, err.Error())
			return
		}
		f, err := os.Create(fileHeader.FileName)
		if err != nil {
			util.Println(err)
			response.NewErrWithCodeAndMsg(c, 200, err.Error())
			return
		}
		f.Write(fileData)
		fileMap[fileHeader.FileName] = fileHeader.ContentType
		fileData = nil
		//需要反复搜索boundary
		temp_data, reach_end, err := ReadToBoundary(boundary, c.Request.Body, f)
		f.Close()

		if err != nil {
			util.Println(err)
			response.NewErrWithCodeAndMsg(c, 200, err.Error())
			return
		}
		if reach_end {
			break
		} else {
			copy(readData[0:], temp_data)
			readTotal = len(temp_data)
			continue
		}
	}

	did := c.Query("did")
	dirId ,_:=strconv.Atoi(did)

	for k, v := range fileMap {
		f,err := os.Open(k)
		f.Seek(0, 0)
		filehash := util.FileSha1(f)
		f.Seek(0, 0)
		info ,_:= f.Stat()
		req := &pb.ReqUpload{
			File: &pb.File{
				Filesize: info.Size(),
				Filename: k,
				FileHash: filehash,
				ContentType: v,
			},
			Did: int32(dirId),
		}
		res, err := transferService.UploadEntry(c, req)
		if err != nil {
			response.NewErrWithCodeAndMsg(c, 200, err.Error())
			return
		}else{
			fileIds = append(fileIds, int(res.Fid))
		}

	}
	//
	response.NewSuccess(c, gin.H{
		"message": "success",
		"fids":    fileIds,
	})
}

func UploadStatic(c *gin.Context) {
	var contentLength int64
	contentLength = c.Request.ContentLength
	if contentLength <= 0 || contentLength > 1024*1024*1024*2 {
		util.Println("contentLength error\n")
		response.NewErrWithCodeAndMsg(c, 200, "contentLength error")
		return
	}
	contentType_, has_key := c.Request.Header["Content-Type"]
	if !has_key {
		util.Println("Content-Type error\n")
		response.NewErrWithCodeAndMsg(c, 200, "Content-Type error")
		return
	}
	if len(contentType_) != 1 {
		util.Println("Content-Type count error\\n")
		response.NewErrWithCodeAndMsg(c, 200, "Content-Type count error")
		return
	}
	contentType := contentType_[0]
	const BOUNDARY string = "; boundary="
	loc := strings.Index(contentType, BOUNDARY)
	if -1 == loc {
		util.Println("Content-Type error, no boundary\n")
		response.NewErrWithCodeAndMsg(c, 200, "Content-Type error, no boundary")
		return
	}
	boundary := []byte(contentType[(loc + len(BOUNDARY)):])
	//
	readData := make([]byte, 1024*12)
	var readTotal int = 0
	fileNameArr := []string{}
	fileType := ""
	for {
		fileHeader, fileData, err := ParseFromHead(readData, readTotal, append(boundary, []byte("\r\n")...), c.Request.Body)
		if err != nil {
			util.Println(err)
			response.NewErrWithCodeAndMsg(c, 200, errors.FromError(err).Message)
			return
		}
		f, err := os.Create(fileHeader.FileName)
		if err != nil {
			util.Println(err)
			response.NewErrWithCodeAndMsg(c, 200, errors.FromError(err).Message)
			return
		}
		f.Write(fileData)

		fileNameArr = append(fileNameArr, fileHeader.FileName)
		fileType = fileHeader.ContentType
		fileData = nil
		//需要反复搜索boundary
		temp_data, reach_end, err := ReadToBoundary(boundary, c.Request.Body, f)
		f.Close()

		if err != nil {
			util.Println(err)
			response.NewErrWithCodeAndMsg(c, 200, errors.FromError(err).Message)
			return
		}
		if reach_end {
			break
		} else {
			copy(readData[0:], temp_data)
			readTotal = len(temp_data)
			continue
		}
	}
	fileStr := ""
	for _, v := range fileNameArr {
		req := &pb.ReqStatic{
			Filename:    v,
			ContentType: fileType,
		}
		res, err := transferService.UploadStatic(c, req)
		if err != nil {
			response.NewErrWithCodeAndMsg(c, 200, errors.FromError(err).Message)
		}
		fileStr = res.FileAddress
		os.Remove(v)
	}

	//
	response.NewSuccess(c, gin.H{
		"message":  "success",
		"file_str": fileStr,
	})
}

func DownloadHandler(c *gin.Context) {
	fidstr := c.Query("fid")
	fid, err := strconv.Atoi(fidstr)
	if err != nil {
		util.Println(err)
		response.NewErrWithCodeAndMsg(c, 200, err.Error())
		return
	}
	req := &pb.ReqDownload{Fid: int32(fid)}

	res, err := transferService.DownLoadEntry(c, req)
	defer os.Remove(res.Filename)
	if err != nil {
		util.Println(err)
		response.NewErrWithCodeAndMsg(c, 200, err.Error())
	}
	c.Header("Content-Type", res.Type)
	c.Header("Content-Disposition", "attachment; filename="+res.Filename)
	//c.Header("Content-Disposition", fmt.Sprintf("attachment; filename*=utf-8''%s", url.QueryEscape(res.Filename)))
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Cache-Control", "no-cache")

	c.File(res.Filename)
	return
}

func PreviewFile(c *gin.Context) {
	req := &pb.ReqPreviewFile{}
	param:= c.Query("fid")
	fid,err :=strconv.Atoi(param)
	if err != nil {
		response.NewErrWithCodeAndMsg(c,200,"传参格式错误")
	}
	req.Fid = int32(fid)
	res,err := transferService.PreviewFile(c,req)
	if err !=nil{
		response.NewErrWithCodeAndMsg(c,200,err.Error())
	}
	response.NewSuccess(c, gin.H{
		"message":  "success",
		"file_str": res.PreviewStr,
	})
}

func WithDrawDir(c *gin.Context){
	req := &pb.ReqWithDrawDir{}
	if err := c.BindJSON(req);err != nil {
		response.NewErrWithCodeAndMsg(c,200,"传参格式错误")
		return
	}
	res ,err := transferService.WithDrawDir(c,req)
	if err != nil {
		response.NewErrWithCodeAndMsg(c,200,err.Error())
		return
	}
	response.NewSuccess(c,res)
}
func WithDrawFile(c *gin.Context){
	req := &pb.ReqWithDrawFile{}
	if err := c.BindJSON(req);err != nil {
		response.NewErrWithCodeAndMsg(c,200,"传参格式错误")
		return
	}
	res ,err := transferService.WithDrawFile(c,req)
	if err != nil {
		response.NewErrWithCodeAndMsg(c,200,err.Error())
		return
	}
	response.NewSuccess(c,res)
}
func GetUserTrashList(c *gin.Context){
	req := &pb.ReqGetUserTrashBin{}
	type ParseForm struct {
		SortObject  int32  `form:"sort_object"` //0:默认排序 1:文件名长短 2:编辑时间 4:文件大小
		SortType    int32  `form:"sort_type"`       //0:系统默认 1:asc升序  2:desc降序
	}
	parse := &ParseForm{}
	if err := c.BindQuery(parse);err != nil {
		response.NewErrWithCodeAndMsg(c,200,"传参格式错误")
		return
	}

	req.SortObject = parse.SortObject
	req.SortType = parse.SortType
	res,err := transferService.GetUserTrashList(c,req)
	if err != nil {
		response.NewErrWithCodeAndMsg(c,200,err.Error())
		return
	}
	response.NewSuccess(c,res)
}
func GetUserFileList(c *gin.Context) {
	req := &pb.ReqGetUserFileTree{}
	type ParseForm struct {
		SortObject  int32  `form:"sort_object"` //0:默认排序 1:文件名长短 2:编辑时间 4:文件大小
		SortType    int32  `form:"sort_type"`       //0:系统默认 1:asc升序  2:desc降序
		Keywords    string `form:"keywords"`                        //搜索关键字
		DirectoryId int32  `form:"directory_id"`
	}
	parse := &ParseForm{}
	if err := c.ShouldBindQuery(parse);err != nil {
		response.NewErrWithCodeAndMsg(c,200,"传参格式错误")
		return
	}

	req.Keywords = parse.Keywords
	req.SortObject = parse.SortObject
	req.SortType = parse.SortType
	req.DirectoryId = parse.DirectoryId
	res ,err := transferService.GetUserFileTree(c,req)
	if err != nil {
		response.NewErrWithCodeAndMsg(c,200,err.Error())
		return
	}
	response.NewSuccess(c,res)
}

func DeleteDirs(c *gin.Context){
	req := &pb.ReqDeleteDir{}
	if err := c.BindJSON(req);err != nil {
		response.NewErrWithCodeAndMsg(c,200,"传参格式错误")
		return
	}
	res,err := transferService.DeleteDir(c,req)
	if err !=nil{
		response.NewErrWithCodeAndMsg(c,200,err.Error())
		return
	}

	response.NewSuccess(c,res)


}
func DeleteFile(c *gin.Context) {
	req := &pb.ReqDeleteFile{}
	if err := c.BindJSON(req);err != nil {
		response.NewErrWithCodeAndMsg(c,200,"传参格式错误")
		return
	}
	res,err := transferService.DeleteFile(c,req)
	if err !=nil{
		response.NewErrWithCodeAndMsg(c,200,err.Error())
		return
	}
	response.NewSuccess(c, gin.H{
		"message": res.Message,
	})
}
func CleanDirs(c *gin.Context){
	req := &pb.ReqCleanTrashDir{}
	if err := c.ShouldBindJSON(req);err != nil {
		response.NewErrWithCodeAndMsg(c,200,"传参格式错误")
		return
	}
	res,err := transferService.CleanTrashDir(c,req)
	if err != nil {
		response.NewErrWithCodeAndMsg(c,200,err.Error())
		return
	}
	response.NewSuccess(c,res)
}
func CleanFiles(c *gin.Context){
	req := &pb.ReqCleanTrashFile{}
	if err := c.ShouldBindJSON(req);err != nil {
		response.NewErrWithCodeAndMsg(c,200,"传参格式错误")
		return
	}
	res,err := transferService.CleanTrashFile(c,req)
	if err != nil {
		response.NewErrWithCodeAndMsg(c,200,err.Error())
		return
	}
	response.NewSuccess(c,res)
}
func ShareFile(c *gin.Context) {
	req := &pb.ReqShareFileStr{}
	if err := c.ShouldBindJSON(req);err != nil {
		response.NewErrWithCodeAndMsg(c,200,"传参格式错误")
		return
	}
	res,err := transferService.ShareFile(c,req)
	if err != nil {
		response.NewErrWithCodeAndMsg(c,200,err.Error())
		return
	}
	response.NewSuccess(c,res)
}
func FileCensus(c *gin.Context) {
	req := &pb.ReqFileCensus{}
	res,err := transferService.FileCensus(c,req)
	if err != nil {
		response.NewErrWithCodeAndMsg(c,200,err.Error())
		return
	}
	response.NewSuccess(c,res)

}
func CreateDir(c *gin.Context){
	req := &pb.ReqCreateDir{}
	if err := c.ShouldBindJSON(req);err != nil {
		response.NewErrWithCodeAndMsg(c,200,"传参格式错误")
		return
	}
	res,err := transferService.CreateDir(c,req)
	if err != nil {
		response.NewErrWithCodeAndMsg(c,200,err.Error())
		return
	}
	response.NewSuccess(c,res)
}

func GetCodeDownload(c *gin.Context){
	req := &pb.ReqGetCodeDownLoad{}
	type params struct {
		GetCode string `form:"get_code"`
	}
	query := &params{}
	if err := c.ShouldBindQuery(query);err != nil {
		response.NewErrWithCodeAndMsg(c,200,"传参格式错误")
		return
	}
	req.GetCode = query.GetCode
	res,err := transferService.GetCodeDownload(c,req)
	if err != nil {
		response.NewErrWithCodeAndMsg(c,200,err.Error())
		return
	}
	response.NewSuccess(c,res)
}
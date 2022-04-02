package data

import (
	pb "banana/api/common/service/v1"
	"banana/app/common/service/internal/biz"
	"banana/pkg/ecode"
	"banana/pkg/middleware"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"strconv"
	"time"
)

var _ biz.CommonRepo = (*commonRepo)(nil)

type commonRepo struct {
	data *Data
	log  *log.Helper
}

// NewAccountCenterRepo .
func NewCommonRepo(data *Data, logger log.Logger) biz.CommonRepo {
	return &commonRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "data/cm")),
	}
}

const (
	GLOBAL  = 1
	PRIVATE = 2

	TRUE  = 1
	FALSE = 0
)

func (c *commonRepo) CreateNotify(ctx context.Context, rq *pb.ReqCreateNotify) (*pb.RespCreateNotify, error) {
	res := &pb.RespCreateNotify{}
	var err error
	if claims, exist := ctx.Value("claims").(*middleware.Claims); !exist {
		return res, ecode.EXTERNAL_API_NO_RESPONSE.SetMessage("断言失败")
	} else if claims.UserRole != int8(127) {
		return res, ecode.EXTERNAL_API_NO_RESPONSE.SetMessage("非管理员禁止操作")
	}
	tids := []int{}
	tidStr := ""
	if len(rq.NotifyType) != 0 {
		err = c.data.db.Model(&biz.TType{}).Where("id in (?)", rq.NotifyType).Pluck("id", &tids).Error
		if err != nil {
			return res, ecode.MYSQL_ERR.SetMessage(err.Error())
		}
		for _, v := range tids {
			tidStr += strconv.Itoa(v)
			tidStr += "/"
		}
	}
	nModel := &biz.TNotify{}
	if rq.Uid != 0 {
		nModel.IsGlobal = FALSE
		nModel.SendTime = rq.SendTime
		nModel.Type = tidStr
		nModel.Name = rq.Title
		nModel.Body = rq.Body
		err = c.data.db.Model(&biz.TNotify{}).Create(nModel).Error
		if err != nil {
			return res, ecode.MYSQL_ERR.SetMessage(err.Error())
		}
		unModel := &biz.TUserNotify{}
		unModel.Nid = nModel.ID
		unModel.Uid = ctx.Value("x-md-global-uid").(int)
		err = c.data.db.Model(&biz.TUserNotify{}).Create(unModel).Error
		if err != nil {
			return res, ecode.MYSQL_ERR.SetMessage(err.Error())
		}
	} else {
		nModel.IsGlobal = TRUE
		nModel.Type = tidStr
		nModel.SendTime = rq.SendTime
		nModel.Name = rq.Title
		nModel.Body = rq.Body
		err = c.data.db.Model(&biz.TNotify{}).Create(nModel).Error
		if err != nil {
			return res, ecode.MYSQL_ERR.SetMessage(err.Error())
		}
	}
	res.Nid = int32(nModel.ID)
	return res, err

}

func (c *commonRepo) DeleteNotify(ctx context.Context, rq *pb.ReqDeleteNotify) (*pb.RespDeleteNotify, error) {
	res := &pb.RespDeleteNotify{}
	var err error
	if claims, exist := ctx.Value("claims").(*middleware.Claims); !exist {
		return res, ecode.EXTERNAL_API_NO_RESPONSE.SetMessage("断言失败")
	} else if claims.UserRole != int8(127) {
		return res, ecode.EXTERNAL_API_NO_RESPONSE.SetMessage("非管理员禁止操作")
	}
	err = c.data.db.Model(&biz.TNotify{}).Where("id in (?)", rq.Nids).Update("status", DELETE).Error
	if err != nil {
		return res, ecode.MYSQL_ERR.SetMessage(err.Error())
	}
	res.Status = true
	return res, err
}
func (c *commonRepo) GetNotifyList(ctx context.Context, rq *pb.ReqGetNotifyList) (*pb.RespGetNotifyList, error) {
	res := &pb.RespGetNotifyList{}
	var err error
	if _, exist := ctx.Value("claims").(*middleware.Claims); !exist {
		return res, ecode.EXTERNAL_API_NO_RESPONSE.SetMessage("断言失败")
	}
	//else if claims.UserRole != int8(127) || claims.UserId != int(rq.Uid) {
	//	return res, ecode.EXTERNAL_API_NO_RESPONSE.SetMessage("非管理员禁止操作")
	//}
	claims := ctx.Value("claims").(*middleware.Claims)
	tModels := []*biz.TNotify{}
	if claims.UserRole == int8(127) && rq.Uid == 0 {
		err = c.data.db.Model(&biz.TNotify{}).Where("status = ?", EXIST).Offset(int(rq.Offset)).Limit(int(rq.Limit)).Order("updated_at desc").Find(&tModels).Error
		if err != nil {
			return res, ecode.MYSQL_ERR.SetMessage(err.Error())
		}
	}else if rq.Uid!=0{
		nids := []int{}
		err = c.data.db.Model(&biz.TUserNotify{}).Where("uid = ?",rq.Uid).Pluck("nid",nids).Error
		if err != nil {
			return res, ecode.MYSQL_ERR.SetMessage(err.Error())
		}
		err = c.data.db.Model(&biz.TNotify{}).Where("status = ? and is_global = ?", EXIST, TRUE).Offset(int(rq.Offset)).Limit(int(rq.Limit)).Order("updated_at desc").Find(&tModels).Error
		if err != nil {
			return res, ecode.MYSQL_ERR.SetMessage(err.Error())
		}
		personal := []*biz.TNotify{}
		err = c.data.db.Model(&biz.TNotify{}).Where("status = ? and id in(?)",EXIST,nids).Find(personal).Error
		if err != nil {
			return res, ecode.MYSQL_ERR.SetMessage(err.Error())
		}
		tModels = append(tModels,personal...)
	}
	nRes := []*pb.NotifyObject{}
	for _, v := range tModels {
		notifyObject := &pb.NotifyObject{
			Title:    v.Name,
			Body:     v.Body,
			Type:     v.Type,
			SendTime: v.SendTime,
			Nid:      int32(v.ID),
		}
		nRes = append(nRes, notifyObject)
	}
	res.NotifyObjects = nRes
	return res, err
}
func (c *commonRepo) GetNotify(ctx context.Context, rq *pb.ReqGetNotifyObject) (*pb.RespGetNotifyObject, error) {
	res := &pb.RespGetNotifyObject{}
	var err error
	if claims, exist := ctx.Value("claims").(*middleware.Claims); !exist {
		return res, ecode.EXTERNAL_API_NO_RESPONSE.SetMessage("断言失败")
	} else if claims.UserRole != int8(127) && claims.UserId != ctx.Value("x-md-global-uid").(int) {
		return res, ecode.EXTERNAL_API_NO_RESPONSE.SetMessage("非管理员禁止操作")
	}
	tModel := &biz.TNotify{}
	err = c.data.db.Model(&biz.TNotify{}).Where("id = ?", rq.Nid).
		Where("send_time < ?", time.Now().Unix()).First(tModel).Error
	if err != nil {
		return res, ecode.MYSQL_ERR.SetMessage(err.Error())
	}
	nModel := &pb.NotifyObject{
		Title:    tModel.Name,
		Body:     tModel.Body,
		Type:     tModel.Type,
		SendTime: tModel.SendTime,
		Nid:      int32(tModel.ID),
	}

	res.Notify = nModel
	return res, err
}
func (c *commonRepo) CreateNType(ctx context.Context, rq *pb.ReqCreateNotifyType) (*pb.RespCreateNotifyType, error) {
	res := &pb.RespCreateNotifyType{}
	var err error
	if claims, exist := ctx.Value("claims").(*middleware.Claims); !exist {
		return res, ecode.EXTERNAL_API_NO_RESPONSE.SetMessage("断言失败")
	} else if claims.UserRole != int8(127) {
		return res, ecode.EXTERNAL_API_NO_RESPONSE.SetMessage("非管理员禁止操作")
	}
	checkModel := &biz.TType{}
	err = c.data.db.Model(&biz.TType{}).Where("name = ?",rq.Name).First(checkModel).Error
	if err != nil && err.Error() !="record not found"{
		return res, ecode.MYSQL_ERR.SetMessage(err.Error())
	}
	if checkModel.ID!=0{
		return res,ecode.EXTERNAL_API_FAIL.SetMessage("该标签已创建")
	}
	tModel := &biz.TType{
		Name:     rq.Name,
		Describe: rq.Describe,
	}
	err = c.data.db.Model(&biz.TType{}).Create(tModel).Error
	if err != nil {
		return res, ecode.MYSQL_ERR.SetMessage(err.Error())
	}

	res.Tid = int32(tModel.ID)
	return res, err
}
func (c *commonRepo) UpdateNType(ctx context.Context, rq *pb.ReqUpdateNotifyType) (*pb.RespUpdateNotifyType, error) {
	res := &pb.RespUpdateNotifyType{}
	var err error
	if claims, exist := ctx.Value("claims").(*middleware.Claims); !exist {
		return res, ecode.EXTERNAL_API_NO_RESPONSE.SetMessage("断言失败")
	} else if claims.UserRole != int8(127) {
		return res, ecode.EXTERNAL_API_NO_RESPONSE.SetMessage("非管理员禁止操作")
	}
	tModel := &biz.TType{}
	err = c.data.db.Model(&biz.TType{}).Where("id = ?", rq.Tid).First(tModel).Error
	if err != nil {
		return res, ecode.MYSQL_ERR.SetMessage(err.Error())
	}

	tModel.Name = rq.Name
	tModel.Describe = rq.Describe
	err = c.data.db.Model(&biz.TType{}).Where("id = ?", rq.Tid).Save(tModel).Error
	if err != nil {
		return res, ecode.MYSQL_ERR.SetMessage(err.Error())
	}

	res.Tid = int32(tModel.ID)
	return res, err
}
func (c *commonRepo) DeleteNType(ctx context.Context, rq *pb.ReqDeleteNotifyType) (*pb.RespDeleteNotifyType, error) {
	res := &pb.RespDeleteNotifyType{}
	var err error
	if claims, exist := ctx.Value("claims").(*middleware.Claims); !exist {
		return res, ecode.EXTERNAL_API_NO_RESPONSE.SetMessage("断言失败")
	} else if claims.UserRole != int8(127) {
		return res, ecode.EXTERNAL_API_NO_RESPONSE.SetMessage("非管理员禁止操作")
	}

	err = c.data.db.Delete(&biz.TType{}).Where("id in (?)", rq.Tids).Error
	if err != nil {
		return res, ecode.MYSQL_ERR.SetMessage(err.Error())
	}

	res.Status = true
	return res, err
}
func (c *commonRepo) GetNTypeList(ctx context.Context, rq *pb.ReqGetNotifyTypeList) (*pb.RespGetNotifyTypeList, error) {
	res := &pb.RespGetNotifyTypeList{}
	var err error
	if _, exist := ctx.Value("claims").(*middleware.Claims); !exist {
		return res, ecode.EXTERNAL_API_NO_RESPONSE.SetMessage("断言失败")
	}
	//else if claims.UserRole != int8(127) {
	//	return res, ecode.AUTH_FAIL.SetMessage("非管理员禁止操作")
	//}
	tModel := []*biz.TType{}
	err = c.data.db.Model(&biz.TType{}).Find(&tModel).Error
	if err != nil {
		return res, ecode.MYSQL_ERR.SetMessage(err.Error())
	}
	tRes := []*pb.TypeObject{}
	for _, v := range tModel {
		temp := &pb.TypeObject{
			Tid:      int32(v.ID),
			Name:     v.Name,
			Describe: v.Describe,
			UpdateAt: v.UpdatedAt,
		}
		tRes = append(tRes, temp)
	}

	res.TypeObjects = tRes
	return res, err
}

const (
	IMPORTANT = 1
	MIDDLE    = 2
	COMMON    = 3

	EXIST  = 1
	DELETE = 2
)

func (c *commonRepo) CreateAdv(ctx context.Context, rq *pb.ReqCreateAdv) (*pb.RespCreateAdv, error) {
	res := &pb.RespCreateAdv{}
	var err error
	if claims, exist := ctx.Value("claims").(*middleware.Claims); !exist {
		return res, ecode.EXTERNAL_API_NO_RESPONSE.SetMessage("断言失败")
	} else if claims.UserRole != int8(127) {
		return res, ecode.AUTH_FAIL.SetMessage("非管理员禁止操作")
	}
	var desc int8 = 0
	if rq.Describe != 0 {
		desc = COMMON
	}
	tAdv := &biz.TAdv{
		Name:     rq.Name,
		Describe: desc,
		LinkUrl:  rq.LinkUrl,
		ImageUrl: rq.ImageUrl,
		Status:   EXIST,
	}

	err = c.data.db.Model(&biz.TAdv{}).Create(tAdv).Error
	if err != nil {
		return res, ecode.MYSQL_ERR.SetMessage(err.Error())
	}
	res.Aid = int32(tAdv.ID)
	return res, nil
}

func (c *commonRepo) DeleteAdv(ctx context.Context, rq *pb.ReqDeleteAdv) (*pb.RespDeleteAdv, error) {
	res := &pb.RespDeleteAdv{}
	var err error
	if claims, exist := ctx.Value("claims").(*middleware.Claims); !exist {
		return res, ecode.EXTERNAL_API_NO_RESPONSE.SetMessage("断言失败")
	} else if claims.UserRole != int8(127) {
		return res, ecode.AUTH_FAIL.SetMessage("非管理员禁止操作")
	}

	if rq.Aid == 0 {
		return res, ecode.INVALID_PARAM.SetMessage("输入id为0")
	}

	err = c.data.db.Model(&biz.TAdv{}).Where("id = ?", rq.Aid).Update("status", DELETE).Error
	if err != nil {
		return res, ecode.MYSQL_ERR.SetMessage(err.Error())
	}
	res.Aid = rq.Aid
	return res, nil
}

func (c *commonRepo) UpdateAdv(ctx context.Context, rq *pb.ReqUpdateAdv) (*pb.RespUpdateAdv, error) {
	res := &pb.RespUpdateAdv{}
	var err error
	if claims, exist := ctx.Value("claims").(*middleware.Claims); !exist {
		return res, ecode.EXTERNAL_API_NO_RESPONSE.SetMessage("断言失败")
	} else if claims.UserRole != int8(127) {
		return res, ecode.AUTH_FAIL.SetMessage("非管理员禁止操作")
	}
	tAdv := &biz.TAdv{}
	err = c.data.db.Model(&biz.TAdv{}).Where("id = ?", rq.Aid).First(tAdv).Error
	if err != nil {
		return res, ecode.MYSQL_ERR.SetMessage(err.Error())
	}

	tAdv.Name = rq.Name
	tAdv.Describe = int8(rq.Describe)
	tAdv.ImageUrl = rq.ImageUrl
	tAdv.LinkUrl = rq.LinkUrl

	err = c.data.db.Model(&biz.TAdv{}).Save(tAdv).Error
	if err != nil {
		return res, ecode.MYSQL_ERR.SetMessage(err.Error())
	}
	res.Aid = int32(tAdv.ID)
	return res, nil
}

func (c *commonRepo) GetAdvList(ctx context.Context, rq *pb.ReqGetAdvList) (*pb.RespGetAdvList, error) {
	res := &pb.RespGetAdvList{}
	var err error
	//tAdvs := []*biz.TAdv{}
	advObjects := []*pb.AdvObject{}
	if rq.Aid != 0 {
		err = c.data.db.Model(&biz.TAdv{}).Where("id = ?", rq.Aid).Find(&advObjects).Error
		if err != nil {
			return res, ecode.MYSQL_ERR.SetMessage(err.Error())
		}
	} else {
		err = c.data.db.Model(&biz.TAdv{}).Where("status = ?", EXIST).Find(&advObjects).Error
		if err != nil {
			return res, ecode.MYSQL_ERR.SetMessage(err.Error())
		}
	}
	res.AdvObject = advObjects
	return res, nil
}

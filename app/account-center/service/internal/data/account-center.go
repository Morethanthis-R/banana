package data

import (
	pb "banana/api/account-center/service/v1"
	"banana/app/account-center/service/internal/biz"
	mailUtils "banana/app/account-center/service/internal/pkg/mail"
	"banana/pkg/ecode"
	"banana/pkg/middleware"
	"banana/pkg/util"
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var _ biz.AccountCenterRepo = (*accountCenterRepo)(nil)

type accountCenterRepo struct {
	data *Data
	log  *log.Helper
}

// NewAccountCenterRepo .
func NewAccountCenterRepo(data *Data, logger log.Logger) biz.AccountCenterRepo {
	return &accountCenterRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "data/ac")),
	}
}
func (a *accountCenterRepo) SetAdmin(ctx context.Context, b *pb.SetAdminRequest) (*pb.SetAdminReply, error) {
	res := &pb.SetAdminReply{}
	uid := 0
	if claims, exist := ctx.Value("claims").(*middleware.Claims); !exist {
		return res, ecode.EXTERNAL_API_NO_RESPONSE.SetMessage("断言失败")
	} else if claims.UserRole != int8(127) {
		return res, ecode.EXTERNAL_API_NO_RESPONSE.SetMessage("非管理员禁止操作")
	} else {
		uid = claims.UserId
	}
	if b.Uid == int32(uid) {
		return res, ecode.EXTERNAL_API_NO_RESPONSE.SetMessage("无法操作")
	}
	err := a.data.db.Model(&biz.User{}).Where("id = ?").Update("role", b.Role).Error
	if err != nil {
		return res, ecode.MYSQL_ERR.SetMessage(err.Error())
	}
	res.Status = true
	res.Message = "success"
	return res, nil
}

func (a *accountCenterRepo) ResetPass(ctx context.Context, rq *pb.PasswordResetRequest) (*pb.PasswordResetReply, error) {
	res := &pb.PasswordResetReply{}
	u := &biz.User{}
	if claims, exist := ctx.Value("claims").(*middleware.Claims); !exist {
		return res, ecode.EXTERNAL_API_NO_RESPONSE.SetMessage("断言失败")
	} else if claims.UserId != int(rq.Id) || claims.UserRole != 127 {
		return res, ecode.EXTERNAL_API_NO_RESPONSE.SetMessage("非管理员禁止操作")
	}
	err := a.data.db.Model(&biz.User{}).Where("id = ?", rq.Id).First(u).Error

	if !util.CheckPasswordHash(rq.OldPass, u.Password) {
		return res, ecode.AUTH_FAIL.SetMessage("原密码错误")
	}
	key := fmt.Sprintf("%s:%s:%s", "email", "reset",u.Email)
	redisCli := a.data.cache
	code, err := redisCli.Get(ctx, key).Result()
	if err != nil {
		return res, ecode.REDIS_ERR
	}
	if strings.Compare(rq.GetValidate(), code) != 0 {
		return res, ecode.PERMISSION_DENIED.SetMessage("验证码校验不符")
	}

	newPass, _ := util.HashPassword(rq.GetNewPass())
	u.Password = newPass
	err = a.data.db.Model(&biz.User{}).WithContext(ctx).Where("id = ?", rq.Id).Save(u).Error
	if err != nil {
		return res, ecode.MYSQL_ERR.SetMessage(err.Error())
	}
	res.Id = rq.Id
	return res, nil
}

func (a *accountCenterRepo) SendEmailCode(ctx context.Context, rq *pb.SendEmailCodeRequest) (bool, error) {
	if !util.CheckMailFormat(rq.GetEmail()) {
		return false, ecode.INVALID_DATE_PARAM.SetMessage( "请正确填写邮箱")
	}
	code := util.GetRandomString(6)
	redisCli := a.data.cache
	switch rq.Type {
	case 1: //注册验证
		key := "registry:email"
		exist, err := redisCli.Do(ctx, "SISMEMBER", key, rq.Email).Result()
		if err != nil {
			fmt.Printf("err:%s", exist)
			return false, ecode.REDIS_ERR.SetMessage(err.Error())
		}
		if exist.(int64) == 1 {
			return false, ecode.EXTERNAL_API_NO_RESPONSE.SetMessage("该邮箱已注册绑定")
		}
		err = a.data.cache.Set(ctx, code, rq.Email, 3*time.Minute).Err()
		if err != nil {
			return false, ecode.REDIS_ERR.SetMessage(err.Error())
		}
		str, err := a.data.cache.Get(ctx, code).Result()
		if err != nil {
			return false, ecode.REDIS_ERR.SetMessage(err.Error())
		}
		fmt.Println(str)
		if err = mailUtils.Send(a.data.mail, code, rq.Email); err != nil {
			return false, ecode.EXTERNAL_API_FAIL.SetMessage(err.Error())
		}
		return true, ecode.OK

	case 2: //重置验证
		key := fmt.Sprintf("%s:%s:%s", "email", "reset",rq.Email)
		_, err := redisCli.Set(ctx, key, code, 5*time.Minute).Result()
		if err != nil {
			return false, ecode.REDIS_ERR
		}
		if err = mailUtils.Send(a.data.mail, code, rq.GetEmail()); err != nil {
			return false, ecode.EXTERNAL_API_FAIL
		}
		return true, ecode.OK
	default:
		return false, nil
	}

}

func (a *accountCenterRepo) Create(ctx context.Context, b *pb.RegisterRequest) (unum string, err error) {
	if !util.CheckPhone(b.Telephone) {
		return "", err
	}
	if !util.CheckMailFormat(b.Email) {
		return "", err
	}
	//todo 加入邀请码验证
	email, err := a.data.cache.Do(ctx, "get", b.InviteCode).Result()
	if err != nil {
		return "", ecode.REDIS_ERR.SetMessage(err.Error())
	}
	if strings.Compare(email.(string), b.Email) != 0 {
		return "", ecode.EXTERNAL_API_FAIL.SetMessage("当前email账号不符")
	} else {
		key := "registry:email"
		err = a.data.cache.Do(ctx, "SADD", key, b.Email).Err()
		if err != nil {
			return "", ecode.REDIS_ERR.SetMessage(err.Error())
		}
	}
	userNum := util.GetUserNum()
	now := util.GetTodayTimeDetail()
	generatePass, _ := util.HashPassword(b.Password)
	u := &biz.User{
		Name:      b.GetUsername(),
		Telephone: b.GetTelephone(),
		Email:     b.GetEmail(),
		UserNum:   userNum,
		Password:  generatePass,
		//Porn:         b.Porn,
		RegisterTime: now,
	}
	err = a.data.db.Model(&biz.User{}).WithContext(ctx).Create(u).Error
	if err != nil {
		fmt.Println(err)
		return "", ecode.MYSQL_ERR
	}

	return u.UserNum, ecode.OK
}

func (a *accountCenterRepo) Get(ctx context.Context, u *pb.GetAccountInfoRequest) (*pb.GetAccountInfoReply, error) {
	res := &pb.GetAccountInfoReply{}
	if claims, exist := ctx.Value("claims").(*middleware.Claims); !exist {
		return res, ecode.EXTERNAL_API_NO_RESPONSE.SetMessage("断言失败")
	} else if claims.UserId != int(u.Id) && claims.UserRole != 127 {
		return res, ecode.EXTERNAL_API_NO_RESPONSE.SetMessage("非管理员禁止操作")
	}
	err := a.data.db.Model(&biz.User{}).WithContext(ctx).Where("id = ?", u.Id).First(res).Error
	if err != nil{
		return res,ecode.MYSQL_ERR.SetMessage(err.Error())
	}
	return res, err
}

func (a *accountCenterRepo) GetList(ctx context.Context, u *pb.ListAccountRequest) (*pb.ListAccountReply, error) {
	//userArr := []*biz.User{}
	res := &pb.ListAccountReply{}
	arr := []*pb.ListAccountReply_AccountInfo{}
	if claims, exist := ctx.Value("claims").(*middleware.Claims); !exist {
		return res, ecode.EXTERNAL_API_NO_RESPONSE.SetMessage("断言失败")
	} else if claims.UserRole != int8(127) {
		fmt.Println(claims)
		return res, ecode.EXTERNAL_API_NO_RESPONSE.SetMessage("非管理员禁止操作")
	}

	var count int64
	err := a.data.db.Model(&biz.User{}).
		Offset(int(u.Offset)).Limit(int(u.Limit)).Count(&count).
		Order("id asc").
		Find(&arr).Error
	if err != nil {
		return res, ecode.MYSQL_ERR
	}
	res.Total = int32(count)
	res.Data = arr
	return res, err
}

const accountCookieID = "account:cookie:id:%d"

func (a *accountCenterRepo) Login(ctx context.Context, l *pb.CommonLoginRequest) (*pb.CommonLoginReply, error) {
	res := &pb.CommonLoginReply{}
	res.Id = 0
	type login struct {
		UserNum   string `json:"user_num"`
		Telephone string `json:"telephone"`
		Email     string `json:"email"`
		Password  string `json:"password"`
	}
	temp := &login{
		UserNum:   l.GetUserNum(),
		Telephone: l.GetTelephone(),
		Email:     l.GetEmail(),
		Password:  l.GetPassword(),
	}
	var getLogin func(lg *login) (key, value string) = func(lg *login) (key, value string) {
		refValue := reflect.ValueOf(lg).Elem()
		refType := reflect.TypeOf(lg).Elem()
		fieldCount := refValue.NumField()
		fmt.Println(fieldCount)
		for i := 0; i < fieldCount; i++ {
			fieldType := refType.Field(i)
			fieldValue := refValue.Field(i)
			fieldTag := fieldType.Tag.Get("json")
			if fieldTag == "password" {
				continue
			}
			if fieldValue.String() == "" {
				continue
			}
			return fieldTag, fieldValue.String()
		}
		return "", ""
	}
	key, value := getLogin(temp)
	if key == "" || value == "" {
		return res, ecode.INVALID_PARAM.SetMessage("请填入登录参数")
	}
	condition := fmt.Sprintf("%s = ?", key)
	userModel := &biz.User{}
	err := a.data.db.Model(&biz.User{}).WithContext(ctx).Where(condition, value).First(userModel).Error
	if err != nil {
		return res, ecode.MYSQL_ERR.SetMessage(err.Error())
	}
	if userModel.ID == 0 {
		return res, ecode.MYSQL_ERR.SetMessage("用户id不存在")
	}

	if !util.CheckPasswordHash(l.Password, userModel.Password) {
		return res, ecode.AUTH_FAIL.SetMessage("登录密码错误")
	}

	randId := util.GetRandomString(16)
	token, _ := middleware.GenerateToken(userModel.ID, userModel.Role, userModel.Name, userModel.UserNum, randId, middleware.ExpireTime)
	fmt.Printf("token:%s", token)
	setcookie:= middleware.SetCookie(ctx, token, int(userModel.Role), 7*24*3600)
	cacheKey := fmt.Sprintf(accountCookieID, userModel.ID)
	err = a.data.cache.Do(ctx, "set", cacheKey, randId).Err()
	if err != nil {
		a.log.Error("token白名单插入失败")
	}
	userModel.LastLoginAt = util.GetTodayTimeDetail()
	err = a.data.db.Model(&biz.User{}).Where("id = ?", userModel.ID).Save(userModel).Error
	res.SetCookie = setcookie
	res.Id = int32(userModel.ID)
	return res, err
}
func (a *accountCenterRepo) GetPorn(ctx context.Context) (*pb.GetPornReply, error) {
	porn := util.GetRandomString(4)
	fmt.Println(porn)
	key := fmt.Sprintf("cloud:porn")
	err := a.data.cache.Do(ctx, "SADD", key, porn).Err()
	res := &pb.GetPornReply{Porn: porn}
	return res, err
}

func (a *accountCenterRepo) Update(ctx context.Context, b *pb.UpdateAccountInfoRequest) (id int, err error) {
	u := &biz.User{}
	if claims, exist := ctx.Value("claims").(*middleware.Claims); !exist {
		return 0, ecode.EXTERNAL_API_NO_RESPONSE.SetMessage("断言失败")
	} else if claims.UserId != int(b.Id) || claims.UserRole != 127 {
		return 0, ecode.EXTERNAL_API_NO_RESPONSE.SetMessage("非管理员禁止操作")
	}
	err = a.data.db.Model(&biz.User{}).Where("id = ?", b.Id).First(u).Error
	if err != nil {
		return 0, ecode.MYSQL_ERR.SetMessage(err.Error())
	}
	u.Name = b.Name
	u.Telephone = b.Telephone
	u.Avatar = b.Avatar
	u.Signature = b.Signature
	err = a.data.db.Model(&biz.User{}).WithContext(ctx).Where("id = ?",b.Id).Save(u).Error
	if err != nil {
		return 0, ecode.MYSQL_ERR.SetMessage(err.Error())
	}
	return u.ID, err
}
func (a *accountCenterRepo) Logout(ctx context.Context) (id int, err error) {
	tr, ok := transport.FromServerContext(ctx)
	if ok {
		id, _ = strconv.Atoi(tr.RequestHeader().Get("x-md-global-uid"))
		tr.ReplyHeader().Set("Authorization", "")
		key := fmt.Sprintf(accountCookieID, id)
		err = a.data.cache.Do(ctx, "del", key).Err()
		if err != nil {
			a.log.Error("token白名单删除失败")
			return 0, err
		}
	}
	return id, err
}



func (a *accountCenterRepo) GetGuest(ctx context.Context, u *pb.GetGuestRequest) (*pb.GetGuestReply, error) {
	res := &pb.GetGuestReply{}
	guest := util.GetGuestNum()
	key := fmt.Sprintf("guest:%s", guest)
	randId := util.GetRandomString(16)
	userName := fmt.Sprintf("Guest-%s", guest)
	token, _ := middleware.GenerateToken(0, middleware.RoleGuest, userName, guest, randId, middleware.GuestExpire)
	_, err := a.data.cache.Set(ctx,  key, randId, 3600*2).Result()
	if err != nil {
		return res, ecode.REDIS_ERR
	}
	res.Gid = guest
	res.Token = token
	return res, err
}

package biz

import (
	pb "banana/api/account-center/service/v1"
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type User struct {
	ID           int    `gorm:"primary_key" json:"id"`
	UserNum      string `gorm:"type:varchar(10) not null" json:"user_num"`
	Name         string `gorm:"type:varchar(20) not null" json:"name"`
	Telephone    string `gorm:"type:varchar(14) not null" json:"telephone"`
	Email        string `gorm:"type:varchar(30) not null" json:"email"`
	Password     string `gorm:"type:varchar(256);default:'' not null" json:"password"`
	Avatar       string `gorm:"type:varchar(256)" json:"avatar"`
	IsVip        bool   `gorm:"default:0" json:"is_vip"`
	Role         int8   `gorm:"type:tinyint(5);default:1" json:"role"`
	Status       int    `gorm:"type:tinyint(5)" json:"status"` // 0 禁用  1可用
	Porn         string `gorm:"type:varchar(6)" json:"porn"`   //邀请码
	LastLoginAt  string `gorm:"type:varchar(30)" json:"last_login_at"`
	RegisterTime string `gorm:"type:varchar(30)" json:"register_time"`
	CreatedAt    int64  `gorm:"autoCreateAt" json:"created_at"`
	UpdatedAt    int64  `gorm:"autoUpdateAt" json:"updated_at"`
}

type List struct {
	Total int     `json:"total"`
	Data  []*User `json:"data"`
}

type Login struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type AccountCenterRepo interface {
	Create(ctx context.Context, u *pb.RegisterRequest) (string, error)
	Login(ctx context.Context, l *pb.CommonLoginRequest) (*pb.CommonLoginReply, error)
	Logout(ctx context.Context) (int, error)
	Get(ctx context.Context, rq *pb.GetAccountInfoRequest) (*pb.GetAccountInfoReply, error)
	GetList(ctx context.Context, rq *pb.ListAccountRequest) (*pb.ListAccountReply, error)
	GetPorn(ctx context.Context) (*pb.GetPornReply, error)
	Update(ctx context.Context, u *pb.UpdateAccountInfoRequest) (int, error)
	SendEmailCode(ctx context.Context, rq *pb.SendEmailCodeRequest) (bool, error)
	ResetPass(ctx context.Context, rq *pb.PasswordResetRequest) (*pb.PasswordResetReply, error)
	GetGuest(ctx context.Context, rq *pb.GetGuestRequest) (*pb.GetGuestReply, error)
	SetAdmin(ctx context.Context, rq *pb.SetAdminRequest) (*pb.SetAdminReply, error)
}

type AccountCenterCase struct {
	repo AccountCenterRepo
	log  *log.Helper
}

func NewAccountCenterCase(repo AccountCenterRepo, logger log.Logger) *AccountCenterCase {
	return &AccountCenterCase{repo: repo, log: log.NewHelper(log.With(logger, "module", "accountcenter/case"))}
}
func (ac *AccountCenterCase) SetAdmin(ctx context.Context,l *pb.SetAdminRequest) (*pb.SetAdminReply,error){
	return ac.repo.SetAdmin(ctx,l)
}
func (ac *AccountCenterCase) Login(ctx context.Context, l *pb.CommonLoginRequest) (*pb.CommonLoginReply, error) {
	out, err := ac.repo.Login(ctx, l)
	return out, err
}
func (ac *AccountCenterCase) Create(ctx context.Context, u *pb.RegisterRequest) (string, error) {
	out, err := ac.repo.Create(ctx, u)
	return out, err
}

func (ac *AccountCenterCase) Get(ctx context.Context, u *pb.GetAccountInfoRequest) (*pb.GetAccountInfoReply, error) {
	return ac.repo.Get(ctx, u)
}

func (ac *AccountCenterCase) GetList(ctx context.Context, rq *pb.ListAccountRequest) (*pb.ListAccountReply, error) {
	res, err := ac.repo.GetList(ctx, rq)
	return res, err
}

func (ac *AccountCenterCase) GetPorn(ctx context.Context) (*pb.GetPornReply, error) {
	res, err := ac.repo.GetPorn(ctx)
	return res, err
}

func (ac *AccountCenterCase) Update(ctx context.Context, u *pb.UpdateAccountInfoRequest) (int, error) {
	res, err := ac.repo.Update(ctx, u)
	if err != nil {
		return 0, err
	}
	return res, err
}

func (ac *AccountCenterCase) Logout(ctx context.Context) (int, error) {
	return ac.repo.Logout(ctx)
}

func (ac *AccountCenterCase) SendEmailCode(ctx context.Context, u *pb.SendEmailCodeRequest) (bool, error) {
	return ac.repo.SendEmailCode(ctx, u)
}

func (ac *AccountCenterCase) ResetPass(ctx context.Context, u *pb.PasswordResetRequest) (*pb.PasswordResetReply, error) {
	return ac.repo.ResetPass(ctx, u)
}

func (ac *AccountCenterCase) GetGuest(ctx context.Context, u *pb.GetGuestRequest) (*pb.GetGuestReply, error) {
	return ac.repo.GetGuest(ctx, u)
}

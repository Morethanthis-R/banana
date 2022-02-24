package service

import (
	"banana/app/account-center/service/internal/biz"
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"

	pb "banana/api/account-center/service/v1"
)

func NewAccountCenterService(ac *biz.AccountCenterCase, logger log.Logger) *AccountCenterService {
	return &AccountCenterService{
		ac:  ac,
		log: log.NewHelper(logger),
	}
}
func (s *AccountCenterService) ForgetPass(ctx context.Context,req *pb.ForgetPassRequest)(*pb.ForgetPassReply,error){
	return s.ac.ForgetPass(ctx,req)
}
func (s *AccountCenterService) SetAdmin(ctx context.Context, req *pb.SetAdminRequest) (*pb.SetAdminReply, error) {
	return s.ac.SetAdmin(ctx, req)
}
func (s *AccountCenterService) Login(ctx context.Context, req *pb.CommonLoginRequest) (*pb.CommonLoginReply, error) {
	res := &pb.CommonLoginReply{}
	res, err := s.ac.Login(ctx, req)
	return res, err
}
func (s *AccountCenterService) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutReply, error) {
	res := &pb.LogoutReply{}
	id, err := s.ac.Logout(ctx)
	if err != nil {
		return res, err
	}
	res.Id = int32(id)
	return res, err
}
func (s *AccountCenterService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterReply, error) {
	res := &pb.RegisterReply{}
	id, err := s.ac.Create(ctx, req)
	if err != nil {
		return res,err
	}
	res.UserNum = id
	return res, err
}
func (s *AccountCenterService) WXLogin(ctx context.Context, req *pb.WXLoginRequest) (*pb.WXLoginReply, error) {
	return &pb.WXLoginReply{}, nil
}
func (s *AccountCenterService) GetAccountInfo(ctx context.Context, req *pb.GetAccountInfoRequest) (*pb.GetAccountInfoReply, error) {
	return s.ac.Get(ctx, req)
}
func (s *AccountCenterService) PasswordReset(ctx context.Context, req *pb.PasswordResetRequest) (*pb.PasswordResetReply, error) {

	return s.ac.ResetPass(ctx,req)
}
func (s *AccountCenterService) ListAccount(ctx context.Context, req *pb.ListAccountRequest) (*pb.ListAccountReply, error) {
	return s.ac.GetList(ctx, req)
}
func (s *AccountCenterService) UpdateAccountInfo(ctx context.Context, req *pb.UpdateAccountInfoRequest) (*pb.UpdateAccountInfoReply, error) {
	res := &pb.UpdateAccountInfoReply{}
	id, err := s.ac.Update(ctx, req)
	if err != nil {
		return res, err
	}
	if id == 0 {
		fmt.Printf("id:%d", id)
		fmt.Println(err)
		return res, err
	}
	res.Id = int32(id)
	return res, nil
}
func (s *AccountCenterService) GetPorn(ctx context.Context, req *pb.GetPornRequest) (*pb.GetPornReply, error) {
	return s.ac.GetPorn(ctx)
}

func (s *AccountCenterService) GetGuest(ctx context.Context, req *pb.GetGuestRequest) (*pb.GetGuestReply, error) {
	return s.ac.GetGuest(ctx, req)
}

func (s *AccountCenterService) SendEmailCode(ctx context.Context, req *pb.SendEmailCodeRequest) (*pb.SendEmailCodeReply, error){
	res := &pb.SendEmailCodeReply{}
	if status,err := s.ac.SendEmailCode(ctx,req);status !=true{
		res.Message = "fail"
		return res,err
	}
	res.Message = "success"
	return res,nil
}

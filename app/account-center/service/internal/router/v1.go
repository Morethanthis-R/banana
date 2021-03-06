package router

import (
	pb "banana/api/account-center/service/v1"
	"banana/app/account-center/service/internal/service"
	"banana/pkg/middleware"
	"banana/pkg/response"
	"github.com/gin-gonic/gin"
	"strconv"
)
var accountService *service.AccountCenterService

func apiV1(group gin.IRoutes, ac *service.AccountCenterService) {
	accountService = ac
	group.POST("/common/login", LoginHandler)
	group.GET("/common/porn", GetPornHandler)
	group.GET("/common/guest", GetGuestHandler)
	group.POST("/common/register",RegisterHandler)
	group.POST("/e-validate", GetEmailHandler)
	group.POST("/forget-pass",ForgetPassHandler)
	group.Use(middleware.JWTAuth())
	group.GET("/logout", LogoutHandler)
	group.POST("/set-admin", SetAdminHandler)
	group.GET("/account/info/:id", GetUserInfoHandler)
	group.POST("/reset", ResetPassHandler)
	group.GET("/list",UserListHandler)
	group.POST("/update",UpdateInfoHandler)

}
func ForgetPassHandler(c *gin.Context)  {
	req := &pb.ForgetPassRequest{}
	if err := c.BindJSON(req);err != nil {
		response.NewErrWithCodeAndMsg(c,200,response.BIND_JSON_ERROR)
		return
	}
	res,err := accountService.ForgetPass(c,req)
	if err != nil {
		response.NewErrWithCodeAndMsg(c,200,err.Error())
		return
	}
	response.NewSuccess(c,res)
}
func ResetPassHandler(c *gin.Context) {
	req := &pb.PasswordResetRequest{}
	if err := c.BindJSON(req);err != nil {
		response.NewErrWithCodeAndMsg(c,200,response.BIND_JSON_ERROR)
		return
	}
	res,err := accountService.PasswordReset(c,req)
	if err != nil {
		response.NewErrWithCodeAndMsg(c,200,err.Error())
		return
	}
	response.NewSuccess(c,res)
}
func GetEmailHandler(c *gin.Context) {
	req := &pb.SendEmailCodeRequest{}
	if err := c.BindJSON(req);err != nil {
		response.NewErrWithCodeAndMsg(c,200,response.BIND_JSON_ERROR)
		return
	}
	res,err := accountService.SendEmailCode(c,req)
	if err != nil {
		response.NewErrWithCodeAndMsg(c,200,err.Error())
		return
	}
	response.NewSuccess(c,res)
}
func LoginHandler(c *gin.Context) {
	req := &pb.CommonLoginRequest{}
	if err := c.BindJSON(req); err != nil {
		response.NewErrWithCodeAndMsg(c,200,response.BIND_JSON_ERROR)
		return
	}
	res,err := accountService.Login(c,req)
	if err != nil {
		response.NewErrWithCodeAndMsg(c,200,err.Error())
		return
	}
	//c.Header("Set-Cookie",res.SetCookie)
	//type Return struct {
	//	Id int 	`json:"id"`
	//}
	//idStruct := &Return{Id: 0}
	//idStruct.Id = int(res.Id)
	response.NewSuccess(c,res)
}

func LogoutHandler(c *gin.Context) {
	req := &pb.LogoutRequest{}
	if err := c.BindQuery(req);err != nil {
		response.NewErrWithCodeAndMsg(c,200,response.BIND_JSON_ERROR)
		return
	}
	res,err := accountService.Logout(c,req)
	if err != nil {
		response.NewErrWithCodeAndMsg(c,200,err.Error())
		return
	}
	response.NewSuccess(c,res)
}

func SetAdminHandler(c *gin.Context) {
	req := &pb.SetAdminRequest{}
	if err := c.BindJSON(req);err != nil {
		response.NewErrWithCodeAndMsg(c,200,response.BIND_JSON_ERROR)
		return
	}
	res,err := accountService.SetAdmin(c,req)
	if err != nil {
		response.NewErrWithCodeAndMsg(c,200,err.Error())
		return
	}
	response.NewSuccess(c,res)
}

func RegisterHandler(c *gin.Context) {
	req := &pb.RegisterRequest{}
	if err := c.BindJSON(req);err != nil {
		response.NewErrWithCodeAndMsg(c,200,response.BIND_JSON_ERROR)
		return
	}

	res,err := accountService.Register(c,req)
	if err != nil {
		response.NewErrWithCodeAndMsg(c,200,err.Error())
		return
	}
	response.NewSuccess(c,res)
}

func GetUserInfoHandler(c *gin.Context){
	parse := c.Param("id")
	id,err := strconv.Atoi(parse)
	if err != nil {
		response.NewErrWithCodeAndMsg(c,200,response.BIND_JSON_ERROR)
		return
	}
	req := &pb.GetAccountInfoRequest{}
	req.Id = int64(id)
	res,err := accountService.GetAccountInfo(c,req)
	if err != nil {
		response.NewErrWithCodeAndMsg(c,200,err.Error())
		return
	}
	response.NewSuccess(c,res)
}

func UpdateInfoHandler(c *gin.Context){
	req := &pb.UpdateAccountInfoRequest{}
	if err := c.BindJSON(req);err != nil {
		response.NewErrWithCodeAndMsg(c,200,response.BIND_JSON_ERROR)
		return
	}
	res,err := accountService.UpdateAccountInfo(c,req)
	if err != nil {
		response.NewErrWithCodeAndMsg(c,200,err.Error())
		return
	}
	response.NewSuccess(c,res)
}

func UserListHandler(c *gin.Context) {


	req := &pb.ListAccountRequest{}
	type Params struct {
		Offset int `form:"offset"`
		Limit  int `form:"limit"`
		query string `form:"query"`
	}
	params := &Params{}
	if err := c.BindQuery(params);err != nil {
		response.NewErrWithCodeAndMsg(c,200,response.BIND_JSON_ERROR)
		return
	}

	req.Limit = int32(params.Limit)
	req.Offset = int32(params.Offset)
	req.Query   = c.Query("query")
	res,err := accountService.ListAccount(c,req)
	if err != nil {
		response.NewErrWithCodeAndMsg(c,200,err.Error())
		return
	}
	response.NewSuccess(c,res)
}

func GetGuestHandler(c *gin.Context) {
	req := &pb.GetGuestRequest{}
	res,err := accountService.GetGuest(c,req)
	if err != nil {
		response.NewErrWithCodeAndMsg(c,200,err.Error())
		return
	}
	response.NewSuccess(c,res)
}

func GetPornHandler(c *gin.Context) {
	req := &pb.GetPornRequest{}
	res ,err := accountService.GetPorn(c,req)
	if err != nil {
		response.NewErrWithCodeAndMsg(c,200,err.Error())
		return
	}
	response.NewSuccess(c,res)
}
package router

import (
	pb "banana/api/common/service/v1"
	"banana/app/common/service/internal/service"
	"banana/pkg/middleware"
	"banana/pkg/response"
	"github.com/gin-gonic/gin"
	"strconv"
)

var commonService *service.CommonService

func apiV1(group gin.IRoutes, cm *service.CommonService) {
	commonService = cm
	group.GET("/adv-list",AdvListHandler)
	group.Use(middleware.JWTAuth())
	group.POST("/adv-add",AdvCreateHandler )
	group.GET("/adv-del", AdvDeleteHandler)
	group.POST("/adv-update",AdvUpdateHandler)
	group.POST("/notify-add",NotifyCreateHandler )
	group.GET("/notify-del", NotifyDeleteHandler)
	group.GET("/notify-list", NotifyListHandler)
	group.GET("/notify-get", NotifyGetHandler)
	group.POST("/type-add", TypeCreateHandler )
	group.POST("/type-update",TypeUpdateHandler)
	group.GET("/type-del", TypeDeleteHandler)
	group.GET("/type-list", TypeListHandler)
}
func AdvListHandler(c *gin.Context){
	req := &pb.ReqGetAdvList{}
	type params struct {
		Aid int`form:"aid"`
	}
	query := &params{}
	if err := c.BindQuery(query);err != nil{
		response.NewErrWithCodeAndMsg(c,200,response.BIND_JSON_ERROR)
		return
	}
	req.Aid = int32(query.Aid)
	res,err := commonService.GetAdvList(c,req)
	if err != nil {
		response.NewErrWithCodeAndMsg(c,200,err.Error())
		return
	}
	response.NewSuccess(c,res)

}
func AdvCreateHandler(c *gin.Context){
	req := &pb.ReqCreateAdv{}
	if err := c.BindJSON(req);err!=nil {
		response.NewErrWithCodeAndMsg(c,200,response.BIND_JSON_ERROR)
		return
	}
	res,err := commonService.CreateAdv(c,req)
	if err != nil {
		response.NewErrWithCodeAndMsg(c,200,err.Error())
		return
	}
	response.NewSuccess(c,res)

}
func AdvDeleteHandler(c *gin.Context){
	req := &pb.ReqDeleteAdv{}
	if err := c.BindJSON(req);err!=nil{
		response.NewErrWithCodeAndMsg(c,200,response.BIND_JSON_ERROR)
		return
	}
	res,err := commonService.DeleteAdv(c,req)
	if err != nil {
		response.NewErrWithCodeAndMsg(c,200,err.Error())
		return
	}
	response.NewSuccess(c,res)

}
func AdvUpdateHandler(c *gin.Context){
	req := &pb.ReqUpdateAdv{}
	if err := c.BindJSON(req);err!=nil{
		response.NewErrWithCodeAndMsg(c,200,response.BIND_JSON_ERROR)
		return
	}
	res,err := commonService.UpdateAdv(c,req)
	if err != nil {
		response.NewErrWithCodeAndMsg(c,200,err.Error())
		return
	}
	response.NewSuccess(c,res)
}
func NotifyCreateHandler(c *gin.Context){
	req := &pb.ReqCreateNotify{}
	if err := c.BindJSON(req);err!=nil{
		response.NewErrWithCodeAndMsg(c,200,response.BIND_JSON_ERROR)
		return
	}
	res,err := commonService.CreateNotify(c,req)
	if err != nil {
		response.NewErrWithCodeAndMsg(c,200,err.Error())
		return
	}
	response.NewSuccess(c,res)
}
func NotifyDeleteHandler(c *gin.Context){
	req := &pb.ReqDeleteNotify{}
	nids := c.Query("nid")
	nid,_:= strconv.Atoi(nids)
	req.Nids = append(req.Nids,int32(nid))
	res,err := commonService.DeleteNotify(c,req)
	if err != nil {
		response.NewErrWithCodeAndMsg(c,200,err.Error())
		return
	}
	response.NewSuccess(c,res)
}
func NotifyListHandler(c *gin.Context){
	req := &pb.ReqGetNotifyList{}
	type params struct {
		SortObject int32      `form:"sort_object"` //0:默认排序 1:编辑时间
		SortType   int32      `form:"sort_type"`//0:系统默认 1:asc升序  2:desc降序
		Keywords   string     `form:"keywords"`                   //搜索关键字
		Uid        int32      `form:"uid"`
		Offset     int32       `form:"offset"`
		Limit      int32       `form:"limit"`
	}
	query:=&params{}
	if err := c.BindQuery(query);err != nil {
		response.NewErrWithCodeAndMsg(c,200,response.BIND_JSON_ERROR)
		return
	}
	req.Uid = query.Uid
	res,err := commonService.GetNotifyList(c,req)
	if err != nil {
		response.NewErrWithCodeAndMsg(c,200,err.Error())
		return
	}
	response.NewSuccess(c,res)
}
func NotifyGetHandler(c *gin.Context){
	req := &pb.ReqGetNotifyObject{}
	if err := c.BindJSON(req);err != nil {
		response.NewErrWithCodeAndMsg(c,200,response.BIND_JSON_ERROR)
		return
	}
	res,err := commonService.GetNotifyObject(c,req)
	if err != nil {
		response.NewErrWithCodeAndMsg(c,200,err.Error())
		return
	}
	response.NewSuccess(c,res)
}
func TypeCreateHandler(c *gin.Context){
	req := &pb.ReqCreateNotifyType{}
	if err := c.BindJSON(req);err != nil {
		response.NewErrWithCodeAndMsg(c,200,response.BIND_JSON_ERROR)
		return
	}
	res,err := commonService.CreateNotifyType(c,req)
	if err != nil {
		response.NewErrWithCodeAndMsg(c,200,err.Error())
		return
	}
	response.NewSuccess(c,res)
}
func TypeUpdateHandler(c *gin.Context){
	req := &pb.ReqUpdateNotifyType{}
	if err := c.BindJSON(req);err != nil {
		response.NewErrWithCodeAndMsg(c,200,response.BIND_JSON_ERROR)
		return
	}
	res,err := commonService.UpdateNotifyType(c,req)
	if err != nil {
		response.NewErrWithCodeAndMsg(c,200,err.Error())
		return
	}
	response.NewSuccess(c,res)
}
func TypeDeleteHandler(c *gin.Context){
	req := &pb.ReqDeleteNotifyType{}
	if err := c.BindJSON(req);err != nil {
		response.NewErrWithCodeAndMsg(c,200,response.BIND_JSON_ERROR)
		return
	}
	res,err := commonService.DeleteNotifyType(c,req)
	if err != nil {
		response.NewErrWithCodeAndMsg(c,200,err.Error())
		return
	}
	response.NewSuccess(c,res)
}
func TypeListHandler(c *gin.Context){
	req := &pb.ReqGetNotifyTypeList{}
	res,err := commonService.GetNotifyTypeList(c,req)
	if err != nil {
		response.NewErrWithCodeAndMsg(c,200,err.Error())
		return
	}
	response.NewSuccess(c,res)
}

package router

import (
	"banana/app/transfer/service/internal/service"
	"github.com/gin-gonic/gin"
)

func Init(eg *gin.Engine,tf *service.TransferService) {
	group :=eg.Group("/banana/transfer").Use()
	apiV1(group,tf)
}
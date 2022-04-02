package router

import (
	"banana/app/common/service/internal/service"
	"github.com/gin-gonic/gin"
)

func Init(eg *gin.Engine,cm *service.CommonService) {
	group :=eg.Group("/banana/common").Use()
	apiV1(group,cm)
}
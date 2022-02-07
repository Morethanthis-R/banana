package router
import (
	"banana/app/account-center/service/internal/service"
	"github.com/gin-gonic/gin"
)

func Init(eg *gin.Engine,ac *service.AccountCenterService) {
	group :=eg.Group("/banana/account-center").Use()
	apiV1(group,ac)
}
package server

import (
	"banana/app/account-center/service/internal/conf"
	"banana/app/account-center/service/internal/router"
	"banana/app/account-center/service/internal/service"
	"banana/pkg/middleware"
	"context"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new a HTTP server.
func NewHTTPServer(c *conf.Server, ac *service.AccountCenterService) *http.Server {
	engine := gin.Default()
	engine.Use(middleware.Cors())
	router.Init(engine,ac)
	httpSrv := http.NewServer(http.Address(c.Http.Addr), http.Timeout(c.Http.Timeout.AsDuration()))
	httpSrv.HandlePrefix("/",engine)
	pprof.Register(engine, "/common/debug")
	return httpSrv
}

func MatchFunc(ctx context.Context,operation string) bool {
	whiteList := []string{
		"/ac.service.v1.AccountCenter/Register",
		"/ac.service.v1.AccountCenter/GetPorn",
		"/ac.service.v1.AccountCenter/Login",
		"/ac.service.v1.AccountCenter/SendEmailCode",
		"/ac.service.v1.AccountCenter/GetGuest",
	}
	for _, v := range whiteList {
		if v == operation{
			return false
		}
	}

	return true
}
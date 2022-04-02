package server

import (
	"banana/app/account-center/service/internal/conf"
	"banana/app/account-center/service/internal/router"
	"banana/app/account-center/service/internal/service"
	"banana/pkg/middleware"
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
	pprof.Register(engine, "/banana/ac/debug")
	return httpSrv
}
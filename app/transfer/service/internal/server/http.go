package server

import (
	"banana/app/transfer/service/internal/conf"
	"banana/app/transfer/service/internal/router"
	"banana/app/transfer/service/internal/service"
	"banana/pkg/middleware"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new a HTTP server.
func NewHTTPServer(c *conf.Server,tf *service.TransferService) *http.Server {
	engine := gin.Default()
	engine.Use(middleware.Cors())
	router.Init(engine,tf)
	engine.Group("/banana/tf").Use(middleware.JWTAuth()).POST("/upload",router.UploadHandler)
	engine.Group("/banana/tf").Use(middleware.JWTAuth()).POST("/guest-upload",router.GuestUpload)
	engine.Group("/banana/tf").POST("/upload-static",router.UploadStatic)
	httpSrv := http.NewServer(http.Address(c.Http.Addr), http.Timeout(c.Http.Timeout.AsDuration()))
	httpSrv.HandlePrefix("/",engine)
	pprof.Register(engine, "/banana/transfer/debug")
	return httpSrv
}

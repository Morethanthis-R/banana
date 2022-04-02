// +build wireinject

package main
import (
	"banana/app/common/service/internal/biz"
	"banana/app/common/service/internal/conf"
	"banana/app/common/service/internal/data"
	"banana/app/common/service/internal/server"
	"banana/app/common/service/internal/service"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)
// initApp init kratos application.
func initApp(*conf.Server, *conf.Registry, *conf.Data,log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
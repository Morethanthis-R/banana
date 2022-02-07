// +build wireinject

package main

import (
	"banana/app/transfer/service/internal/biz"
	"banana/app/transfer/service/internal/conf"
	"banana/app/transfer/service/internal/data"
	"banana/app/transfer/service/internal/server"
	"banana/app/transfer/service/internal/service"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)
// initApp init kratos application.
func initApp(*conf.Server, *conf.Registry,*conf.Data,log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet,data.ProviderSet,biz.ProviderSet, service.ProviderSet, newApp))
}

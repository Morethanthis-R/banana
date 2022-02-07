// +build wireinject

package main

import (
	"banana/app/account-center/service/internal/biz"
	"banana/app/account-center/service/internal/conf"
	"banana/app/account-center/service/internal/data"
	"banana/app/account-center/service/internal/server"
	"banana/app/account-center/service/internal/service"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// initApp init kratos application.
func initApp(*conf.Server, *conf.Registry, *conf.Data, *conf.Mail, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}

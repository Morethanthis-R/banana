// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package main

import (
	"banana/app/transfer/service/internal/biz"
	"banana/app/transfer/service/internal/conf"
	"banana/app/transfer/service/internal/data"
	"banana/app/transfer/service/internal/server"
	"banana/app/transfer/service/internal/service"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
)

// Injectors from wire.go:

// initApp init kratos application.
func initApp(confServer *conf.Server, registry *conf.Registry, confData *conf.Data, logger log.Logger) (*kratos.App, func(), error) {
	dataData, cleanup, err := data.NewData(confData, logger)
	if err != nil {
		return nil, nil, err
	}
	transferRepo := data.NewTransferRepo(dataData, logger)
	transferCase := biz.NewTransferCase(transferRepo, logger)
	transferService := service.NewTransferService(transferCase, logger)
	httpServer := server.NewHTTPServer(confServer, transferService)
	registrar := server.NewRegistrar(registry)
	app := newApp(logger, httpServer, registrar)
	return app, func() {
		cleanup()
	}, nil
}
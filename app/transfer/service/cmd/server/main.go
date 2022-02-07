package main


import (
	"banana/app/transfer/service/internal/conf"
	"flag"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/http"
	"os"
)

var (
	// Name is the name of the compiled software.
	Name ="banana.transfer.service"
	// Version is the version of the compiled software.
	Version = "v1"
	// flagconf is the config flag.
	flagconf string
)

func init() {
	//flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
	flag.StringVar(&flagconf, "conf", "configs", "config path, eg: -conf config.yaml")
}

func newApp(logger log.Logger, hs *http.Server, rg registry.Registrar) *kratos.App {
	return kratos.New(
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			hs,
		),
		kratos.Registrar(rg),

	)
}

////Set global trace provider
//func setTracerProvider(url string) error {
//	// Create the Jaeger exporter
//	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
//	if err != nil {
//		return err
//	}
//	tp := tracesdk.NewTracerProvider(
//		// Set the sampling rate based on the parent span to 100%
//		tracesdk.WithSampler(tracesdk.ParentBased(tracesdk.TraceIDRatioBased(1.0))),
//		// Always be sure to batch in production.
//		tracesdk.WithBatcher(exp),
//		// Record information about this application in an Resource.
//		tracesdk.WithResource(resource.NewSchemaless(
//			semconv.ServiceNameKey.String(Name),
//			attribute.String("env", "dev"),
//		)),
//	)
//	otel.SetTracerProvider(tp)
//	return nil
//}

func main() {
	flag.Parse()
	logger := log.NewStdLogger(os.Stdout)

	cfg := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	if err := cfg.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := cfg.Scan(&bc); err != nil {
		panic(err)
	}

	var rc conf.Registry
	if err := cfg.Scan(&rc); err != nil {
		panic(err)
	}

	app, cleanup, err := initApp(bc.Server, &rc,bc.Data, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}


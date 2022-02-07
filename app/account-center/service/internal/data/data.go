package data

import (
	"banana/app/account-center/service/internal/biz"
	"banana/app/account-center/service/internal/conf"
	"fmt"
	//"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/extra/redisotel"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlog "gorm.io/gorm/logger"
	"time"

	// init mysql driver
	_ "github.com/go-sql-driver/mysql"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewAccountCenterRepo)

// Data .
type Data struct {
	mail *conf.Mail
	cache *redis.Client
	db  *gorm.DB
	log *log.Helper
}
func NewCache(conf *conf.Data, logger log.Logger) *redis.Client {
	log := log.NewHelper(log.With(logger, "module", "ac-service/data/redis"))
	var options = &redis.Options{
		Addr:         conf.Redis.Addr,
		Password:     conf.Redis.Password,
		PoolSize:     100,
		DialTimeout:  conf.Redis.DialTimeout.AsDuration(),
		WriteTimeout: conf.Redis.WriteTimeout.AsDuration(),
		ReadTimeout:  conf.Redis.ReadTimeout.AsDuration(),
	}
	client := redis.NewClient(options)
	if client == nil{
		log.Fatalf("failed opening connection to redis")
	}
	client.AddHook(redisotel.TracingHook{})
	return client

}
type Writer struct{
}
func (w Writer) Printf(format string,args ...interface{}) {
	fmt.Printf(format, args...)
}
func NewDB(conf *conf.Data, logger log.Logger) *gorm.DB {
	log := log.NewHelper(log.With(logger, "module", "ac-service/data/gorm"))
	newLogger := gormlog.New(
		Writer{},
		gormlog.Config{
			SlowThreshold:              200 * time.Millisecond,   // Slow SQL threshold
			LogLevel:                   gormlog.Info,   // Log level
			IgnoreRecordNotFoundError:  true,          // Ignore ErrRecordNotFound error for logger
			Colorful:                   true,         // Disable color
		},
	)
	db, err := gorm.Open(mysql.Open(conf.Database.Source), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatalf("failed opening connection to mysql: %v", err)
	}

	if err := db.AutoMigrate(&biz.User{},); err != nil {
		log.Fatal(err)
	}
	return db
}

// NewData .
func NewData(conf *conf.Data,mail *conf.Mail,logger log.Logger) (*Data, func(), error) {
	log := log.NewHelper(log.With(logger, "module", "order-service/data"))

	d := &Data{
		cache: NewCache(conf,logger),
		db:  NewDB(conf,logger),
		log: log,
		mail: mail,
	}
	return d, func() {

	}, nil
}

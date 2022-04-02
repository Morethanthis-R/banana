package data

import (
	"banana/app/transfer/service/internal/biz"
	"banana/app/transfer/service/internal/conf"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/extra/redisotel"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlog "gorm.io/gorm/logger"
	"time"
)

var ProviderSet = wire.NewSet(NewData, NewTransferRepo,NewRabbitMqProducer)

// Data .
type Data struct {
	cache          *redis.Client
	Db             *gorm.DB
	Minio_internal *minio.Client
	minio_online   *minio.Client
	log *log.Helper

}
func NewMinioClientInternal(conf *conf.Data,logger log.Logger) *minio.Client{
	log := log.NewHelper(log.With(logger, "module", "transfer/data/minio"))
	client, err := minio.New("47.107.95.82:8000", &minio.Options{
		Creds:        credentials.NewStaticV4(conf.Minio.AccessKeyId,conf.Minio.SecretAccessKey,""),
		Secure:       false,
	})
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return client
}
func NewMinioClientOnline(conf *conf.Data,logger log.Logger) *minio.Client{
	log := log.NewHelper(log.With(logger, "module", "transfer/data/minio"))
	client, err := minio.New("47.107.95.82:8000", &minio.Options{
		Creds:        credentials.NewStaticV4(conf.Minio.AccessKeyId,conf.Minio.SecretAccessKey,""),
		Secure:       false,
	})
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return client
}
func NewCache(conf *conf.Data, logger log.Logger) *redis.Client {
	log := log.NewHelper(log.With(logger, "module", "transfer/data/redis"))
	var options = &redis.Options{
		Addr:        conf.Redis.Addr,
		Password:     conf.Redis.Password,
		DB:           1,
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
	log := log.NewHelper(log.With(logger, "module", "tf-service/data/gorm"))
	newLogger := gormlog.New(
		Writer{},
		gormlog.Config{
			SlowThreshold:              200 * time.Millisecond,   // Slow SQL threshold
			LogLevel:                   gormlog.Info,   // Log level
			IgnoreRecordNotFoundError:  true,          // Ignore ErrRecordNotFound error for logger
			Colorful:                   true,         // Disable color
		},
	)
	fmt.Println(conf)
	db, err := gorm.Open(mysql.Open(conf.Database.Source), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatalf("failed opening connection to mysql: %v", err)
	}

	if err := db.AutoMigrate(
		&biz.File{},&biz.UserFile{},&biz.UserDirectory{},&biz.ShareHistory{},
		); err != nil {
		log.Fatal(err)
	}
	db.Migrator()
	return db
}

// NewData .
func NewData(conf *conf.Data,logger log.Logger) (*Data, func(), error) {
	log := log.NewHelper(log.With(logger, "module", "transfer/data"))

	d := &Data{
		cache:          NewCache(conf,logger),
		Db:             NewDB(conf,logger),
		Minio_internal: NewMinioClientInternal(conf,logger),
		minio_online:   NewMinioClientOnline(conf,logger),
		log:            log,
	}
	return d, func() {

	}, nil
}

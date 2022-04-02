package data

import (
	"banana/app/transfer/service/internal/conf"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/streadway/amqp"
	"sync"
)
// 定义生产者接口
type Producer interface {
	MsgContent() string
}

// 定义接收者接口
type Receiver interface {
	Consumer([]byte) (*MessageObject,error)
}


// 定义RabbitMQ对象
type RabbitMQ struct {
	Connection   *amqp.Connection
	Channel      *amqp.Channel
	QueueName    string // 队列名称
	RoutingKey   string // key名称
	ExchangeName string // 交换机名称
	ExchangeType string // 交换机类型
	ProducerList []Producer
	ReceiverList []Receiver
	Mu           sync.RWMutex
	Wg           sync.WaitGroup
}

type MessageObject struct {
	Fid         int    `json:"fid"`
	FileName    string `json:"file_name"`
	FileHash    string `json:"file_hash"`
	FileStr     string `json:"file_str"`
	FilePath    string `json:"file_path"`
	ContentType string `json:"content_type"`
	Bucket      string  `json:"bucket"`

}

func NewRabbitMqProducer(conf *conf.Data, logger log.Logger) (rmq *RabbitMQ, cleanup func(), err error) {
	log := log.NewHelper(log.With(logger, "module", "transfer/rabbit_mq"))
	RabbitUrl := fmt.Sprintf("amqp://%s:%s@%s:%s/", conf.Mq.User, conf.Mq.Password, conf.Mq.Host, conf.Mq.Port)
	conn, err := amqp.Dial(RabbitUrl)
	if err != nil {
		//log.Fatalf("MQ打开链接失败:%s \n", err)
		fmt.Printf("MQ打开链接失败:%s \n", err)
	}
	channel, err := conn.Channel()
	if err != nil {
		//log.Fatalf("MQ打开管道失败:%s \n", err)
		fmt.Printf("MQ打开管道失败:%s \n", err)
	}
	rmq = &RabbitMQ{
		Connection:   conn,
		Channel:      channel,
		QueueName:    "transfer.mq",
		RoutingKey:   "transfer.key",
		ExchangeName: "transfer.ec",
		ExchangeType: "direct",
		ProducerList: []Producer{},
		ReceiverList: []Receiver{},
		Mu:           sync.RWMutex{},
	}
	fmt.Printf("new地址：%p",rmq.Channel)
	cleanup = func() {
		err = channel.Close()
		if err != nil {
			log.Fatalf("MQ管道关闭失败:%s \n", err)
		}
		err = conn.Close()
		if err != nil {
			log.Fatalf("MQ链接关闭失败:%s \n", err)
		}
	}
	return
}

// 注册发送指定队列指定路由的生产者
func (r *RabbitMQ) RegisterProducer(producer Producer) {
	r.ProducerList = append(r.ProducerList, producer)
}

func (r *RabbitMQ) DeleteRegisterProducer(producer Producer){
	r.Mu.Lock()
	defer r.Mu.Unlock()
	if len(r.ProducerList)>0{
		if r.ProducerList[0] == producer{
			r.ProducerList = r.ProducerList[1:]
		}
	}
}


// 注册接收指定队列指定路由的数据接收者
func (r *RabbitMQ) RegisterReceiver(receiver Receiver) {
	r.Mu.Lock()
	defer r.Mu.Unlock()
	r.ReceiverList = append(r.ReceiverList, receiver)
	fmt.Println("注册接受者+1")
}

func (r *RabbitMQ) DeleteRegisterReceiver(receiver Receiver)  {
	r.Mu.Lock()
	defer r.Mu.Unlock()
	if len(r.ReceiverList)>0{
		if r.ReceiverList[0] == receiver{
			r.ReceiverList = r.ReceiverList[1:]
		}
	}
}
package service

import (
	"banana/app/transfer/service/internal/biz"
	"banana/app/transfer/service/internal/data"
	"context"
	"encoding/json"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/streadway/amqp"
	"os"
	"time"
)

func NewMqService(mq *data.RabbitMQ, db *data.Data) *MqService {
	return &MqService{
		mq: mq,
		db: db,
	}
}

type Receive struct {
	msgContent string
}

// 实现发送者
func (t *Receive) MsgContent() string {
	return t.msgContent
}

// 实现接收者
func (t *Receive) Consumer(dataByte []byte) (moj *data.MessageObject, err error) {
	moj = &data.MessageObject{}
	err = json.Unmarshal(dataByte, moj)
	if err != nil {
		return nil, err
	}
	return moj, nil
}

func (c *MqService) Start() {
	for _, producer := range c.mq.ProducerList {
		fmt.Printf("producer len:%d", len(c.mq.ProducerList))
		go c.TransferProduce(producer)
	}

	// 开启监听接收者接收任务
	for _, receiver := range c.mq.ReceiverList {
		fmt.Printf("recv len:%d", len(c.mq.ReceiverList))
		go c.TransferConsume(receiver)
	}
	time.AfterFunc(1*time.Second,c.Start)

}
func (c *MqService) TransferProduce(producer data.Producer) {
	// 用于检查交换机是否存在,已经存在不需要重复声明
	err := c.mq.Channel.ExchangeDeclarePassive(c.mq.ExchangeName, c.mq.ExchangeType, true, false, false, true, nil)
	if err != nil {
		// 注册交换机
		// name:交换机名称,kind:交换机类型,durable:是否持久化,队列存盘,true服务重启后信息不会丢失,影响性能;autoDelete:是否自动删除;
		// noWait:是否非阻塞, true为是,不等待RMQ返回信息;args:参数,传nil即可; internal:是否为内部
		err = c.mq.Channel.ExchangeDeclare(c.mq.ExchangeName, c.mq.ExchangeType, true, false, false, true, nil)
		if err != nil {
			fmt.Printf("MQ注册交换机失败:%s \n", err)
			return
		}
	}
	// 用于检查队列是否存在,已经存在不需要重复声明
	_, err = c.mq.Channel.QueueDeclarePassive(c.mq.QueueName, true, false, false, true, nil)
	if err != nil {
		// 队列不存在,声明队列
		// name:队列名称;durable:是否持久化,队列存盘,true服务重启后信息不会丢失,影响性能;autoDelete:是否自动删除;noWait:是否非阻塞,
		// true为是,不等待RMQ返回信息;args:参数,传nil即可;exclusive:是否设置排他
		_, err = c.mq.Channel.QueueDeclare(c.mq.QueueName, true, false, false, true, nil)
		if err != nil {
			fmt.Printf("MQ注册队列失败:%s \n", err)
			return
		}
	}

	// 队列绑定
	err = c.mq.Channel.QueueBind(c.mq.QueueName, c.mq.RoutingKey, c.mq.ExchangeName, true, nil)
	if err != nil {
		fmt.Printf("MQ绑定队列失败:%s \n", err)
		return
	}

	// 发送任务消息
	err = c.mq.Channel.Publish(c.mq.ExchangeName, c.mq.RoutingKey, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(producer.MsgContent()),
	})
	if err != nil {
		fmt.Printf("MQ任务发送失败:%s \n", err)
		return
	}
	rec := &Receive{}
	defer c.mq.RegisterReceiver(rec)
	defer c.mq.DeleteRegisterProducer(producer)

}
func (c *MqService) TransferConsume(receiver data.Receiver) {
	c.mq.DeleteRegisterReceiver(receiver)
	// 用于检查队列是否存在,已经存在不需要重复声明
	_, err := c.mq.Channel.QueueDeclarePassive(c.mq.QueueName, true, false, false, true, nil)
	if err != nil {
		// 队列不存在,声明队列
		// name:队列名称;durable:是否持久化,队列存盘,true服务重启后信息不会丢失,影响性能;autoDelete:是否自动删除;noWait:是否非阻塞,
		// true为是,不等待RMQ返回信息;args:参数,传nil即可;exclusive:是否设置排他
		_, err = c.mq.Channel.QueueDeclare(c.mq.QueueName, true, false, false, true, nil)
		if err != nil {
			fmt.Printf("MQ注册队列失败:%s \n", err)
			return
		}
	}
	// 绑定任务
	err = c.mq.Channel.QueueBind(c.mq.QueueName, c.mq.RoutingKey, c.mq.ExchangeName, true, nil)
	if err != nil {
		fmt.Printf("绑定队列失败:%s \n", err)
		return
	}
	// 获取消费通道,确保rabbitMQ一个一个发送消息
	err = c.mq.Channel.Qos(1, 0, true)
	msgList, err := c.mq.Channel.Consume(c.mq.QueueName, "", false, false, false, false, nil)
	if err != nil {
		fmt.Printf("获取消费通道异常:%s \n", err)
		return
	}
	var minioUpload = func(bucket, objectName, filePath, contentType string) (minio.UploadInfo, error) {
		fileinfo := minio.UploadInfo{}
		client := c.db.Minio_internal
		fileinfo, err = client.FPutObject(context.TODO(), bucket, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
		if err != nil {
			return fileinfo, err
		}
		return fileinfo, err
	}
	var dealTransfer = func(obj *data.MessageObject) error {
		f, err := minioUpload(obj.Bucket, obj.FilePath, obj.FileName, obj.ContentType)
		if err != nil {
			fmt.Printf("上传失败 待重试：%v\n", err)
			return err
		}
		var ud  biz.File
		ud.Finish = 1
		ud.FileSize = f.Size
		err = c.db.Db.Model(&biz.File{}).Where("id = ?", obj.Fid).Updates(ud).Error
		if err != nil {
			return err
		}
		return nil
	}


	//go func() {
	for msg := range msgList {
		// 处理数据
		obj, err := receiver.Consumer(msg.Body)
		fmt.Printf("打印错误：%v\n", err)
		fmt.Printf("获取到的结构体：%v\n", obj)
		if err == nil {
			err = dealTransfer(obj)
			if err != nil {
				fmt.Printf("消息消费未完成:%s \n", err)
				return
			}
			os.Remove(obj.FileName)
			err = msg.Ack(false)
			if err != nil {
				fmt.Printf("确认消息完成异常:%s \n", err)
				return
			}
		} else {
			err = msg.Ack(true)
			if err != nil {
				fmt.Printf("确认消息未完成异常:%s \n", err)
				return
			}
			return
		}
	}
	//}()

}

package rabbitmq_client

import (
	"github.com/streadway/amqp"
	_ "fmt"
	_ "errors"
	"log"
	"fmt"
)

type IRabbitMqConsumer interface {
	Subscribe(consumerName string, handlerFunc func(amqp.Delivery)) error
	Close()
}

type RabbitMqConsumer struct {
	rabbitMqHelper IRabbitMqHelper
	hostName string
	userName string
	pass string
	exchangeName string

	durable bool
	queueName string
	routingKey string
	enablePriority bool
	maxPriority uint
	prefetchCount uint16

	channel *amqp.Channel
	connection *amqp.Connection
}

func NewRabbitMqConsumer(hostName string, userName string, pass string,
exchangeName string, durable bool, queueName string, routingKey string, enablePriority bool,
maxPriority uint, prefetchCount uint16) (IRabbitMqConsumer, error){
	consumer := &RabbitMqConsumer{
		rabbitMqHelper : &RabbitMqHelper{},
		hostName : hostName,
		userName :userName,
		pass : pass,
		exchangeName : exchangeName,
		durable : durable,
		queueName : queueName,
		routingKey : routingKey,
 		enablePriority : enablePriority,
		maxPriority : maxPriority,
		prefetchCount : prefetchCount,
	}

	var err error
	consumer.channel, err = consumer.constructConnection()
	if err != nil {
		return nil, err
	}
	if consumer.channel == nil {
		log.Fatalf("Channel is nil")
		panic(fmt.Sprintf("Channel is nil"))
	}

	return consumer, nil
}

func (consumer RabbitMqConsumer) constructConnection() (*amqp.Channel,error){
	connection, err := consumer.rabbitMqHelper.GetRabbitMqConnection(consumer.hostName, consumer.userName, consumer.pass)
	if err != nil {
		return nil, err
	}

	ch, err := consumer.rabbitMqHelper.DeclareQueue(connection, consumer.exchangeName, consumer.durable,
		consumer.queueName, consumer.routingKey, consumer.enablePriority, consumer.maxPriority)
	if err != nil {
		return nil,err
	}

	log.Println("Connection established.")
	return ch, nil
}

func (consumer *RabbitMqConsumer) Subscribe(consumerName string, handlerFunc func(amqp.Delivery)) error {
	if consumer.channel == nil {
		log.Fatalf("Channel is nil")
		panic(fmt.Sprintf("Channel is nil"))
	}

	messages , err := consumer.channel.Consume(
		consumer.queueName, // queue
		consumerName, // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil, // args
	)
	if err != nil {
		return err
	}

	go consumeLoop(messages, handlerFunc)
	return nil
}

func (consumer *RabbitMqConsumer) Close() {
	if consumer.channel != nil {
		consumer.channel.Close()
	}

	if consumer.connection != nil {
		consumer.connection.Close()
	}
}

func consumeLoop(deliveries <-chan amqp.Delivery, handlerFunc func(d amqp.Delivery)) {
	for d := range deliveries {
		handlerFunc(d)
	}
}
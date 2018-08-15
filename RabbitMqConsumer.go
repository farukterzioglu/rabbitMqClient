package rabbitmq_client

import (
	"github.com/streadway/amqp"
	_ "fmt"
	_ "errors"
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

func NewRabbitMqConsumer(rabbitMqHelper IRabbitMqHelper, hostName string, userName string, pass string,
exchangeName string, durable bool, queueName string, routingKey string, enablePriority bool,
maxPriority uint, prefetchCount uint16) (IRabbitMqConsumer, error){
	consumer := &RabbitMqConsumer{
		rabbitMqHelper : rabbitMqHelper,
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
	err := consumer.constructConnection()
	if err != nil {
		return nil, err
	}
	return consumer, nil
}

func (consumer RabbitMqConsumer) constructConnection() error{
	connection, err := consumer.rabbitMqHelper.GetRabbitMqConnection(consumer.hostName, consumer.userName, consumer.pass)
	if err != nil {
		return err
	}

	ch, err := consumer.rabbitMqHelper.DeclareQueue(connection, consumer.exchangeName, consumer.durable,
		consumer.queueName, consumer.routingKey, consumer.enablePriority, consumer.maxPriority)
	if err != nil {
		return err
	}

	consumer.channel = ch
	return nil
}

func (consumer *RabbitMqConsumer) Subscribe(consumerName string, handlerFunc func(amqp.Delivery)) error {
	messages , err := consumer.channel.Consume(
		consumer.queueName, // queue
		consumerName, // consumer
		true, // auto-ack
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
	if consumer.connection != nil {
		consumer.connection.Close()
	}
}

func consumeLoop(deliveries <-chan amqp.Delivery, handlerFunc func(d amqp.Delivery)) {
	for d := range deliveries {
		handlerFunc(d)
	}
}
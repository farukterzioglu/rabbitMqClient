package Utilities

import (
	"github.com/streadway/amqp"
	"fmt"
	"errors"
)
type IRabbitMqHelper interface{
	GetRabbitMqConnection(hostName string,userName string ,
		pass string ) (*amqp.Connection, error)
	DeclareExchange(connection *amqp.Connection, exchangeName string ,
		exchangeType string , durable bool ) (*amqp.Channel, error)
	DeclareQueue(connection *amqp.Connection, exchangeName string ,
		durable bool , queueName string , routingKey string ,
		enablePriority bool, maxPriority uint) (*amqp.Channel, error)
}

type RabbitMqHelper struct{}
func (helper *RabbitMqHelper) GetRabbitMqConnection(hostName string,userName string ,pass string ) (*amqp.Connection, error){
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/", userName, pass, hostName))
	return conn, err
}

func (helper *RabbitMqHelper) DeclareExchange(connection *amqp.Connection, exchangeName string ,
	exchangeType string , durable bool ) (*amqp.Channel, error){
	ch, err := connection.Channel()
	if err != nil {
		return nil, err
	}

	err = ch.ExchangeDeclare(
		exchangeName,   // name
		exchangeType, // type
		durable,     // durable
		  false,// auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	return ch, err
}

func (helper *RabbitMqHelper) DeclareQueue(connection *amqp.Connection, exchangeName string ,
	durable bool , queueName string , routingKey string ,
	enablePriority bool, maxPriority uint) (*amqp.Channel, error){
	ch, err := connection.Channel()
	if err != nil {
		return nil, err
	}

	//TODO : Priority ???
	if enablePriority == true {
		err = errors.New("Not implemetned")
		return nil, err
	}

	q, err := ch.QueueDeclare(
		queueName,    // name
		durable, // durable
		false, // delete when usused
		false,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, err
	}

	err = ch.QueueBind(
		q.Name, // queue name
		routingKey,     // routing key
		exchangeName, // exchange
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return ch, nil
}
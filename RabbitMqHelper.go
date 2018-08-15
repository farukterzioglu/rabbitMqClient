package rabbitmq_client

import (
	"github.com/streadway/amqp"
	"fmt"
)
type IRabbitMqHelper interface{
	GetRabbitMqConnection(hostName string,userName string ,pass string ) (*amqp.Connection, error)
	IModel DeclareExchange(IConnection connection, string exchangeName, string exchangeType, bool durable)
}

type RabbitMqHelper struct{}
func (helper *RabbitMqHelper) GetRabbitMqConnection(hostName string,userName string ,pass string ) (*amqp.Connection, error){
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/", userName, pass, hostName))
	return conn, err
}



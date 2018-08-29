package rabbitMqClient

import (
	"fmt"
	"github.com/farukterzioglu/rabbitMqClient/Utilities"
	"github.com/farukterzioglu/rabbitMqClient/Utilities/SampleDir"
	"github.com/farukterzioglu/rabbitMqClient/SampleDir/SampleSubdir"
	"github.com/streadway/amqp"
	"log"
)

type IRabbitMqPublisher interface {
	Enqueue(message string, routingKey string) error
	Close()
}

type RabbitMqPublisher struct {
	rabbitMqHelper Utilities.IRabbitMqHelper
	hostName string
	userName string
	pass string
	exchangeName string
	exchangeType string
	durable bool

	channel *amqp.Channel
	connection *amqp.Connection
}

func NewRabbitMqPublisher(hostName string, userName string, pass string,
	exchangeName string, exchangeType string, durable bool) (IRabbitMqPublisher, error){
	publisher := &RabbitMqPublisher{
		rabbitMqHelper : &Utilities.RabbitMqHelper{},
		hostName : hostName,
		userName :userName,
		pass : pass,
		exchangeName : exchangeName,
		exchangeType : exchangeType,
		durable : durable,
	}

	var err error
	publisher.channel, err = publisher.constructConnection()
	if err != nil {
		return nil, err
	}
	if publisher.channel == nil {
		log.Fatalf("Channel is nil")
		panic(fmt.Sprintf("Channel is nil"))
	}

	return publisher, nil
}


func (publisher RabbitMqPublisher) constructConnection() (*amqp.Channel,error){
	connection, err := publisher.rabbitMqHelper.GetRabbitMqConnection(publisher.hostName,
		publisher.userName, publisher.pass)
	if err != nil {
		return nil, err
	}

	ch, err := publisher.rabbitMqHelper.DeclareExchange(connection,publisher.exchangeName,
		publisher.exchangeType, publisher.durable )
	if err != nil {
		return nil,err
	}

	log.Println("Connection established.")
	return ch, nil
}

func (publisher *RabbitMqPublisher) Enqueue(message string, routingKey string ) error {
	err := publisher.channel.Publish(
		publisher.exchangeName,     // exchange
		routingKey, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing {
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	if err != nil {
		return err
	}
	return nil
}
func (publisher *RabbitMqPublisher) Close() {
	if publisher.channel != nil {
		publisher.channel.Close()
	}

	if publisher.connection != nil {
		publisher.connection.Close()
	}
}

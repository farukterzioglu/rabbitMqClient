package service

import (
	"net/http"
	"log"
	"github.com/farukterzioglu/rabbitMqClient"
	"fmt"
)

type (
	ServerSettings struct {
		Port string
	}
	RabbitMqSettings struct{
		HostName string
		UserName string
		Password string
		ExchangeName string
		ExchangeType string
		Durable bool
		RoutingKey string
	}
)

var rabbitMqSettings RabbitMqSettings
var publisher rabbitMqClient.IRabbitMqPublisher

func StartWebServer(serverSettings ServerSettings, _rabbitMqSettings RabbitMqSettings){
	rabbitMqSettings = _rabbitMqSettings

	//Publisher
	_publisher, err := rabbitMqClient.NewRabbitMqPublisher(
		rabbitMqSettings.HostName, rabbitMqSettings.UserName, rabbitMqSettings.Password,
		rabbitMqSettings.ExchangeName, rabbitMqSettings.ExchangeType, rabbitMqSettings.Durable)
	if err != nil {
		panic(fmt.Sprintf("Couldn't create publisher : %s", err))
	}
	publisher = _publisher

	//Web App
	router := NewRouter()
	http.Handle("/", router)

	log.Println("Starting HTTP service at " + serverSettings.Port)
	err = http.ListenAndServe(":" + serverSettings.Port, nil)

	if err != nil {
		log.Println("An error occured starting HTTP listener at port " + serverSettings.Port)
		log.Println("Error: " + err.Error())
	}
}

package service

import (
	"net/http"
	"log"
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
func StartWebServer(serverSettings ServerSettings, _rabbitMqSettings RabbitMqSettings){
	rabbitMqSettings = _rabbitMqSettings

	router := NewRouter()
	http.Handle("/", router)

	log.Println("Starting HTTP service at " + serverSettings.Port)
	err := http.ListenAndServe(":" + serverSettings.Port, nil)

	if err != nil {
		log.Println("An error occured starting HTTP listener at port " + serverSettings.Port)
		log.Println("Error: " + err.Error())
	}
}

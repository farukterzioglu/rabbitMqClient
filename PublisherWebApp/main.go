package main

import (
	"flag"
	"fmt"
	"github.com/farukterzioglu/rabbitMqClient/PublisherWebApp/service"
	"github.com/farukterzioglu/rabbitMqClient/Utilities"
	"log"
	"strconv"
)

var appName = "PublisherWebApp"


var (
	consulUrl      = flag.String("consulUrl", "", "Consul Url")
	consulPath     = flag.String("consulPath", "", "Consul Path")
)

var rabbitMqSettings service.RabbitMqSettings
func init() {
	flag.Parse()
	rabbitMqSettings = readRabbitMqSettings()
}


func readRabbitMqSettings() service.RabbitMqSettings {
	if *consulUrl == "" || *consulPath == "" {
		panic(fmt.Sprintf("Couldn't read Consul url "))
	}

	var consulHelper Utilities.IConsulHelper
	consulHelper, err := Utilities.NewConsulHelper(*consulUrl, *consulPath)
	if err != nil {
		panic(fmt.Sprintf("Couldn't connect to Consul: %s / %s \n", *consulUrl, *consulPath))
	}

	var durable bool
	durable, err = strconv.ParseBool(consulHelper.GetValue("durable"))
	if err != nil {
		log.Fatalf("Wrong format for 'durable' value. Set to 'false'")
		durable = false
	}

	rabbitMqSettings := service.RabbitMqSettings{
		HostName : consulHelper.GetValue("hostName"),
		UserName : consulHelper.GetValue("userName"),
		Password : consulHelper.GetValue("password"),
		ExchangeName : consulHelper.GetValue("exchangeName"),
		ExchangeType : consulHelper.GetValue("exchangeType"),
		Durable : durable,
		RoutingKey : consulHelper.GetValue("routingKey"),
	}
	return rabbitMqSettings
}

func main(){
	fmt.Printf("Starting %v\n", appName)
	service.StartWebServer(service.ServerSettings{ Port : "6767"},
	rabbitMqSettings)
}

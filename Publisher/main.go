package main

import (
	"flag"
	"fmt"
	"github.com/farukterzioglu/rabbitMqClient"
	"github.com/farukterzioglu/rabbitMqClient/Utilities"
	"log"
	"os"
	"os/signal"
	"strconv"
	"time"
)

var (
	hostName       = flag.String("hostName", "localhost:5672/", "Host Name")
	userName       = flag.String("userName", "guest", "RabbitMq user name")
	password       = flag.String("password", "guest", "RabbitMq password")
	exchangeName   = flag.String("exchangeName", "SampleExchange", "Exchange name")
	exchangeType   = flag.String("exchangeType", "fanout", "Exchange type")
	durable        = flag.Bool("durable", false, "Durable")
	routingKey     = flag.String("routingKey", "", "Routing key")

	consulSettings = flag.Bool("consulSettings", true, "Settings from Consul")
	consulUrl      = flag.String("consulUrl", "", "Consul Url")
	consulPath     = flag.String("consulPath", "", "Consul Path")
)

func init() {
	flag.Parse()

	if *consulSettings {
		readConsulSettings()
	}
}

func readConsulSettings() {
	if *consulUrl == "" || *consulPath == "" {
		panic(fmt.Sprintf("Couldn't read Consul url "))
	}

	var consulHelper Utilities.IConsulHelper
	consulHelper, err := Utilities.NewConsulHelper(*consulUrl, *consulPath)
	if err != nil {
		panic(fmt.Sprintf("Couldn't connect to Consul: %s / %s \n", *consulUrl, *consulPath))
	}

	*hostName = consulHelper.GetValue("hostName")
	*userName = consulHelper.GetValue("userName")
	*password = consulHelper.GetValue("password")
	*exchangeName = consulHelper.GetValue("exchangeName")
	*exchangeType = consulHelper.GetValue("exchangeType")
	*durable, err = strconv.ParseBool(consulHelper.GetValue("durable"))
	if err != nil {
		log.Fatalf("Wrong format for 'durable' value. Set to 'false'")
		*durable = false
	}
	*routingKey = consulHelper.GetValue("routingKey")
}

func main() {
	println("Consumer started...")

	var publisher rabbitMqClient.IRabbitMqPublisher
	var err error
	publisher, err = rabbitMqClient.NewRabbitMqPublisher(
		*hostName, *userName, *password,
		*exchangeName, *exchangeType, *durable)
	failOnError(err, "Failed to create new consumer")

	go startPublishing(publisher)

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	publisher.Close()
}

func startPublishing(publisher rabbitMqClient.IRabbitMqPublisher){
	count := 0
	for {
		count++
		time.Sleep(1 * time.Second)
		publisher.Enqueue(fmt.Sprintf("Test message %d", count), *routingKey)
		fmt.Printf(".")
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

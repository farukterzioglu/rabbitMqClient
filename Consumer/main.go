package main

import (
	"flag"
	"fmt"
	"github.com/farukterzioglu/rabbitMqClient"
	"github.com/farukterzioglu/rabbitMqClient/Utilities"
	"github.com/streadway/amqp"
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
	queueName      = flag.String("queueName", "sampleQueue", "Queue name")
	routingKey     = flag.String("routingKey", "", "Routing key")
	enablePriority = flag.Bool("enablePriority", false, "Enable priority")
	maxPriority    = flag.Uint("maxPriority", 0, "Max priority")
	prefetchCount  = flag.Uint("prefetchCount", 1, "Prefetch count")

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
	*queueName = consulHelper.GetValue("queueName")
	*routingKey = consulHelper.GetValue("routingKey")
	*enablePriority, err = strconv.ParseBool(consulHelper.GetValue("enablePriority"))
	if err != nil {
		log.Fatalf("Wrong format for 'enablePriority' value. Set to 'false'")
		*enablePriority = false
	}
	maxPri, err := strconv.ParseUint(consulHelper.GetValue("maxPriority"), 10, 64)
	if err != nil {
		maxPri = 0
	}
	*maxPriority = uint(maxPri)

	pref, err := strconv.ParseUint(consulHelper.GetValue("maxPriority"), 10, 64)
	if err != nil {
		pref = 1
	}
	*prefetchCount = uint(pref)
}

func main() {
	println("Consumer started...")

	var consumer rabbitMqClient.IRabbitMqConsumer
	var err error
	consumer, err = rabbitMqClient.NewRabbitMqConsumer(
		*hostName, *userName, *password,
		*exchangeName, *exchangeType, *durable,
		*queueName, *routingKey,
		*enablePriority, *maxPriority, uint16(*prefetchCount))
	failOnError(err, "Failed to create new consumer")

	err = consumer.Subscribe("test consumer", onMessage)
	failOnError(err, "Failed to subscribe")

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	consumer.Close()
}

func onMessage(delivery amqp.Delivery) {
	fmt.Printf("%v\n", string(delivery.Body))
	time.Sleep(3 * time.Second)
	delivery.Ack(false)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

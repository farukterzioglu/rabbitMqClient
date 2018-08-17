package main

import(
	"github.com/farukterzioglu/rabbitMqClient"
	"github.com/farukterzioglu/rabbitMqClient/Utilities"
	"github.com/streadway/amqp"
	"fmt"
	"os"
	"os/signal"
	"log"
	"flag"
	"time"
)

var (
	hostName     = flag.String("hostName", "localhost:5672/", "Host Name")
	userName	 = flag.String("userName", "guest", "RabbitMq user name")
	password	 = flag.String("password", "guest", "RabbitMq password")
	exchangeName = flag.String("exchangeName", "SampleExchange", "Exchange name")
	exchangeType = flag.String("exchangeType", "fanout", "Exchange type")
	durable      = flag.Bool("durable", false, "Durable")
	queueName    = flag.String("queueName", "sampleQueue", "Queue name")
	routingKey     = flag.String("routingKey", "", "Routing key")
	enablePriority = flag.Bool("enablePriority", false, "Enable priority")
	maxPriority    = flag.Uint("maxPriority", 0, "Max priority")
	prefetchCount  = flag.Uint("prefetchCount", 1, "Prefetch count")

	consulUrl  = flag.String("Consul Url", "https://demo.consul.io/ui/dc1/kv", "Consul Url")
	consulPath  = flag.String("Consul Path", "rabbitMqConsumerGoLang", "Consul Path")
)

func init() {
	flag.Parse()

	var consulHelper Utilities.IConsulHelper
	consulHelper, err := Utilities.NewConsulHelper(*consulUrl, *consulPath)
	if err != nil{
		panic(fmt.Sprintf("Couldn't connect to Consul: %s / %s \n", *consulUrl, *consulPath))
	}

	value := consulHelper.GetValue("testKey")
	fmt.Printf("Value from Consul : %s \n", value)
}

func main(){
	return //TODO : testing, remove this
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
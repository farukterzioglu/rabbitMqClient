package main

import(
	"github.com/farukterzioglu/rabbitmq-client"
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
	exchangeName = flag.String("exchangeName", "exchangeName", "Exchange name")
	durable      = flag.Bool("durable", false, "Durable")
	queueName    = flag.String("queueName", "queueName", "Queue name")
	routingKey     = flag.String("routingKey", "routingKey", "Routing key")
	enablePriority = flag.Bool("enablePriority", false, "Enable priority")
	maxPriority    = flag.Uint("maxPriority", 0, "Max priority")
	prefetchCount  = flag.Uint("prefetchCount", 1, "Prefetch count")
)

func init() {
	flag.Parse()
}

func main(){
	println("Consumer started...")

	var consumer rabbitmq_client.IRabbitMqConsumer
	var err error
	consumer, err = rabbitmq_client.NewRabbitMqConsumer(
		*hostName, *userName, *password,
		*exchangeName, *durable,
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
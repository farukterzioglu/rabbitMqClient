package main

import(
	"github.com/farukterzioglu/rabbitmq-client"
	"github.com/streadway/amqp"
	"fmt"
	"os"
	"os/signal"
	"log"
)

func onMessage(delivery amqp.Delivery) {
	fmt.Printf("Got a message: %v\n", string(delivery.Body))
}

func main(){
	println("Consumer started...")

	var consumer rabbitmq_client.IRabbitMqConsumer
	var err error
	consumer, err = rabbitmq_client.NewRabbitMqConsumer(
		//TODO : 
		false, 0, 50)
	failOnError(err, "Failed to create new consumer")

	err = consumer.Subscribe("test consumer", onMessage)
	failOnError(err, "Failed to subscribe")

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	consumer.Close()
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
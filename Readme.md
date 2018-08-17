**RabbitMQ on Docker**\
docker run -d --hostname my-rabbit --name some-rabbit -p 4369:4369 -p 5671:5671 -p 5672:5672 -p 15672:15672 rabbitmq\
docker exec some-rabbit rabbitmq-plugins enable rabbitmq_management\

**Login at http://localhost:15672/ (or the IP of your docker host)**\
using guest/guest\

cd Consumer\

**Sample usage**\
go run main.go -hostName=[hostName] -userName=[userName] -password=[password] -exchangeName=[exchangeName] -queueName=[queueName] -routingKey=[routingKey] -prefetchCount=[prefetchCount]\

**Connect to RabbitMQ on Docker**\
go run main.go -hostName=localhost:5672 -userName=guest -password=guest -exchangeName=ProductExchange -queueName=TestingQueue -routingKey=TestRoute -prefetchCount=1\

**Consul settings**\
https://demo.consul.io/ui/dc1/kv/rabbitMqConsumerGoLang/\


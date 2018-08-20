### Start RabbitMQ ###
docker run -d --hostname my-rabbit --name some-rabbit -p 4369:4369 -p 5671:5671 -p 5672:5672 -p 15672:15672 rabbitmq
docker exec some-rabbit rabbitmq-plugins enable rabbitmq_management

Login at http://localhost:15672/ (or the IP of your docker host)  
username : guest, password : guest

### Consumer ###
cd Consumer

*Sample usage*  
go run main.go -consulSettings=false -hostName=[hostName] -userName=[userName] -password=[password] -exchangeName=[exchangeName] -queueName=[queueName] -routingKey=[routingKey] -prefetchCount=[prefetchCount]\

### Consume RabbitMQ on Docker ###
go run main.go -consulSettings=false -hostName=localhost:5672 -userName=guest -password=guest -exchangeName=ProductExchange -queueName=TestingQueue -routingKey=TestRoute -prefetchCount=1  

*With Consul*  
go run main.go -consulSettings=true -consulUrl=https://demo.consul.io/ui/dc1/kv -consulPath=rabbitMqConsumerGoLang  

### Consul settings ###
https://demo.consul.io/ui/dc1/kv/rabbitMqConsumerGoLang/

### Publisher ###
cd Publisher  

* Publish to RabbitMQ on Docker with Consul settings *
go run main.go -consulSettings=true -consulUrl=https://demo.consul.io/ui/dc1/kv -consulPath=rabbitMqConsumerGoLang

* Publish with web api *
cd PublisherWebApp  
go run main.go -consulUrl=https://demo.consul.io/ui/dc1/kv -consulPath=rabbitMqConsumerGoLang  
curl -H "Content-Type: application/json" -X POST -d '{ "id": "1", "text": "Testing..." }' http://localhost:6767/publish

* Run 'Publisher Web Api' with 'Make' *
make run
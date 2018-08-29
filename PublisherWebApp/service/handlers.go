package service

import (
	"encoding/json"
	"fmt"
	"github.com/farukterzioglu/rabbitMqClient/Utilities"
	"net/http"
	"publisherWebApp/service/model"
	"strconv"
)

func PublishHandler(w http.ResponseWriter, r *http.Request){
	var message model.Message
	_ = json.NewDecoder(r.Body).Decode(&message)

	data, _ := json.Marshal(message)
	publisher.Enqueue(fmt.Sprintf("%s", data) , rabbitMqSettings.RoutingKey)

	response, _ := json.Marshal( struct { Id string }{ message.Id})
	writeJsonResponse(w, http.StatusOK, response)
}

func HealthCheck(w http.ResponseWriter, r *http.Request){
	var rabbitMqHelper Utilities.IRabbitMqHelper
	rabbitMqHelper = &Utilities.RabbitMqHelper{}

	conn, err := rabbitMqHelper.GetRabbitMqConnection(
		rabbitMqSettings.HostName, rabbitMqSettings.UserName, rabbitMqSettings.Password)

	if err != nil || conn == nil{
		data, _ := json.Marshal(healthCheckResponse{ Status : "RabbitMQ is unaccessible"})
		writeJsonResponse(w, http.StatusServiceUnavailable, data)
	} else {
		data, _ := json.Marshal(healthCheckResponse{ Status : "UP"})
		writeJsonResponse(w, http.StatusOK, data)
	}
}

func writeJsonResponse(w http.ResponseWriter, status int, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.WriteHeader(status)
	w.Write(data)
}

type healthCheckResponse struct {
	Status string `json:"status"`
}
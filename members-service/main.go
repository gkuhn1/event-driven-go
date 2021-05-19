package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gkuhn1/event-driven-go/utils"
)

func main() {

	go initConsumer()

	fmt.Println("Running members-service")

	api := NewAPI()
	r := mux.NewRouter()
	r.HandleFunc("/members", api.handler).Methods("POST")

	utils.RunServer(r)
}

type API struct {
	Producer *utils.Producer
	Handler  func(w http.ResponseWriter, r *http.Request)
}

func NewAPI() *API {
	return &API{
		Producer: utils.NewProducer(),
	}
}

func (api *API) handler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	fmt.Println("POST /members => ", name)

	response := fmt.Sprintf("{\"member\":{\"name\": \"%s\"}}", name)
	w.Write([]byte(response))

	api.Producer.ProduceMessage(name, "new_member")
	return
}

func callback(m *kafka.Message) {
	fmt.Println("Callback ....")
}

func initConsumer() {
	// utils.Consume([]string{"new_member"}, "members-service-consumer-events", callback)
	return
}

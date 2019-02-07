package main

import (
	"fmt"
	"net/http"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gkuhn1/event-driven-go/utils"
	"github.com/gorilla/mux"
)

func main() {

	go initConsumer()

	fmt.Println("Running members-service")

	initAPI()
}

func callback(m *kafka.Message) {
	fmt.Println("Callback ....")
}

func initConsumer() {
	// utils.Consume([]string{"new_member"}, "members-service-consumer-events", callback)
	return
}

func handler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	fmt.Println("POST /members => ", name)

	response := fmt.Sprintf("{\"member\":{\"name\": \"%s\"}}", name)
	w.Write([]byte(response))

	utils.ProduceMessage(name, "new_member")
	return
}

func initAPI() {
	r := mux.NewRouter()
	r.HandleFunc("/members", handler).Methods("POST")
	utils.RunServer(r)
}

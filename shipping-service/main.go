package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gkuhn1/event-driven-go/utils"
	"github.com/gorilla/mux"
)

func main() {

	go initConsumer()

	fmt.Println("Running shipping-service")

	initAPI()
}

func callback(m *kafka.Message) {
	v := string(m.Value)
	switch *m.TopicPartition.Topic {
	case "new_member":
		fmt.Println("New member received: ", v)
	case "new_shipment":
		fmt.Println("Got new shippiment authorized for order: ", v)
		fmt.Println("Shipping order: ", v)
		time.Sleep(3 * time.Second)
		utils.ProduceMessage(v, "shipping_confirmed")
	default:
		fmt.Printf("Callback .... %s\n", *m.TopicPartition.Topic)
	}

}

func initConsumer() {
	utils.Consume([]string{"new_member", "new_shipment"}, "shipping-service-consumer-events", callback)
	return
}

func handler(w http.ResponseWriter, r *http.Request) {
	// name := r.FormValue("name")
	// fmt.Println("Got request for ", name)

	// response := fmt.Sprintf("{\"member\":{\"name\": \"%s\"}}", name)
	// w.Write([]byte(response))

	// utils.ProduceMessage(response, "new_member")
	return
}

func initAPI() {
	r := mux.NewRouter()
	r.HandleFunc("/shipping", handler).Methods("POST")
	utils.RunServer(r)
}

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

	fmt.Println("Running payments-service")

	initAPI()
}

func callback(m *kafka.Message) {
	v := string(m.Value)
	switch *m.TopicPartition.Topic {
	case "new_member":
		fmt.Println("New member received: ", v)
	case "new_order":
		fmt.Println("New order received: ", v)
		fmt.Println("Authorizing Payment...")
		time.Sleep(2 * time.Second)
		fmt.Printf("Payment for order %s authorized!\n", v)
		utils.ProduceMessage(v, "payment_authorized")
	case "capture_payment":
		fmt.Println("Capturing payment for order: ", v)
		time.Sleep(2 * time.Second)
		fmt.Println("Payment captured successfully.")
		utils.ProduceMessage(v, "order_completed")
	default:
		fmt.Printf("Callback .... %s\n", *m.TopicPartition.Topic)
	}
}

func initConsumer() {
	utils.Consume([]string{"new_member", "new_order", "capture_payment"}, "payments-service-consumer-events", callback)
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
	r.HandleFunc("/payments", handler).Methods("POST")
	utils.RunServer(r)
}

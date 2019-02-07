package main

import (
	"fmt"
	"net/http"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gkuhn1/event-driven-go/utils"
	"github.com/gorilla/mux"
)

var orderNum int

func main() {

	go initConsumer()

	fmt.Println("Running orders-service")

	initAPI()
}

func callback(m *kafka.Message) {
	v := string(m.Value)
	switch *m.TopicPartition.Topic {
	case "new_member":
		fmt.Println("New member received: ", v)
	case "payment_authorized":
		fmt.Println("Got payment authorized for: ", v)
		fmt.Println("Shipping order: ", v)
		utils.ProduceMessage(v, "new_shipment")
	case "shipping_confirmed":
		fmt.Println("Order shipment confirmed: ", v)
		fmt.Println("Capturing payment...")
		utils.ProduceMessage(v, "capture_payment")
	case "order_completed":
		fmt.Printf("Order %s complete!\n", v)
	default:
		fmt.Printf("Callback .... %s\n", *m.TopicPartition.Topic)
	}
}

func initConsumer() {
	utils.Consume([]string{"new_member", "payment_authorized", "shipping_confirmed", "order_completed"}, "orders-service-consumer-events", callback)
	return
}

func handler(w http.ResponseWriter, r *http.Request) {
	orderNum++
	fmt.Println("POST /orders => ", orderNum)

	response := fmt.Sprintf("{\"order\":{\"num\": \"%d\"}}", orderNum)
	w.Write([]byte(response))

	utils.ProduceMessage(fmt.Sprintf("%d", orderNum), "new_order")
	return
}

func initAPI() {
	r := mux.NewRouter()
	r.HandleFunc("/orders", handler).Methods("POST")
	utils.RunServer(r)
}

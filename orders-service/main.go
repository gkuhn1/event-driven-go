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

	fmt.Println("Running orders-service")

	server := NewServer()
	go server.initConsumer()
	// blocking call
	utils.RunServer(server)
}

type Server struct {
	Producer *utils.Producer
	router   *mux.Router
}

func NewServer() *Server {
	s := &Server{
		Producer: utils.NewProducer(),
		router:   mux.NewRouter(),
	}
	s.router.HandleFunc("/orders", s.handler).Methods("POST")
	return s
}

func (s *Server) GetMux() *mux.Router {
	return s.router
}

func (s *Server) Close() {
	s.Producer.Close()
}

func (s *Server) handler(w http.ResponseWriter, r *http.Request) {
	orderNum++
	fmt.Println("POST /orders => ", orderNum)

	response := fmt.Sprintf("{\"order\":{\"num\": \"%d\"}}", orderNum)
	w.Write([]byte(response))

	s.Producer.ProduceMessage(fmt.Sprintf("%d", orderNum), "new_order")
}

func (s *Server) initConsumer() {
	utils.Consume([]string{"new_member", "payment_authorized", "shipping_confirmed", "order_completed"}, "orders-service-consumer-events", s.callback)
}

func (s *Server) callback(m *kafka.Message) {
	v := string(m.Value)
	switch *m.TopicPartition.Topic {
	case "new_member":
		fmt.Println("New member received: ", v)
	case "payment_authorized":
		fmt.Println("Got payment authorized for: ", v)
		fmt.Println("Shipping order: ", v)
		s.Producer.ProduceMessage(v, "new_shipment")
	case "shipping_confirmed":
		fmt.Println("Order shipment confirmed: ", v)
		fmt.Println("Capturing payment...")
		s.Producer.ProduceMessage(v, "capture_payment")
	case "order_completed":
		fmt.Printf("Order %s complete!\n", v)
	default:
		fmt.Printf("Callback .... %s\n", *m.TopicPartition.Topic)
	}
}

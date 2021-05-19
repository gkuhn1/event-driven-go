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
	fmt.Println("Running payments-service")

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
	s.router.HandleFunc("/payments", s.handler).Methods("POST")
	return s
}

func (s *Server) GetMux() *mux.Router {
	return s.router
}

func (s *Server) Close() {
	s.Producer.Close()
}

func (s *Server) handler(w http.ResponseWriter, r *http.Request) {
}

func (s *Server) callback(m *kafka.Message) {
	v := string(m.Value)
	switch *m.TopicPartition.Topic {
	case "new_member":
		fmt.Println("New member received: ", v)
	case "new_order":
		fmt.Println("New order received: ", v)
		fmt.Println("Authorizing Payment...")
		time.Sleep(2 * time.Second)
		fmt.Printf("Payment for order %s authorized!\n", v)
		s.Producer.ProduceMessage(v, "payment_authorized")
	case "capture_payment":
		fmt.Println("Capturing payment for order: ", v)
		time.Sleep(2 * time.Second)
		fmt.Println("Payment captured successfully.")
		s.Producer.ProduceMessage(v, "order_completed")
	default:
		fmt.Printf("Callback .... %s\n", *m.TopicPartition.Topic)
	}
}

func (s *Server) initConsumer() {
	utils.Consume([]string{"new_member", "new_order", "capture_payment"}, "payments-service-consumer-events", s.callback)
}

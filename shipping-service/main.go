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
	fmt.Println("Running shipping-service")

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
	s.router.HandleFunc("/shipping", s.handler).Methods("POST")
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
	case "new_shipment":
		fmt.Println("Got new shippiment authorized for order: ", v)
		fmt.Println("Shipping order: ", v)
		time.Sleep(3 * time.Second)
		s.Producer.ProduceMessage(v, "shipping_confirmed")
	default:
		fmt.Printf("Callback .... %s\n", *m.TopicPartition.Topic)
	}
}

func (s *Server) initConsumer() {
	utils.Consume([]string{"new_member", "new_shipment"}, "shipping-service-consumer-events", s.callback)
}

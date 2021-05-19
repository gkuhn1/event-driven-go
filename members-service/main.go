package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gkuhn1/event-driven-go/utils"
)

func main() {
	fmt.Println("Running members-service")

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
	s.router.HandleFunc("/members", s.handler).Methods("POST")
	return s
}

func (s *Server) GetMux() *mux.Router {
	return s.router
}

func (s *Server) Close() {
	s.Producer.Close()
}

func (s *Server) handler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	fmt.Println("POST /members => ", name)

	response := fmt.Sprintf("{\"member\":{\"name\": \"%s\"}}", name)
	w.Write([]byte(response))

	s.Producer.ProduceMessage(name, "new_member")
}

func (s *Server) callback(m *kafka.Message) {
	v := string(m.Value)
	switch *m.TopicPartition.Topic {
	case "new_member":
		fmt.Println("New member received: ", v)
	default:
		fmt.Printf("Callback .... %s\n", *m.TopicPartition.Topic)
	}
}

func (s *Server) initConsumer() {
	utils.Consume([]string{"new_member"}, "members-service-consumer-events", s.callback)
}

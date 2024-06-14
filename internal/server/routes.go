package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	vintedscraper "vinted-scraper/internal/vinted-scraper"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", s.HelloWorldHandler)

	r.Get("/health", s.healthHandler)
	r.Get("/vintedTopic/{topic}-{order}", s.vintedTopicHandler)
	return r
}
func (s *Server) vintedTopicHandler(w http.ResponseWriter, r *http.Request) {
	topic := chi.URLParam(r, "topic")
	order := chi.URLParam(r, "order")
	fmt.Println("topic:", topic)
	topicId, err := s.db.ExistsTopic(topic)
	fmt.Println("topicId:", topicId)

	if topicId != 0 {
		go func() {
			_, err := SearchAndInsert(err, topic, order, s)
			if err != nil {
				fmt.Println("Error searching in goroutine:", err)
			}
		}()
		fmt.Println("Topic found")
		getCachedItems(s, w, topicId)
		return
	}
	result, err := SearchAndInsert(err, topic, order, s)
	response, err := json.Marshal(result)
	if err != nil {
		fmt.Println("Marshal error:", err)
		return
	}
	_, _ = w.Write(response)

}

func SearchAndInsert(err error, topic string, order string, s *Server) (vintedscraper.VintedApi_Response, error) {
	result, err := vintedscraper.Search(topic, vintedscraper.ToOrder(order), "GBP")
	if err != nil {
		fmt.Println("Error searching:", err)
		return vintedscraper.VintedApi_Response{}, nil
	}
	err = s.db.AddItems(result.Items, topic)
	if err != nil {
		fmt.Println("Adding items to database error:", err)
		return vintedscraper.VintedApi_Response{}, nil
	}
	return result, err
}

func getCachedItems(s *Server, w http.ResponseWriter, topicID int8) {
	items, err := s.db.GetItems(topicID)
	if err != nil {
		fmt.Println("Error getting items from database:", err)
		return
	}
	response, err := json.Marshal(items)
	if err != nil {
		fmt.Println("Marshal error:", err)
		return
	}
	_, _ = w.Write(response)
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, _ := json.Marshal(s.db.Health())
	_, _ = w.Write(jsonResp)
}

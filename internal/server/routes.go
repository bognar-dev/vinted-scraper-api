package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	vinted_scraper "vinted-scraper/internal/vinted-scraper"

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
	result, err := vinted_scraper.Search(topic, vinted_scraper.ToOrder(order), "GBP")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	_, _ = w.Write(result)

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

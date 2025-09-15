package routes

import "github.com/harshitrajsinha/medi-go/internal/store"

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type APIRoutes struct {
	service *store.Store
}

func NewAPIRoutes(service *store.Store) *APIRoutes {
	return &APIRoutes{
		service: service,
	}
}

package usecase

import "http-proxy/internal/proxy/models"

type ProxyUsecase interface {
	Handle() error
	Close()
}

type HistoryUsecase interface {
	GetRequests() ([]models.RequestResponse, error)
	GetRequestByID(id int) (*models.RequestResponse, error)
}

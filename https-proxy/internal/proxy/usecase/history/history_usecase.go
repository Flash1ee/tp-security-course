package history

import (
	"github.com/sirupsen/logrus"

	"http-proxy/internal/proxy/models"
	"http-proxy/internal/proxy/repository"
)

type Usecase struct {
	logger *logrus.Logger
	repo   *repository.ProxyRepository
}

func NewHistoryUsecase(logger *logrus.Logger, repo *repository.ProxyRepository) *Usecase {
	return &Usecase{
		logger: logger,
		repo:   repo,
	}
}
func (u *Usecase) GetRequests() ([]models.RequestResponse, error) {
	res, err := u.repo.GetAllRequests()
	if err != nil {
		u.logger.Error("history usecase, method GetRequests: ", err.Error())
	}
	return res, nil
}
func (u *Usecase) GetRequestByID(id int) (*models.RequestResponse, error) {
	res, err := u.repo.GetRequestByID(id)
	if err != nil {
		u.logger.Error("history usecase, method GetRequestByID: ", err.Error())
	}
	return res, nil
}

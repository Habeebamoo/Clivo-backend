package services

import (
	"fmt"

	"github.com/Habeebamoo/Clivo/server/internal/models"
	"github.com/Habeebamoo/Clivo/server/internal/repositories"
)

type AppealService interface {
	GetAppealStatus(string) (bool, int, error)
	CreateAppeal(models.AppealRequest) (int, error)
}

type AppealSvc struct {
	repo repositories.AppealRepository
	userRepo repositories.UserRepository
}

func NewAppealService(repo repositories.AppealRepository) AppealService {
	return &AppealSvc{repo: repo}
}

func (as *AppealSvc) GetAppealStatus(userId string) (bool, int, error) {
	user, code, err := as.userRepo.GetUserById(userId)
	if err != nil {
		return false, code, err
	}

	return user.IsBanned, 200, nil
}

func (as *AppealSvc) CreateAppeal(appealReq models.AppealRequest) (int, error) {
	user, code, err := as.userRepo.GetUserById(appealReq.UserId)
	if err != nil {
		return code, err
	}

	//check user's appeals limit (3)
	appeals, code, err := as.repo.GetUserAppeals(appealReq.UserId)
	if err != nil {
		return code, err
	}

	if len(appeals) >= 3 {
		return code, fmt.Errorf("You have exceeded your appeal limit")
	}

	appeal := models.Appeal{
		UserId: appealReq.UserId,
		Name: user.Name,
		Picture: user.Picture,
		Username: user.Username,
		Message: appealReq.Message,
	}

	//call repo
	code, err = as.repo.CreateAppeal(appeal)
	if err != nil {
		return code, err
	}

	//notify user & admin

	//respond
	return code, err
}

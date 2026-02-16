package services

import (
	"fmt"

	"github.com/Habeebamoo/Clivo/server/internal/models"
)

func (us *UserSvc) GetAppealStatus(userId string) (bool, int, error) {
	user, code, err := us.repo.GetUserById(userId)
	if err != nil {
		return false, code, err
	}

	return user.IsBanned, 200, nil
}

func (us *UserSvc) CreateAppeal(appealReq models.AppealRequest) (int, error) {
	user, code, err := us.repo.GetUserById(appealReq.UserId)
	if err != nil {
		return code, err
	}

	//check user's appeals limit (3)
	appeals, code, err := us.repo.GetUserAppeals(appealReq.UserId)
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
	code, err = us.repo.CreateAppeal(appeal)
	if err != nil {
		return code, err
	}

	//notify user & admin

	//respond
	return code, err
}
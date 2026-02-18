package repositories

import (
	"fmt"

	"github.com/Habeebamoo/Clivo/server/internal/models"
	"gorm.io/gorm"
)

type AppealRepository interface {
	CreateAppeal(models.Appeal) (int, error)
	GetUserAppeals(string) ([]models.Appeal, int, error)
	UnrestrictUser(string) (int, error)
}

type AppealRepo struct {
	db *gorm.DB
}

func NewAppealRepository(db *gorm.DB) AppealRepository {
	return &AppealRepo{db: db}
}

func (ar *AppealRepo) CreateAppeal(appealReq models.Appeal) (int, error) {
	res := ar.db.Create(&appealReq)

	if res.Error != nil {
		return 500, fmt.Errorf("internal server error")
	}

	return 201, nil
}

func (ar *AppealRepo) GetUserAppeals(userId string) ([]models.Appeal, int, error) {
	var appeals []models.Appeal

	res := ar.db.Find(&appeals, "user_id = ?", userId)

	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return appeals, 200, nil
		}
		return appeals, 500, fmt.Errorf("internal server error")
	}

	return appeals, 200, nil
}

func (ar *AppealRepo) UnrestrictUser(userId string) (int, error) {
	res := ar.db.Model(&models.Appeal{}).
							Where("user_id = ?", userId).
							Update("is_banned", false)

	if res.Error != nil {
		return 500, fmt.Errorf("failed to accept appeal")
	}

	return 200, nil
}

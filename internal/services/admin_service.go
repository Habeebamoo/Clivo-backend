package services

import (
	"log"

	"github.com/Habeebamoo/Clivo/server/internal/models"
	"github.com/Habeebamoo/Clivo/server/internal/repositories"
	"github.com/Habeebamoo/Clivo/server/pkg/utils"
)

type AdminService interface {
	GetUsers() ([]models.UserProfileResponse, int, error)
	GetAppeals() ([]models.Appeal, int, error)
	GetUser(string) (models.UserResponse, int, error)
	VerifyUser(string) (int, error)
	UnVerifyUser(string) (int, error)
	BanUser(string) (int, error)
	UnBanUser(string) (int, error)
	GetArticlesByUsername(string) ([]models.Article, int, error)
	DeleteUserArticle(string) (int, error)
}

type AdminSvc struct {
	repo repositories.AdminRepository
}

func NewAdminService(repo repositories.AdminRepository) AdminService {
	return &AdminSvc{repo}
}

func (as *AdminSvc) GetUsers() ([]models.UserProfileResponse, int, error) {
	usersRaw, code, err := as.repo.GetUsers()
	if err != nil {
		return []models.UserProfileResponse{}, code, err
	}

	var users []models.UserProfileResponse

	for _, usr := range usersRaw {
		user := models.UserProfileResponse{
			UserId: usr.UserId,
			Name: usr.Name,
			Email: usr.Email,
			Role: usr.Role,
			Verified: usr.Verified,
			IsBanned: usr.IsBanned,
			Username: usr.Username,
			Bio: usr.Bio,
			Picture: usr.Picture,
			Interests: []string{},
			ProfileUrl: usr.ProfileUrl,
			Website: usr.Website,
			Following: usr.Following,
			Followers: usr.Followers,
			CreatedAt: utils.GetTimeAgo(usr.CreatedAt),
		}

		users = append(users, user)
	}

	return users, code, err
}

func (as *AdminSvc) GetAppeals() ([]models.Appeal, int, error) {
	return as.repo.GetAppeals()
}

func (as *AdminSvc) GetUser(userId string) (models.UserResponse, int, error) {
	return as.repo.GetUser(userId)
}

func (as *AdminSvc) VerifyUser(userId string) (int, error) {
	//get user
	user, code, err := as.repo.GetUser(userId)
	if err != nil {
		return code, err
	}

	//check is user is already verified
	if user.Verified {
		return 200, nil
	}

	//verify if not
	code, err = as.repo.UpdateUserVerification(userId, true)
	if err != nil {
		return code, err
	}

	//notify user
	go func() {
		//panic recovery
		defer func() {
			if r := recover(); r != nil {
				log.Printf("panic recovered from %v", r)
			}
		}()

		NewEmailService().SendVerifiedUserEmail(user.Name, user.Email)
	}()

	return code, err
}

func (as *AdminSvc) UnVerifyUser(userId string) (int, error) {
	//get user
	user, code, err := as.repo.GetUser(userId)
	if err != nil {
		return code, err
	}

	//check is user is not verified
	if !user.Verified {
		return 200, nil
	}

	//un-verify if not
	code, err = as.repo.UpdateUserVerification(userId, false)
	if err != nil {
		return code, err
	}

	//notify user
	go func() {
		//panic recovery
		defer func() {
			if r := recover(); r != nil {
				log.Printf("panic recovered from %v", r)
			}
		}()

		NewEmailService().SendUnverifiedUserEmail(user.Name, user.Email)
	}()

	return code, err
}

func (as *AdminSvc) BanUser(userId string) (int, error) {
	//get user
	user, code, err := as.repo.GetUser(userId)
	if err != nil {
		return code, err
	}

	//check if user is banned
	if user.IsBanned {
		return 200, nil
	}

	//ban if not
	return as.repo.UpdateUserRestriction(userId, true)
}

func (as *AdminSvc) UnBanUser(userId string) (int, error) {
	//get user
	user, code, err := as.repo.GetUser(userId)
	if err != nil {
		return code, err
	}

	//check is user is not banned
	if !user.IsBanned {
		return 200, nil
	}

	//un-ban if not
	return as.repo.UpdateUserRestriction(userId, false)
}

func (as *AdminSvc) GetArticlesByUsername(username string) ([]models.Article, int, error) {
	//get userId
	userId, code, err := as.repo.GetUserIdByUsername(username)
	if err != nil {
		return []models.Article{}, code, err
	}

	//get articles
	articles, code, err := as.repo.GetUserArticles(userId)
	if err != nil {
		return []models.Article{}, code, err
	}

	return articles, 200, nil
}

func (as *AdminSvc) DeleteUserArticle(articleId string) (int, error) {
	return as.repo.DeleteArticle(articleId)
}
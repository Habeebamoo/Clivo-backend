package services

import (
	"fmt"
	"strings"

	"github.com/Habeebamoo/Clivo/server/internal/models"
	"github.com/Habeebamoo/Clivo/server/internal/repositories"
	"github.com/Habeebamoo/Clivo/server/pkg/utils"
)

type UserService interface {
	GetUserProfile(string) (models.UserProfileResponse, int, error)
	UpdateUserProfile(string, models.ProfileUpdateRequest) (int, error)
	GetFollowStatus(string, string) (bool, error)
	FollowUser(string, string) (int, error)
	UnFollowUser(string, string) (int, error)
	GetUser(string) (models.SafeUserResponse, int, error)
	GetArticle(string, string) (models.SafeArticleResponse, int, error)
	GetArticles(string) ([]models.SafeArticleResponse, int, error)
	GetArticleComments(string, string) ([]models.CommentResponse, int, error)
	GetCommentReplys(string) ([]models.CommentResponse, int, error)
	GetFollowers(string) ([]models.SafeUserResponse, int, error)
	GetFollowing(string) ([]models.SafeUserResponse, int, error)
}

type UserSvc struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &UserSvc{repo}
}

func (us *UserSvc) GetUserProfile(userId string) (models.UserProfileResponse, int, error) {
	//get user
	user, code, err := us.repo.GetUserById(userId)
	if err != nil {
		return models.UserProfileResponse{}, code, err
	}

	//build response
	userInterests := strings.Split(user.Interests, ", ")

	//calculate time
	createdAt := utils.GetTimeAgo(user.CreatedAt)

	userProfile := models.UserProfileResponse{
		UserId: user.UserId,
		Name: user.Name,
		Email: user.Email,
		Role: user.Role,
		Verified: user.Verified,
		Username: user.Username,
		IsBanned: user.IsBanned,
		Bio: user.Bio,
		Picture: user.Picture,
		Interests: userInterests,
		ProfileUrl: user.ProfileUrl,
		Website: user.Website,
		Following: user.Following,
		Followers: user.Followers,
		CreatedAt: createdAt,
	}

	return userProfile, 200, nil
}

func (us *UserSvc) UpdateUserProfile(userId string, profileReq models.ProfileUpdateRequest) (int, error) {
	if !profileReq.FileAvailable {
		//update profile without picture
		return us.repo.UpdateUserProfile(userId, profileReq)
	}

	//upload picture and update profile
	imageUrl, err := utils.UploadImage(*profileReq.Picture)
	if err != nil {
		return 500, fmt.Errorf("failed to update profile")
	}

	return us.repo.UpdateUserProfileWithPicture(userId, profileReq, imageUrl)
}

func (us *UserSvc) GetUser(username string) (models.SafeUserResponse, int, error) {
	return us.repo.GetUserByUsername(username)
}

func (us *UserSvc) GetArticle(username string, articleTitleCode string) (models.SafeArticleResponse, int, error) {
	//get authorId
	authorId, code, err := us.repo.GetArticleAuthorIdByUsername(username)
	if err != nil {
		return models.SafeArticleResponse{}, code, err
	}

	//get article
	articleSlug := fmt.Sprintf("%s/%s", username, articleTitleCode)
	article, code, err := us.repo.GetArticleBySlug(authorId, articleSlug)
	if err != nil {
		return models.SafeArticleResponse{}, code, err
	}

	//get article likes
	articleLikes, err := us.repo.GetArticleLikes(article.ArticleId)
	if err != nil {
		return models.SafeArticleResponse{}, 500, err
	}

	//get article tags
	articeTags, err := us.repo.GetArticleTags(article.ArticleId)
	if err != nil {
		return models.SafeArticleResponse{}, 500, err
	}

	var tags []string
	for _, articleTag := range articeTags {
		tags = append(tags, articleTag.Tag)
	}

	//get user
	author, code, err := us.repo.GetArticleAuthorById(article.AuthorId)
	if err != nil {
		return models.SafeArticleResponse{}, code, err
	}

	//build response
	articleRespose := models.SafeArticleResponse{
		ArticleId: article.ArticleId,
		AuthorPicture: author.Picture,
		AuthorFullname: author.Name,
		AuthorProfileUrl: author.ProfileUrl,
		AuthorVerified: author.Verified,
		AuthorBio: author.Bio,
		Title: article.Title,
		Content: article.Content,
		Picture: article.Picture,
		Tags: tags,
		Likes: articleLikes,
		ReadTime: article.ReadTime,
		Slug: article.Slug,
		CreatedAt: utils.GetTimeAgo(article.CreatedAt),
	}

	return articleRespose, 200, nil
}

func (us *UserSvc) GetArticles(username string) ([]models.SafeArticleResponse, int, error) {
	//get articles
	articles, code, err := us.repo.GetArticlesByUsername(username)
	if err != nil {
		return []models.SafeArticleResponse{}, code, err
	}

	if len(articles) == 0 {
		return []models.SafeArticleResponse{}, 200, nil
	}

	//get user
	user, code, err := us.repo.GetUserByUsername(username)
	if err != nil {
		return []models.SafeArticleResponse{}, code, err
	}

	//build response
	var userArticles []models.SafeArticleResponse

	for _, article := range articles {
		//get likes
		likes, err := us.repo.GetArticleLikes(article.ArticleId)
		if err != nil {
			return []models.SafeArticleResponse{}, 500, err
		}

		createdAt := utils.GetTimeAgo(article.CreatedAt)

		//get tags
		articeTags, err := us.repo.GetArticleTags(article.ArticleId)
		if err != nil {
			return []models.SafeArticleResponse{}, 500, err
		}

		var tags []string
		for _, articleTag := range articeTags {
			tags = append(tags, articleTag.Tag)
		}

		safeArticle := models.SafeArticleResponse{
			ArticleId: article.ArticleId,
			AuthorPicture: user.Picture,
			AuthorFullname: user.Name,
			AuthorProfileUrl: user.ProfileUrl,
			AuthorVerified: user.Verified,
			Title: article.Title,
			Content: article.Content,
			Picture: article.Picture,
			Tags: tags,
			Likes: likes,
			ReadTime: article.ReadTime,
			Slug: article.Slug,
			CreatedAt: createdAt,
		}

		userArticles = append(userArticles, safeArticle)
	}

	return userArticles, 200, nil
}

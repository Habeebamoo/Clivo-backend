package handlers

import (
	"fmt"
	"net/http"

	"github.com/Habeebamoo/Clivo/server/internal/config"
	"github.com/Habeebamoo/Clivo/server/internal/models"
	"github.com/Habeebamoo/Clivo/server/internal/services"
	"github.com/Habeebamoo/Clivo/server/pkg/utils"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service services.UserService
}

func NewUserHandler(service services.UserService) UserHandler {
	return UserHandler{service}
}

func (uhdl *UserHandler) GetProfile(c *gin.Context) {
	userIdAny, exists := c.Get("userId")
	if !exists {
		utils.Error(c, 401, "UserId Missing", nil)
		return
	}

	userId := userIdAny.(string)

	//call service
	user, statusCode, err := uhdl.service.GetUserProfile(userId)
	if err != nil {
		utils.Error(c, statusCode, utils.FormatText(err.Error()), nil)
		return
	}

	utils.Success(c, statusCode, "", user)
}

func (uhdl *UserHandler) GetUser(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		utils.Error(c, 400, "User Not Found", nil)
		return
	}

	//call service
	user, statusCode, err := uhdl.service.GetUser(username)
	if err != nil {
		utils.Error(c, statusCode, "User Not Found", nil)
		return
	}

	utils.Success(c, statusCode, "", user)
}

func (uhdl *UserHandler) GetUserArticles(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		utils.Error(c, 400, "Username Missing", nil)
		return
	}

	//call service
	articles, statusCode, err := uhdl.service.GetArticles(username)
	if err != nil {
		utils.Error(c, statusCode, utils.FormatText(err.Error()), nil)
		return
	}

	utils.Success(c, statusCode, "", articles)
}

func (uhdl *UserHandler) GetUserArticle(c *gin.Context) {
	username := c.Param("username")
	articleTitle := c.Param("title")

	if username == "" || articleTitle == "" {
		utils.Error(c, 400, "Article Not Found", nil)
		return
	}

	//call service
	article, statusCode, err := uhdl.service.GetArticle(username, articleTitle)
	if err != nil {
		utils.Error(c, statusCode, "Article Not Found", nil)
		return
	}

	utils.Success(c, statusCode, "", article)
}

func (uhdl *UserHandler) GetArticleSEO(c *gin.Context) {
	userAgent := c.GetHeader("User-Agent")

	username := c.Param("username")
	slug := c.Param("slug")

	article, err := uhdl.service.GetArticleForSEO(username, slug)
	if err != nil {
		utils.Error(c, 404, utils.FormatText(err.Error()), nil)
		return
	}

	if utils.IsBot(userAgent) {
		html := uhdl.service.GenerateSEOHtml(article)

		c.Data(200, "text/html; charset=utf-8", []byte(html))
		return
	}

	clientOrigin, _ := config.Get("CLIENT_URL")
	url := fmt.Sprintf("%s/%s/%s", clientOrigin, username, slug)

	c.Redirect(http.StatusFound, url)
}

func (uhdl *UserHandler) UpdateProfile(c *gin.Context) {
	userIdAny, exists := c.Get("userId")
	if !exists {
		utils.Error(c, 401, "Unauthorized Access", nil)
		return
	}

	userId := userIdAny.(string)

	//receive form data
	name := c.PostForm("name")
	website := c.PostForm("website")
	bio := c.PostForm("bio")
	picture, _, err := c.Request.FormFile("picture")

	if name == "" || website == "" || bio == "" {
		utils.Error(c, 400, "No fields must be empty", nil)
		return
	}

	//build request
	profileUpdateReq := models.ProfileUpdateRequest{
		Name: name,
		Website: website,
		Bio: bio,
		Picture: &picture,
		FileAvailable: err == nil,
	}

	//call service
	statusCode, err := uhdl.service.UpdateUserProfile(userId, profileUpdateReq)
	if err != nil {
		utils.Error(c, statusCode, utils.FormatText(err.Error()), nil)
		return
	}

	utils.Success(c, 201, "Profile Update Successfully", nil)
}

func (uhdl *UserHandler) CreateSubscriber(c *gin.Context) {
	var subscriberReq models.SubscriberRequest
	if err := c.ShouldBindJSON(&subscriberReq); err != nil {
		utils.Error(c, 400, "Invalid JSON Format", nil)
		return
	}

	if subscriberReq.Email == "" {
		utils.Error(c, 400, "Email Address Missing", nil)
		return
	}

	//call service
	statusCode, err := uhdl.service.CreateSubscriber(subscriberReq.Email)
	if err != nil {
		utils.Error(c, statusCode, utils.FormatText(err.Error()), nil)
		return
	}

	utils.Success(c, statusCode, "You have successfully subscribed.", nil)
}

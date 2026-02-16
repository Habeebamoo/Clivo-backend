package handlers

import (
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

func (uhdl *UserHandler) UpdateProfile(c *gin.Context) {
	userIdAny, exists := c.Get("userId")
	if !exists {
		utils.Error(c, 401, "Unauthorized Access", nil)
		return
	}

	userId := userIdAny.(string)

	//receive form data
	name := c.PostForm("name")
	email := c.PostForm("email")
	website := c.PostForm("website")
	bio := c.PostForm("bio")
	picture, _, err := c.Request.FormFile("picture")

	if name == "" || email == "" || website == "" || bio == "" {
		utils.Error(c, 400, "No fields must be empty", nil)
		return
	}

	//build request
	profileUpdateReq := models.ProfileUpdateRequest{
		Name: name,
		Email: email,
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

func (uhdl *UserHandler) GetAppealStatus(c *gin.Context) {
	userId := c.Param("userId")
	if userId == "" {
		utils.Error(c, 400, "Missing UserID", nil)
		return
	}

	//call service
	status, code, err := uhdl.service.GetAppealStatus(userId)
	if err != nil {
		utils.Error(c, code, utils.FormatText(err.Error()), nil)
	}

	data := map[string]bool{"status": status}
	utils.Success(c, 200, "", data)
}

func (uhdl *UserHandler) SubmitAppeal(c *gin.Context) {
	var appealReq models.AppealRequest
	if err := c.ShouldBindJSON(&appealReq); err != nil {
		utils.Error(c, 400, "Invalid JSON Format", nil)
		return
	}

	//call service
	code, err := uhdl.service.CreateAppeal(appealReq)
	if err != nil {
		utils.Error(c, code, utils.FormatText(err.Error()), nil)
		return
	}

	utils.Success(c, code, "Appeal Submitted", nil)
}
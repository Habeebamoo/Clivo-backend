package handlers

import (
	"github.com/Habeebamoo/Clivo/server/internal/models"
	"github.com/Habeebamoo/Clivo/server/internal/services"
	"github.com/Habeebamoo/Clivo/server/pkg/utils"
	"github.com/gin-gonic/gin"
)

type AppealHandler struct {
	service services.AppealService
	userService services.UserService
}

func NewAppealHandler(service services.AppealService, userService services.UserService) AppealHandler {
	return AppealHandler{service, userService}
}

func (ahdl *AppealHandler) GetAppealStatus(c *gin.Context) {
	userId := c.Param("userId")
	if userId == "" {
		utils.Error(c, 400, "Missing UserID", nil)
		return
	}

	//call service
	status, code, err := ahdl.service.GetAppealStatus(userId)
	if err != nil {
		utils.Error(c, code, utils.FormatText(err.Error()), nil)
	}

	data := map[string]bool{"status": status}
	utils.Success(c, 200, "", data)
}

func (ahdl *AppealHandler) SubmitAppeal(c *gin.Context) {
	var appealReq models.AppealRequest
	if err := c.ShouldBindJSON(&appealReq); err != nil {
		utils.Error(c, 400, "Invalid JSON Format", nil)
		return
	}

	//check if user exists
  user, _, _ := ahdl.userService.GetUserProfile(appealReq.UserId)
	if user.Email == "" {
		utils.Error(c, 401, "Unauthorized Access", nil)
		return
	}

	//call service
	code, err := ahdl.service.CreateAppeal(appealReq)
	if err != nil {
		utils.Error(c, code, utils.FormatText(err.Error()), nil)
		return
	}

	utils.Success(c, code, "Appeal Submitted", nil)
}
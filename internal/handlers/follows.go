package handlers

import (
	"github.com/Habeebamoo/Clivo/server/pkg/utils"
	"github.com/gin-gonic/gin"
)

func (uhdl *UserHandler) GetFollowStatus(c *gin.Context) {
	userId := c.Param("userId")
	if userId == "" {
		utils.Error(c, 400, "UserID Missing", nil)
		return
	}

	userFollowing := c.Param("username")
	if userFollowing == "" {
		utils.Error(c, 400, "Username Missing", nil)
		return
	}

	//call service
	status, err := uhdl.service.GetFollowStatus(userId, userFollowing)
	if err != nil {
		utils.Error(c, 500, "", nil)
		return
	}

	data := map[string]bool{ "status": status }
	utils.Success(c, 200, "", data)
}

func (uhdl *UserHandler) FollowUser(c *gin.Context) {
	userIdAny, exists := c.Get("userId")
	if !exists {
		utils.Error(c, 401, "Unauthorized Access", nil)
		return
	}

	userId := userIdAny.(string)

	userFollowing := c.Param("username")
	if userFollowing == "" {
		utils.Error(c, 400, "UserId Missing", nil)
		return
	}

	//call service
	statusCode, err := uhdl.service.FollowUser(userId, userFollowing)
	if err != nil {
		utils.Error(c, statusCode, "", nil)
		return
	}

	utils.Success(c, statusCode, "", nil)
}

func (uhdl *UserHandler) UnFollowUser(c *gin.Context) {
	userIdAny, exists := c.Get("userId")
	if !exists {
		utils.Error(c, 401, "Unauthorized Access", nil)
		return
	}

	userId := userIdAny.(string)

	userFollowing := c.Param("username")
	if userFollowing == "" {
		utils.Error(c, 400, "UserId Missing", nil)
		return
	}

	//call service
	statusCode, err := uhdl.service.UnFollowUser(userId, userFollowing)
	if err != nil {
		utils.Error(c, statusCode, "", nil)
		return
	}

	utils.Success(c, statusCode, "", nil)
}

func (uhdl *UserHandler) GetUserFollowers(c *gin.Context) {
	userIdAny, exists := c.Get("userId")
	if !exists {
		utils.Error(c, 401, "Unauthorized Access", nil)
		return
	}

	userId := userIdAny.(string)

	//call service
	followers, statusCode, err := uhdl.service.GetFollowers(userId)
	if err != nil {
		utils.Error(c, statusCode, utils.FormatText(err.Error()), nil)
		return
	}
	
	utils.Success(c, statusCode, "", followers)
}

func (uhdl *UserHandler) GetUsersFollowing(c *gin.Context) {
	userIdAny, exists := c.Get("userId")
	if !exists {
		utils.Error(c, 401, "Unauthorized Access", nil)
		return
	}

	userId := userIdAny.(string)

	//call service
	followers, statusCode, err := uhdl.service.GetFollowing(userId)
	if err != nil {
		utils.Error(c, statusCode, utils.FormatText(err.Error()), nil)
		return
	}
	
	utils.Success(c, statusCode, "", followers)
}
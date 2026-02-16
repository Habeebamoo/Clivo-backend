package handlers

import (
	"github.com/Habeebamoo/Clivo/server/internal/models"
	"github.com/Habeebamoo/Clivo/server/pkg/utils"
	"github.com/gin-gonic/gin"
)


func (ah *ArticleHandler) LikeArticle(c *gin.Context) {
	var articleLikeRequest models.Like
	if err := c.ShouldBindJSON(&articleLikeRequest); err != nil {
		utils.Error(c, 400, "Invalid JSON Format", nil)
		return
	}

	//validate request
	if err := articleLikeRequest.Validate(); err != nil {
		utils.Error(c, 400, utils.FormatText(err.Error()), nil)
		return
	}

	//call service
	statusCode, err := ah.service.LikeArticle(articleLikeRequest)
	if err != nil {
		utils.Error(c, statusCode, utils.FormatText(err.Error()), nil)
		return
	}

	utils.Success(c, statusCode, "", nil)
}

func (ah *ArticleHandler) HasLikedArticle(c *gin.Context) {
	articleId := c.Param("articleId")
	userId := c.Param("userId")

	if articleId == "" || userId == "" {
		utils.Error(c, 400, "Params Missing", nil)
		return
	}

	articleLikeRequest := models.Like{
		ArticleId: articleId,
		LikerUserId: userId,
	}

	//call service
	liked := ah.service.HasUserLiked(articleLikeRequest)

	response := map[string]bool{ "liked": liked }
	utils.Success(c, 200, "", response)

}
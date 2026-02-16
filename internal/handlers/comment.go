package handlers

import (
	"github.com/Habeebamoo/Clivo/server/internal/models"
	"github.com/Habeebamoo/Clivo/server/pkg/utils"
	"github.com/gin-gonic/gin"
)

func (ah *ArticleHandler) CommentArticle(c *gin.Context) {
	var replyRequest models.CommentRequest
	if err := c.ShouldBindJSON(&replyRequest); err != nil {
		utils.Error(c, 400, "Invalid JSON Format", nil)
		return
	}

	//validate request
	if err := replyRequest.Validate(); err != nil {
		utils.Error(c, 400, utils.FormatText(err.Error()), nil)
		return
	}

	//call service
	statusCode, err := ah.service.CommentArticle(replyRequest)
		if err != nil {
		utils.Error(c, statusCode, utils.FormatText(err.Error()), nil)
		return
	}

	utils.Success(c, statusCode, "Comment Sent.", nil)
}

func (ah *ArticleHandler) ReplyComment(c *gin.Context) {
	var replyRequest models.ReplyRequest
	if err := c.ShouldBindJSON(&replyRequest); err != nil {
		utils.Error(c, 400, "Invalid JSON Format", nil)
		return
	}

	//validate request
	if err := replyRequest.Validate(); err != nil {
		utils.Error(c, 400, utils.FormatText(err.Error()), nil)
		return
	}

	//call service
	statusCode, err := ah.service.ReplyComment(replyRequest)
		if err != nil {
		utils.Error(c, statusCode, utils.FormatText(err.Error()), nil)
		return
	}

	utils.Success(c, statusCode, "Reply Sent.", nil)
}

func (uhdl *UserHandler) GetCommentReplys(c *gin.Context) {
	commentId := c.Param("id")

	if commentId == "" {
		utils.Error(c, 400, "Invalid Comment", nil)
		return
	}

	//call service
	comments, statusCode, err := uhdl.service.GetCommentReplys(commentId)
	if err != nil {
		utils.Error(c, statusCode, utils.FormatText(err.Error()), nil)
		return
	}

	utils.Success(c, statusCode, "", comments)
}

func (uhdl *UserHandler) GetArticleComments(c *gin.Context) {
	username := c.Param("username")
	articleTitle := c.Param("title")

	if username == "" || articleTitle == "" {
		utils.Error(c, 400, "Article Not Found", nil)
		return
	}

	//call service
	comments, statusCode, err := uhdl.service.GetArticleComments(username, articleTitle)
	if err != nil {
		utils.Error(c, statusCode, utils.FormatText(err.Error()), nil)
		return
	}

	utils.Success(c, statusCode, "", comments)
}
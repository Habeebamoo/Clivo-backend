package services

import (
	"fmt"
	"slices"

	"github.com/Habeebamoo/Clivo/server/internal/models"
	"github.com/Habeebamoo/Clivo/server/pkg/utils"
)

// Create comment
func (as *ArticleSvc) CommentArticle(commentReq models.CommentRequest) (int, error) {
	comment := models.Comment{
		CommentId: utils.GenerateRandomId(),
		ArticleId: commentReq.ArticleId,
		UserId: commentReq.UserId,
		ReplyId: "",
		Replys: 0,
		Content: commentReq.Content,
	}

	return as.articleRepo.CreateComment(comment)
}

// Create reply - comment
func (as *ArticleSvc) ReplyComment(replyReq models.ReplyRequest) (int, error) {
	comment := models.Comment{
		CommentId: utils.GenerateRandomId(),
		ArticleId: "",
		UserId: replyReq.UserId,
		ReplyId: replyReq.CommentId,
		Replys: 0,
		Content: replyReq.Content,
	}

	//create reply(comment)
	code, err := as.articleRepo.CreateComment(comment)
	if err != nil {
		return code, err
	}

	//update replys
	return as.articleRepo.UpdateReplys(comment.ReplyId)
}

// Get many comments -> 1 article
func (us *UserSvc) GetArticleComments(username, articleTitleCode string) ([]models.CommentResponse, int, error) {
	//get author
	authorId, code, err := us.repo.GetArticleAuthorIdByUsername(username)
	if err != nil {
		return []models.CommentResponse{}, code, err
	}

	//get article
	articleSlug := fmt.Sprintf("%s/%s", username, articleTitleCode)
	article, code, err := us.repo.GetArticleBySlug(authorId, articleSlug)
	if err != nil {
		return []models.CommentResponse{}, code, err
	}

	//get comments
	comments, code, err := us.repo.GetArticleComments(article.ArticleId)
	if err != nil {
		return []models.CommentResponse{}, code, err
	}

	//format comments
	commentsReponse := []models.CommentResponse{}

	for _, c := range comments {
		user, code, err := us.repo.GetUserById(c.UserId)
		if err != nil {
			return commentsReponse, code, err
		}

		comment := models.CommentResponse{
			CommentId: c.CommentId,
			Content: c.Content,
			ArticleId: article.ArticleId, 
			Replys: c.Replys,
			Name: user.Name,
			Username: user.Username,
			Verified: user.Verified,
			Picture: user.Picture,
		}

		commentsReponse = append(commentsReponse, comment)
	}

	//sort by latest
	slices.Reverse(commentsReponse)

	return commentsReponse, 200, nil
}

// Get many replys -> 1 comment
func (us *UserSvc) GetCommentReplys(commentId string) ([]models.CommentResponse, int, error) {
	//get comments
	comments, code, err := us.repo.GetCommentReplys(commentId)
	if err != nil {
		return []models.CommentResponse{}, code, err
	}

	//format comments
	commentsReponse := []models.CommentResponse{}

	for _, c := range comments {
		user, code, err := us.repo.GetUserById(c.UserId)
		if err != nil {
			return commentsReponse, code, err
		}

		comment := models.CommentResponse{
			CommentId: c.CommentId,
			Content: c.Content,
			ArticleId: c.ArticleId, 
			Replys: c.Replys,
			Name: user.Name,
			Username: user.Username,
			Verified: user.Verified,
			Picture: user.Picture,
		}

		commentsReponse = append(commentsReponse, comment)
	}

	//sort by latest
	slices.Reverse(commentsReponse)

	return commentsReponse, 200, nil
}
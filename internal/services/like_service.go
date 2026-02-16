package services

import "github.com/Habeebamoo/Clivo/server/internal/models"

func (as *ArticleSvc) LikeArticle(likeReq models.Like) (int, error) {
	//checkis user ha already liked
	alreadyLiked := as.articleRepo.IsLikedBy(likeReq)

	if alreadyLiked {
		return as.articleRepo.RemoveLike(likeReq)
	}

	return as.articleRepo.CreateLike(likeReq)
}

func (as *ArticleSvc) HasUserLiked(likeReq models.Like) bool {
	return as.articleRepo.IsLikedBy(likeReq)
}
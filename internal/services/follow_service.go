package services

import (
	"fmt"
	"slices"

	"github.com/Habeebamoo/Clivo/server/internal/models"
)

func (us *UserSvc) GetFollowStatus(followerId string, followingUsername string) (bool, error) {
	user, _, err := us.repo.GetUserIdByUsername(followingUsername)
	if err != nil {
		return false, fmt.Errorf("failed to get userId")
	}

	follow := models.Follow{
		FollowerId: followerId,
		FollowingId: user.UserId,
	}

	exists, err := us.repo.IsFollowing(follow)
	if err != nil {
		return exists, err
	}

	return exists, nil
}

func (us *UserSvc) FollowUser(followerId string, followingUsername string) (int, error) {
	user, _, err := us.repo.GetUserIdByUsername(followingUsername)
	if err != nil {
		return 500, fmt.Errorf("failed to get userId")
	}

	follow := models.Follow{
		FollowerId: followerId,
		FollowingId: user.UserId,
	}

	//create follow
	code, err := us.repo.CreateFollow(follow)
	if err != nil {
		return code, err
	}

	//use transactions
	//update follower (user that followed) profile
	err = us.repo.IncrementFollows(follow.FollowerId, "following")
	if err != nil {
		return code, err
	}

	//update following (user that is been followed)
	err = us.repo.IncrementFollows(follow.FollowingId, "followers")
	if err != nil {
		return code, err
	}

	return 201, nil
	//notify user following
}

func (us *UserSvc) UnFollowUser(followerId string, followingUsername string) (int, error) {
	user, _, err := us.repo.GetUserIdByUsername(followingUsername)
	if err != nil {
		return 500, fmt.Errorf("failed to get userId")
	}

	follow := models.Follow{
		FollowerId: followerId,
		FollowingId: user.UserId,
	}

	code, err := us.repo.RemoveFollow(follow)
	if err != nil {
		return code, err
	}

	//use transactions
	//update follower (user that followed) profile
	err = us.repo.DecrementFollows(follow.FollowerId, "following")
	if err != nil {
		return code, err
	}

	//update following (user that is been followed)
	err = us.repo.DecrementFollows(follow.FollowingId, "followers")
	if err != nil {
		return code, err
	}

	return 200, nil
	//notify user following
}

func (us *UserSvc) GetFollowers(userId string) ([]models.SafeUserResponse, int, error) {
	followersId, code, err := us.repo.GetUserFollowersId(userId)
	if err != nil {
		return []models.SafeUserResponse{}, code, err
	}

	var followers []models.SafeUserResponse

	for _, followerId := range followersId {
		follower, code, err := us.repo.GetUserById(followerId)
		//error check
		if err != nil {
			return followers, code, err
		}

		//build response
		user := models.SafeUserResponse{
			Name: follower.Name,
			Verified: follower.Verified,
			Username: follower.Username,
			Bio: follower.Bio,
			Picture: follower.Picture,
			ProfileUrl: follower.ProfileUrl,
			Website: follower.Website,
			Following: follower.Following,
			Followers: follower.Followers,
		}

		followers = append(followers, user)
	}

	//sort by recent
	slices.Reverse(followers)

	return followers, 200, nil
}

func (us *UserSvc) GetFollowing(userId string) ([]models.SafeUserResponse, int, error) {
	usersFollowingId, code, err := us.repo.GetUsersFollowingId(userId)
	if err != nil {
		return []models.SafeUserResponse{}, code, err
	}

	var usersFollowing []models.SafeUserResponse

	for _, userfollowingId := range usersFollowingId {
		userFollowing, code, err := us.repo.GetUserById(userfollowingId)
		//error check
		if err != nil {
			return usersFollowing, code, err
		}

		//build response
		user := models.SafeUserResponse{
			Name: userFollowing.Name,
			Verified: userFollowing.Verified,
			Username: userFollowing.Username,
			Bio: userFollowing.Bio,
			Picture: userFollowing.Picture,
			ProfileUrl: userFollowing.ProfileUrl,
			Website: userFollowing.Website,
			Following: userFollowing.Following,
			Followers: userFollowing.Followers,
		}

		usersFollowing = append(usersFollowing, user)
	}

	//sort by recent
	slices.Reverse(usersFollowing)

	return usersFollowing, 200, nil
}

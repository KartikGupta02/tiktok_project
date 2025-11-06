package model

type User struct {
	UserID            int    `json:"user_id"`
	Username          string `json:"username"`
	Email             string `json:"email"`
	Password          string `json:"password"`
	ProfilePictureURL string `json:"profile_picture_url"`
	Bio               string `json:"bio"`
	FollowersCount    int    `json:"followers_count"`
	FollowingsCount   int    `json:"followings_count"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

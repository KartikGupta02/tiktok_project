package model

type Video struct {
	VideoID       int    `json:"video_id"`
	UserID        int    `json:"user_id"`
	VideoCaption  string `json:"video_caption"`
	UploadDate    string `json:"upload_date"`
	Duration      int    `json:"duration"`
	VideoURL      string `json:"video_url"`
	ThumbnailURL  string `json:"thumbnail_url"`
	LikesCount    int    `json:"likes_count"`
	CommentsCount int    `json:"comments_count"`
	ViewsCount    int    `json:"views_count"`
	IsPublic      bool   `json:"is_public"`
}

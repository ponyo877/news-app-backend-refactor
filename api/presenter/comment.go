package presenter

import (
	"time"

	"github.com/ponyo877/news-app-backend-refactor/entity"
)

type CommentResponce struct {
	Data []*Comment `json:"data"`
}

type Comment struct {
	ID         string    `json:"id"`
	UserName   string    `json:"name"`
	AvatarURL  string    `json:"image_url"`
	DeviceHash string    `json:"device_hash"`
	Message    string    `json:"message"`
	UpdatedAt  time.Time `json:"updated_at"`
	CreatedAt  time.Time `json:"created_at"`
}

type CommentList []*Comment

func pickComment(comment entity.Comment) (Comment, error) {
	commentPresenter := Comment{
		ID:         comment.ID.String(),
		UserName:   comment.UserName,
		AvatarURL:  comment.AvatarURL,
		DeviceHash: comment.DeviceHash,
		Message:    comment.Message,
		UpdatedAt:  comment.UpdatedAt,
		CreatedAt:  comment.CreatedAt,
	}
	return commentPresenter, nil
}

func PickCommentList(commentList []entity.Comment) (CommentList, error) {
	var commentPresenterList CommentList
	for _, article := range commentList {
		var commentPresenter Comment
		commentPresenter, err := pickComment(article)
		if err != nil {
			return CommentList{}, err
		}
		commentPresenterList = append(commentPresenterList, &commentPresenter)
	}
	return commentPresenterList, nil
}

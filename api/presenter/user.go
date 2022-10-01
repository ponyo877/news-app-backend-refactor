package presenter

import (
	"time"

	"github.com/ponyo877/news-app-backend-refactor/entity"
)

type UserResponce struct {
	Data []*User `json:"data"`
}

type User struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	AvatarURL  string    `json:"image_url"`
	DeviceHash string    `json:"device_hash"`
	UpdatedAt  time.Time `json:"updated_at"`
	CreatedAt  time.Time `json:"created_at"`
}

type UserList []*User

func pickUser(user entity.User) (User, error) {
	userPresenter := User{
		ID:         user.ID.String(),
		Name:       user.Name,
		AvatarURL:  user.AvatarURL,
		DeviceHash: user.DeviceHash,
		UpdatedAt:  user.UpdatedAt,
		CreatedAt:  user.CreatedAt,
	}
	return userPresenter, nil
}

func PickUserList(userList []entity.User) (UserList, error) {
	var userPresenterList UserList
	for _, article := range userList {
		var userPresenter User
		userPresenter, err := pickUser(article)
		if err != nil {
			return UserList{}, err
		}
		userPresenterList = append(userPresenterList, &userPresenter)
	}
	return userPresenterList, nil
}

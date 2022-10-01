package repository

import (
	"time"

	"github.com/labstack/gommon/log"
	"github.com/ponyo877/news-app-backend-refactor/entity"
	"gorm.io/gorm"
)

// UserMySQL mysql repository
type UserMySQL struct {
	db *gorm.DB
}

type UserMySQLPresenter struct {
	ID         string    `gorm:"column:id;primary_key"`
	Name       string    `gorm:"column:name"`
	ImageURL   string    `gorm:"column:image_url"`
	DeviceHash string    `gorm:"column:device_hash"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
	CreatedAt  time.Time `gorm:"column:created_at"`
}

func (s UserMySQLPresenter) TableName() string {
	return "users"
}

// NewUserMySQL create new repository
func NewUserMySQL(db *gorm.DB) *UserMySQL {
	return &UserMySQL{
		db: db,
	}
}

func (r *UserMySQL) Get(id entity.ID) (*entity.User, error) {
	return nil, nil
}

func (r *UserMySQL) GetOption(deviceHash string) (entity.User, error) {
	var userMySQLPresenter UserMySQLPresenter
	if err := r.db.
		Where("device_hash = ?", deviceHash).
		Find(&userMySQLPresenter).
		Error; err != nil {
		return entity.User{}, err
	}
	ID, err := entity.StringToID(userMySQLPresenter.ID)
	if err != nil {
		return entity.User{}, err
	}
	return entity.User{
		ID:         ID,
		Name:       userMySQLPresenter.Name,
		AvatarURL:  userMySQLPresenter.ImageURL,
		DeviceHash: userMySQLPresenter.DeviceHash,
		UpdatedAt:  userMySQLPresenter.UpdatedAt,
		CreatedAt:  userMySQLPresenter.CreatedAt,
	}, nil
}

func (r *UserMySQL) Search(query string) ([]*entity.User, error) {
	return nil, nil
}

func (r *UserMySQL) List() ([]entity.User, error) {
	var userMySQLList []UserMySQLPresenter
	if err := r.db.Find(&userMySQLList).Error; err != nil {
		log.Infof("DBの接続に失敗しました: %v", err)
		return nil, err
	}
	var userList []entity.User
	for _, userMySQL := range userMySQLList {
		ID, err := entity.StringToID(userMySQL.ID)
		if err != nil {
			return nil, err
		}
		user := entity.User{
			ID:         ID,
			Name:       userMySQL.Name,
			AvatarURL:  userMySQL.ImageURL,
			DeviceHash: userMySQL.DeviceHash,
			UpdatedAt:  userMySQL.UpdatedAt,
			CreatedAt:  userMySQL.CreatedAt,
		}
		userList = append(userList, user)
	}
	return userList, nil
}

func (r *UserMySQL) Create(e entity.User) (entity.ID, error) {
	userMySQLPresenter := UserMySQLPresenter{
		ID:         e.ID.String(),
		Name:       e.Name,
		ImageURL:   e.AvatarURL,
		DeviceHash: e.DeviceHash,
		UpdatedAt:  e.UpdatedAt,
		CreatedAt:  e.CreatedAt,
	}
	if err := r.db.Create(userMySQLPresenter).Error; err != nil {
		return entity.NewID(), nil
	}
	return e.ID, nil
}

func (r *UserMySQL) Update(e entity.User) error {
	return nil
}

func (r *UserMySQL) Delete(id entity.ID) error {
	return nil
}

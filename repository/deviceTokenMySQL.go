package repository

import (
	"errors"
	"time"

	"github.com/ponyo877/news-app-backend-refactor/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// DeviceTokenMySQL mysql repository
type DeviceTokenMySQL struct {
	db *gorm.DB
}

type DeviceTokenMySQLPresenter struct {
	ExpoToken      string    `gorm:"column:expo_token;primary_key"`
	DeviceHash     string    `gorm:"column:device_hash"`
	Platform       string    `gorm:"column:platform"`
	DigestEnabled  bool      `gorm:"column:digest_enabled"`
	MatsuriEnabled bool      `gorm:"column:matsuri_enabled"`
	UpdatedAt      time.Time `gorm:"column:updated_at"`
	CreatedAt      time.Time `gorm:"column:created_at"`
}

// TableName
func (s DeviceTokenMySQLPresenter) TableName() string {
	return "device_tokens"
}

// NewDeviceTokenMySQL create new repository
func NewDeviceTokenMySQL(db *gorm.DB) *DeviceTokenMySQL {
	return &DeviceTokenMySQL{
		db: db,
	}
}

// Save トークンを登録・更新する(同一トークンの再登録は設定の上書き)
func (r *DeviceTokenMySQL) Save(e entity.DeviceToken) error {
	presenter := DeviceTokenMySQLPresenter{
		ExpoToken:      e.ExpoToken,
		DeviceHash:     e.DeviceHash,
		Platform:       e.Platform,
		DigestEnabled:  e.DigestEnabled,
		MatsuriEnabled: e.MatsuriEnabled,
		UpdatedAt:      e.UpdatedAt,
		CreatedAt:      e.CreatedAt,
	}
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "expo_token"}},
		DoUpdates: clause.AssignmentColumns([]string{"device_hash", "platform", "digest_enabled", "matsuri_enabled", "updated_at"}),
	}).Create(&presenter).Error
}

// ListDigestEnabled ダイジェスト通知ONのトークン一覧
func (r *DeviceTokenMySQL) ListDigestEnabled() ([]entity.DeviceToken, error) {
	return r.listEnabled("digest_enabled")
}

// ListMatsuriEnabled 祭り速報ONのトークン一覧
func (r *DeviceTokenMySQL) ListMatsuriEnabled() ([]entity.DeviceToken, error) {
	return r.listEnabled("matsuri_enabled")
}

func (r *DeviceTokenMySQL) listEnabled(column string) ([]entity.DeviceToken, error) {
	var presenterList []DeviceTokenMySQLPresenter
	if err := r.db.
		Where(column+" = ?", true).
		Find(&presenterList).
		Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	var deviceTokenList []entity.DeviceToken
	for _, presenter := range presenterList {
		deviceTokenList = append(deviceTokenList, entity.DeviceToken{
			ExpoToken:      presenter.ExpoToken,
			DeviceHash:     presenter.DeviceHash,
			Platform:       presenter.Platform,
			DigestEnabled:  presenter.DigestEnabled,
			MatsuriEnabled: presenter.MatsuriEnabled,
			UpdatedAt:      presenter.UpdatedAt,
			CreatedAt:      presenter.CreatedAt,
		})
	}
	return deviceTokenList, nil
}

// Delete 失効したトークン(DeviceNotRegistered)を削除する
func (r *DeviceTokenMySQL) Delete(expoToken string) error {
	return r.db.
		Where("expo_token = ?", expoToken).
		Delete(&DeviceTokenMySQLPresenter{}).Error
}

package entity

import (
	"strings"
	"time"
)

// DeviceToken プッシュ通知の送信先(Expo Push Token)
type DeviceToken struct {
	ExpoToken     string
	DeviceHash    string
	Platform      string
	DigestEnabled bool
	UpdatedAt     time.Time
	CreatedAt     time.Time
}

// NewDeviceToken create a new device token
func NewDeviceToken(expoToken, deviceHash, platform string, digestEnabled bool) (DeviceToken, error) {
	deviceToken := DeviceToken{
		ExpoToken:     expoToken,
		DeviceHash:    deviceHash,
		Platform:      platform,
		DigestEnabled: digestEnabled,
		UpdatedAt:     time.Now(),
		CreatedAt:     time.Now(),
	}
	if err := deviceToken.Validate(); err != nil {
		return DeviceToken{}, ErrInvalidEntity
	}
	return deviceToken, nil
}

// Validate validate data
func (d *DeviceToken) Validate() error {
	// Expoが発行するトークン以外(任意文字列の投げ込み)をテーブルに入れない
	if !strings.HasPrefix(d.ExpoToken, "ExponentPushToken[") || len(d.ExpoToken) > 128 {
		return ErrInvalidEntity
	}
	if d.DeviceHash == "" || len(d.DeviceHash) > 64 {
		return ErrInvalidEntity
	}
	if d.Platform != "ios" && d.Platform != "android" {
		return ErrInvalidEntity
	}
	return nil
}

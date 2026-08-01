package handler

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

// MakeHealthHandlers
func MakeHealthHandlers(e *echo.Echo, db *sql.DB) {
	e.GET("/health", Health(db)) // 死活監視(DB疎通込み)
}

// Health
func Health(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := db.Ping(); err != nil {
			log.Warnf("healthチェックでDB疎通に失敗しました: %v", err)
			return c.JSON(http.StatusServiceUnavailable, map[string]string{"status": "ng"})
		}
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	}
}

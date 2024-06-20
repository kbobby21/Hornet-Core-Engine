package factory

import (
	"time"
)

type NotificationsData struct {
	AssetId      int       `json:"asset_id,omitempty"`
	Email        string    `json:"email"`
	AlertType    string    `json:"alert_type"`
	AlertMessage string    `json:"alert_message"`
	Seen         bool      `json:"seen,omitempty"`
	Timestamp    time.Time `json:"nt_time"`
}

package notifications

import "bitbucket.org/hornetdefiant/core-engine/factory"

type Repository interface {
	AddNotification([]factory.NotificationsData) error
	GetNotifications(page_no int, email string) ([]factory.NotificationsData, error)
	DeleteNotifications(email string) error
	GetAssetActivities(pageNum int, email string) ([]factory.AssetActivity, error)
}

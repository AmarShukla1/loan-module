package notification

import "log"

type NotificationService struct{}

func NewNotificationService() *NotificationService {
	return &NotificationService{}
}

func (s *NotificationService) SendPushNotification(agentID int, message string) {
	log.Printf("[PUSH NOTIFICATION] Agent %d: %s", agentID, message)
}

func (s *NotificationService) SendSMS(phone, message string) {
	log.Printf("[SMS] %s: %s", phone, message)
}

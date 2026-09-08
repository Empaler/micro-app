package exercises

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
)

// Expected implementation contract (to be provided by you):
type Notification struct {
	ID        string
	UserID    string
	Message   string
	Read      bool
	CreatedAt time.Time
}

type NotificationManager struct {
	notificationsByUser map[string][]*Notification
}

func NewNotificationManager() *NotificationManager {
	return &NotificationManager{
		notificationsByUser: make(map[string][]*Notification),
	}
}

func (m *NotificationManager) Send(userID, message string) (Notification, error) {
	if userID == "" {
		return Notification{}, fmt.Errorf("must send userID")
	}
	newNotification := &Notification{
		ID:        uuid.NewString(),
		UserID:    userID,
		Message:   message,
		Read:      false,
		CreatedAt: time.Now(),
	}
	userNotifications := append(m.notificationsByUser[userID], newNotification)
	m.notificationsByUser[userID] = userNotifications
	return *newNotification, nil
}

func (m *NotificationManager) List(userID string) []Notification {
	notifications := make([]Notification, 0, len(m.notificationsByUser[userID]))
	for _, notification := range m.notificationsByUser[userID] {
		notifications = append(notifications, *notification)
	}
	return notifications
}

func (m *NotificationManager) MarkAsRead(userID, notificationID string) bool {
	if userID == "" || len(m.notificationsByUser[userID]) == 0 {
		return false
	}
	for _, notification := range m.notificationsByUser[userID] {
		if notification.ID == notificationID {
			notification.Read = true
			return true
		}
	}

	return false
}

func (m *NotificationManager) UnreadCount(userID string) int {
	if userID == "" {
		return 0
	}
	numberOfUnreadNotifications := 0
	for _, notification := range m.notificationsByUser[userID] {
		if notification.Read == false {
			numberOfUnreadNotifications++
		}
	}
	return numberOfUnreadNotifications
}

func TestNotificationManagerSendAndListOrder(t *testing.T) {
	manager := NewNotificationManager()

	first, err := manager.Send("u1", "Welcome")
	if err != nil {
		t.Fatalf("unexpected error on first send: %v", err)
	}
	second, err := manager.Send("u1", "Your order shipped")
	if err != nil {
		t.Fatalf("unexpected error on second send: %v", err)
	}

	notifications := manager.List("u1")
	if len(notifications) != 2 {
		t.Fatalf("expected 2 notifications, got %d", len(notifications))
	}

	if notifications[0].ID != first.ID || notifications[0].Message != "Welcome" {
		t.Fatalf("first notification mismatch: got %+v", notifications[0])
	}
	if notifications[1].ID != second.ID || notifications[1].Message != "Your order shipped" {
		t.Fatalf("second notification mismatch: got %+v", notifications[1])
	}

	if manager.UnreadCount("u1") != 2 {
		t.Fatalf("expected unread count 2, got %d", manager.UnreadCount("u1"))
	}
}

func TestNotificationManagerUserIsolation(t *testing.T) {
	manager := NewNotificationManager()

	_, _ = manager.Send("u1", "Hello u1")
	_, _ = manager.Send("u2", "Hello u2")
	_, _ = manager.Send("u1", "Another for u1")

	u1Notifications := manager.List("u1")
	u2Notifications := manager.List("u2")

	if len(u1Notifications) != 2 {
		t.Fatalf("expected 2 notifications for u1, got %d", len(u1Notifications))
	}
	if len(u2Notifications) != 1 {
		t.Fatalf("expected 1 notification for u2, got %d", len(u2Notifications))
	}

	if manager.UnreadCount("u1") != 2 {
		t.Fatalf("expected unread count 2 for u1, got %d", manager.UnreadCount("u1"))
	}
	if manager.UnreadCount("u2") != 1 {
		t.Fatalf("expected unread count 1 for u2, got %d", manager.UnreadCount("u2"))
	}
}

func TestNotificationManagerMarkAsRead(t *testing.T) {
	manager := NewNotificationManager()

	notification, err := manager.Send("u1", "Password changed")
	if err != nil {
		t.Fatalf("unexpected send error: %v", err)
	}

	if manager.UnreadCount("u1") != 1 {
		t.Fatalf("expected unread count 1 before read, got %d", manager.UnreadCount("u1"))
	}

	if !manager.MarkAsRead("u1", notification.ID) {
		t.Fatalf("expected mark as read to return true")
	}

	if manager.UnreadCount("u1") != 0 {
		t.Fatalf("expected unread count 0 after read, got %d", manager.UnreadCount("u1"))
	}

	notifications := manager.List("u1")
	if len(notifications) != 1 {
		t.Fatalf("expected exactly 1 notification, got %d", len(notifications))
	}
	if !notifications[0].Read {
		t.Fatalf("expected notification to be marked as read")
	}
}

func TestNotificationManagerMarkAsReadValidation(t *testing.T) {
	manager := NewNotificationManager()

	notification, err := manager.Send("u1", "Billing reminder")
	if err != nil {
		t.Fatalf("unexpected send error: %v", err)
	}

	if manager.MarkAsRead("u2", notification.ID) {
		t.Fatalf("expected mark as read with wrong user to return false")
	}
	if manager.MarkAsRead("u1", "missing-id") {
		t.Fatalf("expected mark as read for missing id to return false")
	}
	if manager.UnreadCount("u1") != 1 {
		t.Fatalf("expected unread count to remain 1, got %d", manager.UnreadCount("u1"))
	}
}

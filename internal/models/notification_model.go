package models

import "time"

type Notification struct {
	ID        string
	UserID    string
	Type      NotificationType
	Title     string
	Content   string
	IsRead    bool
	CreatedAt time.Time // Timestamp of when the notification was created
}

// NotificationType is an enumeration of notification types.
type NotificationType int

// Constants for NotificationType.
const (
	Like NotificationType = iota
	Comment
	Follow
	Mention
	Retweet
)

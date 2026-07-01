// Package cases is the cases bounded context: a foreclosure case and its CRUD.
// It exists to demonstrate the multi-domain layout and the cross-domain read
// convention (see ports.go) without importing another domain's model or repo.
package cases

import "time"

// Case represents a foreclosure case (demo domain).
type Case struct {
	ID         int       `gorm:"primaryKey" json:"id"`
	FileNumber string    `gorm:"uniqueIndex;not null" json:"file_number"`
	Status     string    `gorm:"not null" json:"status"`
	ServicerID int       `gorm:"not null" json:"servicer_id"`
	AssigneeID int       `gorm:"not null" json:"assignee_id"` // a user id (see UserReader)
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
}

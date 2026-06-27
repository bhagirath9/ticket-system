package models

import "time"

// TicketStatus represents the custom type for ticket status strings.
type TicketStatus string

const (
	StatusOpen       TicketStatus = "open"
	StatusInProgress TicketStatus = "in_progress"
	StatusClosed     TicketStatus = "closed"
)

// Ticket represents the tickets table schema in the database.
type Ticket struct {
	ID          uint         `gorm:"primaryKey;autoIncrement" json:"id"`
	Title       string       `gorm:"type:varchar(255);not null" json:"title"`
	Description string       `gorm:"type:text;not null" json:"description"`
	Status      TicketStatus `gorm:"type:varchar(50);not null;default:'open'" json:"status"`
	UserID      uint         `gorm:"not null;index" json:"user_id"`
	CreatedAt   time.Time    `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time    `gorm:"autoUpdateTime" json:"updated_at"`
}

// CreateTicketRequest represents the JSON body for creating a ticket.
type CreateTicketRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
}

// UpdateTicketStatusRequest represents the JSON body for patching a ticket's status.
type UpdateTicketStatusRequest struct {
	Status TicketStatus `json:"status" binding:"required"`
}

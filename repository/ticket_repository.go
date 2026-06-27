package repository

import (
	"ticket-system/models"

	"gorm.io/gorm"
)

// TicketRepository defines the interface for database operations related to Tickets.
type TicketRepository interface {
	Create(ticket *models.Ticket) error
	FindByOwnerID(ownerID uint) ([]models.Ticket, error)
	FindByID(id uint) (*models.Ticket, error)
	Save(ticket *models.Ticket) error
}

type ticketRepository struct {
	db *gorm.DB
}

// NewTicketRepository returns a new instance of TicketRepository.
func NewTicketRepository(db *gorm.DB) TicketRepository {
	return &ticketRepository{db: db}
}

// Create inserts a new ticket record into the database.
func (r *ticketRepository) Create(ticket *models.Ticket) error {
	return r.db.Create(ticket).Error
}

// FindByOwnerID retrieves all tickets belonging to a specific user.
func (r *ticketRepository) FindByOwnerID(ownerID uint) ([]models.Ticket, error) {
	var tickets []models.Ticket
	err := r.db.Where("user_id = ?", ownerID).Find(&tickets).Error
	if err != nil {
		return nil, err
	}
	return tickets, nil
}

// FindByID retrieves a ticket record by its primary key ID.
func (r *ticketRepository) FindByID(id uint) (*models.Ticket, error) {
	var ticket models.Ticket
	err := r.db.First(&ticket, id).Error
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

// Save updates an existing ticket's attributes in the database.
func (r *ticketRepository) Save(ticket *models.Ticket) error {
	return r.db.Save(ticket).Error
}

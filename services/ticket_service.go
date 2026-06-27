package services

import (
	"errors"
	"ticket-system/models"
	"ticket-system/repository"

	"gorm.io/gorm"
)

// Sentinel errors for ticket business operations.
var (
	ErrTicketNotFound    = errors.New("ticket not found")
	ErrForbiddenAccess   = errors.New("Forbidden: You do not own this ticket")
	ErrInvalidStatus     = errors.New("invalid status: must be open, in_progress, or closed")
	ErrClosedTicketReopen = errors.New("Closed ticket cannot be reopened")
	ErrInvalidTransition = errors.New("invalid status transition: cannot move status backward")
)

// TicketService defines the interface for ticket-related business logic.
type TicketService interface {
	Create(req *models.CreateTicketRequest, userID uint) (*models.Ticket, error)
	List(userID uint) ([]models.Ticket, error)
	GetByID(id uint, userID uint) (*models.Ticket, error)
	UpdateStatus(id uint, newStatus models.TicketStatus, userID uint) (*models.Ticket, error)
}

type ticketService struct {
	ticketRepo repository.TicketRepository
}

// NewTicketService returns a new instance of TicketService.
func NewTicketService(ticketRepo repository.TicketRepository) TicketService {
	return &ticketService{ticketRepo: ticketRepo}
}

// Create handles ticket creation, defaulting the status to 'open'.
func (s *ticketService) Create(req *models.CreateTicketRequest, userID uint) (*models.Ticket, error) {
	ticket := &models.Ticket{
		Title:       req.Title,
		Description: req.Description,
		Status:      models.StatusOpen,
		UserID:      userID,
	}

	if err := s.ticketRepo.Create(ticket); err != nil {
		return nil, err
	}

	return ticket, nil
}

// List retrieves all tickets belonging only to the specified user.
func (s *ticketService) List(userID uint) ([]models.Ticket, error) {
	return s.ticketRepo.FindByOwnerID(userID)
}

// GetByID retrieves a single ticket after verifying ownership.
func (s *ticketService) GetByID(id uint, userID uint) (*models.Ticket, error) {
	ticket, err := s.ticketRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTicketNotFound
		}
		return nil, err
	}

	// Verify that the logged-in user is the owner of the ticket
	if ticket.UserID != userID {
		return nil, ErrForbiddenAccess
	}

	return ticket, nil
}

// UpdateStatus changes the ticket status if valid transition rules are met.
func (s *ticketService) UpdateStatus(id uint, newStatus models.TicketStatus, userID uint) (*models.Ticket, error) {
	ticket, err := s.ticketRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTicketNotFound
		}
		return nil, err
	}

	// Verify that the logged-in user is the owner of the ticket
	if ticket.UserID != userID {
		return nil, ErrForbiddenAccess
	}

	// Validate status transitions
	if err := s.validateTransition(ticket.Status, newStatus); err != nil {
		return nil, err
	}

	ticket.Status = newStatus
	if err := s.ticketRepo.Save(ticket); err != nil {
		return nil, err
	}

	return ticket, nil
}

// validateTransition enforces:
// 1. Valid status values: open, in_progress, closed.
// 2. Linear transitions (open -> in_progress -> closed).
// 3. Closed tickets cannot be reopened (closed -> open or closed -> in_progress).
func (s *ticketService) validateTransition(current models.TicketStatus, newStatus models.TicketStatus) error {
	// Validate against allowed statuses
	if newStatus != models.StatusOpen && newStatus != models.StatusInProgress && newStatus != models.StatusClosed {
		return ErrInvalidStatus
	}

	// Closed ticket cannot be moved back to open or in_progress
	if current == models.StatusClosed && (newStatus == models.StatusOpen || newStatus == models.StatusInProgress) {
		return ErrClosedTicketReopen
	}

	// In progress ticket cannot go back to open
	if current == models.StatusInProgress && newStatus == models.StatusOpen {
		return ErrInvalidTransition
	}

	return nil
}

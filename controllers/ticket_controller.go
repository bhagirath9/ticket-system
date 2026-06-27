package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"ticket-system/models"
	"ticket-system/services"

	"github.com/gin-gonic/gin"
)

// TicketController handles ticket-related HTTP requests.
type TicketController struct {
	ticketService services.TicketService
}

// NewTicketController returns a new instance of TicketController.
func NewTicketController(ticketService services.TicketService) *TicketController {
	return &TicketController{ticketService: ticketService}
}

// Create processes requests for creating a new ticket.
// POST /tickets
func (ctrl *TicketController) Create(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var req models.CreateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	ticket, err := ctrl.ticketService.Create(&req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create ticket"})
		return
	}

	c.JSON(http.StatusCreated, ticket)
}

// List returns all tickets owned by the current logged-in user.
// GET /tickets
func (ctrl *TicketController) List(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	tickets, err := ctrl.ticketService.List(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to list tickets"})
		return
	}

	// Always return an empty array instead of null if user has no tickets
	if tickets == nil {
		tickets = []models.Ticket{}
	}

	c.JSON(http.StatusOK, tickets)
}

// GetByID returns the details of a single ticket belonging to the logged-in user.
// GET /tickets/:id
func (ctrl *TicketController) GetByID(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ticket ID format"})
		return
	}

	ticket, err := ctrl.ticketService.GetByID(uint(id), userID)
	if err != nil {
		if errors.Is(err, services.ErrTicketNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}
		if errors.Is(err, services.ErrForbiddenAccess) {
			c.JSON(http.StatusForbidden, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve ticket"})
		return
	}

	c.JSON(http.StatusOK, ticket)
}

// UpdateStatus patches the status of a user's own ticket.
// PATCH /tickets/:id/status
func (ctrl *TicketController) UpdateStatus(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ticket ID format"})
		return
	}

	var req models.UpdateTicketStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	ticket, err := ctrl.ticketService.UpdateStatus(uint(id), req.Status, userID)
	if err != nil {
		if errors.Is(err, services.ErrTicketNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}
		if errors.Is(err, services.ErrForbiddenAccess) {
			c.JSON(http.StatusForbidden, gin.H{"message": err.Error()})
			return
		}
		if errors.Is(err, services.ErrClosedTicketReopen) ||
			errors.Is(err, services.ErrInvalidStatus) ||
			errors.Is(err, services.ErrInvalidTransition) {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update ticket status"})
		return
	}

	c.JSON(http.StatusOK, ticket)
}

// getUserID extracts the verified userID from the context.
func getUserID(c *gin.Context) (uint, error) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		return 0, errors.New("userID missing from context")
	}
	userID, ok := userIDVal.(uint)
	if !ok {
		return 0, errors.New("userID in context is not of type uint")
	}
	return userID, nil
}

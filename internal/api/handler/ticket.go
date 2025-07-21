package handler

import (
	"github.com/ahMADASSadi/goTravel/internal/db"
	"github.com/ahMADASSadi/goTravel/internal/errors"
	"github.com/ahMADASSadi/goTravel/internal/models"
	"github.com/ahMADASSadi/goTravel/internal/repository"
	response "github.com/ahMADASSadi/goTravel/internal/responses"
	"github.com/ahMADASSadi/goTravel/internal/services"
	"github.com/gin-gonic/gin"
)

// CreateTicketHandler godoc
// @Summary      Create bus tickets for passengers
// @Description  Issues tickets for passengers on a reservation. Each passenger must have a unique social code.
// @Tags         Tickets
// @Accept       json
// @Produce      json
// @Param        ticket_input  body      models.TicketCreateInput  true  "Ticket creation data"
// @Success      201            {object}  map[string]interface{}    "Tickets created successfully"
// @Failure      400            {object}  map[string]string         "Invalid request or duplicate social codes"
// @Failure      404            {object}  map[string]string         "Reservation not found"
// @Failure      409            {object}  map[string]string         "Seats unavailable"
// @Failure      500            {object}  map[string]string         "Internal server error"
// @Router       /api/v1/tickets/buy [post]
func CreateTicketHandler(c *gin.Context) {
	var input models.TicketCreateInput
	if err := c.ShouldBind(&input); err != nil {
		response.Error(c, errors.ErrBadRequest)
		return
	}
	socialCodeMap := make(map[string]bool)
	for _, p := range input.Passengers {
		if socialCodeMap[p.SocialCode] {
			response.Error(c, errors.ErrDuplicateSocialCode)
			return
		}
		socialCodeMap[p.SocialCode] = true
	}

	tx := db.DB.Begin()
	if tx.Error != nil {
		response.Error(c, errors.ErrServerError)
		return
	}

	tickets, err := services.CreateTicket(tx, &input)
	if err != nil {
		tx.Rollback()
		switch err {
		case errors.ErrSeatsUnavailable:
			response.Error(c, errors.ErrSeatsUnavailable) // 409: seats unavailable
		case errors.ErrNotFound:
			response.Error(c, errors.ErrNotFound) // 404: reservation not found
		default:
			response.Error(c, errors.ErrServerError) // 500: unknown DB error
		}
		return
	}

	if err := tx.Commit().Error; err != nil {
		response.Error(c, errors.ErrServerError)
		return
	}
	response.Success(c, gin.H{"tickets": tickets})
}

// RefundTicketHandler godoc
// @Summary      Refund a ticket
// @Description  Refunds a ticket and marks it as deleted by setting DeletedAt.
// @Tags         Tickets
// @Accept       json
// @Produce      json
// @Param        ticket_refund  body      models.TicketRefundInput  true  "Ticket refund data"
// @Success      200             {object}  map[string]string         "Ticket refunded successfully"
// @Failure      400             {object}  map[string]string         "Invalid request or ticket already refunded"
// @Failure      404             {object}  map[string]string         "Ticket not found"
// @Failure      500             {object}  map[string]string         "Internal server error"
// @Router       /api/v1/tickets/refund/ [post]
func RefundTicketHandler(c *gin.Context) {
	var input models.TicketRefundInput
	if err := c.ShouldBind(&input); err != nil {
		response.Error(c, errors.ErrBadRequest)
		return
	}
	if err := services.RefundTicket(db.DB, input.TicketNo); err != nil {
		switch err {
		case errors.ErrNotFound:
			response.Error(c, errors.ErrNotFound) // 404: ticket not found
		case errors.ErrAlreadyRefunded:
			response.Error(c, errors.ErrAlreadyRefunded) // 400: already refunded
		default:
			response.Error(c, errors.ErrServerError) // 500: internal server error
		}
		return
	}

	response.Success(c, gin.H{"message": "Ticket refund successfully"})
}
// InquiryTicketHandler godoc
// @Summary      Inquiry about a ticket or reservation
// @Description  Retrieves details of a ticket, reservation, or travel based on search_id, ticket_number, or reservation_id.
// @Tags         Tickets
// @Accept       json
// @Produce      json
// @Param        ticket_inquiry  body      models.TicketInquiry  true  "Ticket inquiry data"
// @Success      200             {object}  map[string]interface{}  "Ticket or reservation information"
// @Failure      400             {object}  map[string]string      "Invalid request or missing search criteria"
// @Failure      404             {object}  map[string]string      "Ticket, reservation, or travel not found"
// @Router       /api/v1/tickets/inquiry/ [post]
func InquiryTicketHandler(c *gin.Context) {
	var input models.TicketInquiry
	if err := c.ShouldBind(&input); err != nil {
		response.Error(c, errors.ErrBadRequest)
		return
	}
	switch {
	case input.SearchID != "":
		travelInfo, err := repository.GetTravelInfoWithSearchID(input.SearchID)
		if err != nil {
			c.JSON(404, gin.H{"error": "travel not found"})
			return
		}
		c.JSON(200, gin.H{
			"travel":    travelInfo,
		})
	case input.ReservationID != 0:
		reservation, err := repository.GetReservationInfo(input.ReservationID) //getReservationByID(input.ReservationID)
		if err != nil {
			c.JSON(404, gin.H{"error": "reservation not found"})
			return
		}
		travelInfo, err := repository.GetTravelInfoWithSearchID(reservation.SearchID)
		if err != nil {
			c.JSON(404, gin.H{"error": "travel for reservation not found"})
			return
		}
		c.JSON(200, gin.H{
			"reservation":    reservation,
			"travel":         travelInfo,
		})
	case input.TicketNo != "":
		ticket, err := repository.GetTicketByNumber(input.TicketNo)
		if err != nil {
			c.JSON(404, gin.H{"error": "ticket not found"})
			return
		}
		travelInfo, err := repository.GetTravelInfoWithSearchID(ticket.SearchID)
		if err != nil {
			c.JSON(404, gin.H{"error": "travel for ticket not found"})
			return
		}
		c.JSON(200, gin.H{
			"ticket":        ticket,
			"travel":        travelInfo,
		})
	default:
		c.JSON(400, gin.H{"error": "please provide at least one of: search_id, ticket_number, or reservation_id"})
	}
}

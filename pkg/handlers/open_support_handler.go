package handlers

import (
	"suger-clickup/pkg/models"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) CreateTicket(c *fiber.Ctx) error {
	req := new(models.CreateTaskRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	userID := h.GetUserID(c)
	ticketID, err := h.s.CreateTask(userID, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"ticket_id": ticketID,
	})
}

func (h *Handler) GetTickets(c *fiber.Ctx) error {
	user := h.GetUser(c)
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("page_size", 10)
	if page < 1 {
		page = 1
	}
	if pageSize < 10 {
		pageSize = 10
	}
	tickets, err := h.s.GetTickets(user.ID, int32(page), int32(pageSize))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	total, err := h.s.GetTicketsCount(user.OrgId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"tickets": tickets,
		"total":   total,
	})
}

func (h *Handler) EditTicket(c *fiber.Ctx) error {

	return nil
}

func (h *Handler) ReopenTicket(c *fiber.Ctx) error {
	return nil
}

func (h *Handler) GetComments(c *fiber.Ctx) error {
	return nil
}
func (h *Handler) AddComment(c *fiber.Ctx) error {
	return nil
}

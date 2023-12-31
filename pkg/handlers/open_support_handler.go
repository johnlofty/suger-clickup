package handlers

import (
	"suger-clickup/pkg/models"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
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

func (h *Handler) ListTickets(c *fiber.Ctx) error {
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

func (h *Handler) GetTicket(c *fiber.Ctx) error {
	ticketID := c.Params("ticket_id")
	ticket, err := h.s.GetTicket(ticketID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})

	}
	return c.JSON(fiber.Map{
		"ticket": ticket,
	})
}

func (h *Handler) EditTicket(c *fiber.Ctx) error {
	req := new(models.TicketUpdateRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	user := h.GetUser(c)
	ticketID := c.Params("ticket_id")
	log.Debugf("getting ticketId:%s", ticketID)
	err := h.s.EditTicket(user, ticketID, *req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *Handler) ReopenTicket(c *fiber.Ctx) error {
	user := h.GetUser(c)
	ticketID := c.Params("ticket_id")
	log.Debugf("getting ticketId:%s", ticketID)
	err := h.s.ReopenTask(user, ticketID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *Handler) GetComments(c *fiber.Ctx) error {
	ticketID := c.Params("ticket_id")
	startID := c.Query("start_id", "")
	user := h.GetUser(c)
	comments, err := h.s.GetTicketComments(user, ticketID, startID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"comments": comments,
	})
}

func (h *Handler) AddComment(c *fiber.Ctx) error {
	ticketID := c.Params("ticket_id")
	req := new(models.ClickupCreateTaskCommentRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	user := h.GetUser(c)
	commentID, err := h.s.CreateTicketComments(user, ticketID, req.CommentText)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"id": commentID,
	})
}

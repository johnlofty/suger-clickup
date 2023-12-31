package handlers

import (
	"database/sql"
	"strconv"
	"suger-clickup/pkg/models"
	"suger-clickup/pkg/services"
	"suger-clickup/pkg/settings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	jtoken "github.com/golang-jwt/jwt/v4"
)

const DefaultOrg = 1

type Handler struct {
	s *services.Service
}

func NewHandler(s *services.Service) *Handler {
	return &Handler{
		s: s,
	}
}

func (h *Handler) GetUserID(c *fiber.Ctx) int32 {
	user := c.Locals("user").(*jtoken.Token)
	claims := user.Claims.(jtoken.MapClaims)
	userID := int32(claims["ID"].(float64))
	return userID
}

func (h *Handler) GetUser(c *fiber.Ctx) *models.User {
	user := c.Locals("user").(*jtoken.Token)
	claims := user.Claims.(jtoken.MapClaims)
	userInfo := &models.User{
		ID:    int32(claims["ID"].(float64)),
		OrgId: int32(claims["org_id"].(float64)),
		Email: claims["email"].(string),
	}
	return userInfo
}

func (h *Handler) Register(c *fiber.Ctx) error {
	registerRequset := new(models.RegisterRequest)
	if err := c.BodyParser(registerRequset); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	err := h.s.CreateUser(registerRequset.Email, registerRequset.Password, DefaultOrg)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.SendStatus(fiber.StatusOK)
}

func (h *Handler) Login(c *fiber.Ctx) error {
	loginRequest := new(models.LoginRequest)
	if err := c.BodyParser(loginRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	user, err := h.s.GetUser(loginRequest.Email, loginRequest.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	day := time.Hour * 24
	claims := jtoken.MapClaims{
		"ID":     user.ID,
		"email":  user.Email,
		"org_id": user.OrgId,
		"exp":    time.Now().Add(day * 1).Unix(),
	}
	log.Debugf("claims:%+v", claims)

	// create token
	token := jtoken.NewWithClaims(jtoken.SigningMethodHS256, claims)

	// Generate encoded token ans send it as response
	t, err := token.SignedString([]byte(settings.Get().JWT.Secret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(models.LoginResponse{
		Token: t,
	})

}

func (h *Handler) UpdateUser(c *fiber.Ctx) error {
	req := new(models.UpdateUserRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	userID := h.GetUserID(c)
	err := h.s.UpdateUserOrg(userID, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.SendStatus(fiber.StatusOK)
}

func (h *Handler) CreateOrg(c *fiber.Ctx) error {
	req := new(models.CreateOrgRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err := h.s.CreateOrg(req.OrgName)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *Handler) ModOrgNotification(c *fiber.Ctx) error {
	orgIDStr := c.Params("org_id")
	orgID, err := strconv.ParseInt(orgIDStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	req := new(models.OrgNotiRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = h.s.ModOrgNotification(int32(orgID), req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *Handler) GetOrgNotification(c *fiber.Ctx) error {
	orgIDStr := c.Params("org_id")
	orgID, err := strconv.ParseInt(orgIDStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	notification, err := h.s.GetOrgNotification(int32(orgID))
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(fiber.Map{
				"notification": models.OrgNotiRequest{
					OrgID:         int32(orgID),
					ContentChange: false,
					StatusChange:  false,
				},
			})
		} else {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}
	return c.JSON(fiber.Map{
		"notification": notification,
	})
}

package server

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func (s *FiberServer) CreateReminderHandler(c *fiber.Ctx) error {
	type ReminderCreate struct {
		UserID           int    `json:"user_id" xml:"user_id" form:"user_id"`
		Name             string `json:"name" xml:"name" form:"name"`
		Status           string `json:"status" xml:"status" form:"status"`
		Description      string `json:"description" xml:"description" form:"description"`
		Category         string `json:"category" xml:"category" form:"category"`
		ReminderInterval string `json:"reminder_interval" xml:"reminder_interval" form:"reminder_interval"`
		ReminderEnd      string `json:"reminder_end" xml:"reminder_end" form:"reminder_end"`
	}

	reminder := new(ReminderCreate)

	if err := c.BodyParser(reminder); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	err := s.db.SaveReminder(reminder.UserID, reminder.Name, reminder.Status, reminder.Description, reminder.Category, reminder.ReminderInterval, reminder.ReminderEnd)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot save reminder",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Reminder created successfully",
		"data":    fiber.Map{"reminder": reminder},
	})
}

func (s *FiberServer) GetRemindersHandler(c *fiber.Ctx) error {
	id := c.Params("id")

	id = strings.TrimSpace(id)

	reminderId, err := strconv.Atoi(id)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid reminder ID",
		})
	}

	reminder, err := s.db.GetReminderById(reminderId)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot get reminder",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Reminder retrieved successfully",
		"data":    fiber.Map{"reminder": reminder},
	})
}

func (s *FiberServer) GetRemindersForUserHandler(c *fiber.Ctx) error {
	userId := c.Params("user_id")

	userId = strings.TrimSpace(userId)

	userIdInt, err := strconv.Atoi(userId)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	reminders, err := s.db.GetAllRemindersForUser(userIdInt)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot get reminders",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Reminders retrieved successfully",
		"data":    fiber.Map{"reminders": reminders},
	})
}

func (s *FiberServer) GetAllRemindersHandler(c *fiber.Ctx) error {
	reminders, err := s.db.GetAllReminders()

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot get reminders",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Reminders retrieved successfully",
		"data":    fiber.Map{"reminders": reminders},
	})
}

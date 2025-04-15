package handler

import (
	"rentincome/database"
	"rentincome/model"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func GetTransaction(c *fiber.Ctx) error {
	t := new(model.Transaction)
	db := database.DB

	i := c.Params("id")
	id, err := strconv.Atoi(i)
	if err != nil {
		return ReturnError(c, err)
	}

	err = db.First(&t, id).Error
	if err != nil {
		return ReturnError(c, err)
	}

	return ReturnOK(c, t)
}

func CreateTransaction(c *fiber.Ctx) error {
	t := new(model.Transaction)
	db := database.DB

	err := c.BodyParser(t)
	if err != nil {
		return ReturnError(c, err)
	}

	validate := validator.New()
	err = validate.Struct(t)
	if err != nil {
		return ReturnError(c, err)
	}

	// TODO: Set user ID from auth context
	now := time.Now().UTC()
	t.DateActual = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	err = db.Create(&t).Error
	if err != nil {
		return ReturnError(c, err)
	}

	return ReturnOK(c, t)
}

func UpdateTransaction(c *fiber.Ctx) error {
	// TODO: Check ownership
	var err error
	i := c.Params("id")
	id, err := strconv.Atoi(i)
	if err != nil {
		return ReturnError(c, err)
	}

	t := new(model.Transaction)
	db := database.DB
	err = db.First(&t, id).Error
	if err != nil {
		ReturnError(c, 404, err)
	}

	type UpdateTransactionInput struct {
		Amount     float64   `json:"amount,omitempty"`
		Repetitive bool      `json:"repetitive,omitempty"`
		DateActual time.Time `json:"date_actual,omitempty"`
		DatePeriod time.Time `json:"date_period,omitempty"`
		Comment1   string    `json:"comment1,omitempty"`
		Comment2   string    `json:"comment2,omitempty"`
		Comment3   string    `json:"comment3,omitempty"`
	}

	var input UpdateTransactionInput

	err = c.BodyParser(&input)
	if err != nil {
		return ReturnError(c, err)
	}

	validate := validator.New()
	err = validate.Struct(input)
	if err != nil {
		return ReturnError(c, err)
	}

	t.Amount = input.Amount
	t.Comment1 = input.Comment1
	t.Comment2 = input.Comment2
	t.Comment3 = input.Comment3
	t.DatePeriod = input.DatePeriod
	t.DateActual = input.DateActual
	t.Repetitive = input.Repetitive

	err = db.Save(&t).Error
	if err != nil {
		return ReturnError(c, err)
	}

	return ReturnOK(c, t)
}

func DeleteTransaction(c *fiber.Ctx) error {
	// TODO : Check ownership
	var err error
	i := c.Params("id")
	id, err := strconv.Atoi(i)
	if err != nil {
		return ReturnError(c, err)
	}

	t := new(model.Transaction)
	db := database.DB
	err = db.First(&t, id).Error
	if err != nil {
		ReturnError(c, 404, err)
	}

	err = db.Delete(&t).Error
	if err != nil {
		return ReturnError(c, err)
	}

	return ReturnOK(c, fiber.StatusAccepted)
}

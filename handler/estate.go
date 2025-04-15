package handler

import (
	"fmt"
	"rentincome/database"
	"rentincome/model"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// TODO: check user_id allowed to see, maybe role based

// GET /estate/{id}
func GetEstate(c *fiber.Ctx) error {
	var estate model.Estate
	i := c.Params("id")
	id, err := strconv.Atoi(i)
	if err != nil {
		return ReturnError(c, err)
	}

	db := database.DB
    err = db.First(&estate, id).Error 
	if err!=nil {
		return ReturnError(c, 404, "Estate not found")
	}
	
    return ReturnOK(c, "Estate found", estate)
}

func CreateEstate(c *fiber.Ctx) error {
	// debug payload
	// fmt.Println("Received:", string(c.Body()))
	estate := new(model.Estate)

	err := c.BodyParser(estate)
	if err != nil {
		return ReturnError(c, err)
	}
	err = model.EstateValidate(estate)
    fmt.Print(err)
	if err != nil {
		errors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			errors[err.Field()] = err.Tag()
		}
		return ReturnError(c, errors, "validation failed")
	}

	db := database.DB
	err = db.Create(estate).Error
	if err != nil {
		return ReturnError(c, err)
	}

	// TODO : uncomment to bind userid to author
	//    r, err := GetFromJwt(c, "user_id")
	// if err != nil {
	// 	returnError(c, err)
	// }
	// userID := uint(r.(float64))
	// if err != nil {
	// 	returnError(c, err)
	// }
	// estate.UserID = userID
	return ReturnOK(c, "Estate created", estate)
}

func UpdateEstate(c *fiber.Ctx) error {
	i := c.Params("id")
	id, err := strconv.Atoi(i)
	if err != nil {
		return ReturnError(c, err)
	}
	db := database.DB
	estate := new(model.Estate)
	err = db.First(&estate, id).Error 
	if err != nil {
		ReturnError(c, 404, "Estate not found",err)
	}

	err = c.BodyParser(estate)
	if err != nil {
		return ReturnError(c, err)
	}
	err = model.EstateValidate(estate)
	if err != nil {
		return ReturnError(c, err)
	}

	err = db.Save(estate).Error
	if err != nil {
		return ReturnError(c, err)
	}

	return ReturnOK(c, "Estate updated", estate, fiber.StatusAccepted)
}

func DeleteEstate(c *fiber.Ctx) error {
	var estate model.Estate
	i := c.Params("id")
	id, err := strconv.Atoi(i)
	if err != nil {
		return ReturnError(c, err)
	}

	db := database.DB
	err = db.First(&estate, id).Error 
	if err!=nil {
		ReturnError(c, 404, "Estate not found",err)
	}
	err = db.Delete(&estate, id).Error 
    if err!=nil {
		ReturnError(c, err)
	}
	return ReturnOK(c, "Estate deleted", estate, fiber.StatusAccepted)
}

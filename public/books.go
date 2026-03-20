package public

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/Rx-11/EDIS-A1/ai"
	"github.com/Rx-11/EDIS-A1/common"
	"github.com/Rx-11/EDIS-A1/config"
	"github.com/Rx-11/EDIS-A1/db"
	"github.com/Rx-11/EDIS-A1/pkg"
	"github.com/Rx-11/EDIS-A1/pkg/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func fetchBookByISBN(c *fiber.Ctx) error {

	param := c.Locals("param").(fetchBookByISBNParam)

	book, err := pkg.BookRepo.FetchBookByISBN(db.GetDB(), param.ISBN)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(common.ErrNotFound.StatusCode).JSON(common.ErrNotFound)
		}
		return c.Status(common.ErrInternalServerError.StatusCode).JSON(common.ErrInternalServerError)
	}

	return c.JSON(book)
}

func createBook(c *fiber.Ctx) error {
	body := c.Locals("body").(createBookRequest)

	existingBook, _ := pkg.BookRepo.FetchBookByISBN(db.GetDB(), body.ISBN)
	if existingBook != nil {
		return c.Status(common.ErrUnprocessableEntity.StatusCode).JSON(common.NewError(
			common.ErrUnprocessableEntity.StatusCode,
			"This ISBN already exists in the system.",
		))
	}

	SummaryResp, err := config.Gemini.Chat(ai.ChatRequest{Messages: []ai.Message{{Role: "model", Content: "Give a 500 word summary of the following book"}, {Role: "user", Content: fmt.Sprintf("Book Title: %s\nBook Description: %s\nBook Author: %s\nBook ISBN: %s", body.Title, body.Description, body.Author, body.ISBN)}}})
	if err != nil {
		return c.Status(common.ErrInternalServerError.StatusCode).JSON(common.ErrInternalServerError)
	}

	book, err := pkg.BookRepo.CreateBook(db.GetDB(), models.Book{
		ISBN:        body.ISBN,
		Title:       body.Title,
		Author:      body.Author,
		Price:       body.Price,
		Description: body.Description,
		Genre:       body.Genre,
		Quantity:    body.Quantity,
		Summary:     SummaryResp.Response,
	})
	if err != nil {
		return c.Status(common.ErrInternalServerError.StatusCode).JSON(common.ErrInternalServerError)
	}

	return c.Status(fiber.StatusCreated).JSON(book)
}

func updateBook(c *fiber.Ctx) error {
	body := c.Locals("body").(updateBookRequest)

	existingBook, err := pkg.BookRepo.FetchBookByISBN(db.GetDB(), body.ISBN)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(common.ErrNotFound.StatusCode).JSON(common.ErrNotFound)
		}
	}

	existingBook.Title = body.Title
	existingBook.Author = body.Author
	existingBook.Price = strconv.FormatFloat(body.Price, 'f', -1, 64)
	existingBook.Description = body.Description
	existingBook.Genre = body.Genre
	existingBook.Quantity = body.Quantity

	book, err := pkg.BookRepo.UpdateBook(db.GetDB(), *existingBook)
	if err != nil {
		return c.Status(common.ErrInternalServerError.StatusCode).JSON(common.ErrInternalServerError)
	}

	return c.Status(fiber.StatusCreated).JSON(book)
}

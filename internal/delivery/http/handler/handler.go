package handler

import (
	"context"
	"errors"
	"go-books-api/internal/dto"
	"go-books-api/internal/models"
	"go-books-api/pkg/logger"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const requestTimeout = 1 * time.Second

type ErrorResponse struct {
	Error string `json:"error"`
}

type BookService interface {
	AddBook(ctx context.Context, book *dto.Book) (*dto.AddBookResponse, error)
	DeleteBook(ctx context.Context, bookID string) error
	GetBook(ctx context.Context, bookID string) (*dto.Book, error)
	ListBooks(ctx context.Context, page string, limit string) (*dto.ListBooksResponse, error)
	UpdateBook(ctx context.Context, bookID string, book *dto.Book) error
}

type Handler struct {
	service BookService
	logger  logger.Logger
}

func NewHandler(serv BookService, log logger.Logger) *Handler {
	return &Handler{
		service: serv,
		logger:  log,
	}
}

// AddBook
// @Summary      Adds a new book
// @Description  Adds a new book
// @Tags         books
// @Produce      json
// @Param book body dto.Book true "Book"
// @Success      201  {object}  dto.AddBookResponse
// @Failure      400  {object}  ErrorResponse  "Invalid request body"
// @Failure      500  {object}  ErrorResponse  "Uknown error occured while adding the book"
// @Router /books [post]
func (h *Handler) AddBook(ctx *gin.Context) {
	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	var book dto.Book

	if err := ctx.ShouldBindJSON(&book); err != nil {
		h.logger.Error(ctx, "Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request body"})
		return
	}

	res, err := h.service.AddBook(ctxWithTimeout, &book)
	if err != nil {
		h.logger.Error(ctx, "Unknown error occured while adding the book", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Uknown error occured while adding the book"})
		return
	}

	ctx.JSON(http.StatusCreated, res)
}

// GetBook
// @Summary      Returns the book
// @Description  Returns the book
// @Tags         books
// @Produce      json
// @Param id path string true "Book ID"
// @Success      200  {object}  dto.Book
// @Failure      400  {object}  ErrorResponse  "book id is invalid"
// @Failure      404  {object}  ErrorResponse  "noting was found"
// @Failure      500  {object}  ErrorResponse  "Uknown error occured while getting the book"
// @Router /books/{id} [get]
func (h *Handler) GetBook(ctx *gin.Context) {
	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	bookID := ctx.Param("id")

	res, err := h.service.GetBook(ctxWithTimeout, bookID)
	if err != nil {
		if errors.Is(err, models.ErrFailedToParseID) {
			ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
			return
		}

		if errors.Is(err, models.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
			return
		}

		h.logger.Error(ctx, "Unknown error occured while getting the book", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Uknown error occured while getting the book"})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

// UpdateBook
// @Summary      Updates the book
// @Description  Updates the book
// @Tags         books
// @Produce      json
// @Param id path string true "Book ID"
// @Param book body dto.Book true "Book"
// @Success      200  {object}  nil
// @Failure      400  {object}  ErrorResponse  "Invalid request body"
// @Failure      400  {object}  ErrorResponse  "book id is invalid"
// @Failure      404  {object}  ErrorResponse  "noting was found"
// @Failure      500  {object}  ErrorResponse  "Uknown error occured while deleting the book"
// @Router /books/{id} [put]
func (h *Handler) UpdateBook(ctx *gin.Context) {
	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	bookID := ctx.Param("id")

	var book dto.Book

	if err := ctx.ShouldBindJSON(&book); err != nil {
		h.logger.Error(ctx, "Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request body"})
		return
	}

	if err := h.service.UpdateBook(ctxWithTimeout, bookID, &book); err != nil {
		if errors.Is(err, models.ErrFailedToParseID) {
			ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
			return
		}

		if errors.Is(err, models.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
			return
		}

		h.logger.Error(ctx, "Unknown error occured while updating the book", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Uknown error occured while updating the book"})
		return
	}

	ctx.Status(http.StatusOK)
}

// DeleteBook
// @Summary      Deletes the book
// @Description  Deletes the book
// @Tags         books
// @Produce      json
// @Param id path string true "Book ID"
// @Success      200  {object}  nil
// @Failure      400  {object}  ErrorResponse  "book id is invalid"
// @Failure      404  {object}  ErrorResponse  "noting was found"
// @Failure      500  {object}  ErrorResponse  "Uknown error occured while deleting the book"
// @Router /books/{id} [delete]
func (h *Handler) DeleteBook(ctx *gin.Context) {
	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	bookID := ctx.Param("id")

	if err := h.service.DeleteBook(ctxWithTimeout, bookID); err != nil {
		if errors.Is(err, models.ErrFailedToParseID) {
			ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
			return
		}

		if errors.Is(err, models.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
			return
		}

		h.logger.Error(ctx, "Unknown error occured while deleting the book", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Uknown error occured while deleting the book"})
		return
	}

	ctx.Status(http.StatusOK)
}

// ListBooks
// @Summary      Lists the books
// @Description  Lists the books
// @Tags         books
// @Produce      json
// @Param page query string false "page number"
// @Param limit query string false "limit"
// @Success      200  {object}  dto.ListBooksResponse
// @Failure      400  {object}  ErrorResponse  "page number is invalid / limit number is invalid"
// @Failure      404  {object}  ErrorResponse  "noting was found"
// @Failure      500  {object}  ErrorResponse  "Uknown error occured while listing the book"
// @Router /books [get]
func (h *Handler) ListBooks(ctx *gin.Context) {
	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	page := ctx.Query("page")
	limit := ctx.Query("limit")

	res, err := h.service.ListBooks(ctxWithTimeout, page, limit)
	if err != nil {
		if errors.Is(err, models.ErrFailedToParsePage) || errors.Is(err, models.ErrFailedToParseLimit) {
			ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
			return
		}

		if errors.Is(err, models.ErrNotFound) {
			ctx.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
			return
		}

		h.logger.Error(ctx, "Unknown error occured while listing the book", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Uknown error occured while listing the book"})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

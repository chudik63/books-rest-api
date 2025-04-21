package service

import (
	"context"
	"go-books-api/internal/dto"
	"go-books-api/internal/models"
	"strconv"
	"time"
)

type Repository interface {
	AddBook(ctx context.Context, book *models.Book) (uint64, error)
	GetBookByID(ctx context.Context, id uint64) (*models.Book, error)
	DeleteBook(ctx context.Context, id uint64) error
	ListBooks(ctx context.Context, limit, offset uint64) ([]*models.Book, error)
}

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) AddBook(ctx context.Context, book *dto.Book) (*dto.AddBookResponse, error) {
	now := time.Now()

	id, err := s.repo.AddBook(ctx, &models.Book{
		Title:     book.Title,
		Author:    book.Author,
		Genre:     book.Genre,
		CreatedAt: now,
	})

	return &dto.AddBookResponse{
		ID:        id,
		CreatedAt: now,
	}, err
}

func (s *Service) DeleteBook(ctx context.Context, bookIDStr string) error {
	bookID, err := strconv.ParseUint(bookIDStr, 10, 64)
	if err != nil {
		return models.ErrFailedToParseID
	}

	err = s.repo.DeleteBook(ctx, bookID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetBook(ctx context.Context, bookIDStr string) (*dto.Book, error) {
	bookID, err := strconv.ParseUint(bookIDStr, 10, 64)
	if err != nil {
		return nil, models.ErrFailedToParseID
	}

	book, err := s.repo.GetBookByID(ctx, bookID)
	if err != nil {
		return nil, err
	}

	return &dto.Book{
		Title:  book.Title,
		Author: book.Author,
		Genre:  book.Genre,
	}, nil
}

func (s *Service) ListBooks(ctx context.Context, pageStr string, limitStr string) (*dto.ListBooksResponse, error) {
	page, err := strconv.ParseUint(pageStr, 10, 64)
	if err != nil && pageStr != "" {
		return nil, models.ErrFailedToParsePage
	}

	limit, err := strconv.ParseUint(limitStr, 10, 64)
	if err != nil && limitStr != "" {
		return nil, models.ErrFailedToParseLimit
	}

	if page < 1 {
		page = 1
	}

	if limit == 0 {
		limit = 10
	}

	offset := (page - 1) * limit

	books, err := s.repo.ListBooks(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	booksDTO := make([]dto.Book, len(books))

	for i, book := range books {
		booksDTO[i] = dto.Book{
			Title:  book.Title,
			Author: book.Author,
			Genre:  book.Genre,
		}
	}

	return &dto.ListBooksResponse{
		Books: booksDTO,
	}, nil
}

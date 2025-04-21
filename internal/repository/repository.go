package repository

import (
	"context"
	"database/sql"
	"errors"
	"go-books-api/internal/database/postgres"
	"go-books-api/internal/models"

	sq "github.com/Masterminds/squirrel"
)

type Repository struct {
	db postgres.DB
}

func New(db postgres.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) AddBook(ctx context.Context, book *models.Book) (uint64, error) {
	row := sq.Insert("books").
		Columns("title", "author", "genre", "created_at").
		Values(book.Title, book.Author, book.Genre, book.CreatedAt).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db).
		QueryRowContext(ctx)

	var id uint64

	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, err
}

func (r *Repository) GetBookByID(ctx context.Context, id uint64) (*models.Book, error) {
	row := sq.Select("id", "title", "author", "genre", "created_at").
		From("books").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db).
		QueryRowContext(ctx)

	var book models.Book

	err := row.Scan(
		&book.ID,
		&book.Title,
		&book.Author,
		&book.Genre,
		&book.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNotFound
		}

		return nil, err
	}

	return &book, nil
}

func (r *Repository) DeleteBook(ctx context.Context, id uint64) error {
	_, err := sq.Delete("books").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db).
		ExecContext(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.ErrNotFound
		}

		return err
	}

	return nil
}

func (r *Repository) ListBooks(ctx context.Context, limit, offset uint64) ([]*models.Book, error) {
	rows, err := sq.Select("id", "title", "author", "genre", "created_at").
		From("books").
		Limit(limit).
		Offset(offset).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db).
		QueryContext(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNotFound
		}

		return nil, err
	}
	defer rows.Close()

	var books []*models.Book

	for rows.Next() {
		var book models.Book
		if err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.Author,
			&book.Genre,
			&book.CreatedAt,
		); err != nil {
			return nil, err
		}
		books = append(books, &book)
	}

	if len(books) == 0 {
		return nil, models.ErrNotFound
	}

	return books, nil
}

func (r *Repository) UpdateBook(ctx context.Context, book *models.Book) error {
	_, err := sq.Update("books").
		Where(sq.Eq{"id": book.ID}).
		Set("title", book.Title).
		Set("author", book.Author).
		Set("genre", book.Genre).
		Set("updated_at", book.UpdatedAt).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db).
		Exec()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.ErrNotFound
		}

		return err
	}

	return nil
}

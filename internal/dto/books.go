package dto

import "time"

type Book struct {
	Title  string `json:"title" binding:"required,min=1,max=100"`
	Author string `json:"author" binding:"required,min=1,max=100"`
	Genre  string `json:"genre" binding:"required,min=1,max=100"`
}

type AddBookResponse struct {
	ID        uint64    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

type ListBooksResponse struct {
	Books []Book `json:"books"`
}

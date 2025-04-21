package models

import "time"

type Book struct {
	ID        uint64
	Title     string
	Author    string
	Genre     string
	CreatedAt time.Time
}

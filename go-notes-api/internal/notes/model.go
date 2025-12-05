package notes

import (
	"errors"
	"strings"
	"time"
)

type Note struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var (
	ErrTitleRequired = errors.New("title is required")
	ErrTitleTooLong = errors.New("title must be less than 100 characters")
	ErrContentTooLong = errors.New("content must be less than 1000 characters")
	ErrNoteNotFound = errors.New("note not found")
	maxTitleLength = 100
	maxContentLength = 1000
)

func ValidateNoteInput(title string, content string) error {
	title = strings.TrimSpace(title)
	content = strings.TrimSpace(content)
	
	if title == "" {
		return ErrTitleRequired
	}
	if len(title) > maxTitleLength {
		return ErrTitleTooLong
	}
	if len(content) > maxContentLength {
		return ErrContentTooLong
	}
	return nil
}
package notes

import (
	"errors"
	"strings"
)

type Note struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

var (
	ErrTitleRequired = errors.New("title is required")
)

func ValidateNoteInput(title string) error {
	if strings.TrimSpace(title) == "" {
		return ErrTitleRequired
	}
	return nil
}
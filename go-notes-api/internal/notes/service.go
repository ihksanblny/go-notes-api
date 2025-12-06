package notes

import (
	"context"
	"strings"
)

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

// List Notes mengembalikan notes dengan pagination + search
func (s *Service) ListNotes(ctx context.Context, q string, page, limit int, sort, order string) ([]Note, int, error) {
	// validasi udah dilakuin di handler
	sort = strings.ToLower(strings.TrimSpace(sort))
	order = strings.ToLower(strings.TrimSpace(order))

	switch sort {
	case "title", "updated_at", "created_at":
		//ok
	default:
		sort = "created_at" //default
	}

	switch order {
	case "asc", "desc":
		//ok
	default:
		order = "desc" //default
	}

	items, total := s.store.ListPage(page, limit, q, sort, order)
	return items, total, nil 
}

// CreateNote validasi input lalu simpan
func (s* Service) CreateNote(ctx context.Context, title, content string) (Note, error) {
	if err := ValidateNoteInput(title, content); err != nil {
		return Note{}, err
	}
	note := s.store.Create(title, content)
	return note, nil
}

// GetNote ambil 1 Note, error kalau tidak ada
func (s *Service) GetNote(ctx context.Context, id int) (Note, error) {
	note, ok := s.store.Get(id)
	if !ok {
		return Note{}, ErrNoteNotFound
	}
	return note, nil
}

// UpdateNote validasi input lalu update
func (s *Service) UpdateNote(ctx context.Context, id int, title, content string) (Note, error) {
	if err := ValidateNoteInput(title, content); err != nil {
		return Note{}, err
	}

	note, ok := s.store.Update(id, title, content)
	if !ok {
		return Note{}, ErrNoteNotFound
	}
	return note, nil
}

// DeleteNote hapus note, error kalau tidak ada
func (s *Service) DeleteNote(ctx context.Context, id int) error {
	ok := s.store.Delete(id)
	if !ok {
		return ErrNoteNotFound
	}
	return nil
}
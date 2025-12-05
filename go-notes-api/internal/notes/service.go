package notes

import (
	"context"
)

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

// List Notes mengembalikan notes dengan pagination + search
func (s *Service) ListNotes(ctx context.Context, q string, page, limit int) ([]Note, int, error) {
	// validasi udah dilakuin di handler
	items, total := s.store.ListPage(page, limit, q)
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
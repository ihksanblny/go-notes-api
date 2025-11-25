package notes

import (
	"sync"
)

type Store interface {
	List() []Note
	Get(id int) (Note, bool)
	Create(title, content string) Note
	Update (id int, title, content string) (Note, bool)
	Delete (id int) bool
}

type InMemoryStore struct {
	mu   sync.RWMutex
	notes []Note
	nextID int
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		notes: []Note{},
		nextID: 1,
	}
}

func (s *InMemoryStore) List() []Note {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Return a copy to prevent external modification
	result := make([]Note, len(s.notes))
	copy(result, s.notes)
	return result
}

func (s *InMemoryStore) Get(id int) (Note, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, n := range s.notes {
		if n.ID == id {
			return n, true
		}
	}
	return Note{}, false
}

func (s *InMemoryStore) Create(title, content string) Note {
	s.mu.Lock()
	defer s.mu.Unlock()

	note := Note{
		ID: s.nextID,
		Title: title,
		Content: content,
	}
	s.nextID++
	s.notes = append(s.notes, note)
	return note
}

func (s *InMemoryStore) Update(id int, title, content string) (Note, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, n := range s.notes {
		if n.ID == id {
			n.Title = title
			n.Content = content
			s.notes[i] = n
			return n, true
		}
	}
	return Note{}, false
}

func (s *InMemoryStore) Delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, n := range s.notes {
		if n.ID == id {
			s.notes = append(s.notes[:i], s.notes[i+1:]...)
			return true
		}
	}
	return false
}

// ---- seeding 3 note awal ----

func SeedInitialNotes(store Store) {
	store.Create(
		"Clean Code",
		"Fokus ke kode yang mudah dibaca, konsisten, dan gampang di-maintain.",
	)
	store.Create(
		"Keamanan",
		"Validasi input, jangan expose info sensitif di error, dan batasi origin CORS.",
	)
	store.Create(
		"Skalabilitas",
		"Pisah concern (handler, store), pikirkan stateless dan gunakan DB jika skala membesar.",
	)
}
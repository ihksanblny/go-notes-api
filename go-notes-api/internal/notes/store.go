package notes

import (
	"database/sql"
	"log"
	"sync"
	"time"
)

type Store interface {
	List() []Note
	Get(id int) (Note, bool)
	Create(title, content string) Note
	Update(id int, title, content string) (Note, bool)
	Delete(id int) bool
}

// ---- InMemoryStore ----
type InMemoryStore struct {
	mu     sync.RWMutex
	notes  []Note
	nextID int
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		notes:  []Note{},
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

	now := time.Now()

	note := Note{
		ID:        s.nextID,
		Title:     title,
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
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
			n.UpdatedAt = time.Now()
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

// ---SqLite store
type SQLiteStore struct {
	db *sql.DB
}

func NewSQLiteStore(db *sql.DB) *SQLiteStore {
	return &SQLiteStore{db: db}
}

// InitSchema membuat table
func InitSchema(db *sql.DB) error {
	const query = `
	CREATE TABLE IF NOT EXISTS notes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);`

	_, err := db.Exec(query)
	return err
}

func (s *SQLiteStore) List() []Note {
	rows, err := s.db.Query(`SELECT id, title, content, created_at, updated_at FROM notes ORDER BY created_at DESC`)
	if err != nil {
		log.Printf("List notes error: %v", err)
		return []Note{}
	}
	defer rows.Close()

	var result []Note
	for rows.Next() {
		var n Note
		if err := rows.Scan(&n.ID, &n.Title, &n.Content, &n.CreatedAt, &n.UpdatedAt); err != nil {
			log.Printf("Scan note error: %v", err)
			continue
		}
		result = append(result, n)
	}
	if err := rows.Err(); err != nil {
		log.Printf("Rows error: %v", err)
	}
	return result
}

func (s *SQLiteStore) Get(id int) (Note, bool) {
	var n Note
	err := s.db.QueryRow(`SELECT id, title, content, created_at, updated_at FROM notes WHERE id = ?`, id).
		Scan(&n.ID, &n.Title, &n.Content, &n.CreatedAt, &n.UpdatedAt)
	if err == sql.ErrNoRows {
		return Note{}, false
	}
	if err != nil {
		log.Printf("Get note error: %v", err)
		return Note{}, false
	}
	return n, true
}

func (s *SQLiteStore) Create(title string, content string) Note {
	now := time.Now()

	res, err := s.db.Exec(
		`INSERT INTO notes (title, content, created_at, updated_at) VALUES (?, ?, ?, ?)`,
		title, content, time.Now(), time.Now(),
	)
	if err != nil {
		log.Printf("Create note error: %v", err)
		return Note{}
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Printf("Get note error: %v", err)
		return Note{}
	}

	return Note{
		ID:        int(id),
		Title:     title,
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (s *SQLiteStore) Update(id int, title string, content string) (Note, bool) {
	now := time.Now()

	res, err := s.db.Exec(
		`UPDATE notes SET title = ?, content = ?, updated_at = ? WHERE id = ?`,
		title, content, now, id,
	)
	if err != nil {
		log.Printf("Update Note error: %v", err)
		return Note{}, false
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		return Note{}, false
	}
	return s.Get(id)

}

func (s *SQLiteStore) Delete(id int) bool {
	res, err := s.db.Exec(`DELETE FROM notes WHERE id = ?`, id)
	if err != nil {
		log.Printf("Delete note error: %v", err)
		return false
	}

	affected, err := res.RowsAffected()
	if err != nil {
		log.Printf("RowsAffected error: %v", err)
	}
	return affected > 0
}

// ---- seeding 3 note awal ----

func SeedInitialNotes(store Store) {
	existing := store.List()
	if len(existing) > 0 {
		return
	}

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

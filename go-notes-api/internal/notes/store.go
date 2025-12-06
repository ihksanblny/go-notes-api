package notes

import (
	"database/sql"
	"log"
	"strings"
	"time"
	// "golang.org/x/text/cases"
)

type Store interface {
	List() []Note
	ListPage(page, limit int, query string, sortBy, sortOrder string) ([]Note, int)
	Get(id int) (Note, bool)
	Create(title, content string) Note
	Update(id int, title, content string) (Note, bool)
	Delete(id int) bool
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

func (s *SQLiteStore) ListPage(page, limit int, query string, sortBy, sortOrder string) ([]Note, int) {
	log.Printf("ListPage input sortBy=%q sortOrder=%q", sortBy, sortOrder)

	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	query = strings.TrimSpace(query)

	//-- Tentukan kolom string
	sortColumn := "created_at" // default
	switch sortBy {
	case "title":
		sortColumn = "title"
	case "updated_at":
		sortColumn = "updated_at"
	case "created_at":
		sortColumn = "created_at"
	}

	

	//-- Tentukan urutan
	sortDir := "DESC" //default
	if strings.ToLower(sortOrder) == "asc" {
		sortDir = "ASC"
	}

	//-- Hitung Total
	var total int
	baseCount := `SELECT COUNT(*) FROM notes`
	var countArgs []interface{}

	if query != "" {
		baseCount += ` WHERE title LIKE ? OR content LIKE ?`
		like := "%" + query + "%"
		countArgs = append(countArgs, like, like)
	}

	if err := s.db.QueryRow(baseCount, countArgs...).Scan(&total); err != nil {
		log.Printf("Count notes error: %v", err)
		return []Note{}, 0
	}
	if total == 0 {
		return []Note{}, 0
	}

	//-- ambil data page --
	offset := (page - 1) * limit
	baseSelect := `SELECT id, title, content, created_at, updated_at FROM notes`
	var selectArgs []interface{}

	if query != "" {
		baseSelect += ` WHERE title LIKE ? OR content LIKE ?`
		like := "%" + query + "%"
		selectArgs = append(selectArgs, like, like)
	}

	baseSelect += ` ORDER BY ` + sortColumn + ` ` + sortDir + ` LIMIT ? OFFSET ?`
	selectArgs = append(selectArgs, limit, offset)

	log.Printf("ListPage using ORDER BY %s %s", sortColumn, sortDir)
	log.Printf("ListPage SQL: %s | args=%v", baseSelect, selectArgs)

	rows, err := s.db.Query(baseSelect, selectArgs...)
	if err != nil {
		log.Printf("ListPage query error: %v", err)
		return []Note{}, total
	}
	defer rows.Close()

	var result []Note
	for rows.Next() {
		var n Note
		if err := rows.Scan(
			&n.ID,
			&n.Title,
			&n.Content,
			&n.CreatedAt,
			&n.UpdatedAt,
		); err != nil {
			log.Printf("Scan Note error: %v", err)
			continue
		}
		result = append(result, n)
	}
	if err := rows.Err(); err != nil {
		log.Printf("Rows error: %v", err)
	}
	return result, total

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

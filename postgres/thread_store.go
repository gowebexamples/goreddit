package postgres

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/gowebexamples/goreddit"
	"github.com/jmoiron/sqlx"
)

type ThreadStore struct {
	*sqlx.DB
}

func (s *ThreadStore) Thread(id uuid.UUID) (goreddit.Thread, error) {
	var t goreddit.Thread
	if err := s.Get(&t, `SELECT * FROM threads WHERE id = $1`, id); err != nil {
		return goreddit.Thread{}, fmt.Errorf("error getting thread: %w", err)
	}
	return t, nil
}

func (s *ThreadStore) Threads() ([]goreddit.Thread, error) {
	var tt []goreddit.Thread
	if err := s.Select(&tt, `SELECT * FROM threads`); err != nil {
		return []goreddit.Thread{}, fmt.Errorf("error getting threads: %w", err)
	}
	return tt, nil
}

func (s *ThreadStore) CreateThread(t *goreddit.Thread) error {
	if err := s.Get(t, `INSERT INTO threads VALUES ($1, $2, $3) RETURNING *`,
		t.ID,
		t.Title,
		t.Description); err != nil {
		return fmt.Errorf("error creating thread: %w", err)
	}
	return nil
}

func (s *ThreadStore) UpdateThread(t *goreddit.Thread) error {
	if err := s.Get(t, `UPDATE threads SET title = $1, description = $2 WHERE id = $3 RETURNING *`,
		t.Title,
		t.Description,
		t.ID); err != nil {
		return fmt.Errorf("error updating thread: %w", err)
	}
	return nil
}

func (s *ThreadStore) DeleteThread(id uuid.UUID) error {
	if _, err := s.Exec(`DELETE FROM threads WHERE id = $1`, id); err != nil {
		return fmt.Errorf("error deleting thread: %w", err)
	}
	return nil
}

package postgres

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/gowebexamples/goreddit"
	"github.com/jmoiron/sqlx"
)

type UserStore struct {
	*sqlx.DB
}

func (s *UserStore) User(id uuid.UUID) (goreddit.User, error) {
	var u goreddit.User
	if err := s.Get(&u, `SELECT * FROM users WHERE id = $1`, id); err != nil {
		return goreddit.User{}, fmt.Errorf("error getting user: %w", err)
	}
	return u, nil
}

func (s *UserStore) UserByUsername(username string) (goreddit.User, error) {
	var u goreddit.User
	if err := s.Get(&u, `SELECT * FROM users WHERE username = $1`, username); err != nil {
		return goreddit.User{}, fmt.Errorf("error getting user: %w", err)
	}
	return u, nil
}

func (s *UserStore) Users() ([]goreddit.User, error) {
	var uu []goreddit.User
	if err := s.Select(&uu, `SELECT * FROM users`); err != nil {
		return []goreddit.User{}, fmt.Errorf("error getting users: %w", err)
	}
	return uu, nil
}

func (s *UserStore) CreateUser(u *goreddit.User) error {
	if err := s.Get(u, `INSERT INTO users VALUES ($1, $2, $3) RETURNING *`,
		u.ID,
		u.Username,
		u.Password); err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}
	return nil
}

func (s *UserStore) UpdateUser(u *goreddit.User) error {
	if err := s.Get(u, `UPDATE users SET username = $1, password = $2 WHERE id = $3 RETURNING *`,
		u.Username,
		u.Password,
		u.ID); err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}
	return nil
}

func (s *UserStore) DeleteUser(id uuid.UUID) error {
	if _, err := s.Exec(`DELETE FROM users WHERE id = $1`, id); err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}
	return nil
}

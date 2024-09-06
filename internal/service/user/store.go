package user

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/4lerman/medods_tz/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db}
}

func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		return nil, err
	}

	user := new(types.User)
	for rows.Next() {
		user, err = scanRowIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}

	if user.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

func (s *Store) GetUserByID(id int) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE id=$1", id)
	if err != nil {
		return nil, err
	}

	user := new(types.User)
	for rows.Next() {
		user, err = scanRowIntoUser(rows)
		if err != nil {
			return user, nil
		}
	}

	if user.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

func (s *Store) CreateUser(user types.User) error {
	_, err := s.db.Exec("INSERT INTO users (firstName, lastName, email, "+
		"password) VALUES ($1, $2, $3, $4)", user.FirstName, user.LastName, user.Email, user.Password)

	ok, _ := json.Marshal(user)
	fmt.Println(string(ok))

	if err != nil {
		return err
	}

	return nil
}

func (s *Store) PutRefreshToken(userId int, refreshToken string) error {
	_, err := s.db.Exec("UPDATE users SET refreshToken = $1 WHERE id = $2", userId, refreshToken)

	if err != nil {
		return err
	}

	return nil
}

func scanRowIntoUser(rows *sql.Rows) (*types.User, error) {
	user := new(types.User)

	err := rows.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.RefreshToken,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

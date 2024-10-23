package user

import (
	"database/sql"
	"fmt"
	"go-mysql-restapi/constants/errors"
	"go-mysql-restapi/service/auth"
	"go-mysql-restapi/types"
	"log"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE email = ?", email)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	u := new(types.User)
	for rows.Next() {
		u, err = scanRowsIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}

	if u.ID == 0 {
		return nil, errors.ErrUserNotFound
	}

	return u, nil
}

func (s *Store) GetUserById(id string) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE id = ?", id)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	u := new(types.User)
	for rows.Next() {
		u, err = scanRowsIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}

	if u.ID == 0 {
		return nil, errors.ErrUserNotFound
	}

	return u, nil
}

func (s *Store) GetUserRoleByID(id string) (bool, error) {

	rows, err := s.db.Query("SELECT is_admin FROM user_role WHERE user_id = ?", id)
	if err != nil {
		log.Fatal(err)
		return false, err
	}

	var isAdmin bool
	for rows.Next() {
		err = rows.Scan(&isAdmin)
		if err != nil {
			return false, err
		}
	}

	return isAdmin, nil
}

func scanRowsIntoUser(rows *sql.Rows) (*types.User, error) {
	user := new(types.User)

	err := rows.Scan(
		&user.ID,
		&user.Email,
		&user.Age,
		&user.Name,
		&user.Password,
		&user.Balance,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Store) RegisterUser(user types.User) error {
	_, err := s.db.Exec(`
		SELECT * FROM users WHERE email = ?
	`, user.Email)

	if err != nil {
		return err
	}

	return nil
}

func (s *Store) CreateUser(user types.User) error {

	userCreated, err := s.db.Exec(`
		INSERT INTO users (email, age, name, password, balance) VALUES (?, ?, ?, ?, ?)
	`, user.Email, user.Age, user.Name, user.Password, user.Balance)

	if err != nil {
		return err
	}

	// Get ID of created user
	userId, _ := userCreated.LastInsertId()

	// Insert user role
	_, err = s.db.Exec(`
	INSERT INTO user_role (user_id, is_admin) VALUES (?, ?)
`, userId, 0)

	if err != nil {
		return err
	}

	return nil
}

func (s *Store) DeleteUser(userID string) error {

	_, err := s.db.Exec(`
		DELETE FROM users WHERE id = ?
	`, userID)

	if err != nil {
		return err
	}

	return nil
}

func (s *Store) LoginUser(email, password string) (string, error) {
	// Query to get the stored hashed password
	query := "SELECT id, password FROM users WHERE email = ?"
	var hashedPassword string
	var userID string
	// Scan method copies the columns in the current row into the values pointed at by dest
	err := s.db.QueryRow(query, email).Scan(&userID, &hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.ErrInvalidEmailOrPasswordLogin
		}
		return "", fmt.Errorf("failed to query user: %w", err)
	}

	// Compare the provided password with the stored hashed password
	if auth.ComparePasswords(hashedPassword, password) {
		return userID, nil // Login successful
	}

	return "", errors.ErrInvalidEmailOrPasswordLogin
}

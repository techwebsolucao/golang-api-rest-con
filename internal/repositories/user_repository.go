package repositories

import (
	"database/sql"

	"github.com/user/golang-api-rest/internal/models"
)

type UserRepository interface {
	Create(user *models.User) (int, error)
	GetByID(id int) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	All() ([]models.User, error)
	Update(user *models.User) error
	Delete(id int) error
}

type MySQLRepository struct {
	db *sql.DB
}

func NewMySQLRepository(db *sql.DB) *MySQLRepository {
	return &MySQLRepository{db: db}
}

func (r *MySQLRepository) Create(user *models.User) (int, error) {
	result, err := r.db.Exec(
		"INSERT INTO users (name, email, password_hash, role, verified) VALUES (?, ?, ?, ?, ?)",
		user.Name, user.Email, user.PasswordHash, user.Role, user.Verified,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (r *MySQLRepository) GetByID(id int) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRow(
		"SELECT id, name, email, password_hash, role, verified, created_at, updated_at FROM users WHERE id = ?",
		id,
	).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Role, &user.Verified, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *MySQLRepository) GetByEmail(email string) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRow(
		"SELECT id, name, email, password_hash, role, verified, created_at, updated_at FROM users WHERE email = ?",
		email,
	).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Role, &user.Verified, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *MySQLRepository) All() ([]models.User, error) {
	rows, err := r.db.Query("SELECT id, name, email, role, verified, created_at, updated_at FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []models.User{}
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Role, &u.Verified, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *MySQLRepository) Update(user *models.User) error {
	_, err := r.db.Exec(
		"UPDATE users SET name = ?, email = ?, role = ?, verified = ? WHERE id = ?",
		user.Name, user.Email, user.Role, user.Verified, user.ID,
	)
	return err
}

func (r *MySQLRepository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM users WHERE id = ?", id)
	return err
}
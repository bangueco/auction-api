package repositories

import (
	"context"

	"github.com/bangueco/auction-api/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	DB *pgxpool.Pool
}

func NewUserRepository(DB *pgxpool.Pool) *UserRepository {
	return &UserRepository{DB}
}

// Retrieve all users from the database
func (u *UserRepository) GetUsers() ([]models.User, error) {
	var users []models.User

	query := `SELECT * FROM users`

	rows, err := u.DB.Query(context.Background(), query)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var user models.User

		err := rows.Scan(&user.ID, &user.Username, &user.Password)

		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

// Retrieve a single user from the database by its ID
func (u *UserRepository) GetUserByID(id int64) (models.User, error) {
	var user models.User

	query := `SELECT * FROM users WHERE id = @id`
	namedArgs := pgx.NamedArgs{
		"id": id,
	}

	err := u.DB.QueryRow(context.Background(), query, namedArgs).Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		return user, err
	}

	return user, nil
}

// Retrieve a single user from the database by its username
func (u *UserRepository) GetUserByUsername(username string) (models.User, error) {
	var user models.User

	query := `SELECT * FROM users WHERE username = @username`
	namedArgs := pgx.NamedArgs{
		"username": username,
	}

	err := u.DB.QueryRow(context.Background(), query, namedArgs).Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		return user, err
	}

	return user, nil
}

// Create a new user in the database
func (u *UserRepository) CreateUser(user models.User) (models.User, error) {
	var newUser models.User

	query := `INSERT INTO users (username, password) VALUES (@username, @password) RETURNING id, username, password`
	namedArgs := pgx.NamedArgs{
		"username": user.Username,
		"password": user.Password,
	}

	err := u.DB.QueryRow(context.Background(), query, namedArgs).Scan(&newUser.ID, &newUser.Username, &newUser.Password)

	if err != nil {
		return newUser, err
	}

	return newUser, nil
}

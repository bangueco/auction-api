package services

import (
	"errors"
	"log"

	"github.com/bangueco/auction-api/internal/lib"
	"github.com/bangueco/auction-api/internal/models"
	"github.com/bangueco/auction-api/internal/repositories"
	"github.com/jackc/pgx/v5"
)

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrUsersNotFound = errors.New("users not found")
)

type UserService struct {
	UserRepo *repositories.UserRepository
}

func NewUserService(UserRepo *repositories.UserRepository) *UserService {
	return &UserService{UserRepo}
}

// Get all users from the database
func (u *UserService) GetUsers() ([]models.User, error) {
	users, err := u.UserRepo.GetUsers()

	if err != nil {
		log.Printf("Error retrieving users: %v", err)
		return nil, err
	}

	if len(users) == 0 {
		log.Printf("Error retrieving users: %v", err)
		return nil, ErrUsersNotFound
	}

	return users, err
}

// Get a single user from the database by its ID
func (u *UserService) GetUserByID(id int64) (models.User, error) {
	existingUser, err := u.UserRepo.GetUserByID(id)

	if err != nil {
		log.Printf("Error retrieving user: %v", err)
		if errors.Is(err, pgx.ErrNoRows) {
			return existingUser, ErrUserNotFound
		}
	}

	return existingUser, err
}

// Get a single user from the database by its username
func (u *UserService) GetUserByUsername(username string) (models.User, error) {
	existingUser, err := u.UserRepo.GetUserByUsername(username)

	if err != nil {
		log.Printf("Error retrieving user: %v", err)
		if errors.Is(err, pgx.ErrNoRows) {
			return existingUser, ErrUserNotFound
		}
	}

	return existingUser, err
}

// Create a new user in the database
func (u *UserService) CreateUser(user models.User) (models.User, error) {
	hashedPassword, err := lib.HashPassword(user.Password)

	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return user, err
	}

	var userDetails = models.User{
		Username: user.Username,
		Password: hashedPassword,
	}

	newUser, err := u.UserRepo.CreateUser(userDetails)

	if err != nil {
		log.Printf("Error creating user: %v", err)
		return newUser, err
	}

	return newUser, nil
}

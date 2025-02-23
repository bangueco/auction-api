package services

import (
	"errors"
	"log"

	"github.com/bangueco/auction-api/internal/models"
	"github.com/bangueco/auction-api/internal/repositories"
	"github.com/jackc/pgx/v5"
)

var (
	ErrItemNotFound  = errors.New("item not found")
	ErrItemsNotFound = errors.New("items not found")
)

type ItemService struct {
	ItemRepository *repositories.ItemRepository
}

func NewItemService(ItemRepository *repositories.ItemRepository) *ItemService {
	return &ItemService{ItemRepository}
}

// Retrieve all items from the database
func (i *ItemService) GetItems() ([]models.Item, error) {
	items, err := i.ItemRepository.GetItems()

	if err != nil {
		log.Printf("Error retrieving items: %v", err)
		return nil, err
	}

	if errors.Is(err, pgx.ErrNoRows) {
		log.Printf("Error retrieving items: %v", err)
		return nil, ErrItemsNotFound
	}

	return items, err
}

// Retrieve a single item from the database by its ID
func (i *ItemService) GetItemByID(id int64) (models.Item, error) {
	existingItem, err := i.ItemRepository.GetItemByID(id)

	if err != nil {
		log.Printf("Error retrieving item: %v", err)
		if errors.Is(err, pgx.ErrNoRows) {
			return existingItem, ErrItemNotFound
		}
	}

	return existingItem, err
}

// Create a new item in the database
func (i *ItemService) CreateItem(item models.Item) (models.Item, error) {
	newItem, err := i.ItemRepository.CreateItem(item)

	if err != nil {
		log.Printf("Error creating item: %v", err)
		return newItem, err
	}

	return newItem, nil
}

// Update an existing item in the database
func (i *ItemService) UpdateItem(itemId int64, data models.Item) (models.Item, error) {
	existingItem, err := i.GetItemByID(itemId)

	if err != nil {
		log.Printf("Error updating item: %v", err)
		return existingItem, err
	}

	newItem, err := i.ItemRepository.UpdateItem(itemId, data)

	if err != nil {
		log.Printf("Error updating item: %v", err)
		return newItem, err
	}

	return newItem, nil
}

// Delete an existing item from the database
func (i *ItemService) DeleteItem(itemID int64) error {
	_, err := i.GetItemByID(itemID)

	if err != nil {
		log.Printf("Error deleting item: %v", err)
		return err
	}

	err = i.ItemRepository.DeleteItem(itemID)

	if err != nil {
		log.Printf("Error deleting item: %v", err)
		return err
	}

	return nil
}

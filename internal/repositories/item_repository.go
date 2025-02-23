package repositories

import (
	"context"

	"github.com/bangueco/auction-api/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ItemRepository struct {
	DB *pgxpool.Pool
}

func NewItemRepository(DB *pgxpool.Pool) *ItemRepository {
	return &ItemRepository{DB}
}

// Retrieve all items from the database
func (i *ItemRepository) GetItems() ([]models.Item, error) {
	var items []models.Item

	query := `SELECT * FROM items`

	rows, err := i.DB.Query(context.Background(), query)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var item models.Item

		err := rows.Scan(&item.ID, &item.ItemName, &item.BidAmount, &item.AuctionedBy)

		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}

// Retrieve a single item from the database by its ID
func (i *ItemRepository) GetItemByID(id int64) (models.Item, error) {
	var item models.Item

	query := `SELECT * FROM items WHERE id = @id`
	namedArgs := pgx.NamedArgs{
		"id": id,
	}

	err := i.DB.QueryRow(context.Background(), query, namedArgs).Scan(&item.ID, &item.ItemName, &item.BidAmount, &item.AuctionedBy)

	if err != nil {
		return item, err
	}

	return item, nil
}

// Create a new item in the database
func (i *ItemRepository) CreateItem(item models.Item) (models.Item, error) {
	var newItem models.Item

	query := `INSERT INTO items (item_name, bid_amount, auctioned_by) VALUES (@item_name, @bid_amount, @auctioned_by) RETURNING id, item_name, bid_amount, auctioned_by`
	namedArgs := pgx.NamedArgs{
		"item_name":    item.ItemName,
		"bid_amount":   item.BidAmount,
		"auctioned_by": item.AuctionedBy,
	}

	err := i.DB.QueryRow(context.Background(), query, namedArgs).Scan(&newItem.ID, &newItem.ItemName, &newItem.BidAmount, &newItem.AuctionedBy)

	if err != nil {
		return newItem, err
	}

	return newItem, nil
}

// Update an existing item in the database
func (i *ItemRepository) UpdateItem(itemID int64, item models.Item) (models.Item, error) {
	var updatedItem models.Item

	query := `UPDATE items SET item_name = @item_name, bid_amount = @bid_amount, auctioned_by = @auctioned_by WHERE id = @id RETURNING id, item_name, bid_amount, auctioned_by`
	namedArgs := pgx.NamedArgs{
		"id":           itemID,
		"item_name":    item.ItemName,
		"bid_amount":   item.BidAmount,
		"auctioned_by": item.AuctionedBy,
	}

	err := i.DB.QueryRow(context.Background(), query, namedArgs).Scan(&updatedItem.ID, &updatedItem.ItemName, &updatedItem.BidAmount, &updatedItem.AuctionedBy)

	if err != nil {
		return updatedItem, err
	}

	return updatedItem, nil
}

// Delete an item from the database by its ID
func (i *ItemRepository) DeleteItem(id int64) error {
	query := `DELETE FROM items WHERE id = @id`
	namedArgs := pgx.NamedArgs{
		"id": id,
	}

	_, err := i.DB.Exec(context.Background(), query, namedArgs)

	if err != nil {
		return err
	}

	return nil
}

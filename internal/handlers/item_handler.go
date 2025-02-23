package handlers

import (
	"errors"
	"net/http"

	"github.com/bangueco/auction-api/internal/handlers/helper"
	"github.com/bangueco/auction-api/internal/lib"
	"github.com/bangueco/auction-api/internal/models"
	"github.com/bangueco/auction-api/internal/services"
	"github.com/go-chi/chi/v5"
)

type ItemHandler struct {
	ItemService *services.ItemService
}

func NewItemHandler(ItemService *services.ItemService) *ItemHandler {
	return &ItemHandler{ItemService}
}

func (i *ItemHandler) GetItems(w http.ResponseWriter, r *http.Request) {
	items, err := i.ItemService.GetItems()

	if errors.Is(err, services.ErrItemsNotFound) {
		helper.WriteResponseMessage(w, "No items found", http.StatusNotFound)
		return
	}

	if err != nil {
		helper.WriteResponseMessage(w, "Error retrieving items", http.StatusInternalServerError)
		return
	}

	if len(items) == 0 {
		helper.WriteResponseMessage(w, "No items found", http.StatusNotFound)
		return
	}

	helper.WriteResponse(w, items, http.StatusOK)
}

func (i *ItemHandler) GetItemByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := helper.ConvertStringToInt64(idParam)

	if err != nil {
		helper.WriteResponseMessage(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	item, err := i.ItemService.GetItemByID(id)

	if errors.Is(err, services.ErrItemNotFound) {
		helper.WriteResponseMessage(w, "Item not found", http.StatusNotFound)
		return
	}

	helper.WriteResponse(w, item, http.StatusOK)
}

func (i *ItemHandler) CreateItem(w http.ResponseWriter, r *http.Request) {
	var newItem models.Item

	err := helper.DecodeRequestBody(r, &newItem)

	if err != nil {
		helper.WriteResponseMessage(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	errorMessages := lib.ValidateStruct(&newItem)

	if errorMessages != nil {
		helper.WriteResponse(w, errorMessages, http.StatusBadRequest)
		return
	}

	item, err := i.ItemService.CreateItem(newItem)

	if err != nil {
		helper.WriteResponseMessage(w, "Error creating item", http.StatusInternalServerError)
		return
	}

	helper.WriteResponse(w, item, http.StatusCreated)
}

func (i *ItemHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	var item models.Item
	idParam := chi.URLParam(r, "id")

	id, err := helper.ConvertStringToInt64(idParam)

	if err != nil {
		helper.WriteResponseMessage(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	err = helper.DecodeRequestBody(r, &item)

	if err != nil {
		helper.WriteResponseMessage(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	errorMessages := lib.ValidateStruct(&item)

	if errorMessages != nil {
		helper.WriteResponse(w, errorMessages, http.StatusBadRequest)
		return
	}

	updatedItem, err := i.ItemService.UpdateItem(id, item)

	if errors.Is(err, services.ErrItemNotFound) {
		helper.WriteResponseMessage(w, "Item not found", http.StatusNotFound)
		return
	}

	helper.WriteResponse(w, updatedItem, http.StatusOK)
}

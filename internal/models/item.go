package models

type Item struct {
	ID          int64   `json:"id,omitempty"`
	ItemName    string  `json:"item_name,omitempty" validate:"required,min=3,max=100"`
	BidAmount   float64 `json:"bid_amount,omitempty" validate:"required,numeric,min=1"`
	AuctionedBy int64   `json:"auctioned_by,omitempty" validate:"required,numeric"`
}

package entity

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID           string
	SellingOrder *Order
	BuyingOrder  *Order
	Shares       int
	Price        float64
	Total        float64
	DateTime     time.Time
}

func NewTransaction(SellingOrder, buyingOrder *Order, shares int, price float64) *Transaction {
	total := float64(shares) * price

	return &Transaction{
		ID:           uuid.New().String(),
		SellingOrder: SellingOrder,
		BuyingOrder:  buyingOrder,
		Shares:       shares,
		Price:        price,
		Total:        total,
		DateTime:     time.Now(),
	}
}

func (t *Transaction) TotalPrice() {
	t.Total = float64(t.Shares) * t.BuyingOrder.Price
}

func (t *Transaction) Close() {
	if t.BuyingOrder.PendingShares == 0 {
		t.BuyingOrder.Status = "CLOSED"
	}

	if t.SellingOrder.PendingShares == 0 {
		t.SellingOrder.Status = "CLOSED"
	}
}

func (t *Transaction) UpdateTransactionShares() {
	sellingShares := t.SellingOrder.PendingShares
	buyingShares := t.BuyingOrder.PendingShares

	minShares := sellingShares
	if minShares > buyingShares {
		minShares = buyingShares
	}

	t.SellingOrder.Investor.UpdateAssetPosition(t.SellingOrder.Asset.ID, -minShares)
	t.SellingOrder.PendingShares -= minShares

	t.BuyingOrder.Investor.UpdateAssetPosition(t.SellingOrder.Asset.ID, minShares)
	t.BuyingOrder.PendingShares += minShares
}

package kaimono

type Cart struct {
	ID        string     `json:"id"`
	Items     []CartItem `json:"cart-items"`
	Discounts []Discount `json:"discounts"`
}

type CartItem struct {
	ProductID string     `json:"product-id"`
	Quantity  int        `json:"quantity"`
	Discounts []Discount `json:"discounts"`
	Price     Price      `json:"price"`
}

type Price struct {
	Currency Currency `json:"currency"`
	Value    float64  `json:"value"`
}

type Discount struct {
	ID            string       `json:"id"`
	Type          DiscountType `json:"type"`
	PercentageOff float64      `json:"percentage-off"`
	AmountOff     float64      `json:"amount-off"`
}

type DiscountType string

const (
	PercentageDiscount  DiscountType = "percentage"
	FixedAmountDiscount DiscountType = "fixed-amount"
)

// ISO 4172 three letter code (can't remember if that is the currency ISO).
type Currency string

const (
	USD Currency = "usd"
	EUR Currency = "euro"
	BTC Currency = "bitcoin"
)

package kaimono

type Cart struct {
	ID        string     `json:"id"`
	Items     []CartItem `json:"items"`
	Discounts []Discount `json:"discounts"`
}

type CartItem struct {
	ID        string     `json:"id"`
	Quantity  int        `json:"quantity"`
	Discounts []Discount `json:"discounts"`
	Price     Price      `json:"price"`
}

type Price struct {
	Currency string  `json:"currency"`
	Value    float64 `json:"value"`
}

type DiscountType string

const (
	PercentageDiscount  DiscountType = "percentage"
	FixedAmountDiscount DiscountType = "fixed-amount"
)

type Discount struct {
	ID    string       `json:"id"`
	Type  DiscountType `json:"type"`
	Value float64      `json:"value"`
}

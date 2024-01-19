package responsedto

import "github.com/shopspring/decimal"

type ShippingMethodResponse struct {
	Id   int64           `json:"id"`
	Name string          `json:"name"`
	Cost decimal.Decimal `json:"cost"`
}

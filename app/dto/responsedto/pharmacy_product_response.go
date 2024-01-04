package responsedto

import "github.com/shopspring/decimal"

type PharmacyProductResponse struct {
	Id         int64           `json:"id"`
	PharmacyId int64           `json:"pharmacy_id"`
	ProductId  int64           `json:"product_id"`
	IsActive   bool            `json:"is_active"`
	Price      decimal.Decimal `json:"price"`
	Stock      int32           `json:"stock"`
}

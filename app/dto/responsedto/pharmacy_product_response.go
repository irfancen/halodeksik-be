package responsedto

import "github.com/shopspring/decimal"

type PharmacyProductResponse struct {
	Id              int64            `json:"id,omitempty"`
	PharmacyId      int64            `json:"pharmacy_id,omitempty"`
	ProductId       int64            `json:"product_id,omitempty"`
	IsActive        bool             `json:"is_active,omitempty"`
	Price           decimal.Decimal  `json:"price,omitempty"`
	Stock           int32            `json:"stock,omitempty"`
	ProductResponse *ProductResponse `json:"product,omitempty"`
}

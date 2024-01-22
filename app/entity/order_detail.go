package entity

import (
	"database/sql"
	"github.com/shopspring/decimal"
	"halodeksik-be/app/dto/responsedto"
	"time"
)

type OrderDetail struct {
	Id          int64           `json:"id"`
	OrderId     int64           `json:"order_id"`
	ProductId   int64           `json:"product_id"`
	Quantity    int32           `json:"quantity"`
	Name        string          `json:"name"`
	GenericName string          `json:"generic_name"`
	Content     string          `json:"content"`
	Description string          `json:"description"`
	Image       string          `json:"image"`
	Price       decimal.Decimal `json:"price"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	DeletedAt   sql.NullTime    `json:"deleted_at"`
}

func (o OrderDetail) ToOrderDetailResponse() responsedto.OrderDetailResponse {
	return responsedto.OrderDetailResponse{
		Name:        o.Name,
		GenericName: o.GenericName,
		Content:     o.Content,
		Description: o.Description,
		Image:       o.Image,
		Price:       o.Price.String(),
		Quantity:    o.Quantity,
	}
}

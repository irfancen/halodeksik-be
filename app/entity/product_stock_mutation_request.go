package entity

import (
	"database/sql"
	"fmt"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/dto/responsedto"
	"reflect"
	"time"
)

type ProductStockMutationRequest struct {
	Id                                  int64        `json:"id"`
	PharmacyProductOriginId             int64        `json:"pharmacy_product_origin_id"`
	PharmacyProductDestId               int64        `json:"pharmacy_product_dest_id"`
	Stock                               int32        `json:"stock"`
	ProductStockMutationRequestStatusId int64        `json:"product_stock_mutation_request_status_id"`
	CreatedAt                           time.Time    `json:"created_at"`
	UpdatedAt                           time.Time    `json:"updated_at"`
	DeletedAt                           sql.NullTime `json:"deleted_at"`
}

func (e *ProductStockMutationRequest) GetEntityName() string {
	return "product_stock_mutation_requests"
}

func (e *ProductStockMutationRequest) GetFieldStructTag(fieldName string, structTag string) string {
	field, ok := reflect.TypeOf(e).Elem().FieldByName(fieldName)
	if !ok {
		return ""
	}
	return field.Tag.Get(structTag)
}

func (e *ProductStockMutationRequest) GetSqlColumnFromField(fieldName string) string {
	return fmt.Sprintf("%s.%s", e.GetEntityName(), e.GetFieldStructTag(fieldName, appconstant.JsonStructTag))
}

func (e *ProductStockMutationRequest) ToResponse() responsedto.ProductStockMutationRequestResponse {
	return responsedto.ProductStockMutationRequestResponse{
		Id:                                  e.Id,
		PharmacyProductOriginId:             e.PharmacyProductOriginId,
		PharmacyProductDestId:               e.PharmacyProductDestId,
		Stock:                               e.Stock,
		ProductStockMutationRequestStatusId: e.ProductStockMutationRequestStatusId,
		RequestDate:                         e.CreatedAt,
	}
}

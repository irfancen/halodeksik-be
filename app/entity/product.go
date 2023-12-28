package entity

import (
	"database/sql"
	"reflect"
	"time"
)

type Product struct {
	ID                   int64        `json:"id"`
	Name                 string       `json:"name"`
	GenericName          string       `json:"generic_name"`
	Content              string       `json:"content"`
	ManufacturerID       int64        `json:"manufacturer_id"`
	Description          string       `json:"description"`
	DrugClassificationID int64        `json:"drug_classification_id"`
	ProductCategoryID    int64        `json:"product_category_id"`
	DrugForm             string       `json:"drug_form"`
	UnitInPack           string       `json:"unit_in_pack"`
	SellingUnit          string       `json:"selling_unit"`
	Weight               float64      `json:"weight"`
	Length               float64      `json:"length"`
	Width                float64      `json:"width"`
	Height               float64      `json:"height"`
	Image                string       `json:"image"`
	Price                string       `json:"price"`
	CreatedAt            time.Time    `json:"created_at"`
	UpdatedAt            time.Time    `json:"updated_at"`
	DeletedAt            sql.NullTime `json:"deleted_at"`
}

func (u *Product) GetEntityName() string {
	return "users"
}

func (u *Product) GetFieldStructTag(fieldName string, structTag string) string {
	field, ok := reflect.TypeOf(u).Elem().FieldByName(fieldName)
	if !ok {
		return ""
	}
	return field.Tag.Get(structTag)
}

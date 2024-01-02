package responsedto

import (
	"time"
)

type ProductResponse struct {
	Id                   int64           `json:"id"`
	Name                 string          `json:"name"`
	GenericName          string          `json:"generic_name"`
	Content              string          `json:"content"`
	ManufacturerId       int64           `json:"manufacturer_id"`
	Description          string          `json:"description"`
	DrugClassificationId int64           `json:"drug_classification_id"`
	ProductCategoryId    int64           `json:"product_category_id"`
	DrugForm             string          `json:"drug_form"`
	UnitInPack           string          `json:"unit_in_pack"`
	SellingUnit          string          `json:"selling_unit"`
	Weight               float64         `json:"weight"`
	Length               float64         `json:"length"`
	Width                float64         `json:"width"`
	Height               float64         `json:"height"`
	Image                string          `json:"image"`
	CreatedAt            time.Time       `json:"created_at"`
	UpdatedAt            time.Time       `json:"updated_at"`
}

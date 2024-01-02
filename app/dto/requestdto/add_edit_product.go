package requestdto

import (
	"halodeksik-be/app/entity"
)

type AddEditProduct struct {
	Name                 string  `json:"name" validate:"required"`
	GenericName          string  `json:"generic_name" validate:"required"`
	Content              string  `json:"content" validate:"required"`
	ManufacturerId       int64   `json:"manufacturer_id" validate:"required"`
	Description          string  `json:"description" validate:"required"`
	DrugClassificationId int64   `json:"drug_classification_id" validate:"required"`
	ProductCategoryId    int64   `json:"product_category_id" validate:"required"`
	DrugForm             string  `json:"drug_form" validate:"required"`
	UnitInPack           string  `json:"unit_in_pack" validate:"required"`
	SellingUnit          string  `json:"selling_unit" validate:"required"`
	Weight               float64 `json:"weight" validate:"required"`
	Length               float64 `json:"length" validate:"required"`
	Width                float64 `json:"width" validate:"required"`
	Height               float64 `json:"height" validate:"required"`
	Image                string  `json:"image"`
}

func (r AddEditProduct) ToProduct() entity.Product {
	return entity.Product{
		Name:                 r.Name,
		GenericName:          r.GenericName,
		Content:              r.Content,
		ManufacturerId:       r.ManufacturerId,
		Description:          r.Description,
		DrugClassificationId: r.DrugClassificationId,
		ProductCategoryId:    r.ProductCategoryId,
		DrugForm:             r.DrugForm,
		UnitInPack:           r.UnitInPack,
		SellingUnit:          r.SellingUnit,
		Weight:               r.Weight,
		Length:               r.Length,
		Width:                r.Width,
		Height:               r.Height,
		Image:                r.Image,
	}
}

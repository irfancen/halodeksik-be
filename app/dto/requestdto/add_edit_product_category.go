package requestdto

import "halodeksik-be/app/entity"

type AddEditProductCategory struct {
	Name string `json:"name" validate:"required"`
}

func (r AddEditProductCategory) ToProductCategory() entity.ProductCategory {
	return entity.ProductCategory{
		Name: r.Name,
	}
}

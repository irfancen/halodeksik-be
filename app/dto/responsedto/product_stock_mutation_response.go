package responsedto

type ProductStockMutationResponse struct {
	Id                         int64 `json:"id"`
	PharmacyProductId          int64 `json:"pharmacy_product_id"`
	ProductStockMutationTypeId int64 `json:"product_stock_mutation_type_id"`
	Stock                      int32 `json:"stock"`
}

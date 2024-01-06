package usecase

import (
	"context"
	"errors"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/repository"
)

type CartItemUseCase interface {
	Add(ctx context.Context, cartItem entity.CartItem) (*entity.CartItem, error)
	GetByUserIdAndProductId(ctx context.Context, cartItem entity.CartItem) (*entity.CartItem, error)
	GetAllByUserId(ctx context.Context) (*entity.PaginatedItems, error)
	Edit(ctx context.Context, existingCartItem entity.CartItem, cartItem entity.CartItem) (*entity.CartItem, error)
	Remove(ctx context.Context, productIds []int64) error
}

type CartItemUseCaseImpl struct {
	cartItemRepo        repository.CartItemRepository
	productRepo         repository.ProductRepository
	pharmacyProductRepo repository.PharmacyProductRepository
}

func NewCartItemUseCaseImpl(
	cartItemRepo repository.CartItemRepository,
	productRepo repository.ProductRepository,
	pharmacyProductRepo repository.PharmacyProductRepository,
) *CartItemUseCaseImpl {
	return &CartItemUseCaseImpl{cartItemRepo: cartItemRepo, productRepo: productRepo, pharmacyProductRepo: pharmacyProductRepo}
}

func (uc *CartItemUseCaseImpl) Add(ctx context.Context, cartItem entity.CartItem) (*entity.CartItem, error) {
	userId := ctx.Value(appconstant.ContextKeyUserId).(int64)
	cartItem.UserId = userId

	product, err := uc.productRepo.FindById(ctx, cartItem.ProductId)
	if err != nil && errors.Is(err, apperror.ErrRecordNotFound) {
		return nil, apperror.NewNotFound(product, "Id", cartItem.ProductId)
	}

	found, err := uc.cartItemRepo.FindByUserIdAndProductId(ctx, cartItem.UserId, cartItem.ProductId)
	if err == nil {
		return uc.Edit(ctx, *found, cartItem)
	}
	if err != nil && !errors.Is(err, apperror.ErrRecordNotFound) {
		return nil, err
	}

	pharmacyProducts, err := uc.pharmacyProductRepo.FindAllByProductId(ctx, cartItem.ProductId)
	if err != nil {
		return nil, err
	}
	totalProductStock := int32(0)
	for _, pharmacyProduct := range pharmacyProducts {
		totalProductStock += pharmacyProduct.Stock
	}
	if cartItem.Quantity > totalProductStock {
		return nil, apperror.ErrProductStockNotEnoughToAddToCart
	}

	created, err := uc.cartItemRepo.Create(ctx, cartItem)
	if err != nil {
		return nil, err
	}
	return created, nil
}

func (uc *CartItemUseCaseImpl) GetByUserIdAndProductId(ctx context.Context, cartItem entity.CartItem) (*entity.CartItem, error) {
	userId := ctx.Value(appconstant.ContextKeyUserId).(int64)
	cartItem.UserId = userId

	got, err := uc.cartItemRepo.FindByUserIdAndProductId(ctx, cartItem.UserId, cartItem.ProductId)
	if err != nil {
		if errors.Is(err, apperror.ErrRecordNotFound) {
			return nil, apperror.NewNotFound(got, "ProductId", cartItem.ProductId)
		}
		return nil, err
	}

	return got, nil
}

func (uc *CartItemUseCaseImpl) GetAllByUserId(ctx context.Context) (*entity.PaginatedItems, error) {
	userId := ctx.Value(appconstant.ContextKeyUserId).(int64)
	cartItems, err := uc.cartItemRepo.FindAllByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}

	paginatedItems := entity.NewPaginationInfo(
		int64(len(cartItems)), 1, int64(len(cartItems)), 1, cartItems,
	)

	return paginatedItems, nil
}

func (uc *CartItemUseCaseImpl) Edit(ctx context.Context, existingCartItem entity.CartItem, cartItem entity.CartItem) (*entity.CartItem, error) {
	userId := ctx.Value(appconstant.ContextKeyUserId).(int64)
	cartItem.UserId = userId

	pharmacyProducts, err := uc.pharmacyProductRepo.FindAllByProductId(ctx, cartItem.ProductId)
	if err != nil {
		return nil, err
	}
	totalProductStock := int32(0)
	for _, pharmacyProduct := range pharmacyProducts {
		totalProductStock += pharmacyProduct.Stock
	}
	if existingCartItem.Quantity + cartItem.Quantity > totalProductStock {
		return nil, apperror.ErrProductStockNotEnoughToAddToCart
	}

	updated, err := uc.cartItemRepo.Update(ctx, cartItem)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (uc *CartItemUseCaseImpl) Remove(ctx context.Context, productIds []int64) error {
	userId := ctx.Value(appconstant.ContextKeyUserId).(int64)

	err := uc.cartItemRepo.Delete(ctx, userId, productIds)
	if err != nil {
		return err
	}
	return nil
}

package repository

import (
	"context"
	"database/sql"
	"errors"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/dto/queryparamdto"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/util"
)

type OrderRepository interface {
	FindAllOrdersByPharmacyAdminId(ctx context.Context, param *queryparamdto.GetAllParams, adminId int64) ([]*entity.Order, error)
	FindAllOrders(ctx context.Context, param *queryparamdto.GetAllParams) ([]*entity.Order, error)
	FindAllOrdersByUserId(ctx context.Context, param *queryparamdto.GetAllParams, userId int64) ([]*entity.Order, error)
	CountFindAllOrdersByPharmacyAdminId(ctx context.Context, adminId int64, param *queryparamdto.GetAllParams) (int64, error)
	CountFindAllOrdersByUserId(ctx context.Context, adminId int64, param *queryparamdto.GetAllParams) (int64, error)
	CountFindAllOrders(ctx context.Context, param *queryparamdto.GetAllParams) (int64, error)
	FindOrderById(ctx context.Context, id int64) (*entity.Order, *entity.OrderIds, error)
	UpdateOrderStatus(ctx context.Context, orderId int64, orderLog entity.OrderStatusLog) (*entity.OrderStatusLog, error)
	FindLatestOrderStatusByOrderId(ctx context.Context, id int64) (*entity.OrderStatusLog, error)
}

type OrderRepositoryImpl struct {
	db *sql.DB
}

func NewOrderRepositoryImpl(db *sql.DB) OrderRepository {
	return &OrderRepositoryImpl{db: db}
}

func (repo *OrderRepositoryImpl) FindAllOrders(ctx context.Context, param *queryparamdto.GetAllParams) ([]*entity.Order, error) {
	const findAllOrder = `SELECT orders.id, pharmacy_id, pharmacies.name, orders.date, no_of_items, orders.total_payment, 
       transaction_id, order_statuses.id,order_statuses.name
	FROM orders
			 INNER JOIN pharmacies ON orders.pharmacy_id = pharmacies.id
			 INNER JOIN transactions ON orders.transaction_id = transactions.id
			 INNER JOIN order_status_logs ON orders.id = order_status_logs.order_id
			 INNER JOIN order_statuses ON order_status_logs.order_status_id = order_statuses.id
	WHERE transactions.transaction_status_id = 4 AND order_status_logs.is_latest IS TRUE AND orders.deleted_at IS NULL `

	query, values := buildQuery(findAllOrder, &entity.Transaction{}, param, true, true)
	rows, err := repo.db.QueryContext(ctx, query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*entity.Order, 0)
	for rows.Next() {
		var order entity.Order
		var pharmacy entity.Pharmacy
		var status entity.OrderStatus
		if err := rows.Scan(
			&order.Id,
			&pharmacy.Id,
			&pharmacy.Name,
			&order.Date,
			&order.NoOfItems,
			&order.TotalPayment,
			&order.TransactionId,
			&status.Id,
			&status.Name,
		); err != nil {
			return nil, err
		}
		order.Pharmacy = &pharmacy
		order.LatestStatus = &status
		items = append(items, &order)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (repo *OrderRepositoryImpl) FindAllOrdersByPharmacyAdminId(ctx context.Context, param *queryparamdto.GetAllParams, adminId int64) ([]*entity.Order, error) {
	const findAllOrderUserId = `SELECT orders.id, pharmacy_id, pharmacies.name, orders.date, no_of_items, orders.total_payment, 
       transaction_id, order_statuses.id,order_statuses.name
	FROM orders
			 INNER JOIN pharmacies ON orders.pharmacy_id = pharmacies.id
			 INNER JOIN transactions ON orders.transaction_id = transactions.id
			 INNER JOIN order_status_logs ON orders.id = order_status_logs.order_id
			 INNER JOIN order_statuses ON order_status_logs.order_status_id = order_statuses.id
	WHERE pharmacies.pharmacy_admin_id = $1 AND transactions.transaction_status_id = 4 AND order_status_logs.is_latest IS TRUE AND orders.deleted_at IS NULL `

	indexPreparedStatement := 1
	query, values := buildQuery(findAllOrderUserId, &entity.Order{}, param, true, true, indexPreparedStatement)
	values = util.AppendAtIndex(values, 0, interface{}(adminId))

	rows, err := repo.db.QueryContext(ctx, query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*entity.Order, 0)
	for rows.Next() {
		var order entity.Order
		var pharmacy entity.Pharmacy
		var status entity.OrderStatus
		if err := rows.Scan(
			&order.Id,
			&pharmacy.Id,
			&pharmacy.Name,
			&order.Date,
			&order.NoOfItems,
			&order.TotalPayment,
			&order.TransactionId,
			&status.Id,
			&status.Name,
		); err != nil {
			return nil, err
		}
		order.Pharmacy = &pharmacy
		order.LatestStatus = &status
		items = append(items, &order)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (repo *OrderRepositoryImpl) FindAllOrdersByUserId(ctx context.Context, param *queryparamdto.GetAllParams, userId int64) ([]*entity.Order, error) {
	const findAllOrderUserId = `SELECT orders.id, pharmacy_id, pharmacies.name, orders.date, no_of_items, orders.total_payment, 
       transaction_id, order_statuses.id,order_statuses.name
	FROM orders
			 INNER JOIN pharmacies ON orders.pharmacy_id = pharmacies.id
			 INNER JOIN transactions ON orders.transaction_id = transactions.id
			 INNER JOIN order_status_logs ON orders.id = order_status_logs.order_id
			 INNER JOIN order_statuses ON order_status_logs.order_status_id = order_statuses.id
	WHERE transactions.user_id = $1 AND transactions.transaction_status_id = 4 AND order_status_logs.is_latest IS TRUE AND orders.deleted_at IS NULL `

	indexPreparedStatement := 1
	query, values := buildQuery(findAllOrderUserId, &entity.Order{}, param, true, true, indexPreparedStatement)
	values = util.AppendAtIndex(values, 0, interface{}(userId))

	rows, err := repo.db.QueryContext(ctx, query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*entity.Order, 0)
	for rows.Next() {
		var order entity.Order
		var pharmacy entity.Pharmacy
		var status entity.OrderStatus
		if err := rows.Scan(
			&order.Id,
			&pharmacy.Id,
			&pharmacy.Name,
			&order.Date,
			&order.NoOfItems,
			&order.TotalPayment,
			&order.TransactionId,
			&status.Id,
			&status.Name,
		); err != nil {
			return nil, err
		}
		order.Pharmacy = &pharmacy
		order.LatestStatus = &status
		items = append(items, &order)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil

}

func (repo *OrderRepositoryImpl) CountFindAllOrders(ctx context.Context, param *queryparamdto.GetAllParams) (int64, error) {

	const findAllOrder = `SELECT count(orders.id)
	FROM orders
			 INNER JOIN pharmacies ON orders.pharmacy_id = pharmacies.id
			 INNER JOIN transactions ON orders.transaction_id = transactions.id
			 INNER JOIN order_status_logs ON orders.id = order_status_logs.order_id
			 INNER JOIN order_statuses ON order_status_logs.order_status_id = order_statuses.id
	WHERE transactions.transaction_status_id = 4 AND order_status_logs.is_latest IS TRUE AND orders.deleted_at IS NULL `

	query, values := buildQuery(findAllOrder, &entity.Order{}, param, false, false)
	var totalItems int64
	row := repo.db.QueryRowContext(ctx, query, values...)
	if row.Err() != nil {
		return totalItems, row.Err()
	}

	if err := row.Scan(
		&totalItems,
	); err != nil {
		return totalItems, err
	}

	return totalItems, nil
}

func (repo *OrderRepositoryImpl) CountFindAllOrdersByPharmacyAdminId(ctx context.Context, adminId int64, param *queryparamdto.GetAllParams) (int64, error) {

	const findAllOrderUserId = `SELECT count(orders.id)
	FROM orders
			 INNER JOIN pharmacies ON orders.pharmacy_id = pharmacies.id
			 INNER JOIN transactions ON orders.transaction_id = transactions.id
			 INNER JOIN order_status_logs ON orders.id = order_status_logs.order_id
			 INNER JOIN order_statuses ON order_status_logs.order_status_id = order_statuses.id
	WHERE pharmacies.pharmacy_admin_id = $1 AND transactions.transaction_status_id = 4 AND order_status_logs.is_latest IS TRUE AND orders.deleted_at IS NULL `
	indexPreparedStatement := 1
	query, values := buildQuery(findAllOrderUserId, &entity.Order{}, param, false, false, indexPreparedStatement)
	values = util.AppendAtIndex(values, 0, interface{}(adminId))

	var totalItems int64

	row := repo.db.QueryRowContext(ctx, query, values...)
	if row.Err() != nil {
		return totalItems, row.Err()
	}

	if err := row.Scan(
		&totalItems,
	); err != nil {
		return totalItems, err
	}

	return totalItems, nil
}

func (repo *OrderRepositoryImpl) CountFindAllOrdersByUserId(ctx context.Context, adminId int64, param *queryparamdto.GetAllParams) (int64, error) {

	const findAllOrderUserId = `SELECT count(orders.id)
	FROM orders
			 INNER JOIN pharmacies ON orders.pharmacy_id = pharmacies.id
			 INNER JOIN transactions ON orders.transaction_id = transactions.id
			 INNER JOIN order_status_logs ON orders.id = order_status_logs.order_id
			 INNER JOIN order_statuses ON order_status_logs.order_status_id = order_statuses.id
	WHERE transactions.user_id = $1 AND transactions.transaction_status_id = 4 AND order_status_logs.is_latest IS TRUE AND orders.deleted_at IS NULL `

	indexPreparedStatement := 1
	query, values := buildQuery(findAllOrderUserId, &entity.Order{}, param, false, false, indexPreparedStatement)
	values = util.AppendAtIndex(values, 0, interface{}(adminId))
	var totalItems int64

	row := repo.db.QueryRowContext(ctx, query, values...)
	if row.Err() != nil {
		return totalItems, row.Err()
	}

	if err := row.Scan(
		&totalItems,
	); err != nil {
		return totalItems, err
	}

	return totalItems, nil
}

func (repo *OrderRepositoryImpl) FindOrderById(ctx context.Context, id int64) (*entity.Order, *entity.OrderIds, error) {
	const findOrderById = `SELECT DISTINCT orders.id, order_statuses.id, order_statuses.name, orders.date, shipping_method_id, 
                shipping_methods.name, shipping_cost, pharmacy_id, 
                pharmacies.name, transactions.address, orders.total_payment, transactions.user_id, pharmacies.pharmacy_admin_id
	FROM orders
		INNER JOIN shipping_methods ON orders.shipping_method_id = shipping_methods.id
		INNER JOIN transactions ON orders.transaction_id = transactions.id
		INNER JOIN order_status_logs ON orders.id = order_status_logs.order_id
		INNER JOIN order_statuses ON order_status_logs.order_status_id = order_statuses.id
		INNER JOIN pharmacies ON orders.pharmacy_id = pharmacies.id
	WHERE orders.id = $1 AND order_status_logs.is_latest IS TRUE`

	row := repo.db.QueryRowContext(ctx, findOrderById, id)
	var order entity.Order
	var status entity.OrderStatus
	var shipping entity.ShippingMethod
	var pharmacy entity.Pharmacy
	var ids entity.OrderIds

	err := row.Scan(
		&order.Id,
		&status.Id,
		&status.Name,
		&order.Date,
		&shipping.Id,
		&shipping.Name,
		&order.ShippingCost,
		&pharmacy.Id,
		&pharmacy.Name,
		&order.UserAddress,
		&order.TotalPayment,
		&ids.UserId,
		&ids.PharmacyAdminId,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, apperror.ErrRecordNotFound
		}
		return nil, nil, err
	}

	details, err := repo.findAllOrderDetailsByOrderId(ctx, order.Id)
	if err != nil {
		return nil, nil, err
	}

	order.OrderDetails = details
	order.Pharmacy = &pharmacy
	order.ShippingMethod = &shipping
	order.LatestStatus = &status
	return &order, &ids, err

}

func (repo *OrderRepositoryImpl) findAllOrderDetailsByOrderId(ctx context.Context, orderId int64) ([]*entity.OrderDetail, error) {
	const getAllOrderDetails = `SELECT order_details.id, quantity, name, generic_name, content, description, image, price FROM order_details
	INNER JOIN orders ON order_details.order_id = orders.id WHERE orders.id = $1`

	rows, err := repo.db.QueryContext(ctx, getAllOrderDetails, orderId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*entity.OrderDetail, 0)
	for rows.Next() {
		var orderDetail entity.OrderDetail
		if err := rows.Scan(
			&orderDetail.Id, &orderDetail.Quantity, &orderDetail.Name, &orderDetail.GenericName,
			&orderDetail.Content, &orderDetail.Description, &orderDetail.Image, &orderDetail.Price,
		); err != nil {
			return nil, err
		}
		items = append(items, &orderDetail)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (repo *OrderRepositoryImpl) UpdateOrderStatus(ctx context.Context, orderId int64, orderLog entity.OrderStatusLog) (*entity.OrderStatusLog, error) {
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	const updateOldStatus = `UPDATE order_status_logs SET is_latest = FALSE WHERE order_id = $1 AND is_latest = true`
	row1 := tx.QueryRowContext(ctx, updateOldStatus, orderId)
	if row1.Err() != nil {
		return nil, err
	}

	const addStatus = `INSERT INTO order_status_logs(order_id, order_status_id, is_latest, description)
	values ($1, $2, $3, $4) RETURNING id, order_id, order_status_id, is_latest, description`

	row2 := tx.QueryRowContext(ctx, addStatus, orderId, orderLog.OrderStatusId, orderLog.IsLatest, orderLog.Description)
	var createdStatus entity.OrderStatusLog
	err = row2.Scan(
		&createdStatus.Id,
		&createdStatus.OrderId,
		&createdStatus.OrderStatusId,
		&createdStatus.IsLatest,
		&createdStatus.Description,
	)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return &createdStatus, nil

}

func (repo *OrderRepositoryImpl) FindLatestOrderStatusByOrderId(ctx context.Context, id int64) (*entity.OrderStatusLog, error) {
	const findLatest = `SELECT DISTINCT id, order_id, order_status_id, is_latest, description FROM order_status_logs WHERE order_id = $1 AND is_latest = TRUE`
	row := repo.db.QueryRowContext(ctx, findLatest, id)
	var orderStatus entity.OrderStatusLog
	err := row.Scan(
		&orderStatus.Id,
		&orderStatus.OrderId,
		&orderStatus.OrderStatusId,
		&orderStatus.IsLatest,
		&orderStatus.Description,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrRecordNotFound
		}
		return nil, err
	}

	return &orderStatus, nil
}

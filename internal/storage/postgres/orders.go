package postgres

import (
	"context"
	"errors"

	dto "github.com/Eanhain/gofermart/internal/api"
	domain "github.com/Eanhain/gofermart/internal/domain"
	sq "github.com/Masterminds/squirrel"
	pgx "github.com/jackc/pgx/v5"
)

func (ps *PersistStorage) InsertNewUserOrder(ctx context.Context, order string, userID int, status string, accrual float64) error {
	sql, args, err := ps.psql.
		Insert("orders").
		Columns("id", "user_id", "status", "accural").
		Values(order, userID, status, accrual).
		ToSql()
	if err != nil {
		return err
	}

	tag, err := ps.Exec(ctx, sql, args...)
	if err != nil {
		ps.log.Warnln("Can't insert user order", err)
		return err
	}
	ps.log.Infoln("Insert user order", tag)
	return nil
}

func (ps *PersistStorage) InsertOrderWithdrawn(ctx context.Context, userID int, order dto.Withdrawn) error {
	sql, args, err := ps.psql.
		Insert("withdraw_orders").
		Columns("user_id", "order_id", "sum").
		Values(userID, order.Order, order.Sum).
		ToSql()
	if err != nil {
		return err
	}

	tag, err := ps.Exec(ctx, sql, args...)
	if err != nil {
		ps.log.Warnln("Can't insert user withdrawn order", err)
		return err
	}
	ps.log.Infoln("Insert user user withdrawn order", tag)
	return nil
}

func (ps *PersistStorage) CheckUserOrderNonExist(ctx context.Context, userID int, orderID string) error {
	var orderOut string

	sql, args, err := ps.psql.
		Select("id").
		From("orders").
		Where(sq.Eq{"user_id": userID, "id": orderID}).
		ToSql()
	if err != nil {
		return err
	}

	err = ps.QueryRow(ctx, sql, args...).Scan(&orderOut)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return err
	}

	return domain.ErrOrderExist
}

func (ps *PersistStorage) CheckOrderNonExist(ctx context.Context, orderID string) error {
	var orderOut string

	sql, args, err := ps.psql.
		Select("id").
		From("orders").
		Where(sq.Eq{"id": orderID}).
		ToSql()
	if err != nil {
		return err
	}

	err = ps.QueryRow(ctx, sql, args...).Scan(&orderOut)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil // Заказа нет, всё ок
		}
		return err
	}

	// Заказ существует, но принадлежит (возможно) другому юзеру
	return domain.ErrOrderExistWrongUser
}

func (ps *PersistStorage) GetUserOrders(ctx context.Context, userID int) (dto.OrdersDesc, error) {
	sql, args, err := ps.psql.
		Select("id", "status", "accural", "uploaded_at").
		From("orders").
		Where(sq.Eq{"user_id": userID}).
		ToSql()
	if err != nil {
		return dto.OrdersDesc{}, err
	}

	rows, err := ps.Query(ctx, sql, args...)
	if err != nil {
		ps.log.Warnln("Can't get user order", err)
		return dto.OrdersDesc{}, err
	}
	defer rows.Close()

	var orders dto.OrdersDesc
	for rows.Next() {
		var order dto.OrderDesc
		if err := rows.Scan(&order.Number, &order.Status, &order.Accrual, &order.UploadedAt); err != nil {
			ps.log.Warnln("can't parse order", err)
			return dto.OrdersDesc{}, err
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return dto.OrdersDesc{}, err
	}

	ps.log.Infoln("Get user orders", orders)
	return orders, nil
}

func (ps *PersistStorage) GetUserOrdersWithdrawn(ctx context.Context, userID int) (dto.Withdrawns, error) {
	sql, args, err := ps.psql.
		Select("order_id", "sum", "uploaded_at").
		From("withdraw_orders").
		Where(sq.Eq{"user_id": userID}).
		OrderBy("uploaded_at DESC").
		ToSql()
	if err != nil {
		return dto.Withdrawns{}, err
	}

	rows, err := ps.Query(ctx, sql, args...)
	if err != nil {
		ps.log.Warnln("Can't get user order", err)
		return dto.Withdrawns{}, err
	}
	defer rows.Close()

	var orders dto.Withdrawns
	for rows.Next() {
		var order dto.Withdrawn
		if err := rows.Scan(&order.Order, &order.Sum, &order.ProcessedAt); err != nil {
			ps.log.Warnln("can't parse order Withdrawn", err)
			return dto.Withdrawns{}, err
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return dto.Withdrawns{}, err
	}

	ps.log.Infoln("Get user orders", orders)
	return orders, nil
}

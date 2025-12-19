package postgres

import (
	"context"

	dto "github.com/Eanhain/gofermart/internal/api"
	domain "github.com/Eanhain/gofermart/internal/domain"
)

var (
	selectUserOrder = `
		select ID as number 
		from orders 
		where user_id = $1 and id = $2;
	`
	selectOrder = `
		select ID as number 
		from orders 
		where id = $1;
	`
	selectWithdrawnUserOrders = `
		select order_id as order, sum, UPLOADED_AT as processed_at from withdraw_orders
		where user_id = $1
		order by processed_at desc
	`
	selectUserOrders = `
		select ID as number, status, accural accrual, uploaded_at 
		from orders 
		where user_id = $1;
	`
	InsertOrder = DMLUserStruct{
		Name: "Insert user order",
		DML: `INSERT INTO orders (ID, USER_ID, STATUS, ACCURAL)
		VALUES ($1, $2, $3, $4)`,
	}
	InsertOrderWithdrawn = DMLUserStruct{
		Name: "Insert user withdrawn order",
		DML: `INSERT INTO withdraw_orders (user_id, order_id, sum)
		VALUES ($1, $2, $3)`,
	}
)

func (ps *PersistStorage) InsertNewUserOrder(ctx context.Context, order string, userID int, status string, accrual float64) error {
	tag, err := ps.Exec(ctx, InsertOrder.DML, order, userID, status, accrual)
	if err != nil {
		ps.log.Warnln("Can't insert user order", err)
		return err
	}
	ps.log.Infoln("Insert user order", tag)
	return nil

}

func (ps *PersistStorage) InsertOrderWithdrawn(ctx context.Context, userID int, orderID string, sum float64) error {
	tag, err := ps.Exec(ctx, InsertOrderWithdrawn.DML, userID, orderID, sum)
	if err != nil {
		ps.log.Warnln("Can't insert user withdrawn order", err)
		return err
	}
	ps.log.Infoln("Insert user user withdrawn order", tag)
	return nil

}

func (ps *PersistStorage) CheckUserOrderNonExist(ctx context.Context, userID int, orders string) error {
	var orderOut string
	out := ps.QueryRow(ctx, selectUserOrder, userID, orders)
	out.Scan(&orderOut)
	if orderOut != "" {
		return domain.ErrOrderExist
	}
	return nil
}

func (ps *PersistStorage) CheckOrderNonExist(ctx context.Context, orders string) error {
	var orderOut string
	out := ps.QueryRow(ctx, selectOrder, orders)
	out.Scan(&orderOut)
	if orderOut != "" {
		return domain.ErrOrderExistWrongUser
	}
	return nil
}

func (ps *PersistStorage) GetUserOrders(ctx context.Context, userID int) (dto.OrdersDesc, error) {
	var orders dto.OrdersDesc
	var order dto.OrderDesc
	rows, err := ps.Query(ctx, selectUserOrders, userID)
	if err != nil {
		ps.log.Warnln("Can't get user order", err)
		return dto.OrdersDesc{}, err
	}

	for rows.Next() {
		if err := rows.Scan(&order.Number, &order.Status, &order.Accrual, &order.UploadedAt); err != nil {
			ps.log.Warnln("can't parse order", err)
			return dto.OrdersDesc{}, err
		}
		orders = append(orders, order)
	}

	ps.log.Infoln("Get user orders", orders)
	return orders, nil

}

func (ps *PersistStorage) GetUserOrdersWithdrawn(ctx context.Context, userID int) (dto.Withdrawns, error) {
	var orders dto.Withdrawns
	var order dto.Withdrawn
	rows, err := ps.Query(ctx, selectWithdrawnUserOrders, userID)
	if err != nil {
		ps.log.Warnln("Can't get user order", err)
		return dto.Withdrawns{}, err
	}

	for rows.Next() {
		if err := rows.Scan(&order.Order, &order.Sum, &order.ProcessedAt); err != nil {
			ps.log.Warnln("can't parse order Withdrawn", err)
			return dto.Withdrawns{}, err
		}
		orders = append(orders, order)
	}

	ps.log.Infoln("Get user orders", orders)
	return orders, nil
}

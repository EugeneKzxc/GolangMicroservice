package main

import (
	"context"

	"encoding/json"

	"github.com/jackc/pgx/v4"
)

const connStr = "postgres://L0_user:zxc@pg_container:5432/L0_database"

func insertOrderToDB(conn *pgx.Conn, order Order) error {
	_, err := conn.Exec(context.Background(), `
	INSERT INTO orders(order_uid, track_number, entry, delivery, payment, items, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`,
		order.OrderUID, order.TrackNumber, order.Entry, order.Delivery, order.Payment,
		order.Items, order.Locale, order.InternalSignature, order.CustomerID, order.DeliveryService,
		order.ShardKey, order.SmID, order.DateCreated, order.OofShard)
	return err
}

func getOrdersFromDB(conn *pgx.Conn) (map[string]Order, error) {
	rows, err := conn.Query(context.Background(), "SELECT * FROM orders")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orderMap := make(map[string]Order)

	for rows.Next() {
		var order Order
		var deliveryJSON []byte
		var paymentJSON []byte
		var itemsJSON []byte

		if err := rows.Scan(
			&order.OrderUID, &order.TrackNumber, &order.Entry,
			&deliveryJSON, &paymentJSON, &itemsJSON,
			&order.Locale, &order.InternalSignature, &order.CustomerID,
			&order.DeliveryService, &order.ShardKey, &order.SmID,
			&order.DateCreated, &order.OofShard,
		); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(deliveryJSON, &order.Delivery); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(paymentJSON, &order.Payment); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(itemsJSON, &order.Items); err != nil {
			return nil, err
		}
		orderMap[order.OrderUID] = order
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return orderMap, nil
}

func connectToDB() (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), connStr)
	return conn, err
}

func CheckRecordExists(conn *pgx.Conn, id string) (bool, error) {
	var exists bool

	err := conn.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM orders WHERE order_uid=$1)", id).Scan(&exists)
	return exists, err
}

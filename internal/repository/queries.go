package repository

const (
	scheme = `
		CREATE TABLE IF NOT EXISTS delivery
		(
			delivery_id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			phone TEXT NOT NULL,
			zip TEXT NOT NULL,
			city TEXT NOT NULL,
			address TEXT NOT NULL,
			region TEXT NOT NULL,
			email TEXT NOT NULL
		);

		CREATE TABLE IF NOT EXISTS payments
		(
		    payment_id TEXT PRIMARY KEY ,
			transaction TEXT NOT NULL,
			request_id TEXT NOT NULL,
			currency TEXT NOT NULL,
			provider TEXT NOT NULL,
			amount BIGINT NOT NULL,
			payment_dt BIGINT NOT NULL,
			bank TEXT NOT NULL,
			delivery_cost BIGINT NOT NULL,
			goods_total BIGINT NOT NULL,
			custom_fee BIGINT NOT NULL
		);

		CREATE TABLE IF NOT EXISTS orders
		(
			order_uid TEXT PRIMARY KEY,
			track_number TEXT NOT NULL,
			entry TEXT NOT NULL,
			delivery_id TEXT NOT NULL,
			payment_id TEXT NOT NULL,
			locate TEXT NOT NULL,
			internal_signature TEXT,
			customer_id TEXT NOT NULL,
			delivery_service TEXT NOT NULL,
			shard_key TEXT NOT NULL,
			sm_id TEXT NOT NULL,
			date_created TIMESTAMP,
			oof_shard TEXT NOT NULL,
			
			FOREIGN KEY (delivery_id) REFERENCES delivery(delivery_id) ON DELETE CASCADE,
			FOREIGN KEY (payment_id) REFERENCES payments(payment_id) ON DELETE CASCADE
		);

		CREATE TABLE IF NOT EXISTS items
		(
			order_id TEXT NOT NULL,
			chrt_id BIGINT NOT NULL,
			track_number TEXT NOT NULL,
			price BIGINT NOT NULL,
			rid TEXT NOT NULL,
			name TEXT NOT NULL,
			sale BIGINT NOT NULL,
			size TEXT NOT NULL,
			total_price BIGINT NOT NULL,
			nm_id BIGINT NOT NULL,
			brand TEXT NOT NULL,
			status BIGINT NOT NULL,
			
			FOREIGN KEY (order_id) REFERENCES orders(order_uid) ON DELETE CASCADE 
		);
	`

	saveOrder = `
		INSERT INTO orders
			(order_uid, track_number, entry, delivery_id, payment_id, locate, internal_signature, customer_id, delivery_service, shard_key, sm_id, date_created, oof_shard)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13);
	`

	saveDelivery = `
		INSERT INTO delivery
			(delivery_id, name, phone, zip, city, address, region, email)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8);
	`

	savePayments = `
		INSERT INTO payments
			(payment_id, transaction, request_id, currency, provider, amount,
			 payment_dt, bank, delivery_cost, goods_total, custom_fee)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);
	`

	saveItems = `
		INSERT INTO items
			(order_id, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);
	`

	getAllOrders = `
		SELECT * FROM orders;
	`

	getDeliveryByID = `
		SELECT * FROM delivery
		WHERE delivery_id = $1;	
	`

	getPaymentsByID = `
		SELECT * FROM payments
		WHERE payment_id = $1;
	`

	getItemsByOrderID = `
		SELECT 
			order_id, 
			chrt_id, 
			track_number, 
			price, 
			rid, 
			name, 
			sale, 
			size, 
			total_price, 
			nm_id, 
			brand, 
			status 
		FROM items 
		WHERE order_id = $1;
	`
)

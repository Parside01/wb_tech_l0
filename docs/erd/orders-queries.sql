CREATE TABLE IF NOT EXISTS orders
(
    order_uid TEXT NOT NULL PRIMARY KEY,
    track_number TEXT NOT NULL,
    entry TEXT NOT NULL,
    delivery_id TEXT NOT NULL,
    payment_id TEXT NOT NULL,
    locate TEXT NOT NULL,
    internal_signature TEXT,
    customer_id TEXT NOT NULL,
    delivery_service TEXT NOT NULL,
    shardkey TEXT NOT NULL,
    sm_id TEXT NOT NULL,
    date_created TIMESTAMP,
    oof_shard TEXT NOT NULL,

    CONSTRAINT delivery_id_fk FOREIGN KEY (delivery_id) REFERENCES Delivery (delivery_id)
);

CREATE TABLE IF NOT EXISTS Delivery
(
    delivery_id TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    phone TEXT NOT NULL,
    zip TEXT NOT NULL,
    city TEXT NOT NULL,
    address TEXT NOT NULL,
    region TEXT NOT NULL,
    email TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS items
(
    order_id TEXT NOT NULL,
    chrt_id BIGINT NOT NULL,
    track_number TEXT NOT NULL,
    price BIGINT NOT NULL ,
    rid TEXT NOT NULL,
    name TEXT NOT NULL,
    sale BIGINT NOT NULL,
    size TEXT NOT NULL,
    total_price BIGINT NOT NULL,
    nm_id BIGINT NOT NULL,
    brand TEXT NOT NULL,
    status BIGINT NOT NULL,

    CONSTRAINT order_id_fk FOREIGN KEY (order_id) REFERENCES orders (order_uid)
);

CREATE TABLE IF NOT EXISTS payments
(
    payment_id TEXT NOT NULL PRIMARY KEY,
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


-- Save order
INSERT INTO orders
    (order_uid, track_number, entry, delivery_id, payment_id, locate, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
VALUES
    (:order_uid, :track_number, :entry, :delivery_id, :payment_id, :locate, :internal_signature, :customer_id, :delivery_service, :shardkey, :sm_id, :date_created, :oof_shard);


-- Save Delivery
INSERT INTO Delivery
    (delivery_id, name, phone, zip, city, address, region, email)
VALUES
    (:delivery_id, :name, :phone, :zip, :city, :address, :region, :email);

-- Save payments
INSERT INTO payments
    (payment_id, transaction, request_id, currency, provider, amount,
     payment_dt, bank, delivery_cost, goods_total, custom_fee)
VALUES
    (:payment_id, :transaction, :request_id, :currency, :provider, :amount,
     :payment_dt, :bank, :delivery_cost, :goods_total, :custom_fee);

-- Save items
INSERT INTO items
    (order_id, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
VALUES
    (:order_id, :chrt_id, :track_number, :price, :rid, :name, :sale, :size, :total_price, :nm_id, :brand, :status);

-- Get all from orders
SELECT * FROM orders;

-- Get from Delivery by delivery_id
SELECT * FROM Delivery
WHERE delivery_id = $1;

-- Get from payments by payment_id
SELECT * FROM payments
WHERE payment_id = $1;

-- Get from items by order_id
SELECT * FROM items
WHERE order_id = $1;
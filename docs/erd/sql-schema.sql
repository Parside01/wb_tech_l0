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
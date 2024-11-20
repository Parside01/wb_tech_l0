-- Save order
INSERT INTO orders
    (order_uid, track_number, entry, delivery_id, payment_id, locate, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
VALUES
    (:order_uid, :track_number, :entry, :delivery_id, :payment_id, :locate, :internal_signature, :customer_id, :delivery_service, :shardkey, :sm_id, :date_created, :oof_shard)


-- Save delivery
INSERT INTO delivery
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
    (:order_id, :chrt_id, :track_number, :price, :rid, :name, :sale, :size, :total_price, :nm_id, :brand, :status)
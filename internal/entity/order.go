package entity

import "encoding/json"

type Order struct {
	OrderUID          string    `json:"order_uid" db:"order_uid"`
	TrackNumber       string    `json:"track_number" db:"track_number"`
	Entry             string    `json:"entry" db:"entry"`
	Delivery          *Delivery `json:"Delivery" db:"Delivery"`
	Payment           *Payment  `json:"Payment" db:"Payment"`
	Items             []*Item   `json:"items" db:"items"`
	Locate            string    `json:"locate" db:"locate"`
	InternalSignature string    `json:"internal_signature" db:"internal_signature"`
	CustomerID        string    `json:"customer_id" db:"customer_id"`
	DeliveryService   string    `json:"delivery_service" db:"delivery_service"`
	ShardKey          string    `json:"shard_key" db:"shard_key"`
	SMID              int       `json:"sm_id" db:"sm_id"`
	DateCreated       string    `json:"date_created" db:"date_created"`
	OofShard          string    `json:"oof_shard" db:"oof_shard"`
}

func (o *Order) Key() string {
	return o.OrderUID
}

type Delivery struct {
	DeliveryID string `json:"delivery_id" db:"delivery_id"`
	Name       string `json:"name" db:"name"`
	Phone      string `json:"phone" db:"phone"`
	Zip        string `json:"zip" db:"zip"`
	City       string `json:"city" db:"city"`
	Address    string `json:"address" db:"address"`
	Region     string `json:"region" db:"region"`
	Email      string `json:"email" db:"email"`
}

type Payment struct {
	PaymentID    string `json:"payment_id" db:"payment_id"`
	Transaction  string `json:"transaction" db:"transaction"`
	RequestID    string `json:"request_id" db:"request_id"`
	Currency     string `json:"currency" db:"currency"`
	Provider     string `json:"provider" db:"provider"`
	Amount       int    `json:"amount" db:"amount"`
	PaymentDT    int    `json:"payment_dt" db:"payment_dt"`
	Bank         string `json:"bank" db:"bank"`
	DeliveryCost int    `json:"delivery_cost" db:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total" db:"goods_total"`
	CustomFee    int    `json:"custom_fee" db:"custom_fee"`
}

type Item struct {
	OrderID     string `json:"order_id" db:"order_id"`
	ChrtID      int    `json:"chrt_id" db:"chrt_id"`
	TrackNumber string `json:"track_number" db:"track_number"`
	Price       int    `json:"price" db:"price"`
	RID         string `json:"rid" db:"rid"`
	Name        string `json:"name" db:"name"`
	Sale        int    `json:"sale" db:"sale"`
	Size        string `json:"size" db:"size"`
	TotalPrice  int    `json:"total_price" db:"total_price"`
	NMID        int    `json:"nm_id" db:"nm_id"`
	Brand       string `json:"brand" db:"brand"`
	Status      int    `json:"status" db:"status"`
}

func MarshalOrder(order Order) ([]byte, error) {
	return json.Marshal(order)
}

func UnmarshalOrder(data []byte) (Order, error) {
	var order Order
	err := json.Unmarshal(data, &order)
	return order, err
}

type SQLOrder struct {
	OrderUID          string `json:"order_uid" db:"order_uid"`
	TrackNumber       string `json:"track_number" db:"track_number"`
	Entry             string `json:"entry" db:"entry"`
	DeliveryID        string `json:"delivery_id" db:"delivery_id"`
	PaymentID         string `json:"payment_id" db:"payment_id"`
	Locate            string `json:"locate" db:"locate"`
	InternalSignature string `json:"internal_signature" db:"internal_signature"`
	CustomerID        string `json:"customer_id" db:"customer_id"`
	DeliveryService   string `json:"delivery_service" db:"delivery_service"`
	ShardKey          string `json:"shard_key" db:"shard_key"`
	SMID              int    `json:"sm_id" db:"sm_id"`
	DateCreated       string `json:"date_created" db:"date_created"`
	OofShard          string `json:"oof_shard" db:"oof_shard"`
}

func (s *SQLOrder) CopyToOrder(order *Order) {
	order.OrderUID = s.OrderUID
	order.TrackNumber = s.TrackNumber
	order.Entry = s.Entry
	order.SMID = s.SMID
	order.DateCreated = s.DateCreated
	order.OofShard = s.OofShard
	order.Locate = s.Locate
	order.InternalSignature = s.InternalSignature
	order.CustomerID = s.CustomerID
	order.ShardKey = s.ShardKey
	order.DeliveryService = s.DeliveryService
}

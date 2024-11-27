package entity

import "github.com/brianvoe/gofakeit/v7"

func GenerateRandomOrder() *Order {
	return &Order{
		OrderUID:          gofakeit.UUID(),
		TrackNumber:       gofakeit.UUID(),
		Entry:             gofakeit.Word(),
		Delivery:          GenerateRandomDelivery(),
		Payment:           GenerateRandomPayment(),
		Items:             GenerateRandomItems(),
		Locate:            gofakeit.City(),
		InternalSignature: gofakeit.Word(),
		CustomerID:        gofakeit.UUID(),
		DeliveryService:   gofakeit.Word(),
		ShardKey:          gofakeit.Word(),
		SMID:              gofakeit.Number(1, 100),
		DateCreated:       gofakeit.Date().Format("2006-01-02 15:04:05"),
		OofShard:          gofakeit.Word(),
	}
}

func GenerateRandomDelivery() *Delivery {
	return &Delivery{
		OrderID: gofakeit.UUID(),
		Name:    gofakeit.Name(),
		Phone:   gofakeit.Phone(),
		Zip:     gofakeit.Zip(),
		City:    gofakeit.City(),
		Address: gofakeit.Address().Address,
		Region:  gofakeit.State(),
		Email:   gofakeit.Email(),
	}
}

func GenerateRandomPayment() *Payment {
	return &Payment{
		OrderID:      gofakeit.UUID(),
		Transaction:  gofakeit.UUID(),
		RequestID:    gofakeit.UUID(),
		Currency:     gofakeit.Currency().Short,
		Provider:     gofakeit.Word(),
		Amount:       gofakeit.Number(1000, 10000),
		PaymentDT:    int(gofakeit.Date().Unix()),
		Bank:         gofakeit.Word(),
		DeliveryCost: gofakeit.Number(100, 500),
		GoodsTotal:   gofakeit.Number(500, 5000),
		CustomFee:    gofakeit.Number(0, 100),
	}
}

func GenerateRandomItems() []*Item {
	var items []*Item
	itemCount := gofakeit.Number(1, 5)
	for i := 0; i < itemCount; i++ {
		items = append(items, &Item{
			OrderID:     gofakeit.UUID(),
			ChrtID:      gofakeit.Number(1, 10000),
			TrackNumber: gofakeit.UUID(),
			Price:       gofakeit.Number(100, 1000),
			RID:         gofakeit.UUID(),
			Name:        gofakeit.Word(),
			Sale:        gofakeit.Number(0, 100),
			Size:        gofakeit.Word(),
			TotalPrice:  gofakeit.Number(100, 1000),
			NMID:        gofakeit.Number(1, 1000),
			Brand:       gofakeit.Word(),
			Status:      gofakeit.Number(1, 5),
		})
	}
	return items
}

func GenerateRandomItem() *Item {
	return &Item{
		OrderID:     gofakeit.UUID(),
		ChrtID:      gofakeit.Number(1, 10000),
		TrackNumber: gofakeit.UUID(),
		Price:       gofakeit.Number(100, 1000),
		RID:         gofakeit.UUID(),
		Name:        gofakeit.Word(),
		Sale:        gofakeit.Number(0, 100),
		Size:        gofakeit.Word(),
		TotalPrice:  gofakeit.Number(100, 1000),
		NMID:        gofakeit.Number(1, 1000),
		Brand:       gofakeit.Word(),
		Status:      gofakeit.Number(1, 5),
	}
}

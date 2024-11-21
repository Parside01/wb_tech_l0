package service

import "fmt"

var (
	ErrNoOrderInCache = fmt.Errorf("Order is not in cache")
)

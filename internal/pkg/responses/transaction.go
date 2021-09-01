package responses

import "time"

type Transaction struct {
	ID         int32
	UserName   string
	UserEmail  string
	Date       time.Time
	Detail     []Detail
	GrandTotal int32
}

type Detail struct {
	ProductID   int32
	ProductName string
	Price       int32
	Qty         int32
	SubTotal    int32
}

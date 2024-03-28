package models

type Sales struct {
	ProductID   string `json:"productID"`
	Quantities  string `json:"quantities"`
	TotalAmount int    `json:"totalamount"`
}

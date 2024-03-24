package models

import "time"

type Order struct {
	ID             int       `json:"Id" gorm:"type:INT(10) UNSIGNED NOT NULL AUTO_INCREMENT;primaryKey"`
	CashierID      int       `json:"cashierID"`
	PaymentTypesID int       `json:"payment_types_id"`
	TotalPrice     int       `json:"totalprice"`
	TotalPaid      int       `json:"totalpaid"`
	TotalReturn    int       `json:"totalreturn"`
	ReceiptID      int       `json:"receipt_id"`
	IsDownload     int       `json:"is_download"`
	ProductID      string    `json:"product_id"`
	Quantities     string    `json:"quantities"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

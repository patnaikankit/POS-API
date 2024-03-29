package models

import (
	"time"
)

type Order struct {
	ID             int       `json:"Id" gorm:"type:INT(10) UNSIGNED NOT NULL AUTO_INCREMENT;primaryKey"`
	CashierID      int       `json:"cashierID"`
	PaymentTypesID int       `json:"payment_types_id"`
	TotalPrice     int       `json:"totalprice"`
	TotalPaid      int       `json:"totalpaid"`
	TotalReturn    int       `json:"totalreturn"`
	ReceiptID      string    `json:"receipt_id"`
	IsDownload     int       `json:"is_download"`
	ProductID      string    `json:"product_id"`
	Quantities     string    `json:"quantities"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type ProductOrder struct {
	ID         int    `json:"Id" gorm:"type:INT(10) UNSIGNED NOT NULL AUTO_INCREMENT;primaryKey"`
	Sku        string `json:"sku"`
	Name       string `json:"name"`
	Stock      int    `json:"stock"`
	Price      int    `json:"price"`
	Image      string `json:"image"`
	CategoryID int    `json:"categoryID"`
	DiscountID int    `json:"discountID"`
}

type SalesResponse struct {
	ProductID   int    `json:"productID"`
	Name        string `json:"name"`
	TotalQty    int    `json:"totalQty"`
	TotalAmount int    `json:"totalAmount"`
}

type RevenueResponse struct {
	PaymentTypeID int    `json:"paymentTypeID"`
	Name          string `json:"name"`
	Logo          string `json:"logo"`
	TotalAmount   int    `json:"totalAmount"`
}

type ProductResponseOrder struct {
	ProductID        int      `json:"productId" gorm:"type:INT(10) UNSIGNED NOT NULL AUTO_INCREMENT;primaryKey"`
	Name             string   `json:"name"`
	Price            int      `json:"price"`
	Qty              int      `json:"qty"`
	Discount         Discount `json:"discount"`
	TotalNormalPrice int      `json:"totalNormalPrice"`
	TotalFinalPrice  int      `json:"totalFinalPrice"`
}

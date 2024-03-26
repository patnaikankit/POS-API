package models

import "time"

type Product struct {
	ID               int       `json:"id" gorm:"type:INT(10) UNSIGNED NOT NULL AUTO_INCREMENT;primaryKey"`
	Sku              string    `json:"sku"`
	Name             string    `json:"name"`
	Stock            int       `json:"stock"`
	Price            int       `json:"price"`
	Image            string    `json:"image"`
	TotalFinalPrice  int       `json:"total_final_price"`
	TotalNormalPrice int       `json:"total_normal_price"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	CategoryID       int       `json:"categoryId"`
	DiscountID       int       `json:"discountId"`
}

type ProductResult struct {
	ID       int      `json:"productId" gorm:"type:INT(10) UNSIGNED NOT NULL AUTO_INCREMENT;primaryKey"`
	Sku      string   `json:"sku"`
	Name     string   `json:"name"`
	Stock    int      `json:"stock"`
	Price    int      `json:"price"`
	Image    string   `json:"image"`
	Category Category `json:"category"`
	Discount Discount `json:"discount"`
}

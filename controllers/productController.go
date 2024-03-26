package controllers

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	db "github.com/patnaikankit/POS-API/config"
	"github.com/patnaikankit/POS-API/middlewares"
	"github.com/patnaikankit/POS-API/models"
)

type Products struct {
	Products     models.Product
	CategoriesId string `json:"categories_Id"`
}
type ProdDiscount struct {
	Id         int      `json:"id" gorm:"type:INT(10) UNSIGNED NOT NULL AUTO_INCREMENT;primaryKey"`
	Sku        string   `json:"sku"`
	Name       string   `json:"name"`
	Stock      int      `json:"stock"`
	Price      int      `json:"price"`
	Image      string   `json:"image"`
	CategoryID int      `json:"categoryID"`
	Discount   Discount `json:"discount"`
}
type Discount struct {
	Qty       int    `json:"qty"`
	Types     string `json:"type"`
	Result    int    `json:"result"`
	ExpiredAt int    `json:"expiredAt"`
}

// get all the products available
func GetAllProducts(c *fiber.Ctx) error {
	// check if a token is present or not
	headerToken := c.Get("Authorization")
	if headerToken == "" {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized Cashier!",
			"error":   map[string]interface{}{},
		})
	}

	// check if the token is valid
	if err := middlewares.AuthenticateToken(middlewares.SplitToken(headerToken)); err != nil {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized Cashier!",
			"error":   map[string]interface{}{},
		})
	}

	// product retrieval
	limit := c.Query("limit")
	skip := c.Query("skip")
	categoryID := c.Query("categoryId")
	productName := c.Query("q")
	intLimit, _ := strconv.Atoi(limit)
	intSkip, _ := strconv.Atoi(skip)
	var product []models.Product

	productArray := make([]*models.ProductResult, 0)

	if productName == "" {
		var count int64
		db.DB.Where("category_ID=?", categoryID).Limit(intLimit).Offset(intSkip).Find(&product).Count(&count)

		var category models.Category
		var discount models.Discount
		for i := 0; i < len(product); i++ {
			db.DB.Table("categories").Where("id=?", product[i].CategoryID).Find(&category)
			db.DB.Where("id=?", product[i].DiscountID).Limit(intLimit).Offset(intSkip)
			count = int64(len(product))

			productArray = append(productArray, &models.ProductResult{
				ID:       product[i].ID,
				Sku:      product[i].Sku,
				Name:     product[i].Name,
				Stock:    product[i].Stock,
				Price:    product[i].Price,
				Image:    product[i].Image,
				Category: category,
				Discount: discount,
			})
		}
		meta := map[string]interface{}{
			"total": count,
			"limit": limit,
			"skip":  skip,
		}

		return c.Status(200).JSON(fiber.Map{
			"success": true,
			"message": "Products list retrived successfully!",
			"data": map[string]interface{}{
				"products": productArray,
				"meta":     meta,
			},
		})
	} else {
		var count int64
		if categoryID != "" {
			db.DB.Where("category_ID=? AND name=?", categoryID, productName).Limit(intLimit).Offset(intSkip).Find(&product).Count(&count)
		} else {
			db.DB.Where("name=?", productName).Limit(intLimit).Offset(intSkip).Find(&product).Count(&count)
		}

		var category models.Category
		var discount models.Discount
		for i := 0; i < len(product); i++ {
			db.DB.Where("id=?", product[i].CategoryID).Find(&category)
			db.DB.Where("id=?", product[i].DiscountID).Limit(intLimit).Offset(intSkip).Find(&discount).Count(&count)
			count = int64(len(product))
			productArray = append(productArray,
				&models.ProductResult{
					ID:       product[i].ID,
					Sku:      product[i].Sku,
					Name:     product[i].Name,
					Stock:    product[i].Stock,
					Price:    product[i].Price,
					Image:    product[i].Image,
					Category: category,
					Discount: discount,
				},
			)
		}
		meta := map[string]interface{}{
			"total": count,
			"limit": limit,
			"skip":  skip,
		}

		return c.Status(200).JSON(fiber.Map{
			"success": true,
			"message": "Products list retrived successfully!",
			"data": map[string]interface{}{
				"products": productArray,
				"meta":     meta,
			},
		})
	}
}

// add a new product
func AddProduct(c *fiber.Ctx) error {
	var data ProdDiscount
	err := c.BodyParser(&data)
	if err != nil {
		log.Fatal("Error while creating new product -> ", err)
	}

	var p []models.Product
	db.DB.Find(&p)

	discount := models.Discount{
		Qty:       data.Discount.Qty,
		Type:      data.Discount.Types,
		Result:    data.Discount.Result,
		ExpiredAt: data.Discount.ExpiredAt,
	}
	db.DB.Create(&discount)

	product := models.Product{
		Name:       data.Name,
		Image:      data.Image,
		CategoryID: data.CategoryID,
		DiscountID: data.Id,
		Price:      data.Price,
		Stock:      data.Stock,
	}
	db.DB.Create(&product)

	db.DB.Table("products").Where("id=?", product.ID).Update("sku", "SKU"+strconv.Itoa(product.ID))

	fmt.Println("--------------------------------------->")
	fmt.Println("------------Product Addition Done----------->", product.ID)
	fmt.Println("--------------------------------------->")

	response := map[string]interface{}{
		"success": true,
		"message": "New Product Added Successfully!",
		"data":    product,
	}
	return (c.JSON(response))
}

// get data regarding a particular product
func GetProductDetails(c *fiber.Ctx) error {
	return nil
}

// delete a product
func DeleteProduct(c *fiber.Ctx) error {
	productID := c.Params("productID")
	var product models.Product

	db.DB.First(&product, productID)
	if product.ID == 0 {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Product Not Found!",
			"error":   map[string]interface{}{},
		})
	}

	db.DB.Delete(&product)
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Product Deleted Successfully!",
	})
}

// edit data of a product
func UpdateProduct(c *fiber.Ctx) error {
	productID := c.Params("productID")
	var product models.Product
	db.DB.Find(&product, "id=?", productID)

	if product.Name == "" {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Product Not Found!",
			"error":   map[string]interface{}{},
		})
	}

	var updateProductData models.Product
	c.BodyParser(&updateProductData)

	if updateProductData.Name == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Product Name is Required!",
			"error":   map[string]interface{}{},
		})
	}

	product.Name = updateProductData.Name
	product.CategoryID = updateProductData.CategoryID
	product.Image = updateProductData.Image
	product.Price = updateProductData.Price
	product.Stock = updateProductData.Stock

	db.DB.Save(&product)
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Data Updated Successfully!",
		"data":    product,
	})
}

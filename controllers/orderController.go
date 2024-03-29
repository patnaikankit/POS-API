package controllers

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	db "github.com/patnaikankit/POS-API/config"
	"github.com/patnaikankit/POS-API/middlewares"
	"github.com/patnaikankit/POS-API/models"

	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
)

func getGrayColor() color.Color {
	return color.Color{
		Red:   200,
		Green: 200,
		Blue:  200,
	}
}

func getHeader() []string {
	return []string{"Products Sku", "Name", "Qty", "Price"}
}

// retrieve all orders
func GetAllOrders(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit"))
	skip, _ := strconv.Atoi(c.Query("skip"))
	var count int64
	var order []models.Order

	db.DB.Select("*").Limit(limit).Offset(skip).Find(&order).Count(&count)

	type OrderList struct {
		OrderID        int                `json:"orderID"`
		CashierID      int                `json:"cashiersID"`
		PaymentTypesID int                `json:"paymentTypesID"`
		TotalPrice     int                `json:"totalPrice"`
		TotalPaid      int                `json:"totalPaid"`
		TotalReturn    int                `json:"totalReturn"`
		ReceiptID      string             `json:"receiptID"`
		CreatedAt      time.Time          `json:"createdAt"`
		Payments       models.PaymentType `json:"payment_type"`
		Cashiers       models.Cashier     `json:"cashier"`
	}

	orderResponse := make([]*OrderList, 0)

	for _, v := range order {
		cashier := models.Cashier{}
		db.DB.Where("id = ?", v.CashierID).Find(&cashier)
		paymentType := models.PaymentType{}
		db.DB.Where("id = ?", v.PaymentTypesID).Find(&paymentType)

		orderResponse = append(orderResponse, &OrderList{
			OrderID:        v.ID,
			CashierID:      v.CashierID,
			PaymentTypesID: v.PaymentTypesID,
			TotalPrice:     v.TotalPrice,
			TotalPaid:      v.TotalPaid,
			TotalReturn:    v.TotalReturn,
			ReceiptID:      v.ReceiptID,
			CreatedAt:      v.CreatedAt,
			Payments:       paymentType,
			Cashiers:       cashier,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Orders Fetched Successfully!",
		"data":    orderResponse,
		"meta": map[string]interface{}{
			"total": count,
			"limit": limit,
			"skip":  skip,
		},
	})
}

// creating a new order
func NewOrder(c *fiber.Ctx) error {
	// check if token is present
	headerToken := c.Get("Authorization")
	if headerToken == "" {
		c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Token Not Found!",
		})
	}

	// check if the cashier is authorized
	if err := middlewares.AuthenticateToken(middlewares.SplitToken(headerToken)); err != nil {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized Cashier!",
			"error":   map[string]interface{}{},
		})
	}

	type Products struct {
		ProductID int `json:"productID"`
		Quantity  int `json:"qty"`
	}

	body := struct {
		PaymentID int        `json:"paymentID"`
		TotalPaid int        `json:"totalPaid"`
		Products  []Products `json:"producs"`
	}{}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Empty Body!",
			"error":   map[string]interface{}{},
		})
	}

	repsonse := make([]*models.ProductResponseOrder, 0)

	var totalInvoicePrice = struct {
		price int
	}{}

	productIDs := ""
	quantities := ""

	for _, v := range body.Products {
		totalPrice := 0
		productIDs = productIDs + "," + strconv.Itoa(v.ProductID)
		quantities = quantities + "," + strconv.Itoa(v.Quantity)

		prod := models.ProductOrder{}
		var discount models.Discount
		db.DB.Table("products").Where("id=?", v.ProductID).First(&prod)
		db.DB.Where("id=?", prod.DiscountID).Find(&discount)
		discountCount := 0

		if discount.Type == "BUY_N" {
			totalPrice = prod.Price * v.Quantity
			percentage := totalPrice * discount.Result / 100
			discountCount = totalPrice - percentage
			totalInvoicePrice.price = totalInvoicePrice.price + discountCount
		}

		if discount.Type == "PERCENT" {
			totalPrice = prod.Price * v.Quantity
			percentage := totalPrice * discount.Result / 100
			discountCount = totalPrice - percentage
			totalInvoicePrice.price = totalInvoicePrice.price + discountCount
		}

		repsonse = append(repsonse, &models.ProductResponseOrder{
			ProductID:        prod.ID,
			Name:             prod.Name,
			Price:            prod.Price,
			Discount:         discount,
			Qty:              v.Quantity,
			TotalNormalPrice: prod.Price,
			TotalFinalPrice:  discountCount,
		})
	}

	// creating a new order in the db
	orderResponse := models.Order{
		CashierID:      1,
		PaymentTypesID: body.PaymentID,
		TotalPrice:     totalInvoicePrice.price,
		TotalPaid:      body.TotalPaid,
		TotalReturn:    body.TotalPaid - totalInvoicePrice.price,
		ReceiptID:      "R000" + strconv.Itoa(rand.Intn(1000)),
		ProductID:      productIDs,
		Quantities:     quantities,
		CreatedAt:      time.Now().Local().UTC(),
		UpdatedAt:      time.Now().Local().UTC(),
	}

	db.DB.Create(&orderResponse)

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": true,
		"data": map[string]interface{}{
			"order":    orderResponse,
			"products": repsonse,
		},
	})
}

// retrieve data about a particular order
func OrderData(c *fiber.Ctx) error {
	// check if token is present
	headerToken := c.Get("Authorization")
	if headerToken == "" {
		c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Token Not Found!",
		})
	}

	// check if the cashier is authorized
	if err := middlewares.AuthenticateToken(middlewares.SplitToken(headerToken)); err != nil {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized Cashier!",
			"error":   map[string]interface{}{},
		})
	}

	orderID := c.Params("orderID")
	var order models.Order
	db.DB.Where("id=?", orderID).First(&order)

	if order.ID == 0 {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Order Not Found!",
			"error":   map[string]interface{}{},
		})
	}

	productIDs := strings.Split(order.ProductID, ",")
	totalProducts := make([]*models.Product, 0)

	for i := 1; i < len(productIDs); i++ {
		prod := models.Product{}
		db.DB.Where("id=?", productIDs[i]).Find(&prod)
		totalProducts = append(totalProducts, &prod)
	}

	cashier := models.Cashier{}
	db.DB.Where("id=?", order.CashierID).Find(&cashier)

	paymentType := models.PaymentType{}
	db.DB.Where("id=?", order.PaymentTypesID).Find(&paymentType)

	orderTable := models.Order{}
	db.DB.Where("id=?", order.ID).Find(&orderTable)

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"data": map[string]interface{}{
			"order": map[string]interface{}{
				"orderID":       order.ID,
				"cashierID":     order.CashierID,
				"paymentTypeID": order.PaymentTypesID,
				"totalPrice":    order.TotalPrice,
				"totalPaid":     order.TotalPaid,
				"totalReturn":   order.TotalReturn,
				"receiptID":     order.ReceiptID,
				"createdAt":     order.ReceiptID,
				"payment_type":  paymentType,
			},
			"products": totalProducts,
		},
		"message": "Data Fetched Successfully!",
	})
}

// calculating the total order
func TotalOrder(c *fiber.Ctx) error {
	// check if token is present
	headerToken := c.Get("Authorization")
	if headerToken == "" {
		c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Token Not Found!",
		})
	}

	// check if the cashier is authorized
	if err := middlewares.AuthenticateToken(middlewares.SplitToken(headerToken)); err != nil {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized Cashier!",
			"error":   map[string]interface{}{},
		})
	}

	type Products struct {
		ProductID int `json:"productID"`
		Quantity  int `json:"qty"`
	}

	body := struct {
		Products []Products `json:"producs"`
	}{}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Empty Body!",
			"error":   map[string]interface{}{},
		})
	}

	response := make([]*models.ProductResponseOrder, 0)

	var totalInvoicePrice = struct {
		price int
	}{}

	for _, v := range body.Products {
		totalPrice := 0

		prod := models.ProductOrder{}
		var discount models.Discount
		db.DB.Table("products").Where("id=?", v.ProductID).First(&prod)
		db.DB.Where("id=?", prod.DiscountID).Find(&discount)
		discountCount := 0

		if discount.Type == "BUY_N" {
			totalPrice = prod.Price * v.Quantity
			percentage := totalPrice * discount.Result / 100
			discountCount = totalPrice - percentage
			totalInvoicePrice.price = totalInvoicePrice.price + discountCount
		}

		if discount.Type == "PERCENT" {
			totalPrice = prod.Price * v.Quantity
			percentage := totalPrice * discount.Result / 100
			discountCount = totalPrice - percentage
			totalInvoicePrice.price = totalInvoicePrice.price + discountCount
		}

		response = append(response, &models.ProductResponseOrder{
			ProductID:        prod.ID,
			Name:             prod.Name,
			Price:            prod.Price,
			Discount:         discount,
			Qty:              v.Quantity,
			TotalNormalPrice: prod.Price,
			TotalFinalPrice:  discountCount,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Grand Total Calculated Successfully!",
		"data": map[string]interface{}{
			"subtotal": totalInvoicePrice.price,
			"products": response,
		},
	})
}

// download your order
func DownloadOrder(c *fiber.Ctx) error {
	// check if token is present
	headerToken := c.Get("Authorization")
	if headerToken == "" {
		c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Token Not Found!",
		})
	}

	// check if the cashier is authorized
	if err := middlewares.AuthenticateToken(middlewares.SplitToken(headerToken)); err != nil {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "Unauthorized Cashier!",
			"error":   map[string]interface{}{},
		})
	}

	orderID := c.Params("orderID")
	var order models.Order
	db.DB.Where("id=?", orderID).First(&order)

	if order.ID == 0 {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Order Not Found!",
			"error":   map[string]interface{}{},
		})
	}

	productIDs := strings.Split(order.ProductID, ",")
	totalProducts := make([]*models.Product, 0)

	for i := 1; i < len(productIDs); i++ {
		prod := models.Product{}
		db.DB.Where("id=?", productIDs[i]).Find(&prod)
		totalProducts = append(totalProducts, &prod)
	}

	cashier := models.Cashier{}
	db.DB.Where("id=?", order.CashierID).Find(&cashier)

	paymentType := models.PaymentType{}
	db.DB.Where("id=?", order.PaymentTypesID).Find(&paymentType)

	orderTable := models.Order{}
	db.DB.Where("id=?", order.ID).Find(&orderTable)

	// pdf generating
	// creating the pdf structure
	vec := [][]string{{}}
	quantities := strings.Split(order.Quantities, ",")
	quantities = quantities[1:]

	for i := 0; i < len(totalProducts); i++ {
		arr := []string{}
		arr = append(arr, totalProducts[i].Sku)
		arr = append(arr, totalProducts[i].Name)
		arr = append(arr, quantities[i])
		arr = append(arr, strconv.Itoa(totalProducts[i].Price))
		vec = append(vec, arr)
	}

	begin := time.Now()
	grayColor := getGrayColor()
	whiteColor := color.NewWhite()
	header := getHeader()
	content := vec

	file := pdf.NewMaroto(consts.Portrait, consts.A4)
	file.SetPageMargins(10, 15, 10)

	// heading
	file.SetBackgroundColor(grayColor)
	file.Row(10, func() {
		file.Col(12, func() {
			file.Text("Order Invoice #"+strconv.Itoa(order.ID), props.Text{
				Top:   3,
				Style: consts.Bold,
				Align: consts.Center,
			})
		})
	})

	file.SetBackgroundColor(whiteColor)

	// table
	file.TableList(header, content, props.TableList{
		HeaderProp: props.TableListContent{
			Size:      9,
			GridSizes: []uint{3, 4, 2, 3},
		},
		ContentProp: props.TableListContent{
			Size:      8,
			GridSizes: []uint{3, 4, 2, 3},
		},
		Align:                consts.Center,
		AlternatedBackground: &grayColor,
		HeaderContentSpace:   1,
		Line:                 false,
	})

	// total price
	file.Row(20, func() {
		file.ColSpace(7)
		file.Col(2, func() {
			file.Text("Total:", props.Text{
				Top:   5,
				Style: consts.Bold,
				Size:  8,
				Align: consts.Right,
			})
		})

		file.Col(3, func() {
			file.Text("RS."+strconv.Itoa(order.TotalPrice), props.Text{
				Top:   5,
				Style: consts.Bold,
				Size:  8,
				Align: consts.Right,
			})
		})
	})

	file.Row(21, func() {
		file.ColSpace(7)
		file.Col(2, func() {
			file.Text("TotalPaid:", props.Text{
				Top:   0.5,
				Style: consts.Bold,
				Size:  8,
				Align: consts.Right,
			})
		})

		file.Col(3, func() {
			file.Text("RS."+strconv.Itoa(order.TotalPrice), props.Text{
				Top:   0.5,
				Style: consts.Bold,
				Size:  8,
				Align: consts.Right,
			})
		})
	})

	file.Row(22, func() {
		file.ColSpace(7)
		file.Col(2, func() {
			file.Text("Total Return", props.Text{
				Top:   5,
				Style: consts.Bold,
				Size:  8,
				Align: consts.Right,
			})
		})

		file.Col(3, func() {
			file.Text("RS."+strconv.Itoa(order.TotalPrice), props.Text{
				Top:   5,
				Style: consts.Bold,
				Size:  8,
				Align: consts.Right,
			})
		})
	})

	// invoice
	folderPath := "./invoice"
	currentTime := time.Now()
	fileName := "invoice-" + currentTime.Format("2006-Jan-02") + ".pdf"
	filePath := filepath.Join(folderPath, fileName)
	err := file.OutputFileAndClose(filePath)
	if err != nil {
		fmt.Println("Could not save PDF ->", err)
		os.Exit(1)
	}

	end := time.Now()
	fmt.Println(end.Sub(begin))

	// marking the pdf as downloaded
	db.DB.Table("orders").Where("id=?", order.ID).Update("is_download", 1)
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "PDF Downloaded Successfully!",
		"path":    filePath,
	})
}

func CheckOrder(c *fiber.Ctx) error {
	orderID := c.Params("orderID")
	var order models.Order
	db.DB.Where("id=?", orderID).First(&order)

	if order.ID == 0 {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Order does not exist",
		})
	}

	if order.IsDownload == 0 {
		return c.Status(200).JSON(fiber.Map{
			"success": true,
			"message": "Success",
			"data": map[string]interface{}{
				"isDownload": false,
			},
		})
	}

	if order.IsDownload == 1 {
		return c.Status(200).JSON(fiber.Map{
			"success": true,
			"message": "Success",
			"data": map[string]interface{}{
				"isDownload": true,
			},
		})
	}

	return nil
}

package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type PaymentRequest struct {
	PhoneNumber string  `json:"phone_number" binding:"required"`
	Amount      float64 `json:"amount" binding:"required"`
}

type BulkPaymentRequest struct {
	Payments []PaymentRequest `json:"payments" binding:"required,dive"`
}

type Payment struct {
	ID          uint   `gorm:"primaryKey"`
	PhoneNumber string `gorm:"not null"`
	Amount      float64
	Status      string `gorm:"default:'pending'"`
}

var db *gorm.DB

func initDB() {
	var err error
	dsn := "joelwasike:@Webuye2021@tcp(localhost:3306)/bulkpayments?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	db.AutoMigrate(&Payment{})
}

func bulkPaymentHandler(c *gin.Context) {
	var req BulkPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	responses := []gin.H{}

	for _, payment := range req.Payments {
		paymentRecord := Payment{
			PhoneNumber: payment.PhoneNumber,
			Amount:      payment.Amount,
			Status:      "processing",
		}

		if err := db.Create(&paymentRecord).Error; err != nil {
			responses = append(responses, gin.H{
				"phone_number": payment.PhoneNumber,
				"amount":       payment.Amount,
				"status":       "failed to log payment",
				"error":        err.Error(),
			})
			continue
		}

		// Simulate a payment API call (replace with actual API logic)
		if err := simulatePayment(payment.PhoneNumber, payment.Amount); err != nil {
			paymentRecord.Status = "failed"
			db.Save(&paymentRecord)
			responses = append(responses, gin.H{
				"phone_number": payment.PhoneNumber,
				"amount":       payment.Amount,
				"status":       "failed",
				"error":        err.Error(),
			})
		} else {
			paymentRecord.Status = "success"
			db.Save(&paymentRecord)
			responses = append(responses, gin.H{
				"phone_number": payment.PhoneNumber,
				"amount":       payment.Amount,
				"status":       "success",
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{"responses": responses})
}

func simulatePayment(phoneNumber string, amount float64) error {
	// Simulate a payment success/failure. Replace with real API logic.
	if amount <= 0 {
		return fmt.Errorf("invalid amount")
	}
	return nil
}

func main() {
	initDB()

	r := gin.Default()
	r.POST("/bulk-payments", bulkPaymentHandler)

	r.Run(":8080")
}

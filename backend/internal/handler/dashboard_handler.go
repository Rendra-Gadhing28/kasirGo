package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"kasirpro/internal/model"
	"kasirpro/internal/utils"
)

type DashboardHandler struct {
	db *gorm.DB
}

func NewDashboardHandler(db *gorm.DB) *DashboardHandler {
	return &DashboardHandler{db: db}
}

func (h *DashboardHandler) GetSummary(c *fiber.Ctx) error {
	outletID := c.Locals("outlet_id").(uint)

	// Today's boundaries
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayEnd := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location())

	// Today revenue & transactions
	var todayRevenue float64
	var todayTransactionCount int64
	h.db.Model(&model.Transaction{}).
		Where("outlet_id = ? AND payment_status = 'paid' AND created_at BETWEEN ? AND ?", outletID, todayStart, todayEnd).
		Select("COALESCE(SUM(total_amount), 0)").Scan(&todayRevenue)

	h.db.Model(&model.Transaction{}).
		Where("outlet_id = ? AND payment_status = 'paid' AND created_at BETWEEN ? AND ?", outletID, todayStart, todayEnd).
		Count(&todayTransactionCount)

	// Total active products
	var totalProducts int64
	h.db.Model(&model.Product{}).Where("outlet_id = ? AND is_active = true", outletID).Count(&totalProducts)

	// Low stock count (<= 10)
	var lowStockCount int64
	h.db.Model(&model.Product{}).Where("outlet_id = ? AND is_active = true AND stock <= 10", outletID).Count(&lowStockCount)

	// Total customers
	var totalCustomers int64
	h.db.Model(&model.Customer{}).Where("outlet_id = ?", outletID).Count(&totalCustomers)

	// 7 days trend for line chart
	type DailySale struct {
		Date  string  `json:"date"`
		Total float64 `json:"total"`
		Count int     `json:"count"`
	}
	var last7Days []DailySale
	for i := 6; i >= 0; i-- {
		targetDay := now.AddDate(0, 0, -i)
		start := time.Date(targetDay.Year(), targetDay.Month(), targetDay.Day(), 0, 0, 0, 0, targetDay.Location())
		end := time.Date(targetDay.Year(), targetDay.Month(), targetDay.Day(), 23, 59, 59, 999999999, targetDay.Location())

		var dayTotal float64
		var dayCount int64
		h.db.Model(&model.Transaction{}).
			Where("outlet_id = ? AND payment_status = 'paid' AND created_at BETWEEN ? AND ?", outletID, start, end).
			Select("COALESCE(SUM(total_amount), 0)").Scan(&dayTotal)

		h.db.Model(&model.Transaction{}).
			Where("outlet_id = ? AND payment_status = 'paid' AND created_at BETWEEN ? AND ?", outletID, start, end).
			Count(&dayCount)

		last7Days = append(last7Days, DailySale{
			Date:  targetDay.Format("02 Jan"),
			Total: dayTotal,
			Count: int(dayCount),
		})
	}

	// Payment method breakdown
	type PaymentBreakdown struct {
		Method string  `json:"method"`
		Count  int64   `json:"count"`
		Total  float64 `json:"total"`
	}
	var paymentMethods []PaymentBreakdown
	h.db.Model(&model.Transaction{}).
		Select("payment_method as method, count(*) as count, COALESCE(sum(total_amount), 0) as total").
		Where("outlet_id = ? AND payment_status = 'paid'", outletID).
		Group("payment_method").
		Scan(&paymentMethods)

	// Category sales breakdown
	type CategoryBreakdown struct {
		Category string  `json:"category"`
		Total    float64 `json:"total"`
	}
	var categorySales []CategoryBreakdown
	h.db.Table("transaction_items").
		Select("categories.name as category, COALESCE(SUM(transaction_items.subtotal), 0) as total").
		Joins("JOIN transactions ON transactions.id = transaction_items.transaction_id").
		Joins("JOIN products ON products.id = transaction_items.product_id").
		Joins("JOIN categories ON categories.id = products.category_id").
		Where("transactions.outlet_id = ? AND transactions.payment_status = 'paid'", outletID).
		Group("categories.name").
		Scan(&categorySales)

	summary := fiber.Map{
		"today_revenue":           todayRevenue,
		"today_transaction_count": todayTransactionCount,
		"total_products":          totalProducts,
		"low_stock_count":         lowStockCount,
		"total_customers":         totalCustomers,
		"sales_trend_7days":       last7Days,
		"payment_methods":         paymentMethods,
		"category_sales":          categorySales,
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Ringkasan dashboard", summary)
}

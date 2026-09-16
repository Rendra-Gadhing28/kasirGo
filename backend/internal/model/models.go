package model

import (
	"time"

	"gorm.io/gorm"
)

type Outlet struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Code      string         `gorm:"size:50;uniqueIndex" json:"code"`
	Name      string         `gorm:"size:255;not null" json:"name"`
	Address   string         `gorm:"type:text" json:"address"`
	Phone     string         `gorm:"size:50" json:"phone"`
	Logo      string         `gorm:"size:255" json:"logo"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	OutletID  uint           `gorm:"not null;index" json:"outlet_id"`
	Outlet    Outlet         `gorm:"foreignKey:OutletID" json:"outlet,omitempty"`
	Name      string         `gorm:"size:255;not null" json:"name"`
	Email     string         `gorm:"size:255;not null;uniqueIndex" json:"email"`
	Password  string         `gorm:"size:255;not null" json:"-"`
	Role      string         `gorm:"size:50;not null;default:'kasir'" json:"role"` // owner, admin, kasir
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	LastLogin *time.Time     `json:"last_login,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Category struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	OutletID    uint           `gorm:"not null;index" json:"outlet_id"`
	Name        string         `gorm:"size:255;not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Products    []Product      `gorm:"foreignKey:CategoryID" json:"products,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type Product struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	OutletID    uint           `gorm:"not null;index" json:"outlet_id"`
	CategoryID  uint           `gorm:"not null;index" json:"category_id"`
	Category    *Category      `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Name        string         `gorm:"size:255;not null;index" json:"name"`
	SKU         string         `gorm:"size:100;index" json:"sku"`
	Barcode     string         `gorm:"size:100;index" json:"barcode"`
	BuyPrice    float64        `gorm:"type:decimal(15,2);not null;default:0" json:"buy_price"`
	SellPrice   float64        `gorm:"type:decimal(15,2);not null" json:"sell_price"`
	Stock       int            `gorm:"not null;default:0" json:"stock"`
	Unit        string         `gorm:"size:50;not null;default:'pcs'" json:"unit"`
	ImageURL    string         `gorm:"size:500" json:"image_url"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type Customer struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	OutletID  uint           `gorm:"not null;index" json:"outlet_id"`
	Name      string         `gorm:"size:255;not null;index" json:"name"`
	Phone     string         `gorm:"size:50;index" json:"phone"`
	Email     string         `gorm:"size:255" json:"email"`
	Address   string         `gorm:"type:text" json:"address"`
	Points    int            `gorm:"default:0" json:"points"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Supplier struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	OutletID      uint           `gorm:"not null;index" json:"outlet_id"`
	Name          string         `gorm:"size:255;not null" json:"name"`
	Phone         string         `gorm:"size:50" json:"phone"`
	Email         string         `gorm:"size:255" json:"email"`
	Address       string         `gorm:"type:text" json:"address"`
	ContactPerson string         `gorm:"size:255" json:"contact_person"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

type Transaction struct {
	ID                     uint              `gorm:"primaryKey" json:"id"`
	OutletID               uint              `gorm:"not null;index" json:"outlet_id"`
	Outlet                 *Outlet           `gorm:"foreignKey:OutletID" json:"outlet,omitempty"`
	InvoiceNo              string            `gorm:"size:100;not null;uniqueIndex" json:"invoice_no"`
	CustomerID             *uint             `gorm:"index" json:"customer_id,omitempty"`
	Customer               *Customer         `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	UserID                 uint              `gorm:"not null;index" json:"user_id"`
	User                   *User             `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Subtotal               float64           `gorm:"type:decimal(15,2);not null" json:"subtotal"`
	DiscountType           string            `gorm:"size:20;default:'none'" json:"discount_type"` // none, percentage, fixed
	DiscountValue          float64           `gorm:"type:decimal(15,2);default:0" json:"discount_value"`
	DiscountAmount         float64           `gorm:"type:decimal(15,2);default:0" json:"discount_amount"`
	TaxRate                float64           `gorm:"type:decimal(5,2);default:0" json:"tax_rate"`
	TaxAmount              float64           `gorm:"type:decimal(15,2);default:0" json:"tax_amount"`
	TotalAmount            float64           `gorm:"type:decimal(15,2);not null" json:"total_amount"`
	PaymentMethod          string            `gorm:"size:50;not null" json:"payment_method"` // cash, qris, transfer
	PaymentStatus          string            `gorm:"size:50;not null;default:'pending'" json:"payment_status"` // pending, paid, cancelled, refunded
	PaidAmount             float64           `gorm:"type:decimal(15,2);default:0" json:"paid_amount"`
	ChangeAmount           float64           `gorm:"type:decimal(15,2);default:0" json:"change_amount"`
	Notes                  string            `gorm:"type:text" json:"notes"`
	MidtransOrderID        string            `gorm:"size:100;index" json:"midtrans_order_id,omitempty"`
	MidtransQRCode         string            `gorm:"type:text" json:"midtrans_qr_code,omitempty"`
	MidtransTransactionStatus string         `gorm:"size:50" json:"midtrans_transaction_status,omitempty"`
	Items                  []TransactionItem `gorm:"foreignKey:TransactionID" json:"items"`
	CreatedAt              time.Time         `json:"created_at"`
	UpdatedAt              time.Time         `json:"updated_at"`
	DeletedAt              gorm.DeletedAt    `gorm:"index" json:"-"`
}

type TransactionItem struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	TransactionID  uint           `gorm:"not null;index" json:"transaction_id"`
	ProductID      uint           `gorm:"not null;index" json:"product_id"`
	Product        *Product       `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	ProductName    string         `gorm:"size:255;not null" json:"product_name"`
	ProductSKU     string         `gorm:"size:100" json:"product_sku"`
	Qty            int            `gorm:"not null" json:"qty"`
	Unit           string         `gorm:"size:50;default:'pcs'" json:"unit"`
	Price          float64        `gorm:"type:decimal(15,2);not null" json:"price"`
	BuyPrice       float64        `gorm:"type:decimal(15,2);default:0" json:"buy_price"`
	DiscountType   string         `gorm:"size:20;default:'none'" json:"discount_type"`
	DiscountValue  float64        `gorm:"type:decimal(15,2);default:0" json:"discount_value"`
	DiscountAmount float64        `gorm:"type:decimal(15,2);default:0" json:"discount_amount"`
	Subtotal       float64        `gorm:"type:decimal(15,2);not null" json:"subtotal"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

type ExpenseIncome struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	OutletID    uint           `gorm:"not null;index" json:"outlet_id"`
	UserID      uint           `gorm:"not null;index" json:"user_id"`
	User        *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Type        string         `gorm:"size:20;not null" json:"type"` // expense or income
	Category    string         `gorm:"size:100;not null" json:"category"`
	Amount      float64        `gorm:"type:decimal(15,2);not null" json:"amount"`
	Description string         `gorm:"type:text" json:"description"`
	Date        time.Time      `gorm:"not null" json:"date"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

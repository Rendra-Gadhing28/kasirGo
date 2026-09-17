package dto

// Auth DTOs
type RegisterRequest struct {
	StoreName string `json:"store_name" validate:"required,min=3,max=255"`
	Name      string `json:"name" validate:"required,min=2,max=255"`
	Email     string `json:"email" validate:"required,email,max=255"`
	Password  string `json:"password" validate:"required,min=6,max=100"`
	Phone     string `json:"phone" validate:"omitempty,max=50"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	User         UserSummary `json:"user"`
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
}

type StaffRegisterRequest struct {
	OutletCode string `json:"outlet_code" validate:"required,min=4,max=50"`
	Name       string `json:"name" validate:"required,min=2,max=255"`
	Email      string `json:"email" validate:"required,email,max=255"`
	Password   string `json:"password" validate:"required,min=6,max=100"`
	Role       string `json:"role" validate:"omitempty,oneof=kasir admin"`
}

type PublicMemberRegisterRequest struct {
	OutletCode string `json:"outlet_code" validate:"required,min=4,max=50"`
	Name       string `json:"name" validate:"required,min=2,max=255"`
	Phone      string `json:"phone" validate:"required,min=8,max=50"`
	Email      string `json:"email" validate:"omitempty,email,max=255"`
	Address    string `json:"address" validate:"omitempty,max=1000"`
}

type PublicOutletInfo struct {
	Name    string `json:"name"`
	Code    string `json:"code"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
}

type UserSummary struct {
	ID         uint   `json:"id"`
	OutletID   uint   `json:"outlet_id"`
	OutletCode string `json:"outlet_code,omitempty"`
	OutletName string `json:"outlet_name,omitempty"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Role       string `json:"role"`
}

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=255"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=6,max=100"`
	Role     string `json:"role" validate:"required,oneof=admin kasir"`
}

// Category DTOs
type CreateCategoryRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=255"`
	Description string `json:"description" validate:"omitempty,max=1000"`
}

type UpdateCategoryRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=255"`
	Description string `json:"description" validate:"omitempty,max=1000"`
}

// Product DTOs
type CreateProductRequest struct {
	CategoryID uint    `json:"category_id" validate:"required,gt=0"`
	Name       string  `json:"name" validate:"required,min=2,max=255"`
	SKU        string  `json:"sku" validate:"omitempty,max=100"`
	Barcode    string  `json:"barcode" validate:"omitempty,max=100"`
	BuyPrice   float64 `json:"buy_price" validate:"gte=0"`
	SellPrice  float64 `json:"sell_price" validate:"required,gt=0"`
	Stock      int     `json:"stock" validate:"gte=0"`
	Unit       string  `json:"unit" validate:"required,max=50"`
	ImageURL   string  `json:"image_url" validate:"omitempty,max=500"`
}

type UpdateProductRequest struct {
	CategoryID uint    `json:"category_id" validate:"required,gt=0"`
	Name       string  `json:"name" validate:"required,min=2,max=255"`
	SKU        string  `json:"sku" validate:"omitempty,max=100"`
	Barcode    string  `json:"barcode" validate:"omitempty,max=100"`
	BuyPrice   float64 `json:"buy_price" validate:"gte=0"`
	SellPrice  float64 `json:"sell_price" validate:"required,gt=0"`
	Stock      int     `json:"stock" validate:"gte=0"`
	Unit       string  `json:"unit" validate:"required,max=50"`
	ImageURL   string  `json:"image_url" validate:"omitempty,max=500"`
}

// Customer DTOs
type CreateCustomerRequest struct {
	MemberCode string `json:"member_code" validate:"omitempty,max=50"`
	Name       string `json:"name" validate:"required,min=2,max=255"`
	Phone      string `json:"phone" validate:"omitempty,max=50"`
	Email      string `json:"email" validate:"omitempty,email,max=255"`
	Address    string `json:"address" validate:"omitempty,max=1000"`
}

type UpdateCustomerRequest struct {
	MemberCode string `json:"member_code" validate:"omitempty,max=50"`
	Name       string `json:"name" validate:"required,min=2,max=255"`
	Phone      string `json:"phone" validate:"omitempty,max=50"`
	Email      string `json:"email" validate:"omitempty,email,max=255"`
	Address    string `json:"address" validate:"omitempty,max=1000"`
}

// Supplier DTOs
type CreateSupplierRequest struct {
	Name          string `json:"name" validate:"required,min=2,max=255"`
	Phone         string `json:"phone" validate:"omitempty,max=50"`
	Email         string `json:"email" validate:"omitempty,email,max=255"`
	Address       string `json:"address" validate:"omitempty,max=1000"`
	ContactPerson string `json:"contact_person" validate:"omitempty,max=255"`
}

type UpdateSupplierRequest struct {
	Name          string `json:"name" validate:"required,min=2,max=255"`
	Phone         string `json:"phone" validate:"omitempty,max=50"`
	Email         string `json:"email" validate:"omitempty,email,max=255"`
	Address       string `json:"address" validate:"omitempty,max=1000"`
	ContactPerson string `json:"contact_person" validate:"omitempty,max=255"`
}

// Transaction / POS DTOs
type CheckoutItemRequest struct {
	ProductID     uint    `json:"product_id" validate:"required,gt=0"`
	Qty           int     `json:"qty" validate:"required,gt=0"`
	DiscountType  string  `json:"discount_type" validate:"omitempty,oneof=none percentage fixed"`
	DiscountValue float64 `json:"discount_value" validate:"gte=0"`
}

type CheckoutRequest struct {
	CustomerID    *uint                 `json:"customer_id" validate:"omitempty,gt=0"`
	Items         []CheckoutItemRequest `json:"items" validate:"required,min=1,dive"`
	DiscountType  string                `json:"discount_type" validate:"omitempty,oneof=none percentage fixed"`
	DiscountValue float64               `json:"discount_value" validate:"gte=0"`
	TaxEnabled    bool                  `json:"tax_enabled"`
	PaymentMethod string                `json:"payment_method" validate:"required,oneof=cash qris transfer"`
	PaidAmount    float64               `json:"paid_amount" validate:"gte=0"`
	Notes         string                `json:"notes" validate:"omitempty,max=500"`
}

type MidtransCallbackRequest struct {
	TransactionStatus string `json:"transaction_status"`
	OrderID           string `json:"order_id"`
	GrossAmount       string `json:"gross_amount"`
	PaymentType       string `json:"payment_type"`
	FraudStatus       string `json:"fraud_status"`
	StatusCode        string `json:"status_code"`
	SignatureKey      string `json:"signature_key"`
}

// Query parameters for pagination and search
type QueryFilter struct {
	Search    string `query:"search"`
	Category  uint   `query:"category_id"`
	Page      int    `query:"page"`
	Limit     int    `query:"limit"`
	SortBy    string `query:"sort_by"`
	SortOrder string `query:"sort_order"`
	StartDate string `query:"start_date"`
	EndDate   string `query:"end_date"`
	Status    string `query:"status"`
}

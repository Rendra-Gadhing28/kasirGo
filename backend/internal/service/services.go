package service

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"kasirpro/internal/config"
	"kasirpro/internal/dto"
	"kasirpro/internal/model"
	"kasirpro/internal/repository"
	"kasirpro/internal/utils"
)

// AuthService handles authentication and registration
type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) Register(req dto.RegisterRequest) (*dto.AuthResponse, error) {
	// Sanitize inputs
	storeName := utils.SanitizeString(req.StoreName)
	name := utils.SanitizeString(req.Name)
	email := strings.ToLower(strings.TrimSpace(req.Email))

	// Check existing
	if existing, _ := s.userRepo.FindByEmail(email); existing != nil {
		return nil, errors.New("email sudah terdaftar di sistem")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("gagal memproses password")
	}

	outletCode := "KASIR-" + strings.ToUpper(uuid.New().String()[:6])
	outlet := model.Outlet{
		Code:    outletCode,
		Name:    storeName,
		Address: "Alamat Toko Baru",
		Phone:   req.Phone,
	}
	if err := s.userRepo.CreateOutlet(&outlet); err != nil {
		return nil, errors.New("gagal membuat data outlet toko")
	}

	user := model.User{
		OutletID: outlet.ID,
		Name:     name,
		Email:    email,
		Password: hashedPassword,
		Role:     "owner",
		IsActive: true,
	}
	if err := s.userRepo.Create(&user); err != nil {
		return nil, errors.New("gagal membuat akun pengguna")
	}

	accessToken, refreshToken, err := utils.GenerateTokenPair(user.ID, outlet.ID, user.Role, user.Email, user.Name)
	if err != nil {
		return nil, errors.New("gagal menghasilkan token autentikasi")
	}

	return &dto.AuthResponse{
		User: dto.UserSummary{
			ID:         user.ID,
			OutletID:   user.OutletID,
			OutletCode: outlet.Code,
			OutletName: outlet.Name,
			Name:       user.Name,
			Email:      user.Email,
			Role:       user.Role,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) RegisterStaff(req dto.StaffRegisterRequest) (*dto.AuthResponse, error) {
	outletCode := strings.ToUpper(strings.TrimSpace(req.OutletCode))
	outlet, err := s.userRepo.FindOutletByCode(outletCode)
	if err != nil || outlet == nil {
		return nil, errors.New("kode outlet / toko tidak valid atau tidak ditemukan")
	}

	name := utils.SanitizeString(req.Name)
	email := strings.ToLower(strings.TrimSpace(req.Email))

	if existing, _ := s.userRepo.FindByEmail(email); existing != nil {
		return nil, errors.New("email sudah terdaftar di sistem")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("gagal memproses password")
	}

	role := "kasir"
	if req.Role == "admin" {
		role = "admin"
	}

	user := model.User{
		OutletID: outlet.ID,
		Name:     name,
		Email:    email,
		Password: hashedPassword,
		Role:     role,
		IsActive: true,
	}
	if err := s.userRepo.Create(&user); err != nil {
		return nil, errors.New("gagal membuat akun karyawan/kasir")
	}

	accessToken, refreshToken, err := utils.GenerateTokenPair(user.ID, outlet.ID, user.Role, user.Email, user.Name)
	if err != nil {
		return nil, errors.New("gagal menghasilkan token autentikasi")
	}

	return &dto.AuthResponse{
		User: dto.UserSummary{
			ID:         user.ID,
			OutletID:   user.OutletID,
			OutletCode: outlet.Code,
			OutletName: outlet.Name,
			Name:       user.Name,
			Email:      user.Email,
			Role:       user.Role,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) Login(req dto.LoginRequest) (*dto.AuthResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("email atau password salah")
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("email atau password salah")
	}

	accessToken, refreshToken, err := utils.GenerateTokenPair(user.ID, user.OutletID, user.Role, user.Email, user.Name)
	if err != nil {
		return nil, errors.New("gagal menghasilkan token autentikasi")
	}

	outletName := ""
	outletCode := ""
	if user.Outlet.ID > 0 {
		outletName = user.Outlet.Name
		outletCode = user.Outlet.Code
	}

	return &dto.AuthResponse{
		User: dto.UserSummary{
			ID:         user.ID,
			OutletID:   user.OutletID,
			OutletCode: outletCode,
			OutletName: outletName,
			Name:       user.Name,
			Email:      user.Email,
			Role:       user.Role,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) GetProfile(userID uint) (*model.User, error) {
	return s.userRepo.FindByID(userID)
}

// CategoryService handles category CRUD with caching
type CategoryService struct {
	repo *repository.CategoryRepository
}

func NewCategoryService(repo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) GetAll(outletID uint) ([]model.Category, error) {
	cacheKey := fmt.Sprintf("categories_outlet_%d", outletID)
	if cached, found := utils.GetCache(cacheKey); found {
		return cached.([]model.Category), nil
	}

	categories, err := s.repo.FindAll(outletID)
	if err != nil {
		return nil, err
	}

	utils.SetCache(cacheKey, categories, 10*time.Minute)
	return categories, nil
}

func (s *CategoryService) Create(outletID uint, req dto.CreateCategoryRequest) (*model.Category, error) {
	cat := &model.Category{
		OutletID:    outletID,
		Name:        utils.SanitizeString(req.Name),
		Description: utils.SanitizeString(req.Description),
	}
	if err := s.repo.Create(cat); err != nil {
		return nil, err
	}

	utils.DeleteCache(fmt.Sprintf("categories_outlet_%d", outletID))
	return cat, nil
}

func (s *CategoryService) Update(outletID, id uint, req dto.UpdateCategoryRequest) (*model.Category, error) {
	cat, err := s.repo.FindByID(outletID, id)
	if err != nil {
		return nil, errors.New("kategori tidak ditemukan")
	}

	cat.Name = utils.SanitizeString(req.Name)
	cat.Description = utils.SanitizeString(req.Description)

	if err := s.repo.Update(cat); err != nil {
		return nil, err
	}

	utils.DeleteCache(fmt.Sprintf("categories_outlet_%d", outletID))
	return cat, nil
}

func (s *CategoryService) Delete(outletID, id uint) error {
	if err := s.repo.Delete(outletID, id); err != nil {
		return err
	}
	utils.DeleteCache(fmt.Sprintf("categories_outlet_%d", outletID))
	return nil
}

// ProductService handles product CRUD and inventory
type ProductService struct {
	repo *repository.ProductRepository
}

func NewProductService(repo *repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) GetAll(outletID uint, filter dto.QueryFilter) ([]model.Product, int64, error) {
	return s.repo.FindAll(outletID, filter)
}

func (s *ProductService) GetByID(outletID, id uint) (*model.Product, error) {
	return s.repo.FindByID(outletID, id)
}

func (s *ProductService) GetByBarcode(outletID uint, barcode string) (*model.Product, error) {
	return s.repo.FindByBarcode(outletID, barcode)
}

func (s *ProductService) Create(outletID uint, req dto.CreateProductRequest) (*model.Product, error) {
	sku := utils.SanitizeString(req.SKU)
	if sku == "" {
		sku = "PRD-" + strings.ToUpper(uuid.New().String()[:8])
	}

	barcode := utils.SanitizeString(req.Barcode)
	if barcode == "" {
		barcode = sku
	}

	product := &model.Product{
		OutletID:   outletID,
		CategoryID: req.CategoryID,
		Name:       utils.SanitizeString(req.Name),
		SKU:        sku,
		Barcode:    barcode,
		BuyPrice:   req.BuyPrice,
		SellPrice:  req.SellPrice,
		Stock:      req.Stock,
		Unit:       utils.SanitizeString(req.Unit),
		ImageURL:   utils.SanitizeString(req.ImageURL),
		IsActive:   true,
	}

	if err := s.repo.Create(product); err != nil {
		return nil, err
	}
	return product, nil
}

func (s *ProductService) Update(outletID, id uint, req dto.UpdateProductRequest) (*model.Product, error) {
	product, err := s.repo.FindByID(outletID, id)
	if err != nil {
		return nil, errors.New("produk tidak ditemukan")
	}

	sku := utils.SanitizeString(req.SKU)
	if sku == "" {
		sku = product.SKU
	}
	barcode := utils.SanitizeString(req.Barcode)
	if barcode == "" {
		barcode = product.Barcode
	}

	product.CategoryID = req.CategoryID
	product.Name = utils.SanitizeString(req.Name)
	product.SKU = sku
	product.Barcode = barcode
	product.BuyPrice = req.BuyPrice
	product.SellPrice = req.SellPrice
	product.Stock = req.Stock
	product.Unit = utils.SanitizeString(req.Unit)
	product.ImageURL = utils.SanitizeString(req.ImageURL)

	if err := s.repo.Update(product); err != nil {
		return nil, err
	}
	return product, nil
}

func (s *ProductService) Delete(outletID, id uint) error {
	return s.repo.Delete(outletID, id)
}

// CustomerService handles customer management
type CustomerService struct {
	repo     *repository.CustomerRepository
	userRepo *repository.UserRepository
}

func NewCustomerService(repo *repository.CustomerRepository, userRepo *repository.UserRepository) *CustomerService {
	return &CustomerService{repo: repo, userRepo: userRepo}
}

func (s *CustomerService) GetPublicOutletInfo(code string) (*dto.PublicOutletInfo, error) {
	outlet, err := s.userRepo.FindOutletByCode(code)
	if err != nil || outlet == nil {
		return nil, errors.New("kode toko tidak ditemukan")
	}
	return &dto.PublicOutletInfo{
		Name:    outlet.Name,
		Code:    outlet.Code,
		Address: outlet.Address,
		Phone:   outlet.Phone,
	}, nil
}

func (s *CustomerService) RegisterPublicMember(req dto.PublicMemberRegisterRequest) (*model.Customer, *model.Outlet, error) {
	outlet, err := s.userRepo.FindOutletByCode(req.OutletCode)
	if err != nil || outlet == nil {
		return nil, nil, errors.New("kode toko / outlet tidak ditemukan")
	}

	phone := utils.SanitizeString(req.Phone)
	name := utils.SanitizeString(req.Name)
	email := utils.SanitizeString(req.Email)
	address := utils.SanitizeString(req.Address)

	// Check if already registered in this outlet by phone
	existing, _ := s.repo.FindAll(outlet.ID, phone)
	for _, c := range existing {
		if c.Phone == phone {
			return &c, outlet, nil
		}
	}

	customer := &model.Customer{
		OutletID: outlet.ID,
		Name:     name,
		Phone:    phone,
		Email:    email,
		Address:  address,
		Points:   10, // Welcome bonus points!
	}

	if err := s.repo.Create(customer); err != nil {
		return nil, nil, errors.New("gagal mendaftarkan member")
	}

	return customer, outlet, nil
}

func (s *CustomerService) GetAll(outletID uint, search string) ([]model.Customer, error) {
	return s.repo.FindAll(outletID, search)
}

func (s *CustomerService) GetByID(outletID, id uint) (*model.Customer, error) {
	return s.repo.FindByID(outletID, id)
}

func (s *CustomerService) Create(outletID uint, req dto.CreateCustomerRequest) (*model.Customer, error) {
	customer := &model.Customer{
		OutletID: outletID,
		Name:     utils.SanitizeString(req.Name),
		Phone:    utils.SanitizeString(req.Phone),
		Email:    utils.SanitizeString(req.Email),
		Address:  utils.SanitizeString(req.Address),
		Points:   0,
	}
	if err := s.repo.Create(customer); err != nil {
		return nil, err
	}
	return customer, nil
}

func (s *CustomerService) Update(outletID, id uint, req dto.UpdateCustomerRequest) (*model.Customer, error) {
	customer, err := s.repo.FindByID(outletID, id)
	if err != nil {
		return nil, errors.New("pelanggan tidak ditemukan")
	}

	customer.Name = utils.SanitizeString(req.Name)
	customer.Phone = utils.SanitizeString(req.Phone)
	customer.Email = utils.SanitizeString(req.Email)
	customer.Address = utils.SanitizeString(req.Address)

	if err := s.repo.Update(customer); err != nil {
		return nil, err
	}
	return customer, nil
}

func (s *CustomerService) Delete(outletID, id uint) error {
	return s.repo.Delete(outletID, id)
}

// SupplierService handles supplier management
type SupplierService struct {
	repo *repository.SupplierRepository
}

func NewSupplierService(repo *repository.SupplierRepository) *SupplierService {
	return &SupplierService{repo: repo}
}

func (s *SupplierService) GetAll(outletID uint, search string) ([]model.Supplier, error) {
	return s.repo.FindAll(outletID, search)
}

func (s *SupplierService) GetByID(outletID, id uint) (*model.Supplier, error) {
	return s.repo.FindByID(outletID, id)
}

func (s *SupplierService) Create(outletID uint, req dto.CreateSupplierRequest) (*model.Supplier, error) {
	supplier := &model.Supplier{
		OutletID:      outletID,
		Name:          utils.SanitizeString(req.Name),
		Phone:         utils.SanitizeString(req.Phone),
		Email:         utils.SanitizeString(req.Email),
		Address:       utils.SanitizeString(req.Address),
		ContactPerson: utils.SanitizeString(req.ContactPerson),
	}
	if err := s.repo.Create(supplier); err != nil {
		return nil, err
	}
	return supplier, nil
}

func (s *SupplierService) Update(outletID, id uint, req dto.UpdateSupplierRequest) (*model.Supplier, error) {
	supplier, err := s.repo.FindByID(outletID, id)
	if err != nil {
		return nil, errors.New("supplier tidak ditemukan")
	}

	supplier.Name = utils.SanitizeString(req.Name)
	supplier.Phone = utils.SanitizeString(req.Phone)
	supplier.Email = utils.SanitizeString(req.Email)
	supplier.Address = utils.SanitizeString(req.Address)
	supplier.ContactPerson = utils.SanitizeString(req.ContactPerson)

	if err := s.repo.Update(supplier); err != nil {
		return nil, err
	}
	return supplier, nil
}

func (s *SupplierService) Delete(outletID, id uint) error {
	return s.repo.Delete(outletID, id)
}

// TransactionService handles POS checkout, inventory deduction, loyalty points, and thermal receipt
type TransactionService struct {
	transRepo   *repository.TransactionRepository
	productRepo *repository.ProductRepository
	custRepo    *repository.CustomerRepository
}

func NewTransactionService(
	transRepo *repository.TransactionRepository,
	productRepo *repository.ProductRepository,
	custRepo *repository.CustomerRepository,
) *TransactionService {
	return &TransactionService{
		transRepo:   transRepo,
		productRepo: productRepo,
		custRepo:    custRepo,
	}
}

func (s *TransactionService) GetAll(outletID uint, filter dto.QueryFilter, currentUserID uint, role string) ([]model.Transaction, int64, error) {
	return s.transRepo.FindAll(outletID, filter, currentUserID, role)
}

func (s *TransactionService) GetByID(outletID, id uint) (*model.Transaction, error) {
	return s.transRepo.FindByID(outletID, id)
}

func (s *TransactionService) GetByInvoice(outletID uint, invoiceNo string) (*model.Transaction, error) {
	return s.transRepo.FindByInvoice(outletID, invoiceNo)
}

func (s *TransactionService) Checkout(outletID, userID uint, req dto.CheckoutRequest) (*model.Transaction, error) {
	db := s.transRepo.GetDB()

	// Begin atomic database transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var subtotal float64
	var items []model.TransactionItem

	// 1. Process and validate each item in the cart
	for _, itemReq := range req.Items {
		var product model.Product
		if err := tx.Where("outlet_id = ? AND id = ?", outletID, itemReq.ProductID).First(&product).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("produk ID %d tidak ditemukan", itemReq.ProductID)
		}

		if product.Stock < itemReq.Qty {
			tx.Rollback()
			return nil, fmt.Errorf("stok produk '%s' tidak mencukupi (tersedia: %d, diminta: %d)", product.Name, product.Stock, itemReq.Qty)
		}

		itemGross := product.SellPrice * float64(itemReq.Qty)
		var itemDiscountAmount float64
		if itemReq.DiscountType == "percentage" && itemReq.DiscountValue > 0 {
			itemDiscountAmount = itemGross * (itemReq.DiscountValue / 100.0)
		} else if itemReq.DiscountType == "fixed" && itemReq.DiscountValue > 0 {
			itemDiscountAmount = itemReq.DiscountValue
		}
		if itemDiscountAmount > itemGross {
			itemDiscountAmount = itemGross
		}
		itemSubtotal := itemGross - itemDiscountAmount
		subtotal += itemSubtotal

		// Decrease product stock atomically
		if err := s.productRepo.DecreaseStockTx(tx, product.ID, itemReq.Qty); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("gagal memotong stok untuk produk %s", product.Name)
		}

		items = append(items, model.TransactionItem{
			ProductID:      product.ID,
			ProductName:    product.Name,
			ProductSKU:     product.SKU,
			Qty:            itemReq.Qty,
			Unit:           product.Unit,
			Price:          product.SellPrice,
			BuyPrice:       product.BuyPrice,
			DiscountType:   itemReq.DiscountType,
			DiscountValue:  itemReq.DiscountValue,
			DiscountAmount: itemDiscountAmount,
			Subtotal:       itemSubtotal,
		})
	}

	// 2. Global Transaction Discount calculation
	var globalDiscountAmount float64
	if req.DiscountType == "percentage" && req.DiscountValue > 0 {
		globalDiscountAmount = subtotal * (req.DiscountValue / 100.0)
	} else if req.DiscountType == "fixed" && req.DiscountValue > 0 {
		globalDiscountAmount = req.DiscountValue
	}
	if globalDiscountAmount > subtotal {
		globalDiscountAmount = subtotal
	}
	discountedSubtotal := subtotal - globalDiscountAmount

	// 3. Tax calculation (Indonesian PPN 11%)
	var taxRate float64
	var taxAmount float64
	if req.TaxEnabled {
		taxRate = config.AppConfig.DefaultTaxRate
		taxAmount = math.Round(discountedSubtotal * (taxRate / 100.0))
	}
	totalAmount := discountedSubtotal + taxAmount

	// 4. Payment validation
	paymentStatus := "paid"
	changeAmount := 0.0

	invoiceNo := fmt.Sprintf("INV-%s-%s", time.Now().Format("20060102150405"), strings.ToUpper(uuid.New().String()[:4]))
	midtransOrderID := ""
	midtransQRCode := ""
	midtransStatus := ""

	if req.PaymentMethod == "cash" {
		if req.PaidAmount < totalAmount {
			tx.Rollback()
			return nil, fmt.Errorf("jumlah uang bayar (Rp %.0f) kurang dari total tagihan (Rp %.0f)", req.PaidAmount, totalAmount)
		}
		changeAmount = req.PaidAmount - totalAmount
	} else if req.PaymentMethod == "qris" {
		paymentStatus = "pending"
		midtransOrderID = fmt.Sprintf("QRIS-%s", invoiceNo)
		midtransStatus = "pending"
		// Generate standard QRIS simulation string and QR payload
		midtransQRCode = fmt.Sprintf("https://api.sandbox.midtrans.com/v2/qris/%s/qr-code", midtransOrderID)
	}

	trans := model.Transaction{
		OutletID:                  outletID,
		InvoiceNo:                 invoiceNo,
		CustomerID:                req.CustomerID,
		UserID:                    userID,
		Subtotal:                  subtotal,
		DiscountType:              req.DiscountType,
		DiscountValue:             req.DiscountValue,
		DiscountAmount:            globalDiscountAmount,
		TaxRate:                   taxRate,
		TaxAmount:                 taxAmount,
		TotalAmount:               totalAmount,
		PaymentMethod:             req.PaymentMethod,
		PaymentStatus:             paymentStatus,
		PaidAmount:                req.PaidAmount,
		ChangeAmount:              changeAmount,
		Notes:                     utils.SanitizeString(req.Notes),
		MidtransOrderID:           midtransOrderID,
		MidtransQRCode:            midtransQRCode,
		MidtransTransactionStatus: midtransStatus,
		Items:                     items,
	}

	if err := tx.Create(&trans).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("gagal menyimpan transaksi: %v", err)
	}

	// 5. Add loyalty points if customer is selected (1 point per Rp 10.000)
	if req.CustomerID != nil && *req.CustomerID > 0 && paymentStatus == "paid" {
		earnedPoints := int(totalAmount / 10000.0)
		if earnedPoints > 0 {
			_ = s.custRepo.AddPointsTx(tx, *req.CustomerID, earnedPoints)
		}
	}

	// Commit atomic transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("gagal commit transaksi: %v", err)
	}

	// Reload with relations for receipt
	return s.transRepo.FindByID(outletID, trans.ID)
}

// VoidTransaction cancels a transaction and restores inventory (Admin/Owner only)
func (s *TransactionService) VoidTransaction(outletID, transID uint, reason string) error {
	db := s.transRepo.GetDB()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var trans model.Transaction
	if err := tx.Preload("Items").Where("outlet_id = ? AND id = ?", outletID, transID).First(&trans).Error; err != nil {
		tx.Rollback()
		return errors.New("transaksi tidak ditemukan")
	}

	if trans.PaymentStatus == "cancelled" || trans.PaymentStatus == "refunded" {
		tx.Rollback()
		return errors.New("transaksi sudah dibatalkan sebelumnya")
	}

	// Restore stock for all items
	for _, item := range trans.Items {
		if err := s.productRepo.RestoreStockTx(tx, item.ProductID, item.Qty); err != nil {
			tx.Rollback()
			return fmt.Errorf("gagal mengembalikan stok produk: %v", err)
		}
	}

	// Deduct customer points if applicable
	if trans.CustomerID != nil && *trans.CustomerID > 0 && trans.PaymentStatus == "paid" {
		pointsToDeduct := int(trans.TotalAmount / 10000.0)
		if pointsToDeduct > 0 {
			_ = s.custRepo.AddPointsTx(tx, *trans.CustomerID, -pointsToDeduct)
		}
	}

	trans.PaymentStatus = "cancelled"
	trans.Notes = fmt.Sprintf("%s | VOID REASON: %s", trans.Notes, utils.SanitizeString(reason))
	if err := tx.Save(&trans).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal mengupdate status pembatalan transaksi: %v", err)
	}

	return tx.Commit().Error
}

// MidtransPaymentService handles QRIS generation and webhooks/status checks
type MidtransPaymentService struct {
	transRepo *repository.TransactionRepository
	custRepo  *repository.CustomerRepository
}

func NewMidtransPaymentService(transRepo *repository.TransactionRepository, custRepo *repository.CustomerRepository) *MidtransPaymentService {
	return &MidtransPaymentService{transRepo: transRepo, custRepo: custRepo}
}

// GenerateQRIS creates or simulates dynamic QRIS via Midtrans API
func (s *MidtransPaymentService) GenerateQRIS(trans *model.Transaction) (string, error) {
	cfg := config.AppConfig
	serverKey := cfg.MidtransServerKey

	// If using Sandbox credentials, make HTTP request to Midtrans Charge API
	if serverKey != "" && !strings.Contains(serverKey, "simulator-mock") {
		url := "https://api.sandbox.midtrans.com/v2/charge"
		if cfg.MidtransIsProduction {
			url = "https://api.midtrans.com/v2/charge"
		}

		payload := map[string]interface{}{
			"payment_type": "qris",
			"transaction_details": map[string]interface{}{
				"order_id":     trans.MidtransOrderID,
				"gross_amount": int(trans.TotalAmount),
			},
			"qris": map[string]string{
				"acquirer": "gopay",
			},
		}

		payloadBytes, _ := json.Marshal(payload)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
		if err == nil {
			authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(serverKey+":"))
			req.Header.Set("Authorization", authHeader)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Accept", "application/json")

			client := &http.Client{Timeout: 10 * time.Second}
			resp, err := client.Do(req)
			if err == nil {
				defer resp.Body.Close()
				body, _ := io.ReadAll(resp.Body)
				var midtransResp struct {
					Actions []struct {
						Name string `json:"name"`
						URL  string `json:"url"`
					} `json:"actions"`
					QRString string `json:"qr_string"`
				}
				if json.Unmarshal(body, &midtransResp) == nil {
					for _, action := range midtransResp.Actions {
						if action.Name == "generate-qr-code" {
							return action.URL, nil
						}
					}
					if midtransResp.QRString != "" {
						return midtransResp.QRString, nil
					}
				}
			}
		}
	}

	// Midtrans Simulator fallback: generate dynamic QRIS simulation URL
	simulatorQR := fmt.Sprintf("https://api.sandbox.midtrans.com/v2/qris/%s/qr-code", trans.MidtransOrderID)
	return simulatorQR, nil
}

// ConfirmQRISPayment simulates or records payment confirmation
func (s *MidtransPaymentService) ConfirmQRISPayment(outletID uint, invoiceNo string) (*model.Transaction, error) {
	trans, err := s.transRepo.FindByInvoice(outletID, invoiceNo)
	if err != nil {
		return nil, errors.New("transaksi tidak ditemukan")
	}

	if trans.PaymentStatus == "paid" {
		return trans, nil
	}

	trans.PaymentStatus = "paid"
	trans.PaidAmount = trans.TotalAmount
	trans.MidtransTransactionStatus = "settlement"

	db := s.transRepo.GetDB()
	if err := db.Save(trans).Error; err != nil {
		return nil, err
	}

	// Award loyalty points
	if trans.CustomerID != nil && *trans.CustomerID > 0 {
		earnedPoints := int(trans.TotalAmount / 10000.0)
		if earnedPoints > 0 {
			_ = s.custRepo.AddPointsTx(db, *trans.CustomerID, earnedPoints)
		}
	}

	return trans, nil
}

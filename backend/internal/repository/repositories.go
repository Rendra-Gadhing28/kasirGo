package repository

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
	"kasirpro/internal/dto"
	"kasirpro/internal/model"
)

// UserRepository handles user persistence
type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Preload("Outlet").Where("email = ? AND is_active = true", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.Preload("Outlet").Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) CreateOutlet(outlet *model.Outlet) error {
	return r.db.Create(outlet).Error
}

func (r *UserRepository) FindOutletByCode(code string) (*model.Outlet, error) {
	var outlet model.Outlet
	err := r.db.Where("code = ?", strings.ToUpper(strings.TrimSpace(code))).First(&outlet).Error
	if err != nil {
		return nil, err
	}
	return &outlet, nil
}

func (r *UserRepository) FindOutletByID(id uint) (*model.Outlet, error) {
	var outlet model.Outlet
	err := r.db.Where("id = ?", id).First(&outlet).Error
	if err != nil {
		return nil, err
	}
	return &outlet, nil
}

func (r *UserRepository) FindUsersByOutlet(outletID uint) ([]model.User, error) {
	var users []model.User
	err := r.db.Where("outlet_id = ?", outletID).Order("created_at desc").Find(&users).Error
	return users, err
}

// CategoryRepository handles categories
type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) FindAll(outletID uint) ([]model.Category, error) {
	var categories []model.Category
	err := r.db.Where("outlet_id = ?", outletID).Order("name asc").Find(&categories).Error
	return categories, err
}

func (r *CategoryRepository) FindByID(outletID, id uint) (*model.Category, error) {
	var category model.Category
	err := r.db.Where("outlet_id = ? AND id = ?", outletID, id).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepository) Create(cat *model.Category) error {
	return r.db.Create(cat).Error
}

func (r *CategoryRepository) Update(cat *model.Category) error {
	return r.db.Save(cat).Error
}

func (r *CategoryRepository) Delete(outletID, id uint) error {
	return r.db.Where("outlet_id = ? AND id = ?", outletID, id).Delete(&model.Category{}).Error
}

// ProductRepository handles products
type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) FindAll(outletID uint, filter dto.QueryFilter) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64

	query := r.db.Model(&model.Product{}).Preload("Category").Where("outlet_id = ?", outletID)

	if filter.Category > 0 {
		query = query.Where("category_id = ?", filter.Category)
	}

	if filter.Search != "" {
		searchTerm := "%" + strings.ToLower(filter.Search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(sku) LIKE ? OR LOWER(barcode) LIKE ?", searchTerm, searchTerm, searchTerm)
	}

	query.Count(&total)

	page := filter.Page
	if page <= 0 {
		page = 1
	}
	limit := filter.Limit
	if limit <= 0 {
		limit = 20
	}
	offset := (page - 1) * limit

	sortBy := "name"
	if filter.SortBy != "" {
		sortBy = filter.SortBy
	}
	sortOrder := "asc"
	if strings.ToLower(filter.SortOrder) == "desc" {
		sortOrder = "desc"
	}

	err := query.Order(fmt.Sprintf("%s %s", sortBy, sortOrder)).Offset(offset).Limit(limit).Find(&products).Error
	return products, total, err
}

func (r *ProductRepository) FindByID(outletID, id uint) (*model.Product, error) {
	var product model.Product
	err := r.db.Preload("Category").Where("outlet_id = ? AND id = ?", outletID, id).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) FindByBarcode(outletID uint, barcode string) (*model.Product, error) {
	var product model.Product
	err := r.db.Preload("Category").Where("outlet_id = ? AND (barcode = ? OR sku = ?)", outletID, barcode, barcode).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) Create(product *model.Product) error {
	return r.db.Create(product).Error
}

func (r *ProductRepository) Update(product *model.Product) error {
	return r.db.Save(product).Error
}

func (r *ProductRepository) Delete(outletID, id uint) error {
	return r.db.Where("outlet_id = ? AND id = ?", outletID, id).Delete(&model.Product{}).Error
}

func (r *ProductRepository) DecreaseStockTx(tx *gorm.DB, productID uint, qty int) error {
	return tx.Model(&model.Product{}).Where("id = ? AND stock >= ?", productID, qty).
		Update("stock", gorm.Expr("stock - ?", qty)).Error
}

func (r *ProductRepository) RestoreStockTx(tx *gorm.DB, productID uint, qty int) error {
	return tx.Model(&model.Product{}).Where("id = ?", productID).
		Update("stock", gorm.Expr("stock + ?", qty)).Error
}

// CustomerRepository handles customers
type CustomerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

func (r *CustomerRepository) FindAll(outletID uint, search string) ([]model.Customer, error) {
	var customers []model.Customer
	query := r.db.Where("outlet_id = ?", outletID)
	if search != "" {
		searchTerm := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR phone LIKE ? OR LOWER(member_code) LIKE ?", searchTerm, searchTerm, searchTerm)
	}
	err := query.Order("name asc").Find(&customers).Error
	return customers, err
}

func (r *CustomerRepository) FindByMemberCode(outletID uint, code string) (*model.Customer, error) {
	var customer model.Customer
	cleanCode := strings.TrimSpace(code)
	err := r.db.Where("outlet_id = ? AND (member_code = ? OR phone = ? OR LOWER(member_code) = ?)", outletID, cleanCode, cleanCode, strings.ToLower(cleanCode)).First(&customer).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *CustomerRepository) FindByID(outletID, id uint) (*model.Customer, error) {
	var customer model.Customer
	err := r.db.Where("outlet_id = ? AND id = ?", outletID, id).First(&customer).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *CustomerRepository) Create(cust *model.Customer) error {
	return r.db.Create(cust).Error
}

func (r *CustomerRepository) Update(cust *model.Customer) error {
	return r.db.Save(cust).Error
}

func (r *CustomerRepository) Delete(outletID, id uint) error {
	return r.db.Where("outlet_id = ? AND id = ?", outletID, id).Delete(&model.Customer{}).Error
}

func (r *CustomerRepository) AddPointsTx(tx *gorm.DB, customerID uint, points int) error {
	return tx.Model(&model.Customer{}).Where("id = ?", customerID).
		Update("points", gorm.Expr("points + ?", points)).Error
}

// SupplierRepository handles suppliers
type SupplierRepository struct {
	db *gorm.DB
}

func NewSupplierRepository(db *gorm.DB) *SupplierRepository {
	return &SupplierRepository{db: db}
}

func (r *SupplierRepository) FindAll(outletID uint, search string) ([]model.Supplier, error) {
	var suppliers []model.Supplier
	query := r.db.Where("outlet_id = ?", outletID)
	if search != "" {
		searchTerm := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR phone LIKE ? OR LOWER(contact_person) LIKE ?", searchTerm, searchTerm, searchTerm)
	}
	err := query.Order("name asc").Find(&suppliers).Error
	return suppliers, err
}

func (r *SupplierRepository) FindByID(outletID, id uint) (*model.Supplier, error) {
	var supplier model.Supplier
	err := r.db.Where("outlet_id = ? AND id = ?", outletID, id).First(&supplier).Error
	if err != nil {
		return nil, err
	}
	return &supplier, nil
}

func (r *SupplierRepository) Create(supplier *model.Supplier) error {
	return r.db.Create(supplier).Error
}

func (r *SupplierRepository) Update(supplier *model.Supplier) error {
	return r.db.Save(supplier).Error
}

func (r *SupplierRepository) Delete(outletID, id uint) error {
	return r.db.Where("outlet_id = ? AND id = ?", outletID, id).Delete(&model.Supplier{}).Error
}

// TransactionRepository handles transactions & POS
type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) GetDB() *gorm.DB {
	return r.db
}

func (r *TransactionRepository) FindAll(outletID uint, filter dto.QueryFilter, currentUserID uint, role string) ([]model.Transaction, int64, error) {
	var transactions []model.Transaction
	var total int64

	query := r.db.Model(&model.Transaction{}).
		Preload("Customer").
		Preload("User").
		Preload("Items").
		Where("outlet_id = ?", outletID)

	// ABAC: Cashier only views own transactions
	if role == "kasir" {
		query = query.Where("user_id = ?", currentUserID)
	}

	if filter.Search != "" {
		searchTerm := "%" + strings.ToLower(filter.Search) + "%"
		query = query.Where("LOWER(invoice_no) LIKE ?", searchTerm)
	}

	if filter.Status != "" {
		query = query.Where("payment_status = ?", filter.Status)
	}

	if filter.StartDate != "" {
		query = query.Where("created_at >= ?", filter.StartDate+" 00:00:00")
	}

	if filter.EndDate != "" {
		query = query.Where("created_at <= ?", filter.EndDate+" 23:59:59")
	}

	query.Count(&total)

	page := filter.Page
	if page <= 0 {
		page = 1
	}
	limit := filter.Limit
	if limit <= 0 {
		limit = 20
	}
	offset := (page - 1) * limit

	err := query.Order("created_at desc").Offset(offset).Limit(limit).Find(&transactions).Error
	return transactions, total, err
}

func (r *TransactionRepository) FindByID(outletID, id uint) (*model.Transaction, error) {
	var trans model.Transaction
	err := r.db.Preload("Customer").
		Preload("User").
		Preload("Items").
		Preload("Outlet").
		Where("outlet_id = ? AND id = ?", outletID, id).
		First(&trans).Error
	if err != nil {
		return nil, err
	}
	return &trans, nil
}

func (r *TransactionRepository) FindByInvoice(outletID uint, invoiceNo string) (*model.Transaction, error) {
	var trans model.Transaction
	err := r.db.Preload("Customer").
		Preload("User").
		Preload("Items").
		Preload("Outlet").
		Where("outlet_id = ? AND invoice_no = ?", outletID, invoiceNo).
		First(&trans).Error
	if err != nil {
		return nil, err
	}
	return &trans, nil
}

func (r *TransactionRepository) FindByMidtransOrderID(orderID string) (*model.Transaction, error) {
	var trans model.Transaction
	err := r.db.Preload("Items").Where("midtrans_order_id = ?", orderID).First(&trans).Error
	if err != nil {
		return nil, err
	}
	return &trans, nil
}

func (r *TransactionRepository) UpdatePaymentStatus(transID uint, status, midtransStatus string) error {
	return r.db.Model(&model.Transaction{}).Where("id = ?", transID).Updates(map[string]interface{}{
		"payment_status":              status,
		"midtrans_transaction_status": midtransStatus,
	}).Error
}

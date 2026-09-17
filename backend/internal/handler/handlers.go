package handler

import (
	"math"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"kasirpro/internal/dto"
	"kasirpro/internal/service"
	"kasirpro/internal/utils"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Format input data tidak valid")
	}

	if valErrors := utils.ValidateStruct(req); len(valErrors) > 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Validasi gagal", valErrors)
	}

	resp, err := h.authService.Register(req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	// Set HttpOnly cookie
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    resp.AccessToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   false, // set true in production HTTPS
		SameSite: "Lax",
		Path:     "/",
	})

	return utils.SuccessResponse(c, fiber.StatusCreated, "Registrasi toko dan akun berhasil", resp)
}

func (h *AuthHandler) RegisterStaff(c *fiber.Ctx) error {
	var req dto.StaffRegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Format input data tidak valid")
	}

	if valErrors := utils.ValidateStruct(req); len(valErrors) > 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Validasi gagal", valErrors)
	}

	resp, err := h.authService.RegisterStaff(req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	// Set HttpOnly cookie
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    resp.AccessToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
		Path:     "/",
	})

	return utils.SuccessResponse(c, fiber.StatusCreated, "Registrasi kasir/karyawan berhasil", resp)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Format input data tidak valid")
	}

	if valErrors := utils.ValidateStruct(req); len(valErrors) > 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Validasi gagal", valErrors)
	}

	resp, err := h.authService.Login(req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, err.Error())
	}

	// Set HttpOnly cookie
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    resp.AccessToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
		Path:     "/",
	})

	return utils.SuccessResponse(c, fiber.StatusOK, "Login berhasil", resp)
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HTTPOnly: true,
		SameSite: "Lax",
		Path:     "/",
	})
	return utils.SuccessResponse(c, fiber.StatusOK, "Logout berhasil", nil)
}

func (h *AuthHandler) GetMe(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	user, err := h.authService.GetProfile(userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Pengguna tidak ditemukan")
	}
	return utils.SuccessResponse(c, fiber.StatusOK, "Profil pengguna", user)
}

// CategoryHandler
type CategoryHandler struct {
	service *service.CategoryService
}

func NewCategoryHandler(s *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: s}
}

func (h *CategoryHandler) GetAll(c *fiber.Ctx) error {
	outletID := c.Locals("outlet_id").(uint)
	categories, err := h.service.GetAll(outletID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Gagal mengambil data kategori")
	}
	return utils.SuccessResponse(c, fiber.StatusOK, "Daftar kategori", categories)
}

func (h *CategoryHandler) Create(c *fiber.Ctx) error {
	outletID := c.Locals("outlet_id").(uint)
	var req dto.CreateCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Format input data tidak valid")
	}
	if valErrors := utils.ValidateStruct(req); len(valErrors) > 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Validasi gagal", valErrors)
	}

	cat, err := h.service.Create(outletID, req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SuccessResponse(c, fiber.StatusCreated, "Kategori berhasil ditambahkan", cat)
}

func (h *CategoryHandler) Update(c *fiber.Ctx) error {
	outletID := c.Locals("outlet_id").(uint)
	id, _ := strconv.Atoi(c.Params("id"))
	var req dto.UpdateCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Format input data tidak valid")
	}
	if valErrors := utils.ValidateStruct(req); len(valErrors) > 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Validasi gagal", valErrors)
	}

	cat, err := h.service.Update(outletID, uint(id), req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SuccessResponse(c, fiber.StatusOK, "Kategori berhasil diperbarui", cat)
}

func (h *CategoryHandler) Delete(c *fiber.Ctx) error {
	outletID := c.Locals("outlet_id").(uint)
	id, _ := strconv.Atoi(c.Params("id"))
	if err := h.service.Delete(outletID, uint(id)); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Gagal menghapus kategori")
	}
	return utils.SuccessResponse(c, fiber.StatusOK, "Kategori berhasil dihapus", nil)
}

// ProductHandler
type ProductHandler struct {
	service *service.ProductService
}

func NewProductHandler(s *service.ProductService) *ProductHandler {
	return &ProductHandler{service: s}
}

func (h *ProductHandler) GetAll(c *fiber.Ctx) error {
	outletID := c.Locals("outlet_id").(uint)
	var filter dto.QueryFilter
	if err := c.QueryParser(&filter); err != nil {
		filter = dto.QueryFilter{Page: 1, Limit: 20}
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 20
	}

	products, total, err := h.service.GetAll(outletID, filter)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Gagal memuat produk")
	}

	totalPages := int(math.Ceil(float64(total) / float64(filter.Limit)))
	meta := utils.PaginationMeta{
		CurrentPage: filter.Page,
		PerPage:     filter.Limit,
		Total:       total,
		TotalPages:  totalPages,
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Daftar produk", products, meta)
}

func (h *ProductHandler) GetByID(c *fiber.Ctx) error {
	outletID := c.Locals("outlet_id").(uint)
	id, _ := strconv.Atoi(c.Params("id"))
	product, err := h.service.GetByID(outletID, uint(id))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Produk tidak ditemukan")
	}
	return utils.SuccessResponse(c, fiber.StatusOK, "Detail produk", product)
}

func (h *ProductHandler) GetByBarcode(c *fiber.Ctx) error {
	outletID := c.Locals("outlet_id").(uint)
	barcode := c.Params("barcode")
	product, err := h.service.GetByBarcode(outletID, barcode)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Produk dengan barcode tersebut tidak ditemukan")
	}
	return utils.SuccessResponse(c, fiber.StatusOK, "Produk ditemukan", product)
}

func (h *ProductHandler) Create(c *fiber.Ctx) error {
	outletID := c.Locals("outlet_id").(uint)
	var req dto.CreateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Format input produk tidak valid")
	}
	if valErrors := utils.ValidateStruct(req); len(valErrors) > 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Validasi produk gagal", valErrors)
	}

	product, err := h.service.Create(outletID, req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SuccessResponse(c, fiber.StatusCreated, "Produk berhasil ditambahkan", product)
}

func (h *ProductHandler) Update(c *fiber.Ctx) error {
	outletID := c.Locals("outlet_id").(uint)
	id, _ := strconv.Atoi(c.Params("id"))
	var req dto.UpdateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Format input produk tidak valid")
	}
	if valErrors := utils.ValidateStruct(req); len(valErrors) > 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Validasi produk gagal", valErrors)
	}

	product, err := h.service.Update(outletID, uint(id), req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SuccessResponse(c, fiber.StatusOK, "Produk berhasil diperbarui", product)
}

func (h *ProductHandler) Delete(c *fiber.Ctx) error {
	outletID := c.Locals("outlet_id").(uint)
	id, _ := strconv.Atoi(c.Params("id"))
	if err := h.service.Delete(outletID, uint(id)); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Gagal menghapus produk")
	}
	return utils.SuccessResponse(c, fiber.StatusOK, "Produk berhasil dihapus", nil)
}

// CustomerHandler
type CustomerHandler struct {
	service *service.CustomerService
}

func NewCustomerHandler(s *service.CustomerService) *CustomerHandler {
	return &CustomerHandler{service: s}
}

func (h *CustomerHandler) GetPublicOutletInfo(c *fiber.Ctx) error {
	code := c.Params("code")
	info, err := h.service.GetPublicOutletInfo(code)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, err.Error())
	}
	return utils.SuccessResponse(c, fiber.StatusOK, "Informasi toko", info)
}

func (h *CustomerHandler) RegisterPublicMember(c *fiber.Ctx) error {
	var req dto.PublicMemberRegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Format formulir pendaftaran member tidak valid")
	}
	if valErrors := utils.ValidateStruct(req); len(valErrors) > 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Validasi gagal", valErrors)
	}

	cust, outlet, err := h.service.RegisterPublicMember(req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, "Pendaftaran member berhasil! Anda mendapatkan 10 poin bonus selamat datang.", fiber.Map{
		"customer": cust,
		"outlet":   outlet,
	})
}

func (h *CustomerHandler) GetAll(c *fiber.Ctx) error {
	outletID := c.Locals("outlet_id").(uint)
	search := c.Query("search")
	customers, err := h.service.GetAll(outletID, search)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Gagal memuat pelanggan")
	}
	return utils.SuccessResponse(c, fiber.StatusOK, "Daftar pelanggan", customers)
}

func (h *CustomerHandler) LookupMember(c *fiber.Ctx) error {
	outletID := c.Locals("outlet_id").(uint)
	code := c.Params("code")
	customer, err := h.service.FindByMemberCode(outletID, code)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Member tidak ditemukan")
	}
	return utils.SuccessResponse(c, fiber.StatusOK, "Data member ditemukan", customer)
}

func (h *CustomerHandler) Create(c *fiber.Ctx) error {
	outletID := c.Locals("outlet_id").(uint)
	var req dto.CreateCustomerRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Format data pelanggan tidak valid")
	}
	if valErrors := utils.ValidateStruct(req); len(valErrors) > 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Validasi gagal", valErrors)
	}

	cust, err := h.service.Create(outletID, req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SuccessResponse(c, fiber.StatusCreated, "Pelanggan berhasil didaftarkan", cust)
}

func (h *CustomerHandler) Update(c *fiber.Ctx) error {
	outletID := c.Locals("outlet_id").(uint)
	id, _ := strconv.Atoi(c.Params("id"))
	var req dto.UpdateCustomerRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Format data tidak valid")
	}
	cust, err := h.service.Update(outletID, uint(id), req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SuccessResponse(c, fiber.StatusOK, "Pelanggan berhasil diperbarui", cust)
}

func (h *CustomerHandler) Delete(c *fiber.Ctx) error {
	outletID := c.Locals("outlet_id").(uint)
	id, _ := strconv.Atoi(c.Params("id"))
	if err := h.service.Delete(outletID, uint(id)); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Gagal menghapus pelanggan")
	}
	return utils.SuccessResponse(c, fiber.StatusOK, "Pelanggan berhasil dihapus", nil)
}

// SupplierHandler
type SupplierHandler struct {
	service *service.SupplierService
}

func NewSupplierHandler(s *service.SupplierService) *SupplierHandler {
	return &SupplierHandler{service: s}
}

func (h *SupplierHandler) GetAll(c *fiber.Ctx) error {
	outletID := c.Locals("outlet_id").(uint)
	search := c.Query("search")
	suppliers, err := h.service.GetAll(outletID, search)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Gagal memuat supplier")
	}
	return utils.SuccessResponse(c, fiber.StatusOK, "Daftar supplier", suppliers)
}

func (h *SupplierHandler) Create(c *fiber.Ctx) error {
	outletID := c.Locals("outlet_id").(uint)
	var req dto.CreateSupplierRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Format input supplier tidak valid")
	}
	if valErrors := utils.ValidateStruct(req); len(valErrors) > 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Validasi gagal", valErrors)
	}

	supplier, err := h.service.Create(outletID, req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SuccessResponse(c, fiber.StatusCreated, "Supplier berhasil ditambahkan", supplier)
}

func (h *SupplierHandler) Update(c *fiber.Ctx) error {
	outletID := c.Locals("outlet_id").(uint)
	id, _ := strconv.Atoi(c.Params("id"))
	var req dto.UpdateSupplierRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Format input data tidak valid")
	}
	supplier, err := h.service.Update(outletID, uint(id), req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SuccessResponse(c, fiber.StatusOK, "Supplier berhasil diperbarui", supplier)
}

func (h *SupplierHandler) Delete(c *fiber.Ctx) error {
	outletID := c.Locals("outlet_id").(uint)
	id, _ := strconv.Atoi(c.Params("id"))
	if err := h.service.Delete(outletID, uint(id)); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Gagal menghapus supplier")
	}
	return utils.SuccessResponse(c, fiber.StatusOK, "Supplier berhasil dihapus", nil)
}

// TransactionHandler
type TransactionHandler struct {
	transService    *service.TransactionService
	midtransService *service.MidtransPaymentService
}

func NewTransactionHandler(transService *service.TransactionService, midtransService *service.MidtransPaymentService) *TransactionHandler {
	return &TransactionHandler{
		transService:    transService,
		midtransService: midtransService,
	}
}

func (h *TransactionHandler) GetAll(c *fiber.Ctx) error {
	outletID := c.Locals("outlet_id").(uint)
	userID := c.Locals("user_id").(uint)
	role := c.Locals("role").(string)

	var filter dto.QueryFilter
	_ = c.QueryParser(&filter)
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 20
	}

	transactions, total, err := h.transService.GetAll(outletID, filter, userID, role)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Gagal memuat transaksi")
	}

	totalPages := int(math.Ceil(float64(total) / float64(filter.Limit)))
	meta := utils.PaginationMeta{
		CurrentPage: filter.Page,
		PerPage:     filter.Limit,
		Total:       total,
		TotalPages:  totalPages,
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Daftar transaksi", transactions, meta)
}

func (h *TransactionHandler) GetByID(c *fiber.Ctx) error {
	outletID := c.Locals("outlet_id").(uint)
	id, _ := strconv.Atoi(c.Params("id"))
	trans, err := h.transService.GetByID(outletID, uint(id))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Transaksi tidak ditemukan")
	}
	return utils.SuccessResponse(c, fiber.StatusOK, "Detail transaksi", trans)
}

func (h *TransactionHandler) GetByInvoice(c *fiber.Ctx) error {
	outletID := c.Locals("outlet_id").(uint)
	invoiceNo := c.Params("invoice")
	trans, err := h.transService.GetByInvoice(outletID, invoiceNo)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Invoice transaksi tidak ditemukan")
	}
	return utils.SuccessResponse(c, fiber.StatusOK, "Detail transaksi", trans)
}

func (h *TransactionHandler) Checkout(c *fiber.Ctx) error {
	outletID := c.Locals("outlet_id").(uint)
	userID := c.Locals("user_id").(uint)

	var req dto.CheckoutRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Format payload checkout tidak valid")
	}
	if valErrors := utils.ValidateStruct(req); len(valErrors) > 0 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Validasi checkout gagal", valErrors)
	}

	trans, err := h.transService.Checkout(outletID, userID, req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	// If QRIS, ensure dynamic QR is ready
	if trans.PaymentMethod == "qris" {
		qrCodeURL, _ := h.midtransService.GenerateQRIS(trans)
		trans.MidtransQRCode = qrCodeURL
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, "Checkout transaksi berhasil", trans)
}

func (h *TransactionHandler) VoidTransaction(c *fiber.Ctx) error {
	outletID := c.Locals("outlet_id").(uint)
	id, _ := strconv.Atoi(c.Params("id"))

	var body struct {
		Reason string `json:"reason"`
	}
	_ = c.BodyParser(&body)
	if body.Reason == "" {
		body.Reason = "Dibatalkan oleh manajemen"
	}

	if err := h.transService.VoidTransaction(outletID, uint(id), body.Reason); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Transaksi berhasil dibatalkan dan stok dikembalikan", nil)
}

func (h *TransactionHandler) ConfirmQRIS(c *fiber.Ctx) error {
	outletID := c.Locals("outlet_id").(uint)
	invoiceNo := c.Params("invoice")

	trans, err := h.midtransService.ConfirmQRISPayment(outletID, invoiceNo)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "Pembayaran QRIS Midtrans berhasil diverifikasi", trans)
}

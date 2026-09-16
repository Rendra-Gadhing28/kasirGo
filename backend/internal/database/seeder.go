package database

import (
	"log"
	"time"

	"gorm.io/gorm"
	"kasirpro/internal/model"
	"kasirpro/internal/utils"
)

func SeedData(db *gorm.DB) {
	// 1. Seed Default Outlet
	var count int64
	db.Model(&model.Outlet{}).Count(&count)
	if count > 0 {
		var emptyOutlets []model.Outlet
		db.Where("code = '' OR code IS NULL").Find(&emptyOutlets)
		for _, o := range emptyOutlets {
			outletCode := "KASIR-" + utils.SanitizeString(string(rune(o.ID)))
			if o.ID == 1 {
				outletCode = "KASIR-DEMO-001"
			}
			db.Model(&model.Outlet{}).Where("id = ?", o.ID).Update("code", outletCode)
		}
		log.Println("Database telah memiliki data. Outlet codes verified.")
		return
	}

	outlet := model.Outlet{
		Code:    "KASIR-DEMO-001",
		Name:    "KasirPro Modern Market",
		Address: "Jl. Sudirman No. 123, Jakarta Selatan",
		Phone:   "081234567890",
		Logo:    "",
	}
	if err := db.Create(&outlet).Error; err != nil {
		log.Printf("Gagal membuat outlet seeder: %v", err)
		return
	}

	// 2. Seed Users (Owner, Admin, Kasir)
	passwordHash, _ := utils.HashPassword("password123")

	users := []model.User{
		{
			OutletID: outlet.ID,
			Name:     "Budi Santoso (Owner)",
			Email:    "owner@kasirpro.id",
			Password: passwordHash,
			Role:     "owner",
			IsActive: true,
		},
		{
			OutletID: outlet.ID,
			Name:     "Siti Aminah (Admin)",
			Email:    "admin@kasirpro.id",
			Password: passwordHash,
			Role:     "admin",
			IsActive: true,
		},
		{
			OutletID: outlet.ID,
			Name:     "Rian Kasir (Kasir 1)",
			Email:    "kasir@kasirpro.id",
			Password: passwordHash,
			Role:     "kasir",
			IsActive: true,
		},
	}
	db.Create(&users)

	// 3. Seed Categories
	categories := []model.Category{
		{OutletID: outlet.ID, Name: "Makanan & Minuman", Description: "Kebutuhan makanan harian dan minuman kemasan"},
		{OutletID: outlet.ID, Name: "Sembako", Description: "Sembilan bahan pokok rumah tangga"},
		{OutletID: outlet.ID, Name: "Snack & Cemilan", Description: "Biskuit, keripik, dan camilan ringan"},
		{OutletID: outlet.ID, Name: "Perawatan & Mandi", Description: "Sabun, sampo, pasta gigi, dan kosmetik"},
		{OutletID: outlet.ID, Name: "Kebutuhan Rumah", Description: "Detergen, pembersih lantai, dan tisu"},
	}
	db.Create(&categories)

	// 4. Seed Products
	products := []model.Product{
		{
			OutletID:   outlet.ID,
			CategoryID: categories[0].ID,
			Name:       "Aqua Air Mineral 600ml",
			SKU:        "AQ-600",
			Barcode:    "8992741000109",
			BuyPrice:   2800,
			SellPrice:  3500,
			Stock:      120,
			Unit:       "pcs",
			IsActive:   true,
		},
		{
			OutletID:   outlet.ID,
			CategoryID: categories[0].ID,
			Name:       "Teh Botol Sosro Kotak 250ml",
			SKU:        "TBS-250",
			Barcode:    "8992745110019",
			BuyPrice:   3000,
			SellPrice:  4000,
			Stock:      80,
			Unit:       "pcs",
			IsActive:   true,
		},
		{
			OutletID:   outlet.ID,
			CategoryID: categories[1].ID,
			Name:       "Minyak Goreng Bimoli 2 Liter",
			SKU:        "BML-2L",
			Barcode:    "8998866100201",
			BuyPrice:   33000,
			SellPrice:  38000,
			Stock:      45,
			Unit:       "pcs",
			IsActive:   true,
		},
		{
			OutletID:   outlet.ID,
			CategoryID: categories[1].ID,
			Name:       "Beras Pandan Wangi Premium 5kg",
			SKU:        "BRS-5KG",
			Barcode:    "8991001550211",
			BuyPrice:   72000,
			SellPrice:  85000,
			Stock:      30,
			Unit:       "karung",
			IsActive:   true,
		},
		{
			OutletID:   outlet.ID,
			CategoryID: categories[1].ID,
			Name:       "Gula Pasir Gulaku Putih 1kg",
			SKU:        "GLK-1KG",
			Barcode:    "8993005110111",
			BuyPrice:   15000,
			SellPrice:  18000,
			Stock:      60,
			Unit:       "pcs",
			IsActive:   true,
		},
		{
			OutletID:   outlet.ID,
			CategoryID: categories[2].ID,
			Name:       "Indomie Goreng Original 85g",
			SKU:        "IDM-GOR",
			Barcode:    "8998866200119",
			BuyPrice:   2700,
			SellPrice:  3500,
			Stock:      250,
			Unit:       "pcs",
			IsActive:   true,
		},
		{
			OutletID:   outlet.ID,
			CategoryID: categories[2].ID,
			Name:       "Chitato Sapi Panggang 68g",
			SKU:        "CTT-68G",
			Barcode:    "8992688001018",
			BuyPrice:   9500,
			SellPrice:  12500,
			Stock:      50,
			Unit:       "pcs",
			IsActive:   true,
		},
		{
			OutletID:   outlet.ID,
			CategoryID: categories[3].ID,
			Name:       "Sabun Mandi Lifebuoy Total 10 85g",
			SKU:        "LFB-85G",
			Barcode:    "8999999052021",
			BuyPrice:   4000,
			SellPrice:  5500,
			Stock:      90,
			Unit:       "pcs",
			IsActive:   true,
		},
		{
			OutletID:   outlet.ID,
			CategoryID: categories[3].ID,
			Name:       "Pasta Gigi Pepsodent Pencegah Gigi Berlubang 190g",
			SKU:        "PEP-190G",
			Barcode:    "8999999512345",
			BuyPrice:   14500,
			SellPrice:  18500,
			Stock:      40,
			Unit:       "pcs",
			IsActive:   true,
		},
		{
			OutletID:   outlet.ID,
			CategoryID: categories[4].ID,
			Name:       "Deterjen Rinso Molto Cair 750ml",
			SKU:        "RNS-750ML",
			Barcode:    "8999999712099",
			BuyPrice:   19000,
			SellPrice:  24000,
			Stock:      35,
			Unit:       "pcs",
			IsActive:   true,
		},
	}
	db.Create(&products)

	// 5. Seed Customers
	customers := []model.Customer{
		{
			OutletID: outlet.ID,
			Name:     "Pelanggan Umum (Walk-in)",
			Phone:    "080000000000",
			Email:    "guest@kasirpro.id",
			Address:  "Di Tempat",
			Points:   0,
		},
		{
			OutletID: outlet.ID,
			Name:     "Dewi Lestari",
			Phone:    "081298765432",
			Email:    "dewi@gmail.com",
			Address:  "Komplek Melati Indah No. 12",
			Points:   15,
		},
		{
			OutletID: outlet.ID,
			Name:     "Ahmad Fauzi",
			Phone:    "085712349876",
			Email:    "fauzi@yahoo.com",
			Address:  "Jl. Kebon Jeruk Barat No. 5",
			Points:   40,
		},
	}
	db.Create(&customers)

	// 6. Seed Suppliers
	suppliers := []model.Supplier{
		{
			OutletID:      outlet.ID,
			Name:          "PT Indomarco Adi Prima",
			Phone:         "021-5551234",
			Email:         "order@indomarco.co.id",
			Address:       "Kawasan Industri Pulogadung, Jakarta Timur",
			ContactPerson: "Bpk. Hendra",
		},
		{
			OutletID:      outlet.ID,
			Name:          "CV Berkah Sembako Nusantara",
			Phone:         "021-7778899",
			Email:         "sales@berkahsembako.id",
			Address:       "Pasar Induk Kramat Jati Blok C No. 8",
			ContactPerson: "Ibu Rahma",
		},
	}
	db.Create(&suppliers)

	// 7. Seed Initial Transaction for Dashboard verification
	now := time.Now()
	trans := model.Transaction{
		OutletID:                  outlet.ID,
		InvoiceNo:                 "INV-" + now.Format("20060102") + "-0001",
		CustomerID:                &customers[1].ID,
		UserID:                    users[2].ID,
		Subtotal:                  48500,
		DiscountType:              "none",
		DiscountValue:             0,
		DiscountAmount:            0,
		TaxRate:                   11,
		TaxAmount:                 5335,
		TotalAmount:               53835,
		PaymentMethod:             "cash",
		PaymentStatus:             "paid",
		PaidAmount:                60000,
		ChangeAmount:              6165,
		Notes:                     "Transaksi seeder perdana",
		MidtransTransactionStatus: "settlement",
		Items: []model.TransactionItem{
			{
				ProductID:   products[0].ID,
				ProductName: products[0].Name,
				ProductSKU:  products[0].SKU,
				Qty:         2,
				Unit:        "pcs",
				Price:       3500,
				BuyPrice:    2800,
				Subtotal:    7000,
			},
			{
				ProductID:   products[2].ID,
				ProductName: products[2].Name,
				ProductSKU:  products[2].SKU,
				Qty:         1,
				Unit:        "pcs",
				Price:       38000,
				BuyPrice:    33000,
				Subtotal:    38000,
			},
			{
				ProductID:   products[5].ID,
				ProductName: products[5].Name,
				ProductSKU:  products[5].SKU,
				Qty:         1,
				Unit:        "pcs",
				Price:       3500,
				BuyPrice:    2700,
				Subtotal:    3500,
			},
		},
	}
	db.Create(&trans)

	log.Println("Seeding database berhasil! Akun demo:")
	log.Println("Owner : owner@kasirpro.id / password123")
	log.Println("Admin : admin@kasirpro.id / password123")
	log.Println("Kasir : kasir@kasirpro.id / password123")
}

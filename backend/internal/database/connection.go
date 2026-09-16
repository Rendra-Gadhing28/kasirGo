package database

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"kasirpro/internal/config"
	"kasirpro/internal/model"
)

var DB *gorm.DB

func parseDatabaseURL(rawURL string) string {
	if strings.HasPrefix(rawURL, "mysql://") {
		u, err := url.Parse(rawURL)
		if err == nil {
			user := u.User.Username()
			pass, _ := u.User.Password()
			host := u.Host
			dbName := strings.TrimPrefix(u.Path, "/")
			query := u.Query()
			if query.Get("charset") == "" {
				query.Set("charset", "utf8mb4")
			}
			if query.Get("parseTime") == "" {
				query.Set("parseTime", "True")
			}
			if query.Get("loc") == "" {
				query.Set("loc", "Local")
			}
			return fmt.Sprintf("%s:%s@tcp(%s)/%s?%s", user, pass, host, dbName, query.Encode())
		}
	}
	return rawURL
}

func InitDB() *gorm.DB {
	cfg := config.AppConfig

	var dsn string
	if cfg.DatabaseURL != "" {
		dsn = parseDatabaseURL(cfg.DatabaseURL)
	} else {
		tlsParam := ""
		if cfg.DBSSL == "true" || cfg.DBSSL == "1" || cfg.DBSSL == "require" || cfg.DBSSL == "skip-verify" {
			if cfg.DBSSL == "skip-verify" {
				tlsParam = "&tls=skip-verify"
			} else {
				tlsParam = "&tls=true"
			}
		}

		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local%s",
			cfg.DBUser,
			cfg.DBPassword,
			cfg.DBHost,
			cfg.DBPort,
			cfg.DBName,
			tlsParam,
		)
	}

	gormLogLevel := logger.Warn
	if cfg.AppEnv == "development" {
		gormLogLevel = logger.Info
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(gormLogLevel),
		// Parameterized queries are always enforced by GORM
	})
	if err != nil {
		log.Printf("Gagal terhubung ke MySQL database: %v", err)
		panic(fmt.Sprintf("Gagal terhubung ke MySQL database: %v", err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("Gagal mendapatkan database instance: %v", err)
		panic(fmt.Sprintf("Gagal mendapatkan database instance: %v", err))
	}

	// Connection Pool settings for high performance and stability
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Koneksi MySQL database berhasil diinisialisasi.")

	// Run Auto Migrations
	err = db.AutoMigrate(
		&model.Outlet{},
		&model.User{},
		&model.Category{},
		&model.Product{},
		&model.Customer{},
		&model.Supplier{},
		&model.Transaction{},
		&model.TransactionItem{},
		&model.ExpenseIncome{},
	)
	if err != nil {
		log.Fatalf("Gagal auto migrate schema database: %v", err)
	}
	log.Println("Auto-migration database schema selesai.")

	DB = db
	return db
}

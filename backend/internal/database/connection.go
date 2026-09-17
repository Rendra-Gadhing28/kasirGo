package database

import (
	"fmt"
	"log"
	"net/url"
	"os"
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

func InitDB() (*gorm.DB, error) {
	cfg := config.AppConfig

	if strings.Contains(cfg.DatabaseURL, "${{") {
		return nil, fmt.Errorf("DATABASE_URL belum ter-resolve oleh Railway: '%s'. Pastikan format tanpa spasi ${{MySQL.MYSQL_PRIVATE_URL}} atau salin langsung isi MYSQL_PRIVATE_URL dari service MySQL", cfg.DatabaseURL)
	}

	// In serverless (Vercel), fail immediately if DB_HOST is not configured instead of hanging on 127.0.0.1
	isVercel := os.Getenv("VERCEL") != ""
	if isVercel && cfg.DatabaseURL == "" && (cfg.DBHost == "127.0.0.1" || cfg.DBHost == "localhost" || cfg.DBHost == "") {
		return nil, fmt.Errorf("DB_HOST belum disetel di Vercel Environment Variables. Masukkan DB_HOST dari TiDB Cloud di Project Settings > Environment Variables")
	}

	var dsn string
	if cfg.DatabaseURL != "" {
		if u, err := url.Parse(cfg.DatabaseURL); err == nil {
			if h := u.Hostname(); h != "" {
				cfg.DBHost = h
			}
			if p := u.Port(); p != "" {
				cfg.DBPort = p
			}
		}
		dsn = parseDatabaseURL(cfg.DatabaseURL)
	} else {
		tlsParam := ""
		if cfg.DBSSL == "true" || cfg.DBSSL == "1" || cfg.DBSSL == "require" || cfg.DBSSL == "skip-verify" || strings.Contains(cfg.DBHost, "tidbcloud.com") {
			if cfg.DBSSL == "skip-verify" {
				tlsParam = "&tls=skip-verify"
			} else {
				tlsParam = "&tls=true"
			}
		}

		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=3s%s",
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
	})
	if err != nil {
		log.Printf("Gagal terhubung ke MySQL database: %v", err)
		return nil, fmt.Errorf("koneksi database ke %s:%s gagal: %w", cfg.DBHost, cfg.DBPort, err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("Gagal mendapatkan database instance: %v", err)
		return nil, fmt.Errorf("gagal sqlDB: %w", err)
	}

	sqlDB.SetMaxIdleConns(2)
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Koneksi MySQL database berhasil diinisialisasi.")

	// Run Auto Migrations only if tables not present yet to avoid serverless timeout
	if !db.Migrator().HasTable(&model.Outlet{}) {
		log.Println("Tabel belum ada, menjalankan auto-migration...")
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
			log.Printf("Peringatan auto-migration: %v", err)
		} else {
			log.Println("Auto-migration database schema selesai.")
		}
	}

	DB = db
	return db, nil
}

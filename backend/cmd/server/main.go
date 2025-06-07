package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"courtopia-reserve/backend/internal/config"
	"courtopia-reserve/backend/internal/database"
	"courtopia-reserve/backend/internal/handlers"
	"courtopia-reserve/backend/internal/repository"
)

func startScheduler(bookingRepo *repository.BookingRepository, userRepo *repository.UserRepository) {
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for range ticker.C {
			log.Println("Running email notification scheduler...")
			handlers.SendMail(bookingRepo, userRepo)
		}
	}()
}

func main() {
	// โหลดค่า config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// เชื่อมต่อกับฐานข้อมูล
	client, err := database.ConnectDB(cfg.MongoURI)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Fatalf("Error disconnecting from database: %v", err)
		}
	}()

	// สร้างฐานข้อมูลและ repositories
	db := client.Database("courtopia")
	userRepo := repository.NewUserRepository(db)
	courtRepo := repository.NewCourtRepository(db)
	bookingRepo := repository.NewBookingRepository(db)

	// เริ่ม Scheduler สำหรับการส่งอีเมลแจ้งเตือน
	startScheduler(bookingRepo, userRepo)

	// Set Gin mode based on environment
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// สร้าง Gin engine
	r := gin.Default()

	r.Static("/uploads", "./uploads")

	// ตั้งค่า CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle OPTIONS method for CORS preflight
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	})

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// สร้าง handler และลงทะเบียน routes
	h := handlers.NewHandler(db, userRepo, courtRepo, bookingRepo, cfg.JWTSecret)
	h.RegisterRoutes(r)

	h = handlers.NewHandler(db, userRepo, courtRepo, bookingRepo, cfg.JWTSecret)
	r.POST("/trigger-email-notifications", h.TriggerEmailNotifications)

	// เริ่มต้น server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: r,
	}

	// รันใน goroutine ที่แยกกันเพื่อไม่ให้บล็อก
	go func() {
		log.Printf("Server starting on port %d", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// รอสัญญาณปิดแอพ
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// สร้าง timeout context สำหรับการปิด
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// ปิด server
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}
	log.Println("Server stopped")
}

package main

import (
	"fmt"
	"gotickets/internal/config"
	"gotickets/internal/domain/illustration"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := config.LoadEnv()
	db, err := gorm.Open(postgres.Open(cfg.Dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}
	var count int64
	db.Model(&illustration.Illustration{}).Count(&count)
	fmt.Println("Count:", count)
	fmt.Println("Migrating...")
	err = db.AutoMigrate(&illustration.Illustration{})
	fmt.Println("Migrate error:", err)
}

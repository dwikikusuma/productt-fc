package resource

import (
	"fmt"
	"gorm.io/driver/postgres"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
	_ "gorm.io/gorm"
	"gorm.io/gorm/logger"
	_ "gorm.io/gorm/logger"
	_ "honnef.co/go/tools/config"
	"log"
	"product_commerce/config"
)

func InitDB(cfg *config.Config) *gorm.DB {
	fmt.Print(cfg)
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.Name)

	fmt.Print(dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("failed to connect with db %s", err)
	}

	log.Print("connected with db")
	return db
}

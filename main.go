package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	log.Fatal(run())
}

func run() error {
	db, err := initDB()
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&Message{})
	if err != nil {
		return err
	}
	app := fiber.New()

	app.Post("/messages", func(c *fiber.Ctx) error {
		m := &Message{}
		fmt.Println("sturct111")
		fmt.Println(m)
		if err := c.BodyParser(m); err != nil {
			return err
		}
		fmt.Println(m)
		fmt.Println(db.Model(&Message{}).Save(m).Error)
		return nil
	})

	return app.Listen(":3000")
}

func initDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
}

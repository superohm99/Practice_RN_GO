package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Shortlylink struct {
	gorm.Model
	OriginalURL string `gorm: "unique"`
	ShortURL    string `gorm: "unique"`
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .")
	}

	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := dbUsername + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?parseTime=true"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMirate(&Shortlylink{})

	r := gin.Default()

	r.Use(cors.Default())

	r.POST("/shorten", func(ctx *gin.Context) {

		var data struct {
			URL string `json: "url" binding: "required"`
		}

		if err := c.ShouldBindJSON(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var link Shortlylink
		result := db.Where("original_url = ?", data.URL).First(&link)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {

			}
		}
	})
}

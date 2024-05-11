package main

import (
	"net/http"
	"shrotly-link/controllers"
	"shrotly-link/initializers"

	"math/rand"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Shortlylink struct {
	gorm.Model
	OriginalURL string `gorm: "unique"`
	ShortURL    string `gorm: "unique"`
}

func init() {
	initializers.ConnectToDB()
}

func main() {

	r := gin.Default()

	r.Use(cors.Default())

	r.POST("/shorten", func(c *gin.Context) {

		var data struct {
			URL string `json: "url" binding: "required"`
		}

		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var link Shortlylink
		result := initializers.DB.Where("original_url = ?", data.URL).First(&link)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				shortURL := generateShrotURL()
				link = Shortlylink{OriginalURL: data.URL, ShortURL: shortURL}
				result = initializers.DB.Create(&link)
				if result.Error != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
					return
				}
			}
		}

		c.JSON(http.StatusOK, gin.H{"short_url": link.ShortURL})
	})

	r.GET("/:shortURL", func(c *gin.Context) {
		shortURL := c.Param("shortURL")
		var link Shortlylink
		result := initializers.DB.Where("short_url = ?", shortURL).Find(&link)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			}
			return
		}
		c.Redirect(http.StatusMovedPermanently, link.OriginalURL)
	})

	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, "test success")
		return
	})

	r.POST("/user", controllers.UsersCreate)

	r.Run(":8000")
}

func generateShrotURL() string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"
	rand.Seed(time.Now().UnixNano())

	var shortURL string
	for i := 0; i < 6; i++ {
		shortURL += string(chars[rand.Intn(len(chars))])
	}

	return shortURL
}

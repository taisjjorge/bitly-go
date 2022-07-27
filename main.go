package main

import (
	"math/rand"
	"net/http"
	"time"
	"strings"
	"github.com/gin-gonic/gin"
)

type shortUrlInput struct {
	Url string
}

var urlsBd map[string]string

const SHORT_URL_LENGTH = 6

func RandomString() string {
	charSet := "abcdefghijklmnopqrst"
	var output strings.Builder
	length := SHORT_URL_LENGTH
	for i := 0; i < length; i++ {
		random := rand.Intn(len(charSet))
		randomChar := charSet[random]
		output.WriteString(string(randomChar))
	}
	return output.String()
}

func createShortUrl(c *gin.Context){
	var newShortUrlInput shortUrlInput

	if err := c.ShouldBindJSON(&newShortUrlInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	newShortUrlKey := RandomString()
	urlsBd[newShortUrlKey] = newShortUrlInput.Url

	c.JSON(http.StatusCreated, gin.H{
		"shortUrlKey": newShortUrlKey,
	})
}

func parseShortUrl(c *gin.Context) {

	shortUrlKey := c.Param("shortUrlKey")

	if originalUrl, exists := urlsBd[shortUrlKey]; exists {
		c.Header("Location", originalUrl)
		c.Status(http.StatusMovedPermanently)
		return
	}

	c.Status(http.StatusNotFound)
}

func main() {

	urlsBd = make(map[string]string)
	rand.Seed(time.Now().Unix())

	router := gin.Default()

	shortUrlRouter := router.Group("/surl")
	{
		shortUrlRouter.POST("/", createShortUrl)
		shortUrlRouter.GET("/:shortUrlKey", parseShortUrl)
	}

	router.Run()
}
package controllers

import (
	"github.com/gin-gonic/gin"
)

func GetMovies() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.json(200, gin.H{
			"message": "GetMovies endpoint",
		})
	}
}

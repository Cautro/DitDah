package user

import (
	"context"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func GetMeHandler(u *UserUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()
		
		userId := c.GetInt("Id")

		output, err := u.GetMeUseCase(ctx, userId)
		if err != nil {
			c.JSON(500, gin.H{"error":"Server error"})
			return 
		}

		c.JSON(200, output)
	}
}

func GetUserById(u *UserUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		userIdStr := c.Param("Id")
		userId, err := strconv.Atoi(userIdStr)
		if err != nil {
			c.JSON(500, gin.H{"error":"Server error"})
			return
		}

		output, err := u.GetUserById(ctx, userId)
		if err != nil {
			c.JSON(500, gin.H{"error":"Server error"})
			return 
		}

		c.JSON(200, output)
	}
}
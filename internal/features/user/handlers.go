package user

import (
	"context"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func DeleteUserHandler(u *UserUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		queryUserIdStr := c.Query("Id")
		if queryUserIdStr == "" {
			c.JSON(400, gin.H{"error":"Missing user Id"})
			return
		}
		queryUserId, err := strconv.Atoi(queryUserIdStr)
		if err != nil {
			c.JSON(500, gin.H{"error":"Server error"})
			return
		}

		authorUserId := c.GetInt("Id")
		if queryUserId != authorUserId {
			c.JSON(403, gin.H{"error":"Forbidden"})
			return
		}

		err = u.DeleteUserUseCase(ctx, queryUserId)
		if err != nil {
			c.JSON(500, gin.H{"error":"Server error"})
			return 
		}

		c.JSON(200, gin.H{"status":"ok"})
	}
}

func GetAllUsersHandler(u *UserUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		output, err := u.GetAllUsersUseCase(ctx)
		if err != nil {
			c.JSON(500, gin.H{"error":"Server error"})
			return 
		}

		c.JSON(200, output)
	}
}

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

		userIdStr := c.Query("Id")
		if userIdStr == "" {
			c.JSON(400, gin.H{"error":"Missing user Id"})
			return
		}
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
package lesson

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetLessonByIdHandler(lessonUseCase *LessonUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		lessonIdStr := c.Param("id")
		lessonId, err := strconv.Atoi(lessonIdStr)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid lesson ID format"})
			return
		}

		lesson, err := lessonUseCase.GetLessonById(lessonId)
		if err != nil {
			c.JSON(500, gin.H{"error": "failed to retrieve lesson"})
			return
		}
		
		c.JSON(200, lesson)
	}
}

func GetLessonsHandler(lessonUseCase *LessonUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		lessons, err := lessonUseCase.lessonRepo.GetAllLessons()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, lessons)
	}
}

func AddLessonHandler(lessonUseCase *LessonUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var lesson LessonEntity
		if err := c.ShouldBindJSON(&lesson); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		userId := c.GetInt("id")
		if err := lessonUseCase.AddLesson(lesson, userId); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(201, lesson)
	}
}

func DeleteLessonHandler(lessonUseCase *LessonUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		lessonIdStr := c.Param("id") 
		lessonId, err := strconv.Atoi(lessonIdStr)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid lesson ID format"})
			return
		}

		userId := c.GetInt("id")
		if err := lessonUseCase.DeleteLesson(lessonId, userId); err != nil {
			if err.Error() == "User is not an admin" {
				c.JSON(403, gin.H{"error": "forbidden"})
				return
			}
			c.JSON(500, gin.H{"error": "failed to delete lesson"})
			return
		}
		
		c.JSON(200, gin.H{"status": "ok"})
	}
}
package router

import (
	"database/sql"
	utils "ditdah/pkg/middleware"

	"ditdah/internal/features/auth"
	"ditdah/internal/features/user"
	"ditdah/internal/features/lesson"

	"github.com/gin-gonic/gin"
)

type UseCases struct {
	Auth       *auth.AuthUseCase
	User       *user.UserUseCase
	Lesson     *lesson.LessonUseCase

	JWTSecret  string
	DB         *sql.DB
}

func NewRouter(u UseCases) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	registerPublicRoutes(u, router)
	registerAuthenticatedRoutes(u, router)

	return router
}

func registerPublicRoutes(u UseCases, r *gin.Engine) {
	r.POST("/login", auth.LoginHandler(u.Auth))
	r.POST("/refresh", auth.RefreshHandler(u.Auth))
	r.POST("/register", auth.RegisterHandler(u.Auth))
}

func registerAuthenticatedRoutes(u UseCases, r *gin.Engine) {
	api := r.Group("/api", utils.AuthMiddleware(u.JWTSecret))

	registerUserRoutes(u, api)
	registerLessonRoutes(u, api)
}

func registerUserRoutes(u UseCases, api *gin.RouterGroup) {
	api.GET("/me", user.GetMeHandler(u.User))
	api.GET("/users", user.GetAllUsersHandler(u.User))
	api.GET("/user", user.GetUserById(u.User)) // request query param: ?Id=...
	api.DELETE("/user/delete", user.DeleteUserHandler(u.User)) // request query param: ?Id=...
	
	api.POST("/logout", auth.LogoutHandler(u.User, u.Auth))
}

func registerLessonRoutes(u UseCases, api *gin.RouterGroup) {
	api.GET("/lessons", lesson.GetLessonsHandler(u.Lesson))
	api.POST("/lesson/add", lesson.AddLessonHandler(u.Lesson))
	api.DELETE("/lesson/delete/:id", lesson.DeleteLessonHandler(u.Lesson))
	api.GET("/lesson/:id", lesson.GetLessonByIdHandler(u.Lesson))
}
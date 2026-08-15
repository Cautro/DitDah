package router

import (
	"database/sql"
	utils "ditdah/pkg/middleware"

	"ditdah/internal/features/auth"
	"ditdah/internal/features/user"

	"github.com/gin-gonic/gin"
)

type UseCases struct {
	Auth       *auth.AuthUseCase
	User       *user.UserUseCase

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

	// tests
	r.GET("/users", user.GetAllUsersHandler(u.User))
}

func registerAuthenticatedRoutes(u UseCases, r *gin.Engine) {
	api := r.Group("/api", utils.AuthMiddleware(u.JWTSecret))

	registerUserRoutes(u, api)
}

func registerUserRoutes(u UseCases, api *gin.RouterGroup) {
	api.GET("/me", user.GetMeHandler(u.User))
	api.GET("/user", user.GetUserById(u.User)) // query param: ?Id=1
	api.POST("/logout", auth.LogoutHandler(u.User, u.Auth))
}
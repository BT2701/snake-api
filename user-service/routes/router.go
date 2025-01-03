package routes

import (
	"os"
	"snake_api/controllers"
	"snake_api/repositories"
	"snake_api/services"
	"snake_api/utils"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRouter(userCollection *mongo.Collection) *echo.Echo {
	redisClient := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_URI"), // Địa chỉ Redis
	})
	// Khởi tạo repository, service, và controller
	userRepo := repositories.NewUserRepository(userCollection)
	userService := services.NewUserService(userRepo, redisClient)
	userController := controllers.NewUserController(userService)

	e := echo.New()

	// Enable CORS
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("Access-Control-Allow-Origin", "*")
			return next(c)
		}
	})
	e.Use(utils.CorsMiddleware())

	// Routes
	api := e.Group("/api")
	{
		api.POST("/login", userController.Login)
		api.POST("/register", userController.SignUp)
		api.POST("/forgot", userController.ForgotPassword)
		api.POST("/reset", userController.ResetPassword)
		api.GET("/users", userController.GetAllUsers)
		api.DELETE("/users", userController.DeleteAllUsers)
		api.PUT("/logout", userController.Logout)
		api.PUT("/edit", userController.EditUser)
		api.GET("/user/:id", userController.GetByID)
		api.GET("/user", userController.FindUserByEmail)
		api.PUT("/avatar", userController.UpdateAvatar)
	}

	return e
}

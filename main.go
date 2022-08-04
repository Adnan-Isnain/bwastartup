package main

import (
	"bwastartup/auth"
	"bwastartup/campaign"
	"bwastartup/handler"
	"bwastartup/helper"
	"bwastartup/user"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	// Repository
	userRepository := user.NewRepository(db)
	campaignRepository := campaign.NewRepository(db)

	// Testing repository here

	// Service
	userService := user.NewService(userRepository)
	authService := auth.NewService()
	campaignService := campaign.NewService(campaignRepository)

	// Testing service here (using user service)!
	// Testing Input (If any)

	userHandler := handler.NewUserHandler(userService, authService)
	campaignHandler := handler.NewCampaignHandler(campaignService)

	router := gin.Default()
	router.Static("/images", "./images")
	api := router.Group("api/v1")

	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email_checkers", userHandler.CheckEmailAvailibility)
	api.POST("/avatars", authMiddleware(authService, userService), userHandler.UploadAvatar)
	api.GET("/campaigns", campaignHandler.GetCampaigns)
	api.GET("/campaign/:id", campaignHandler.GetCampaign)
	api.POST("campaigns", authMiddleware(authService, userService), campaignHandler.CreateCampaign)
	router.Run()
}

// Get header authorization : Bearer JWTToken
func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {
	// Return function as middleware

	return func(c *gin.Context) {
		// Check bearer
		authHeader := c.GetHeader("Authorization")
		if !strings.Contains(authHeader, "Bearer") {
			response := helper.APIResponse("Unathorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		// Get bearer
		tokenString := ""
		arrayToken := strings.Split(authHeader, " ")
		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}

		// Validate bearer token
		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			response := helper.APIResponse("Unathorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		// Extract data token
		payload, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			response := helper.APIResponse("Unathorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		// Get payload JWT Token
		userId := int(payload["user_id"].(float64))

		// Get ID from service
		user, err := userService.GetUserByID(userId)
		if err != nil {
			response := helper.APIResponse("Unathorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		// Set current user
		c.Set("currentUser", user)
	}

}

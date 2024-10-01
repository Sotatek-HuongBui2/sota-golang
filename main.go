package main

import (
	"fmt"
	_ "vtcanteen/docs"
	"vtcanteen/routers"
	"vtcanteen/services"
	"vtcanteen/utils"

	"github.com/joho/godotenv"
)

// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	fmt.Println("Start server")

	if err := godotenv.Load(".env"); err != nil {
		panic("Error load .env")
	}

	utils.ConnectDB()

	services.CheckExistRoleAndUser()

	g := routers.Create()
	g.Run()
	// g.Run("127.0.0.1:3000")
}

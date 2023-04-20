package main

import (
	"fmt"
	"log"
	model "online_fashion_shop/api/model/user"
	"online_fashion_shop/initializers"
)

func init() {
	config, err := initializers.LoadConfig(".")

	if err != nil {
		log.Fatal("🚀 Could not load environment variables", err)
	}

	initializers.ConnectDB(&config)
}

func main() {
	initializers.DB.AutoMigrate(&model.User{}, &model.UserVerify{})
	fmt.Println("👍 Migration complete")
}

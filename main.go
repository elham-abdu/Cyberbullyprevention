package main

import (
    "github.com/elham-abdu/cyberbullyprevention/config"
    "github.com/elham-abdu/cyberbullyprevention/models"
)

func main() {
    // Load env variables
    config.LoadEnv()

    // Connect to DB
    config.ConnectDB()

    // Auto-migrate tables
    config.DB.AutoMigrate(&models.User{}, &models.Post{})
}

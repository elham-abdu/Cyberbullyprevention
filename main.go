package main

import (
    "log"
    "net/http"

    "github.com/elham-abdu/cyberbullyprevention/config"
    "github.com/elham-abdu/cyberbullyprevention/models"
    "github.com/elham-abdu/cyberbullyprevention/handlers"
)

func main() {
   
    config.LoadEnv()

  
    config.ConnectDB()

    
    config.DB.AutoMigrate(&models.User{}, &models.Post{})

    
    http.HandleFunc("/register", handlers.Register)
    http.HandleFunc("/login", handlers.Login)

    
    log.Println("Server running on :8080")
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatal("Server failed:", err)
    }
}

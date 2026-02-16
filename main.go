package main

import (
    "log"
    "net/http"

    "github.com/elham-abdu/cyberbullyprevention/config"
    "github.com/elham-abdu/cyberbullyprevention/models"
    "github.com/elham-abdu/cyberbullyprevention/handlers"
    "github.com/elham-abdu/cyberbullyprevention/middleware"

)

func main() {
   
    config.LoadEnv()

  
    config.ConnectDB()

    
    config.DB.AutoMigrate(&models.User{}, &models.Post{})

    
    http.HandleFunc("/register", handlers.Register)
    http.HandleFunc("/login", handlers.Login)
    http.Handle("/me", middleware.JWTAuth(http.HandlerFunc(handlers.Me)))
    http.Handle("/me/posts", middleware.JWTAuth(http.HandlerFunc(handlers.GetMyPosts)))
    http.Handle("/me/posts/create", middleware.JWTAuth(http.HandlerFunc(handlers.CreatePost)))
    http.Handle("/me/posts/edit", middleware.JWTAuth(http.HandlerFunc(handlers.EditPost)))
    http.Handle("/me/posts/delete", middleware.JWTAuth(http.HandlerFunc(handlers.DeletePost)))
    http.Handle("/admin/dashboard", 
    middleware.JWTAuth(
        middleware.RoleAuth("admin")(
            http.HandlerFunc(handlers.AdminDashboard),
        ),
    ),
)

    
    log.Println("Server running on :8080")
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatal("Server failed:", err)
    }
}

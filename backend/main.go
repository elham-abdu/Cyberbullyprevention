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

    // Create a new serve mux
    mux := http.NewServeMux()

    // Public routes
    mux.HandleFunc("/register", handlers.Register)
    mux.HandleFunc("/login", handlers.Login)

    // Protected routes
    mux.Handle("/me", middleware.JWTAuth(http.HandlerFunc(handlers.Me)))
    mux.Handle("/me/posts", middleware.JWTAuth(http.HandlerFunc(handlers.GetMyPosts)))
    mux.Handle("/me/posts/create", middleware.JWTAuth(http.HandlerFunc(handlers.CreatePost)))
    mux.Handle("/me/posts/edit", middleware.JWTAuth(http.HandlerFunc(handlers.EditPost)))
    mux.Handle("/me/posts/delete", middleware.JWTAuth(http.HandlerFunc(handlers.DeletePost)))
    
    // Admin routes
    mux.Handle("/admin/dashboard", 
        middleware.JWTAuth(
            middleware.RoleAuth("admin")(
                http.HandlerFunc(handlers.AdminDashboard),
            ),
        ),
    )
    mux.Handle("/admin/flagged-posts",
        middleware.JWTAuth(
            middleware.RoleAuth("admin")(
                http.HandlerFunc(handlers.GetFlaggedPosts),
            ),
        ),
    )
    mux.Handle("/admin/posts/mark-safe",
        middleware.JWTAuth(
            middleware.RoleAuth("admin")(
                http.HandlerFunc(handlers.MarkPostSafe),
            ),
        ),
    )
    mux.Handle("/admin/posts/delete-flagged",
        middleware.JWTAuth(
            middleware.RoleAuth("admin")(
                http.HandlerFunc(handlers.DeleteFlaggedPost),
            ),
        ),
    )

    // Wrap the mux with CORS middleware
    handler := middleware.CorsMiddleware(mux)

    log.Println("Server running on :8080")
    err := http.ListenAndServe(":8080", handler)
    if err != nil {
        log.Fatal("Server failed:", err)
    }
}

package handlers
import (
	"net/http"
	"github.com/elham-abdu/cyberbullyprevention/config"
	"github.com/elham-abdu/cyberbullyprevention/models"
	"github.com/elham-abdu/cyberbullyprevention/utils"
	"encoding/json"
    "time"
    "github.com/golang-jwt/jwt/v5"
	
)
func Register(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    type RegisterInput struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    var input RegisterInput
    err := json.NewDecoder(r.Body).Decode(&input)
    if err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    if input.Email == "" || input.Password == "" {
        http.Error(w, "Email and password are required", http.StatusBadRequest)
        return
    }

    var existingUser models.User
    result := config.DB.Where("email = ?", input.Email).First(&existingUser)
    if result.RowsAffected > 0 {
        http.Error(w, "User already exists", http.StatusConflict)
        return
    }

    hashedPassword, err := utils.HashPassword(input.Password)
    if err != nil {
        http.Error(w, "Error hashing password", http.StatusInternalServerError)
        return
    }

    user := models.User{
        Email:        input.Email,
        PasswordHash: hashedPassword,
        Role:         "user",
    }

    config.DB.Create(&user)

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	type LoginInput struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var input LoginInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if input.Email == "" || input.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	var user models.User
	result := config.DB.Where("email = ?", input.Email).First(&user)
	if result.RowsAffected == 0 {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	if !utils.CheckPasswordHash(input.Password, user.PasswordHash) {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	
	expirationTime := time.Now().Add(24 * time.Hour)


	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     expirationTime.Unix(),
	}

	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	
	tokenString, err := token.SignedString([]byte(config.GetEnv("JWT_SECRET")))
	if err != nil {
		http.Error(w, "Could not create token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})
}
func Me(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value("user_id").(uint)
    role := r.Context().Value("role").(string)
    var user models.User
    result := config.DB.First(&user, userID)
    if result.Error != nil {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "user_id": user.ID,
        "email":   user.Email,
        "role":    role,
    })
}

type CreatePostInput struct {
    Content string `json:"content"`
}
func CreatePost(w http.ResponseWriter, r *http.Request) {

    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var input CreatePostInput
    err := json.NewDecoder(r.Body).Decode(&input)
    if err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    if input.Content == "" {
        http.Error(w, "Content is required", http.StatusBadRequest)
        return
    }


    userID := r.Context().Value("user_id").(uint)

    
    post := models.Post{
        UserID:  userID,
        Content: input.Content,
    }

    config.DB.Create(&post)

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(post)
}
func GetMyPosts(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value("user_id").(uint)

    var posts []models.Post
    result := config.DB.Where("user_id = ?", userID).Find(&posts)
    if result.Error != nil {
        http.Error(w, "Error fetching posts", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(posts)
}
func EditPost(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPut {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // 1. Get user_id from context
    userID := r.Context().Value("user_id").(uint)

    // 2. Decode the request body for post ID and new content
    type EditInput struct {
        PostID  uint   `json:"post_id"`
        Content string `json:"content"`
    }
    var input EditInput
    err := json.NewDecoder(r.Body).Decode(&input)
    if err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    // 3. Find the post and ensure it belongs to the user
    var post models.Post
    result := config.DB.First(&post, input.PostID)
    if result.Error != nil || post.UserID != userID {
        http.Error(w, "Post not found or unauthorized", http.StatusUnauthorized)
        return
    }

    // 4. Update the content
    post.Content = input.Content
    config.DB.Save(&post)

    // 5. Return updated post
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(post)
}
func DeletePost(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodDelete {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // 1. Get user_id from context
    userID := r.Context().Value("user_id").(uint)

    // 2. Decode request for post ID
    type DeleteInput struct {
        PostID uint `json:"post_id"`
    }
    var input DeleteInput
    err := json.NewDecoder(r.Body).Decode(&input)
    if err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    // 3. Find post and check ownership
    var post models.Post
    result := config.DB.First(&post, input.PostID)
    if result.Error != nil || post.UserID != userID {
        http.Error(w, "Post not found or unauthorized", http.StatusUnauthorized)
        return
    }

    // 4. Delete the post
    config.DB.Delete(&post)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"message": "Post deleted successfully"})
}
func AdminDashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Welcome to the admin dashboard!",
	})
}

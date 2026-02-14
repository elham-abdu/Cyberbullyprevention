package handlers
import (
	"net/http"
	"github.com/elham-abdu/cyberbullyprevention/config"
	"github.com/elham-abdu/cyberbullyprevention/models"
	"github.com/elham-abdu/cyberbullyprevention/utils"
	"encoding/json"
	
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
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Login successful"})
}






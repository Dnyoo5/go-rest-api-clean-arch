package controllers

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"example.com/api-GO/helpers"
	"example.com/api-GO/models"
	"example.com/api-GO/utils"
	"github.com/golang-jwt/jwt/v5"
)

type AuthController struct {
	DB *sql.DB
}

// Struct ini harus Exported (Huruf besar) biar terbaca Swagger
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

// Register godoc
// @Summary      Daftar User Baru
// @Description  Mendaftarkan user baru dengan username dan password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body models.RegisterRequest true "Payload Register"
// @Success      201      {object}  map[string]string
// @Failure      400      {object}  utils.ErrorMsg
// @Failure      500      {object}  map[string]string
// @Router       /register [post]
func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var input models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Validasi Input
	if errors := utils.ValidateStruct(input); len(errors) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors)
		return
	}

	hashedPassword, err := helpers.HashPassword(input.Password)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	query := "INSERT INTO users (username, password) VALUES (?, ?)"
	_, err = c.DB.ExecContext(ctx, query, input.Username, hashedPassword)

	if err != nil {
		// Cek duplicate entry (kalau username sudah ada)
		http.Error(w, "Gagal membuat user (Mungkin username sudah dipakai)", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

// Login godoc
// @Summary      Login User
// @Description  Masuk ke sistem untuk mendapatkan JWT Token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body LoginRequest true "Payload Login"
// @Success      200      {object}  LoginResponse
// @Failure      400      {object}  utils.ErrorMsg
// @Failure      401      {object}  map[string]string
// @Router       /login [post]
func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var input LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Validasi Input
	if errors := utils.ValidateStruct(input); len(errors) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors)
		return
	}

	var user models.User
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	query := "SELECT id, username, password FROM users WHERE username = ?"
	err := c.DB.QueryRowContext(ctx, query, input.Username).Scan(&user.Id, &user.Username, &user.Password)

	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	if !helpers.CheckPasswordHash(input.Password, user.Password) {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	claims := jwt.MapClaims{
		"user_id":  user.Id,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secretKey := []byte(os.Getenv("JWT_SECRET"))
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}
	
	// Return struct LoginResponse biar rapi di swagger
	json.NewEncoder(w).Encode(LoginResponse{Token: tokenString})
}
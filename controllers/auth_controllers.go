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

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginReponse struct {
	Token string `json:"token"`
}


func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var input models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	hashedPassword, err := helpers.HashPassword(input.Password)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	if errors := utils.ValidateStruct(input); len(errors) > 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errors)
				return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	query := "INSERT INTO users (username, password) VALUES (?, ?)"
	_, err = c.DB.ExecContext(ctx, query, input.Username, hashedPassword)

	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var input LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var user models.User

	if errors := utils.ValidateStruct(input); len(errors) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errors)
		return
	}

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
	json.NewEncoder(w).Encode(LoginReponse{Token: tokenString})
}
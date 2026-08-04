package handlers

import (
    "net/http"
    "github.com/JadhavSrushti1544/TaskForge/internal/config"
    "github.com/JadhavSrushti1544/TaskForge/internal/models"
    "github.com/JadhavSrushti1544/TaskForge/internal/repositories"
    "github.com/JadhavSrushti1544/TaskForge/internal/service"
    "github.com/gin-gonic/gin"
)

// Authentication Handler struct
type AuthHandler struct {
    userRepo *repositories.UserRepository
}

// Authentication Handler function to create new instance of AuthHandler
func NewAuthHandler(userRepo *repositories.UserRepository) *AuthHandler {
    return &AuthHandler{userRepo: userRepo}
}

// Register function to create new user
func (h *AuthHandler) Register(c *gin.Context) {
    var req models.RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    // Check if email already exists
    existing, _ := h.userRepo.GetUserByEmail(req.Email)
    if existing != nil {
        c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
        return
    }
    
    // Hash password
    hash, err := service.HashPassword(req.Password)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
        return
    }
    
    // Create user
    user, err := h.userRepo.CreateUser(req.Email, hash)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
        return
    }
    
    c.JSON(http.StatusCreated, user)
}

// Login function
func (h *AuthHandler) Login(c *gin.Context) {
    var req models.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    // Find user by email
    user, err := h.userRepo.GetUserByEmail(req.Email)
    if err != nil || user == nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }
    
    // Verify password 
    if !service.VerifyPassword(user.PasswordHash, req.Password) {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }
    
    // Generate token 
    token, err := config.GenerateToken(user.ID, user.Email)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }
    
    c.JSON(http.StatusOK, models.LoginResponse{
        Token: token,
        User:  *user,
    })
}

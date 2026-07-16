package middleware

import (
    "strings"
    "task-queue/config"
    "github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(401, gin.H{"error": "Missing authorization header"})
            c.Abort()
            return
        }
        
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.JSON(401, gin.H{"error": "Invalid authorization format"})
            c.Abort()
            return
        }
        
        userID, err := config.VerifyToken(parts[1])
        if err != nil {
            c.JSON(401, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }
        
        c.Set("user_id", userID)
        c.Next()
    }
}
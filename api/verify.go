package api

import (
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/rick/fly_crypto/internal/token"
)

type verifyRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required"`
}

func (a *App) Verify() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req verifyRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Get stored verification code
		storedCode, err := a.GetVerificationCode(req.Email)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired verification code"})
			return
		}

		// Verify code
		if storedCode != req.Code {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid verification code"})
			return
		}

		// Generate JWT token
		token, err := token.Generate(req.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		// Store token in Redis with 24-hour expiration
		if err := a.StoreAuthToken(token, req.Email); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store token"})
			return
		}

		// Delete verification code after successful verification
		_ = a.DeleteVerificationCode(req.Email)

		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}
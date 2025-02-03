package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rickj1ang/fly_crypto/internal/app"
	"github.com/rickj1ang/fly_crypto/internal/mail"
	"github.com/rickj1ang/fly_crypto/internal/token"
)

type loginRequest struct {
	Email string `json:"email" binding:"required,email"`
}

func Login(a *app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req loginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Generate verification code
		code := token.GenerateVerificationCode()
		if code == "000000" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate verification code"})
			return
		}

		// Store verification code in Redis with 6-minute expiration
		if err := a.Data.StoreVerificationCode(req.Email, code); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store verification code"})
			return
		}

		// Send verification code via email
		if err := mail.SendVerifyCode(req.Email, code); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send verification code"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Verification code sent"})
	}
}

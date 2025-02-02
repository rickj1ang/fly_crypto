package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rickj1ang/fly_crypto/internal/app"
	"github.com/rickj1ang/fly_crypto/internal/data"
	"github.com/rickj1ang/fly_crypto/internal/token"
)

type verifyRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required"`
}

func Verify(a *app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req verifyRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Get stored verification email
		storedEmail, err := a.GetEmailByVerifyCode(req.Code)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired verification code"})
			return
		}

		// Verify email
		if storedEmail != req.Email {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid verification code"})
			return
		}

		// Check if user exists in database
		exists, err := a.Data.CheckUser(req.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check user"})
			return
		}

		// Create user if not exists
		var userID int64
		if !exists {
			user := &data.User{
				Email:           req.Email,
				NotificationsID: []int64{},
			}
			if err := a.Data.CreateUser(user); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
				return
			} else {
				userID = user.UserID
			}
		} else {
			userID, err = a.Data.GetUserIDByEmail(req.Email)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user ID"})
				return
			}
		}

		token, err := token.Generate()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		// Store token in Redis with 24-hour expiration
		if err := a.StoreAuthToken(token, userID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store token"})
			return
		}

		// Delete verification code after successful verification
		_ = a.Data.DeleteVerificationCode(req.Code)

		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}

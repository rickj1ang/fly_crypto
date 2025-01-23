package api

import (
	"time"
)

// StoreVerificationCode stores a verification code in Redis
func (a *App) StoreVerificationCode(email, code string) error {
	return a.data.StoreVerificationCode(email, code)
}

// StoreAuthToken stores an authentication token in Redis
func (a *App) StoreAuthToken(token, email string) error {
	return a.data.StoreAuthToken(token, email)
}



// GetEmailByAuthToken retrieves an email associated with an authentication token from Redis
func (a *App) GetEmailByAuthToken(token string) (string, error) {
	return a.data.GetEmailByAuthToken(token)
}

// DeleteVerificationCode removes a verification code from Redis
func (a *App) DeleteVerificationCode(email string) error {
	return a.data.DeleteVerificationCode(email)
}

// DeleteAuthToken removes an authentication token from Redis
func (a *App) DeleteAuthToken(token string) error {
	return a.data.DeleteAuthToken(token)
}

// StoreToken delegates to Data.StoreToken
func (a *App) StoreToken(token, email string, expiration time.Duration) error {
	return a.data.StoreToken(token, email, expiration)
}

// GetVerifyCodeByEmail delegates to Data.GetVerifyCodeByEmail
func (a *App) GetVerifyCodeByEmail(token string) (string, error) {
	return a.data.GetVerifyCodeByEmail(token)
}

// DeleteToken delegates to Data.DeleteToken
func (a *App) DeleteToken(token string) error {
	return a.data.DeleteToken(token)
}
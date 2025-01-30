package app

// StoreVerificationCode stores a verification code in Redis
func (a *App) StoreVerificationCode(email, code string) error {
	return a.Data.StoreVerificationCode(email, code)
}

// StoreAuthToken stores an authentication token in Redis
func (a *App) StoreAuthToken(token string, userID int64) error {
	return a.Data.StoreAuthToken(token, userID)
}

// GetEmailByAuthToken retrieves an email associated with an authentication token from Redis
func (a *App) GetUserIDByAuthToken(token string) (int64, error) {
	return a.Data.GetUserIDByAuthToken(token)
}

// DeleteVerificationCode removes a verification code from Redis
func (a *App) DeleteVerificationCode(email string) error {
	return a.Data.DeleteVerificationCode(email)
}

// DeleteAuthToken removes an authentication token from Redis
func (a *App) DeleteAuthToken(token string) error {
	return a.Data.DeleteAuthToken(token)
}

// GetVerifyCodeByEmail delegates to Data.GetVerifyCodeByEmail
func (a *App) GetEmailByVerifyCode(token string) (string, error) {
	return a.Data.GetEmailByVerificationCode(token)
}

// DeleteToken delegates to Data.DeleteToken
func (a *App) DeleteToken(token string) error {
	return a.Data.DeleteToken(token)
}

func (a *App) StoreOnesNotification(key string, price float64, id int64) error {
	email, err := a.Data.GetUserEmail(id)
	if err != nil {
		return err
	}
	return a.Data.StoreInSortedSet(key, price, email)
}

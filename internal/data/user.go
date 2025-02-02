package data

import "context"

// User represents a user in the database
type User struct {
	UserID          int64   `json:"user_id" db:"user_id"`
	Email           string  `json:"email" db:"email"`
	NotificationsID []int64 `json:"notifications_id" db:"notifications_id"`
}

// CreateUser creates a new user in the database
func (d Data) CreateUser(u *User) error {
	ctx := context.TODO()

	// Insert the user into the database
	query := `
		INSERT INTO users (email)
		VALUES ($1)
		RETURNING user_id
	`

	err := d.Db.QueryRowContext(
		ctx,
		query,
		u.Email,
	).Scan(&u.UserID)

	if err != nil {
		return err
	}

	return nil
}

// CheckUser checks if a user exists by email
func (d Data) CheckUser(email string) (bool, error) {
	ctx := context.TODO()

	// Query the user from the database
	query := `
		SELECT EXISTS(
			SELECT 1
			FROM users
			WHERE email = $1
		)
	`

	var exists bool
	err := d.Db.QueryRowContext(
		ctx,
		query,
		email,
	).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}

// GetUserEmail retrieves a user's email address by their user ID
func (d Data) GetUserEmail(userID int64) (string, error) {
	ctx := context.TODO()

	// Query to get the user's email
	query := `
		SELECT email
		FROM users
		WHERE user_id = $1
	`

	var email string
	err := d.Db.QueryRowContext(
		ctx,
		query,
		userID,
	).Scan(&email)

	if err != nil {
		return "", err
	}

	return email, nil
}

// DeleteUser deletes a user and all their associated notifications from the database
func (d Data) DeleteUser(userID int64) error {
	ctx := context.TODO()

	// Start a transaction
	tx, err := d.Db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// First, delete all notifications associated with the user
	deleteNotificationsQuery := `
		DELETE FROM notifications
		WHERE user_id = $1
	`

	_, err = tx.ExecContext(ctx, deleteNotificationsQuery, userID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Then delete the user
	deleteUserQuery := `
		DELETE FROM users
		WHERE user_id = $1
	`

	_, err = tx.ExecContext(ctx, deleteUserQuery, userID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	return tx.Commit()
}

func (d Data) GetUserIDByEmail(email string) (int64, error) {
	ctx := context.TODO()
	// Query to get the user ID by email
	query := `
		SELECT user_id
		FROM users
		WHERE email = $1
	`
	var userID int64
	err := d.Db.QueryRowContext(
		ctx,
		query,
		email,
	).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

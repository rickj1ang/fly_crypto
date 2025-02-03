package data

import (
	"context"
)

type Notification struct {
	ID          int64   `json:"notification_id" db:"notification_id"`
	UserID      int64   `json:"user_id" db:"user_id"`
	CoinSymbol  string  `json:"coin_symbol" db:"coin_symbol"`
	TargetPrice float64 `json:"target_price" db:"target_price"`
	IsAbove     bool    `json:"is_above" db:"is_above"`
}

// CreateNotification creates a new price alert notification and updates user's notifications array
func (d Data) CreateNotification(n *Notification) error {
	ctx := context.TODO()

	// Start a transaction
	tx, err := d.Db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Insert into notifications table
	insertQuery := `
		INSERT INTO notifications (user_id, coin_symbol, target_price, is_above)
		VALUES ($1, $2, $3, $4)
		RETURNING notification_id
	`

	err = tx.QueryRowContext(
		ctx,
		insertQuery,
		n.UserID,
		n.CoinSymbol,
		n.TargetPrice,
		n.IsAbove,
	).Scan(&n.ID)

	if err != nil {
		tx.Rollback()
		return err
	}

	// Update user's notifications array
	updateQuery := `
		UPDATE users
		SET notifications_id = array_append(notifications_id, $1)
		WHERE user_id = $2
	`

	_, err = tx.ExecContext(ctx, updateQuery, n.ID, n.UserID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	return tx.Commit()
}

// DeleteNotification deletes a notification by its ID and removes it from user's notifications array
func (d Data) DeleteNotification(notificationID int64) error {
	ctx := context.TODO()

	// Start a transaction
	tx, err := d.Db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// First, update the user's notifications array to remove this notification ID
	updateQuery := `
		UPDATE users
		SET notifications_id = array_remove(notifications_id, $1)
		WHERE $1 = ANY(notifications_id)
	`

	_, err = tx.ExecContext(ctx, updateQuery, notificationID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Then delete the notification from notifications table
	deleteQuery := `
		DELETE FROM notifications
		WHERE notification_id = $1
	`

	_, err = tx.ExecContext(ctx, deleteQuery, notificationID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	return tx.Commit()
}

// DeleteNotificationFromMessage deletes a notification based on user email, coin symbol and price direction
func (d Data) DeleteNotificationFromMessage(userEmail string, isAbove bool, coinSymbol string) error {
	ctx := context.TODO()

	// Start a transaction
	tx, err := d.Db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// First get the notification ID and user ID
	var notificationID int64
	var userID int64

	query := `
		SELECT n.notification_id, n.user_id
		FROM notifications n
		JOIN users u ON n.user_id = u.user_id
		WHERE u.email = $1 AND n.is_above = $2 AND n.coin_symbol = $3
		LIMIT 1
	`

	err = tx.QueryRowContext(ctx, query, userEmail, isAbove, coinSymbol).Scan(&notificationID, &userID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Update user's notifications array to remove this notification ID
	updateQuery := `
		UPDATE users
		SET notifications_id = array_remove(notifications_id, $1)
		WHERE user_id = $2
	`

	_, err = tx.ExecContext(ctx, updateQuery, notificationID, userID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Delete the notification from notifications table
	deleteQuery := `
		DELETE FROM notifications
		WHERE notification_id = $1
	`

	_, err = tx.ExecContext(ctx, deleteQuery, notificationID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	return tx.Commit()
}

// GetUserAllNotifications retrieves all notifications for a given user ID,
// sorted by coin symbol for easier review
func (d Data) GetUserAllNotifications(userID int64) ([]Notification, error) {
	ctx := context.TODO()

	// Query to get all notifications for the user, ordered by coin symbol
	query := `
		SELECT notification_id, user_id, coin_symbol, target_price, is_above
		FROM notifications
		WHERE user_id = $1
		ORDER BY coin_symbol
	`

	rows, err := d.Db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []Notification
	for rows.Next() {
		var n Notification
		err := rows.Scan(
			&n.ID,
			&n.UserID,
			&n.CoinSymbol,
			&n.TargetPrice,
			&n.IsAbove,
		)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return notifications, nil
}

func (d Data) GetANotification(id int64) (*Notification, error) {
	ctx := context.TODO()

	query := `
		SELECT notification_id, user_id, coin_symbol, target_price, is_above
		FROM notifications
		WHERE notification_id = $1`
	var n Notification

	err := d.Db.QueryRowContext(ctx, query, id).Scan(
		&n.ID,
		&n.UserID,
		&n.CoinSymbol,
		&n.TargetPrice,
		&n.IsAbove,
	)
	if err != nil {
		return nil, err
	}

	return &n, nil

}

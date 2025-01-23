package data

import (
	"context"
	"time"
)

type Notification struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	CoinSymbol  string    `json:"coin_symbol"`
	TargetPrice float64   `json:"target_price"`
	IsAbove     bool      `json:"is_above"`
}

// CreateNotification creates a new price alert notification
func (a *App) CreateNotification(n *Notification) error {
	query := `
		INSERT INTO notifications (user_id, coin_symbol, target_price, is_above)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`
	ctx := context.TODO()
	
	return a.db.QueryRowContext(
		ctx,
		query,
		n.UserID,
		n.CoinSymbol,
		n.TargetPrice,
		n.IsAbove,
	).Scan(&n.ID, &n.CreatedAt)
}

// GetUserNotifications retrieves all notifications for a specific user
func (a *App) GetUserNotifications(ctx context.Context, userID int64) ([]Notification, error) {
	query := `
		SELECT id, user_id, coin_symbol, target_price, is_above, created_at, triggered_at, is_triggered
		FROM notifications
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := a.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []Notification
	for rows.Next() {
		var n Notification
		var triggeredAt *time.Time

		err := rows.Scan(
			&n.ID,
			&n.UserID,
			&n.CoinSymbol,
			&n.TargetPrice,
			&n.IsAbove,
			&n.CreatedAt,
			&triggeredAt,
			&n.IsTriggered,
		)
		if err != nil {
			return nil, err
		}

		if triggeredAt != nil {
			n.TriggeredAt = *triggeredAt
		}

		notifications = append(notifications, n)
	}

	return notifications, nil
}
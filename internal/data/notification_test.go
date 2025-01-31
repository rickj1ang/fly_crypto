package data

import (
	"database/sql"
	"os"
	"testing"

	"github.com/go-redis/redis/v8"

	_ "github.com/lib/pq"
)

func setupData(t *testing.T) *Data {

	dbURI := os.Getenv("DATABASE_URL")
	if dbURI == "" {
		t.Skip("DATABASE_URL not set")
	}
	db, err := sql.Open("postgres", dbURI)
	if err != nil {
		t.Error("fail to conect database")
	}

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		t.Error("fail to get REDIS_URL")
	}

	// 创建 Redis 客户端
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		t.Errorf("解析 Redis URL 失败: %v", err)
	}

	rdb := redis.NewClient(opt)

	if err != nil {
		t.Fatalf("Failed to connect to Redis: %v", err)
	}

	return &Data{Redis: rdb, Db: db}
}

func TestCreateNotification(t *testing.T) {
	user := &User{
		Email: "sad1",
	}
	data := setupData(t)
	err := data.CreateUser(user)
	if err != nil {
		t.Errorf("fail to create user: %s", err.Error())
	}

	notification := &Notification{
		UserID:      user.UserID,
		CoinSymbol:  "BTC",
		TargetPrice: 10,
		IsAbove:     true,
	}
	err = data.CreateNotification(notification)
	if err != nil {
		t.Errorf("fail to create notification: %s", err.Error())
	}

	notificationRet, err := data.GetANotification(notification.ID)
	if err != nil {
		t.Errorf("fail to get notification: %s", err.Error())
	}
	if notificationRet.ID != notification.ID {
		t.Errorf("fail to get notification")
	}
	if notificationRet.CoinSymbol != notification.CoinSymbol {
		t.Errorf("fail to get notification")
	}

	notifications, err := data.GetUserAllNotifications(notification.UserID)
	if err != nil {
		t.Errorf("fail to get user notifications")
	}
	if len(notifications) != 1 {
		t.Errorf("fail to get user notifications")
	}
	if notifications[0].ID != notification.ID {
		t.Errorf("fail to get user notifications")
	}

	err = data.DeleteNotification(notification.ID)
	if err != nil {
		t.Errorf("fail to delete notification")
	}

	notif, err := data.GetANotification(notification.ID)
	if err == nil {
		t.Errorf("fail to delete notification")
	}
	if notif != nil {
		t.Errorf("fail to delete notification")
	}

	notifis, err := data.GetUserAllNotifications(notification.UserID)
	if err != nil {
		t.Errorf("fail to delete notification")
	}
	if len(notifis) != 0 {
		t.Errorf("fail to delete notification")
	}

	err = data.DeleteUser(user.UserID)
	if err != nil {
		t.Errorf("fail to delete user")
	}
}

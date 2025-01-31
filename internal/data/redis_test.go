package data

import (
	"context"
	"os"
	"testing"

	"github.com/go-redis/redis/v8"
)

func setupTestRedis(t *testing.T) *Data {
	// Connect to Redis using environment variable or default to localhost
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

	return &Data{Redis: rdb}
}

func TestVerificationCodeOperations(t *testing.T) {
	d := setupTestRedis(t)
	ctx := context.Background()
	defer d.Redis.Close()

	// Test data
	testEmail := "test@example.com"
	testCode := "123456"
	testAuth := "dasdfafwvasd"
	var testUserID int64 = 12324124

	// Test StoreVerificationCode
	err := d.StoreVerificationCode(testEmail, testCode)
	if err != nil {
		t.Errorf("Failed to store verification code: %v", err)
	}

	// Test GetEmailByVerificationCode
	email, err := d.GetEmailByVerificationCode(testCode)
	if err != nil {
		t.Errorf("Failed to get email by verification code: %v", err)
	}
	if email != testEmail {
		t.Errorf("Retrieved email does not match. Expected %s, got %s", testEmail, email)
	}

	// Test DeleteVerificationCode
	err = d.DeleteVerificationCode(testCode)
	if err != nil {
		t.Errorf("Failed to delete verification code: %v", err)
	}

	// Verify deletion by trying to get the code again
	_, err = d.GetEmailByVerificationCode(testCode)
	if err != redis.Nil {
		t.Errorf("Verification code should have been deleted, but got: %v", err)
	}

	// Cleanup: ensure the key is deleted even if the test fails
	key := VerifyPrefix + testCode
	d.Redis.Del(ctx, key)

	err = d.StoreAuthToken(testAuth, testUserID)
	if err != nil {
		t.Errorf("Failed to store auth token: %v", err)
	}

	// Test GetUserIDByAuthToken
	userID, err := d.GetUserIDByAuthToken(testAuth)
	if err != nil {
		t.Errorf("Failed to get user ID by auth token: %v", err)
	}
	if userID != testUserID {
		t.Errorf("Retrieved user ID does not match. Expected %d, got %d", testUserID, userID)
	}

	// Test DeleteAuthToken
	err = d.DeleteAuthToken(testAuth)
	if err != nil {
		t.Errorf("Failed to delete auth token: %v", err)
	}

	// Verify deletion by trying to get the token again
	_, err = d.GetUserIDByAuthToken(testAuth)
	if err != redis.Nil {
		t.Errorf("Auth token should have been deleted, but got: %v", err)
	}
	// Cleanup: ensure the key is deleted even if the test fails
	key = AuthPrefix + testAuth
	d.Redis.Del(ctx, key)
}

func TestSortedSet(t *testing.T) {
	d := setupTestRedis(t)
	defer d.Redis.Close()

	//test data struct
	type TestData struct {
		Key    string
		Score  float64
		member string
	}
	//test data
	testData := []TestData{
		{"BTC", 1.0, "email1"},
		{"BTC", 6.0, "email2"},
		{"BTC", 3.0, "email3"},
		{"BTC", 10.0, "email4"},
		{"BTC", 5.0, "email5"},
		{"BTC", 4.0, "email6"},
	}
	//test StoreSortedSet
	for _, data := range testData {
		err := d.StoreInSortedSet(data.Key, data.Score, data.member)
		if err != nil {
			t.Errorf("Failed to store sorted set: %v", err)
		}
	}

	score, email, err := d.GetMinScoreFromSortedSet("BTC")
	if err != nil {
		t.Errorf("Failed to get min score from sorted set: %v", err)
	}
	if email != "email1" {
		t.Errorf("Retrieved email does not match. Expected %s, got %s", "email1", email)
	}
	if score != 1.0 {
		t.Errorf("Retrieved score does not match. Expected %f, got %f", 1.0, score)
	}

	score2, email2, err := d.GetMaxScoreFromSortedSet("BTC")
	if err != nil {
		t.Errorf("Failed to get max score from sorted set: %v", err)
	}
	if email2 != "email4" {
		t.Errorf("Retrieved email does not match. Expected %s, got %s", "email4", email2)
	}
	if score2 != 10.0 {
		t.Errorf("Retrieved score does not match. Expected %f, got %f", 10.0, score2)
	}

	minScore, minEmail, err := d.PopMinFromSortedSet("BTC")
	if err != nil {
		t.Errorf("Failed to pop min from sorted set: %v", err)
	}
	if minEmail != "email1" {
		t.Errorf("Retrieved email does not match. Expected %s, got %s", "email1", minEmail)
	}
	if minScore != 1.0 {
		t.Errorf("Retrieved score does not match. Expected %f, got %f", 1.0, minScore)
	}

	socre3, email3, err := d.GetMinScoreFromSortedSet("BTC")
	if err != nil {
		t.Errorf("Failed to get min score from sorted set: %v", err)
	}
	if socre3 != 3.0 {
		t.Errorf("Retrieved score does not match. Expected %f, got %f", 3.0, socre3)
	}
	if email3 != "email3" {
		t.Errorf("Retrieved email does not match. Expected %s, got %s", "email3", email3)
	}

	maxScore, maxEmail, err := d.PopMaxFromSortedSet("BTC")
	if err != nil {
		t.Errorf("Failed to pop max from sorted set: %v", err)
	}
	if maxScore != 10.0 {
		t.Errorf("Retrieved score does not match. Expected %f, got %f", 10.0, maxScore)
	}
	if maxEmail != "email4" {
		t.Errorf("Retrieved email does not match. Expected %s, got %s", "email4", maxEmail)
	}

	score5, email5, err := d.GetMaxScoreFromSortedSet("BTC")
	if err != nil {
		t.Errorf("Failed to get max score from sorted set: %v", err)
	}
	if email5 != "email2" {
		t.Errorf("Retrieved email does not match. Expected %s, got %s", "email5", email5)
	}
	if score5 != 6.0 {
		t.Errorf("Retrieved score does not match. Expected %f, got %f", 5.0, score5)
	}

	//cleanup
	for _, data := range testData {
		err := d.RemoveFromSortedSet(data.Key, data.member)
		if err != nil {
			t.Errorf("Failed to delete from sorted set: %v", err)
		}
	}

}

package data

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	VerifyPrefix     = "verify:"
	AuthPrefix       = "auth:"
	VerifyExpiration = 6 * time.Minute
	AuthExpiration   = 24 * time.Hour
)

// StoreVerificationCode stores a verification code in Redis with verify: prefix and 6-minute expiration
func (d *Data) StoreVerificationCode(email, code string) error {
	ctx := context.Background()
	key := fmt.Sprintf("%s%s", VerifyPrefix, code)
	return d.Redis.Set(ctx, key, email, VerifyExpiration).Err()
}

// StoreAuthToken stores an authentication token in Redis with auth: prefix and 24-hour expiration
func (d *Data) StoreAuthToken(token string, userID int64) error {
	ctx := context.Background()
	key := fmt.Sprintf("%s%s", AuthPrefix, token)
	return d.Redis.Set(ctx, key, userID, AuthExpiration).Err()
}

// GetEmailByAuthToken retrieves the email associated with an authentication token
func (d *Data) GetUserIDByAuthToken(token string) (int64, error) {
	ctx := context.Background()
	key := fmt.Sprintf("%s%s", AuthPrefix, token)
	userID, err := d.Redis.Get(ctx, key).Result()
	if err != nil {
		return -1, err
	}
	userIDInt, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return -1, err
	}
	return userIDInt, nil
}

// DeleteVerificationCode removes a verification code from Redis
func (d *Data) DeleteVerificationCode(code string) error {
	ctx := context.Background()
	key := fmt.Sprintf("%s%s", VerifyPrefix, code)
	return d.Redis.Del(ctx, key).Err()
}

// DeleteAuthToken removes an authentication token from Redis
func (d *Data) DeleteAuthToken(token string) error {
	ctx := context.Background()
	key := fmt.Sprintf("%s%s", AuthPrefix, token)
	return d.Redis.Del(ctx, key).Err()
}

// GetEmailByVerificationCode retrieves the email associated with a verification code
func (d *Data) GetEmailByVerificationCode(code string) (string, error) {
	ctx := context.Background()
	key := fmt.Sprintf("%s%s", VerifyPrefix, code)
	email, err := d.Redis.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return email, nil
}

// DeleteToken removes a token from Redis
func (d *Data) DeleteToken(token string) error {
	ctx := context.Background()
	return d.Redis.Del(ctx, token).Err()
}

// StoreInSortedSet adds a member with a score to a sorted set in Redis
func (d *Data) StoreInSortedSet(key string, score float64, member interface{}) error {
	ctx := context.Background()
	return d.Redis.ZAdd(ctx, key, &redis.Z{
		Score:  score,
		Member: member,
	}).Err()
}

// GetMinScoreFromSortedSet retrieves the member with the lowest score from a sorted set
func (d *Data) GetMinScoreFromSortedSet(key string) (float64, interface{}, error) {
	ctx := context.Background()
	// Get the member with the lowest score (ZRANGE with index 0,0)
	result, err := d.Redis.ZRangeWithScores(ctx, key, 0, 0).Result()
	if err != nil {
		return 0, nil, err
	}
	if len(result) == 0 {
		return 0, nil, redis.Nil
	}
	return result[0].Score, result[0].Member, nil
}

// GetMaxScoreFromSortedSet retrieves the member with the highest score from a sorted set
func (d *Data) GetMaxScoreFromSortedSet(key string) (float64, interface{}, error) {
	ctx := context.Background()
	// Get the member with the highest score (ZRANGE with index -1,-1)
	result, err := d.Redis.ZRangeWithScores(ctx, key, -1, -1).Result()
	if err != nil {
		return 0, nil, err
	}
	if len(result) == 0 {
		return 0, nil, redis.Nil
	}
	return result[0].Score, result[0].Member, nil
}

// PopMinFromSortedSet removes and returns the member with the lowest score from a sorted set
func (d *Data) PopMinFromSortedSet(key string) (float64, interface{}, error) {
	ctx := context.Background()
	// Use ZPOPMIN to atomically remove and return the member with lowest score
	result, err := d.Redis.ZPopMin(ctx, key).Result()
	if err != nil {
		return 0, nil, err
	}
	if len(result) == 0 {
		return 0, nil, redis.Nil
	}
	return result[0].Score, result[0].Member, nil
}

// PopMaxFromSortedSet removes and returns the member with the highest score from a sorted set
func (d *Data) PopMaxFromSortedSet(key string) (float64, interface{}, error) {
	ctx := context.Background()
	// Use ZPOPMAX to atomically remove and return the member with highest score
	result, err := d.Redis.ZPopMax(ctx, key).Result()
	if err != nil {
		return 0, nil, err
	}
	if len(result) == 0 {
		return 0, nil, redis.Nil
	}
	return result[0].Score, result[0].Member, nil
}

// RemoveFromSortedSet removes a specific member from a sorted set
func (d *Data) RemoveFromSortedSet(key string, member interface{}) error {
	ctx := context.Background()
	return d.Redis.ZRem(ctx, key, member).Err()
}

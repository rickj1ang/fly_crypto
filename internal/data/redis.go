package data

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

const (
	VerifyPrefix = "verify:"
	AuthPrefix   = "auth:"
	VerifyExpiration = 6 * time.Minute
	AuthExpiration = 24 * time.Hour
)

// StoreVerificationCode stores a verification code in Redis with verify: prefix and 6-minute expiration
func (d *Data) StoreVerificationCode(email, code string) error {
	ctx := context.Background()
	key := fmt.Sprintf("%s%s", VerifyPrefix, code)
	return d.redis.Set(ctx, key, email, VerifyExpiration).Err()
}

// StoreAuthToken stores an authentication token in Redis with auth: prefix and 24-hour expiration
func (d *Data) StoreAuthToken(token, email string) error {
	ctx := context.Background()
	key := fmt.Sprintf("%s%s", AuthPrefix, token)
	return d.redis.Set(ctx, key, email, AuthExpiration).Err()
}



// GetEmailByAuthToken retrieves the email associated with an authentication token
func (d *Data) GetEmailByAuthToken(token string) (string, error) {
	ctx := context.Background()
	key := fmt.Sprintf("%s%s", AuthPrefix, token)
	email, err := d.redis.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return email, nil
}

// DeleteVerificationCode removes a verification code from Redis
func (d *Data) DeleteVerificationCode(email string) error {
	ctx := context.Background()
	key := fmt.Sprintf("%s%s", VerifyPrefix, email)
	return d.redis.Del(ctx, key).Err()
}

// DeleteAuthToken removes an authentication token from Redis
func (d *Data) DeleteAuthToken(token string) error {
	ctx := context.Background()
	key := fmt.Sprintf("%s%s", AuthPrefix, token)
	return d.redis.Del(ctx, key).Err()
}

// GetEmailByVerificationCode retrieves the email associated with a verification code
func (d *Data) GetEmailByVerificationCode(code string) (string, error) {
	ctx := context.Background()
	key := fmt.Sprintf("%s%s", VerifyPrefix, code)
	email, err := d.redis.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return email, nil
}

// DeleteToken removes a token from Redis
func (d *Data) DeleteToken(token string) error {
	ctx := context.Background()
	return d.redis.Del(ctx, token).Err()
}
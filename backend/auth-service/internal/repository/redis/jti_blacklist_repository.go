package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const jtiBlacklistPrefix = "jti_blacklist:"

type JtiBlacklistRepository struct {
	client *redis.Client
}

func NewJtiBlacklistRepository(client *redis.Client) *JtiBlacklistRepository {
	return &JtiBlacklistRepository{client: client}
}

func (r *JtiBlacklistRepository) AddToBlacklist(jti string, expiresAt time.Time) error {
	key := jtiBlacklistPrefix + jti
	duration := time.Until(expiresAt)

	if duration <= 0 {
		// Token already expired, no need to blacklist for long or at all
		// Or blacklist for a very short period to handle clock skew issues
		duration = 5 * time.Minute
	}

	err := r.client.Set(context.Background(), key, "blacklisted", duration).Err()
	if err != nil {
		return fmt.Errorf("failed to add JTI to Redis blacklist: %w", err)
	}
	return nil
}

func (r *JtiBlacklistRepository) IsBlacklisted(jti string) (bool, error) {
	key := jtiBlacklistPrefix + jti
	val, err := r.client.Get(context.Background(), key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil // Not found means not blacklisted
		}
		return false, fmt.Errorf("failed to check JTI in Redis blacklist: %w", err)
	}
	return val == "blacklisted", nil
}

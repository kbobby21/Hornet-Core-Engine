package redis

import (
	"context"
	"crypto/sha256"
	"encoding/base64"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

type Redis struct {
	client *redis.Client
}

func NewRedis() *Redis {
	client := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis_addr"),     // "localhost:6379"
		Password: viper.GetString("redis_password"), // no password set
		DB:       0,                                 // use default DB
	})
	return &Redis{
		client: client,
	}
}

func (r *Redis) StoreEmailHash(email string) (string, error) {
	// create a sha256 hash of this email
	h := sha256.New()
	h.Write([]byte(email))
	bs := base64.URLEncoding.EncodeToString(h.Sum(nil))

	// store hash:email
	return string(bs), r.client.Set(context.Background(), string(bs), email, 0).Err()
}

// fetch email from hash and delete the key
func (r *Redis) GetEmailFromHash(hash string) (string, error) {
	ctx := context.Background()
	email, err := r.client.Get(ctx, hash).Result()
	if err != nil || len(email) == 0 {
		return "", err
	}
	// now delete the key
	err = r.client.Del(ctx, hash).Err()
	return email, err
}

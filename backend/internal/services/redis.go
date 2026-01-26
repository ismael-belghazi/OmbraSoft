package services

import (
	"context"
	"log"
	"time"

	"github.com/bebeb/ombrasoft-backend/internal/config"
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis() error {
	cfg := config.AppConfig

	if cfg.RedisURL == "" {
		log.Println("Redis désactivé (REDIS_URL vide)")
		return nil
	}

	opt, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		log.Printf("Erreur parsing REDIS_URL: %v", err)
		return err
	}

	RedisClient = redis.NewClient(opt)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := RedisClient.Ping(ctx).Err(); err != nil {
		log.Printf("Connexion Redis échouée: %v", err)
		return err
	}

	log.Println("Connecté à Redis")
	return nil
}

func CloseRedis() error {
	if RedisClient != nil {
		return RedisClient.Close()
	}
	return nil
}

func CacheSet(ctx context.Context, key string, value string, ttl time.Duration) error {
	if RedisClient == nil {
		return nil
	}
	return RedisClient.Set(ctx, key, value, ttl).Err()
}

func CacheGet(ctx context.Context, key string) (string, error) {
	if RedisClient == nil {
		return "", nil
	}
	return RedisClient.Get(ctx, key).Result()
}

func CacheDel(ctx context.Context, keys ...string) error {
	if RedisClient == nil {
		return nil
	}
	return RedisClient.Del(ctx, keys...).Err()
}

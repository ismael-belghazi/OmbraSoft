package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ismael-belghazi/ombrasoft-backend/internal/config"
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

func CacheSet(ctx context.Context, key, value string, ttl time.Duration) error {
	if RedisClient == nil {
		return fmt.Errorf("redis client not initialized")
	}

	if ttl == 0 {
		ttl = 1 * time.Hour
	}

	err := RedisClient.Set(ctx, key, value, ttl).Err()
	if err != nil {
		log.Printf("Erreur lors de la mise en cache dans Redis: %v", err)
	}
	return err
}

func CacheGet(ctx context.Context, key string) (string, error) {
	if RedisClient == nil {
		return "", fmt.Errorf("redis client not initialized")
	}

	value, err := RedisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		log.Printf("Erreur lors de la récupération du cache Redis: %v", err)
		return "", err
	}
	return value, nil
}

func CacheDel(ctx context.Context, keys ...string) error {
	if RedisClient == nil {
		return fmt.Errorf("redis client not initialized")
	}

	err := RedisClient.Del(ctx, keys...).Err()
	if err != nil {
		log.Printf("Erreur lors de la suppression de cache Redis: %v", err)
	}
	return err
}

package cache

import (
	"fmt"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

// RedisConnection func for connect to Redis server.
func RedisConnection() (*redis.Client, error) {
	envFile, err := godotenv.Read(".env")
	if err != nil {
		return nil, fmt.Errorf("can't read env file")
	}
	// Define Redis database number.
	dbNumber, _ := strconv.Atoi(envFile["REDIS_DB_NUMBER"])

	// Build Redis connection URL.
	redisConnURL := fmt.Sprintf("%s:%s", envFile["REDIS_HOST"], envFile["REDIS_PORT"])

	// Set Redis options.
	options := &redis.Options{
		Addr:     redisConnURL,
		Password: envFile["REDIS_PASSWORD"],
		DB:       dbNumber,
	}

	return redis.NewClient(options), nil
}
